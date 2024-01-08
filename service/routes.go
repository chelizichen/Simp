package service

import (
	handlers "Simp/handlers/http"
	"Simp/utils"
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func Registry(ctx *handlers.SimpHttpServerCtx) {
	G := ctx.Engine
	G.POST("/uploadServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		expectedHash := c.PostForm("hash") // 假设客户端提供了文件的哈希值
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		F, err := c.FormFile("file")
		storagePath := filepath.Join(cwd, utils.PublishPath, serverName, F.Filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "接受文件失败", nil))
			return
		}
		// 保存上传的文件到服务器临时目录
		tempPath := filepath.Join(cwd, "temp", F.Filename)
		if err := c.SaveUploadedFile(F, tempPath); err != nil {
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "保存上传的文件到服务器临时目录失败", nil))
			return
		}
		// 校验文件完整性（这里使用MD5哈希值作为示例）
		actualHash, err := utils.CalculateFileHash(tempPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "计算哈希值失败", nil))
			return
		}
		fmt.Println("计算HASH", actualHash)
		// 比较哈希值
		if actualHash != expectedHash && false {
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "比较哈希值失败", nil))
			return
		}
		// 移动文件到目标目录
		fmt.Println("tempPath", tempPath)
		fmt.Println("storagePath", storagePath)
		if err := os.Rename(tempPath, storagePath); err != nil {
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "移动文件失败", nil))
			return
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "上传成功", nil))
	})

	G.POST("/restartServer", func(c *gin.Context) {
		fileName := c.PostForm("fileName")
		serverName := c.PostForm("serverName")
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		storagePath := filepath.Join(cwd, utils.PublishPath, serverName, fileName)
		storageExEPath := filepath.Join(cwd, utils.PublishPath, serverName, fileName, "server")
		dest := filepath.Join(cwd, utils.PublishPath, serverName)
		err = utils.Unzip(storagePath, dest)
		if err != nil {
			fmt.Println("Error To Unzip", err.Error())
		}
		cmd := exec.Command(storageExEPath)
		err = cmd.Start()
		if err != nil {
			fmt.Println("Error To Err", err.Error())
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	G.POST("/test/restart", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		filePath := filepath.Join(cwd, utils.PublishPath, "SimpTestServer/childservice")
		cmd := exec.Command(filePath)
		stdoutPipe, err := cmd.StdoutPipe()
		err = cmd.Start()
		if err != nil {
			fmt.Println("Error To Err :", err.Error())
		}
		time.Sleep(2 * time.Second)
		go func() {
			scanner := bufio.NewScanner(stdoutPipe)
			for scanner.Scan() {
				outputLine := scanner.Text()
				fmt.Println("Child Process Output:", outputLine)
			}
		}()
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	G.POST("/getServers", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverPath := filepath.Join(cwd, utils.PublishPath)
		fmt.Println("serverPath", serverPath)
		subdirectories, err := utils.GetSubdirectories(serverPath)
		if err != nil {
			fmt.Println("Error To GetSubdirectories")
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", subdirectories))
	})

	G.POST("/createServer", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To GetWd", nil))
			return
		}
		value := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, "publish", value)
		err = os.Mkdir(serverPath, 512)
		if err != nil {
			fmt.Println("Error To Mkdir", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To Mkdir", nil))
			return
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	G.POST("/getServerPackageList", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To GetWd", nil))
			return
		}
		value := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils.PublishPath, value)
		var packages []string
		err = filepath.Walk(serverPath, utils.VisitTgzS(&packages))
		if err != nil {
			fmt.Printf("error walking the path %v: %v\n", serverPath, err)
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", packages))
	})

}
