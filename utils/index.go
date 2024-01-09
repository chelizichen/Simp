package utils

import (
	"archive/zip"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := filepath.Join(dest, file.Name)

		if file.FileInfo().IsDir() {
			// 如果是目录，创建目录
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		// 如果是文件，创建文件并写入内容
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		extractedFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer extractedFile.Close()

		sourceFile, err := file.Open()
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		// 将文件内容复制到目标文件
		if _, err := io.Copy(extractedFile, sourceFile); err != nil {
			return err
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

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 排除当前目录
		if path != directoryPath && info.IsDir() {
			paths := strings.SplitAfter(path, PublishPath)
			subdirectories = append(subdirectories, paths[1])
		}
		return nil
	})

	return subdirectories, err
}

func VisitTgzS(archiveFiles *[]string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err) // can't walk here,
			return nil       // but continue walking elsewhere
		}
		if f.IsDir() {
			return nil // not a file, ignore.
		}
		fmt.Println(path)

		// Check if the file has a .tar.gz extension
		if strings.HasSuffix(path, ".tar.gz") {
			path := strings.SplitAfter(path, PublishPath)

			*archiveFiles = append(*archiveFiles, path[1])
		}
		return nil
	}
}

func AddHashToPackageName(packageName *string, hash string) {
	s := strings.Split(*packageName, ".tar.gz")
	*packageName = s[0] + "_" + hash + s[1]
}

// serverName SimpTestServer
// fileName SimpTestServer_asdh213njonasd.tar.gz
func ConfirmFileName(serverName string, fileName string) bool {
	return strings.HasPrefix(fileName, serverName)
}
