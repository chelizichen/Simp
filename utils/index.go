package utils

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/c4milo/unpackit"
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
	fmt.Println("src", src)
	fmt.Println("dest", dest)
	file, err := os.Open(src)
	if err != nil {
		fmt.Println("Open Error", err.Error())
		return err
	}
	defer file.Close()
	err = unpackit.Unpack(file, dest)
	if err != nil {
		fmt.Println("Unpackit Error", err.Error())
		return err
	}
	return nil
	// tarReader := tar.NewReader(file)

	// for {
	// 	header, err := tarReader.Next()

	// 	if err == io.EOF {
	// 		fmt.Println("header Error", err.Error())
	// 		break
	// 	}

	// 	if err != nil {
	// 		return err
	// 	}

	// 	target := filepath.Join(dest, header.Name)

	// 	switch header.Typeflag {
	// 	case tar.TypeDir:
	// 		if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
	// 			return err
	// 		}
	// 	case tar.TypeReg:
	// 		file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
	// 		if err != nil {
	// 			fmt.Println("OpenFile Error", err.Error())
	// 			return err
	// 		}

	// 		if _, err := io.Copy(file, tarReader); err != nil {
	// 			return err
	// 		}
	// 		file.Close()
	// 	default:
	// 		return fmt.Errorf("Unknown type: %v in %s", header.Typeflag, header.Name)
	// 	}
	// }

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

type ReleasePackageVo struct {
	Hash        string
	PackageName string
}

func VisitTgzS(archiveFiles *[]ReleasePackageVo, serverName string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if f.IsDir() {
			return nil // not a file, ignore.
		}

		// 定义正则表达式，匹配文件名中的哈希值
		re := regexp.MustCompile(`_(\d{4}_\d{2}_\d{2}_\d{2}_\d{2}_\d{2}__[a-zA-Z0-9]+).tar.gz$`)
		// Check if the file has a .tar.gz extension
		if strings.HasSuffix(path, ".tar.gz") {
			// 使用正则表达式查找匹配项
			matches := re.FindStringSubmatch(path)
			// 输出匹配的哈希值
			if len(matches) == 2 {
				a := ReleasePackageVo{
					PackageName: serverName,
					Hash:        matches[1],
				}
				*archiveFiles = append(*archiveFiles, a)
			} else {
				fmt.Println("No hash found in the file path.")
			}

		}
		return nil
	}
}

func AddHashToPackageName(packageName *string, hash string) {
	s := strings.Split(*packageName, ".tar.gz")
	format := time.Now().Format("2006_01_02_15_01_05")
	fmt.Println("packageName", packageName)
	fmt.Println("hash", hash)
	fmt.Println("format", format)
	fmt.Println("s", s)
	*packageName = s[0] + "_" + format + "__" + hash + ".tar.gz"
}

// serverName SimpTestServer
// fileName SimpTestServer_asdh213njonasd.tar.gz
func ConfirmFileName(serverName string, fileName string) bool {
	return strings.HasPrefix(fileName, serverName)
}

func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
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

func MoveAndRemove(sourcePath, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	err = os.Remove(sourcePath)
	if err != nil {
		return err
	}

	return nil
}

func AppendDocToTgz(storagePath, docContent string) {
	// Open the existing .tgz file for reading
	file, err := os.Open(storagePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a new .tgz file for writing
	newFile, err := os.Create(storagePath + "_new")
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	// Create a gzip writer for the new file
	gzipWriter := gzip.NewWriter(newFile)
	defer gzipWriter.Close()

	// Create a tar writer for the gzip writer
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Create a tar reader for the existing file
	tarReader := tar.NewReader(file)
	//defer tarReader.Close()

	// Copy existing contents to the new tar archive
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			panic(err)
		}

		if _, err := io.Copy(tarWriter, tarReader); err != nil {
			panic(err)
		}
	}

	// Add the new file to the archive
	newFileHeader := &tar.Header{
		Name: "doc.txt",
		Mode: 0600,
		Size: int64(len(docContent)),
	}

	if err := tarWriter.WriteHeader(newFileHeader); err != nil {
		panic(err)
	}

	if _, err := tarWriter.Write([]byte(docContent)); err != nil {
		panic(err)
	}

	// Rename the new file to replace the original
	if err := os.Rename(storagePath+"_new", storagePath); err != nil {
		panic(err)
	}
}
