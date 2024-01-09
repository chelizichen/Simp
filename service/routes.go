package service

import (
	"Simp/config"
	handlers "Simp/handlers/http"
	"Simp/utils"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Registry(ctx *handlers.SimpHttpServerCtx) {
	G := ctx.Engine
	G.POST("/uploadServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
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
		utils.AddHashToPackageName(&storagePath, actualHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "计算哈希值失败", nil))
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

	// serverName SimpTestServer
	// fileName SimpTestServer_asdh213njonasd.tar.gz
	G.POST("/restartServer", func(c *gin.Context) {
		fileName := c.PostForm("fileName")
		serverName := c.PostForm("serverName")

		isSame := utils.ConfirmFileName(serverName, fileName)
		if !isSame {
			fmt.Println("Error File!", fileName, "  | ", serverName)
		}
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		storagePath := filepath.Join(cwd, utils.PublishPath, serverName, fileName)
		storageExEPath := filepath.Join(cwd, utils.PublishPath, serverName, "service_go")
		storageYmlEPath := filepath.Join(cwd, utils.PublishPath, serverName, "simp.yaml")

		err = utils.IFExistThenRemove(storageExEPath)
		if err != nil {
			fmt.Println("remote File Error "+storageExEPath, err.Error())
		}
		err = utils.IFExistThenRemove(storageYmlEPath)
		if err != nil {
			fmt.Println("remote File Error "+storageYmlEPath, err.Error())
		}
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
		v := make(map[string]interface{}, 10)

		v["pid"] = cmd.Process.Pid
		v["status"] = true
		utils.ServantAlives[serverName] = cmd.Process.Pid

		c.JSON(http.StatusOK, handlers.Resp(0, "ok", v))
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
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName)
		var packages []string
		err = filepath.Walk(serverPath, utils.VisitTgzS(&packages, serverName))
		if err != nil {
			fmt.Printf("error walking the path %v: %v\n", serverPath, err)
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", packages))
	})

	G.POST("/deleteServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		F := c.PostForm("fileName")
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		storagePath := filepath.Join(cwd, utils.PublishPath, serverName, F)
		err = os.Remove(storagePath)
		if err != nil {
			fmt.Println("Error To RemoveFile", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To RemoveFile", nil))
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	G.POST("/checkServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		pid := c.DefaultPostForm("pid", fmt.Sprint(utils.ServantAlives[serverName]))
		P, err := strconv.Atoi(pid)
		if err != nil {
			fmt.Println("Error to Atoi", err.Error())
		}
		b := utils.IsPidAlive(P, serverName)
		v := make(map[string]interface{}, 10)
		v["status"] = false
		if b == true {
			v["pid"] = pid
			v["status"] = true
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", v))
	})

	G.POST("/checkConfig", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		configPath := filepath.Join(utils.PublishPath, serverName, "simp.yaml")
		sc, err := config.NewConfig(configPath)
		if err != nil {
			fmt.Println("Error To Get NewConfig", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", sc))
	})

}
