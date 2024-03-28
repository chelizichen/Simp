package utils

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	p "github.com/shirou/gopsutil/process"

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

func IFNotExistThenCreate(filePath string) (os.File, error) {
	fmt.Println("Exist ", filePath)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error to Open File", err.Error())
	}
	return *f, err
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
		if pid == pid1 {
			return true // 进程存在
		}
	}

	if os.IsNotExist(err) || err == os.ErrProcessDone {
		return false // 进程不存在
	}

	fmt.Printf("Error signaling process: %v\n", err)
	return false
}
func GetProcessMemoryInfo(pid int) *p.Process {
	process, err := p.NewProcess(int32(pid))
	if err != nil {
		fmt.Println("Error creating new process:", err)
		return nil
	}
	return process
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

func CopyFile(src string, dst string) error {
	// 打开源文件
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 创建目标文件
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// 复制文件内容
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func Join(pre string) func(t string) string {
	return func(target string) string {
		return path.Join(pre, target)
	}
}

func DeleteDirectory(path string) error {
	// 获取目录中的文件和子目录
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		var currentPath string
		if file.IsDir() {
			// 如果是子目录，递归调用 DeleteDirectory 函数
			currentPath = fmt.Sprintf("%s/%s", path, file.Name())
			if err := DeleteDirectory(currentPath); err != nil {
				fmt.Printf("Error deleting directory: %v\n", err)
			}
		} else {
			// 如果是文件，直接删除文件
			currentPath = fmt.Sprintf("%s/%s", path, file.Name())
			err := os.Remove(currentPath)
			if err != nil {
				fmt.Printf("Error deleting file: %v\n", err)
			}
		}
	}

	// 删除目录本身
	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}
