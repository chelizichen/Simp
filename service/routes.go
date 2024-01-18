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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const TOKEN = "e609d00404645feed1c1733835b8c127"

func TOKEN_VALIDATE(ctx *gin.Context) {
	s := ctx.Request.Header.Get("token")
	if s != TOKEN {
		ctx.JSON(http.StatusBadRequest, handlers.Resp(-2, "Token Error", nil))
	}
	ctx.Next()
}

func Registry(ctx *handlers.SimpHttpServerCtx) {
	G := ctx.Engine
	G.GET("/web", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/web/login.html")
	})
	G.POST("/login", func(c *gin.Context) {
		token := c.PostForm("token")
		if token == TOKEN {
			c.JSON(http.StatusOK, handlers.Resp(0, "Ok", nil))
			return
		}
		c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error", nil))
	})
	G.POST("/uploadServer", TOKEN_VALIDATE, func(c *gin.Context) {
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

		if err := utils.MoveAndRemove(tempPath, storagePath); err != nil {
			fmt.Println("Error To Rename", err.Error())
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "移动文件失败", nil))
			return
		}
		releaseDoc := c.PostForm("doc")
		utils.AppendDocToTgz(storagePath, releaseDoc)
		c.JSON(http.StatusOK, handlers.Resp(0, "上传成功", nil))
	})

	// serverName SimpTestServer
	// fileName SimpTestServer_asdh213njonasd.tar.gz
	G.POST("/restartServer", TOKEN_VALIDATE, func(c *gin.Context) {
		fileName := c.PostForm("fileName")
		serverName := c.PostForm("serverName")
		isAlive := utils.ServantAlives[serverName]
		if isAlive != 0 {
			cmd := exec.Command("kill", "-9", strconv.Itoa(isAlive))
			// 执行命令
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error killing process:", err)
				return
			}
		}
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
		storageDocPath := filepath.Join(cwd, utils.PublishPath, serverName, "doc.txt")

		err = utils.IFExistThenRemove(storageExEPath)
		if err != nil {
			fmt.Println("remove File Error "+storageExEPath, err.Error())
		}
		err = utils.IFExistThenRemove(storageYmlEPath)
		if err != nil {
			fmt.Println("remove File Error "+storageYmlEPath, err.Error())
		}
		err = utils.IFExistThenRemove(storageDocPath)
		if err != nil {
			fmt.Println("remove File Error "+storageDocPath, err.Error())
		}
		dest := filepath.Join(cwd, utils.PublishPath, serverName)

		err = utils.Unzip(storagePath, dest)
		if err != nil {
			fmt.Println("Error To Unzip", err.Error())
		}
		cmd := exec.Command(storageExEPath)
		stdoutPipe, err := cmd.StdoutPipe()
		// 设置环境变量
		cmd.Env = append(os.Environ(), "SIMP_PRODUCTION=Yes", "SIMP_CONFIG_PATH="+storageYmlEPath, "SIMP_SERVER_PATH="+dest)
		sm, err := utils.NewSimpMonitor(serverName, "")
		err = cmd.Start()
		// 启动一个协程，用于读取并打印命令的输出
		go func() {
			for {
				// 读取输出
				buf := make([]byte, 1024)
				n, err := stdoutPipe.Read(buf)
				if err != nil {
					break
				}
				// 打印输出
				content := "ServerName " + serverName + " || " + string(buf[:n]) + "\n"
				sm.AppendLogger(content)
			}
		}()
		if err != nil {
			fmt.Println("Error To Err", err.Error())
		}
		v := make(map[string]interface{}, 10)
		fmt.Println("v", v)
		v["pid"] = cmd.Process.Pid
		v["status"] = true
		utils.ServantAlives[serverName] = cmd.Process.Pid

		c.JSON(http.StatusOK, handlers.Resp(0, "ok", v))
	})

	G.POST("/test/restart", TOKEN_VALIDATE, func(c *gin.Context) {
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

	G.POST("/getServers", TOKEN_VALIDATE, func(c *gin.Context) {
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

	G.POST("/createServer", TOKEN_VALIDATE, func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To GetWd", nil))
			return
		}
		value := c.PostForm("serverName")
		fmt.Println("createServer | serverName ", value)
		serverPath := filepath.Join(cwd, utils.PublishPath, value)
		utils.AutoCreateLoggerFile(value)
		err = os.Mkdir(serverPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error To Mkdir", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To Mkdir", nil))
			return
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	G.POST("/getServerPackageList", TOKEN_VALIDATE, func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To GetWd", nil))
			return
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName)
		fmt.Println("serverPath", serverPath)
		var packages []utils.ReleasePackageVo
		err = filepath.Walk(serverPath, utils.VisitTgzS(&packages, serverName))
		if err != nil {
			fmt.Printf("error walking the path %v: %v\n", serverPath, err)
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", packages))
	})

	G.POST("/deletePackage", TOKEN_VALIDATE, func(c *gin.Context) {
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

	G.POST("/checkServer", TOKEN_VALIDATE, func(c *gin.Context) {
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

	G.POST("/checkConfig", TOKEN_VALIDATE, func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		configPath := filepath.Join(utils.PublishPath, serverName, "simp.yaml")
		sc, err := config.NewConfig(configPath)
		if err != nil {
			fmt.Println("Error To Get NewConfig", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", sc))
	})

	G.POST("/coverConfig", TOKEN_VALIDATE, func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		uploadConfig := c.PostForm("uploadConfig")
		conf, err := config.ParseConfig(uploadConfig)
		if err != nil {
			fmt.Println("Error To Get NewConfig", err.Error())
			c.JSON(200, handlers.Resp(-1, "Error To ParseConfig", nil))
			return
		}
		configPath := filepath.Join(utils.PublishPath, serverName, "simp.yaml")
		config.CoverConfig(conf, configPath)
		c.JSON(200, handlers.Resp(0, "ok", nil))
	})

	G.POST("/deleteAllPackage", TOKEN_VALIDATE, func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		ErrorToRemoveAll := "Error To Remove All"
		ErrorToMakeAServer := "Error To Make A Sever"
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName)
		err = os.RemoveAll(serverPath)
		if err != nil {
			fmt.Println(ErrorToRemoveAll, err.Error())
			c.JSON(200, handlers.Resp(-1, ErrorToRemoveAll, nil))
			return
		}
		err = os.Mkdir(serverPath, os.ModePerm)
		if err != nil {
			fmt.Println(ErrorToMakeAServer, err.Error())
			c.JSON(200, handlers.Resp(-1, ErrorToMakeAServer, nil))
			return
		}
		c.JSON(200, handlers.Resp(0, "ok", nil))
	})

	G.POST("/shutdownServer", TOKEN_VALIDATE, func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		pid := utils.ServantAlives[serverName]
		if pid == 0 {
			c.JSON(200, handlers.Resp(-1, "暂无PID", nil))
			return
		}
		fmt.Println("shoutDown server", serverName, "pid is ", pid)
		cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
		// 执行命令
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error killing process:", err)
			return
		}
		utils.ServantAlives[serverName] = 0
		c.JSON(200, handlers.Resp(0, "ok", nil))
	})

	G.POST("/getServerLog", TOKEN_VALIDATE, func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		fileName := c.PostForm("fileName")
		pattern := c.DefaultPostForm("pattern", "")
		sm, err := utils.NewSearchLogMonitor(serverName, fileName)
		if err != nil {
			fmt.Println("Error To New SimMonitor", err.Error())
		}
		s, err := sm.GetLogger(pattern)
		if err != nil {
			fmt.Println("Error To GetLogger", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", s))
	})

	G.POST("/getApiJson", TOKEN_VALIDATE, func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName, "API.json")
		Content, err := os.ReadFile(serverPath)
		if err != nil {
			fmt.Println("Error To ReadFile", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", string(Content)))
	})

	G.POST("/getDoc", TOKEN_VALIDATE, func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName, "doc.txt")
		Content, err := os.ReadFile(serverPath)
		if err != nil {
			fmt.Println("Error To ReadFile", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", string(Content)))
	})

	G.POST("/getLogList", TOKEN_VALIDATE, func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName)
		D, err := os.ReadDir(serverPath)
		if err != nil {
			fmt.Println("Error To ReadDir", err.Error())
		}
		var loggers []string
		for i := 0; i < len(D); i++ {
			de := D[i]
			s := de.Name()
			b := strings.HasSuffix(s, ".log")
			if b {
				loggers = append(loggers, s)
			}
		}
		c.JSON(200, handlers.Resp(0, "ok", loggers))
	})

}
