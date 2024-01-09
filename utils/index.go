package utils

import (
	"archive/tar"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error Defer", err.Error())
		}
	}(file)

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// Unzip
// zipFilePath := "path/to/your/file.zip"
// destinationDir := "path/to/your/destination"
//
//	if err := unzip(zipFilePath, destinationDir); err != nil {
//		fmt.Println("Error:", err)
//	} else {
//		fmt.Println("Unzip successful.")
//	}
func Unzip(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		target := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
			file.Close()
		default:
			return fmt.Errorf("Unknown type: %v in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

// GetSubdirectories
// directoryPath := "/path/to/your/directory"
// subdirectories, err := getSubdirectories(directoryPath)
// if err != nil {
// fmt.Println("Error:", err)
// return
// }
// fmt.Println("Subdirectories:")
// for _, subdir := range subdirectories {
// fmt.Println(subdir)
// }
func GetSubdirectories(directoryPath string) ([]string, error) {
	var subdirectories []string

	entries, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subdirectories = append(subdirectories, entry.Name())
		}
	}

	return subdirectories, nil
}

func VisitTgzS(archiveFiles *[]string, serverName string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if f.IsDir() {
			return nil // not a file, ignore.
		}

		// Check if the file has a .tar.gz extension
		if strings.HasSuffix(path, ".tar.gz") {
			s := filepath.Join(PublishPath, serverName)
			fmt.Println("s ", s)
			path := strings.SplitAfter(path, s+"/")
			fmt.Println("path1 ", path[1])

			*archiveFiles = append(*archiveFiles, path[1])
		}
		return nil
	}
}

func AddHashToPackageName(packageName *string, hash string) {
	s := strings.Split(*packageName, ".tar.gz")
	*packageName = s[0] + "_" + hash + ".tar.gz"
}

// serverName SimpTestServer
// fileName SimpTestServer_asdh213njonasd.tar.gz
func ConfirmFileName(serverName string, fileName string) bool {
	return strings.HasPrefix(fileName, serverName)
}

func IFExistThenRemove(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Printf("Path %s does not exist.\n", path)
		return nil
	}

	err = os.Remove(path)
	if err != nil {
		fmt.Printf("Error removing path %s: %v\n", path, err)
		return err
	}

	fmt.Printf("Path %s removed successfully.\n", path)
	return nil
}

func IsPidAlive(pid int, serverName string) bool {

	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("Error finding process: %v\n", err)
		return false
	}

	err = process.Signal(syscall.Signal(0))
	if err == nil {
		pid1 := ServantAlives[serverName]
		// 判断是否为同一个服务
		if pid1 == pid1 {
			return true // 进程存在
		}
	}

	if os.IsNotExist(err) || err == os.ErrProcessDone {
		return false // 进程不存在
	}

	fmt.Printf("Error signaling process: %v\n", err)
	return false
}
