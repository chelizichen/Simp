package service

import (
	"Simp/src/config"
	handlers "Simp/src/http"
	utils2 "Simp/src/utils"
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"gopkg.in/yaml.v2"
)

const TOKEN = "e609d00404645feed1c1733835b8c127"

func TOKEN_VALIDATE(ctx *gin.Context) {
	s := ctx.Request.Header.Get("token")
	fmt.Println("TokenValidate", s)
	fmt.Println("ctx.url", ctx.Request.URL)
	if s != TOKEN {
		if strings.HasPrefix(ctx.Request.URL.Path, "/simpserver/web") {
			ctx.Next()
			return
		}
		if strings.Contains(ctx.Request.URL.Path, "static/source") {
			ctx.Next()
			return
		}
		ctx.Redirect(http.StatusTemporaryRedirect, "/simpserver/web/login")
		return
	} else {
		ctx.Next()
		return
	}
}

type ServerCtx struct {
	ExitSignal *atomic.Bool
	Cron       *cron.Cron
}

func Registry(ctx *handlers.SimpHttpServerCtx, pre string) {

	f := utils2.Join(pre)
	G := ctx.Engine
	var RegistrhServicesCtx = make(map[string]ServerCtx)

	// G.GET(f("/web"), func(c *gin.Context) {
	// 	c.Redirect(http.StatusPermanentRedirect, "/web/login.html")
	// })
	G.POST(f("/login"), func(c *gin.Context) {
		token := c.PostForm("token")
		if token == TOKEN {
			c.JSON(http.StatusOK, handlers.Resp(0, "Ok", nil))
			return
		}
		c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error", nil))
	})

	GROUP := G.Group(pre)
	GROUP.Use(TOKEN_VALIDATE)
	GROUP.POST("/uploadServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		F, err := c.FormFile("file")
		storagePath := filepath.Join(cwd, utils2.PublishPath, serverName, F.Filename)
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
		actualHash, err := utils2.CalculateFileHash(tempPath)
		utils2.AddHashToPackageName(&storagePath, actualHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "计算哈希值失败", nil))
			return
		}
		// 移动文件到目标目录
		fmt.Println("tempPath", tempPath)
		fmt.Println("storagePath", storagePath)

		if err := utils2.MoveAndRemove(tempPath, storagePath); err != nil {
			fmt.Println("Error To Rename", err.Error())
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "移动文件失败", nil))
			return
		}
		releaseDoc := c.PostForm("doc")
		storageDocPath := filepath.Join(cwd, utils2.PublishPath, serverName, "doc.txt")
		E, err := utils2.IFNotExistThenCreate(storageDocPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "打开或创建文件失败"+err.Error(), nil))
		}
		defer E.Close()
		content := storagePath + "\n" + releaseDoc + "\n"
		E.WriteString(content)
		c.JSON(http.StatusOK, handlers.Resp(0, "上传成功", nil))
	})

	// serverName SimpTestServer
	// fileName SimpTestServer_asdh213njonasd.tar.gz
	GROUP.POST("/restartServer", func(c *gin.Context) {
		fileName := c.PostForm("fileName")
		serverName := c.PostForm("serverName")
		// targetPort 为扩容时指定的端口
		targetPort := c.DefaultPostForm("targetPort", "")
		ctxName := serverName + targetPort
		fmt.Println("targetPort | ", targetPort, " | ctxName |", ctxName)
		isAlive := utils2.ServantAlives[ctxName]
		if isAlive != 0 {
			cmd := exec.Command("kill", "-9", strconv.Itoa(isAlive))
			// 执行命令
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error killing process:", err)
				return
			}
			if sc, ok := RegistrhServicesCtx[ctxName]; ok {
				sc.ExitSignal.Store(true)
				sc.Cron.Stop()
			}
		}
		isSame := utils2.ConfirmFileName(serverName, fileName)
		if !isSame {
			fmt.Println("Error File!", fileName, "  | ", serverName)
		}
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		storagePath := filepath.Join(cwd, utils2.PublishPath, serverName, fileName)
		storageExEPath := filepath.Join(cwd, utils2.PublishPath, serverName, "service_go")
		storageNodePath := filepath.Join(cwd, utils2.PublishPath, serverName, "app.js")
		storageYmlEPath := filepath.Join(cwd, utils2.PublishPath, serverName, "simp.yaml")
		storageYmlProdPath := filepath.Join(cwd, utils2.PublishPath, serverName, "simpProd.yaml")
		dest := filepath.Join(cwd, utils2.PublishPath, serverName)
		sc, err := config.NewConfig(storageYmlEPath)
		if err != nil {
			fmt.Println("Error To Get Config")
		}
		s := sc.Server.StaticPath
		storageStaticPath := filepath.Join(cwd, utils2.PublishPath, serverName, s)
		// 判定为普通重启操作条件
		// 1. targetPort 存在为扩容操作，一般重启服务时 targetPort 是不存在的
		// 2. 如果需要重启服务，那么会将几个端口一起重启，那么将会传几个targetPort进来
		// 此时还需要区分主服务和扩容服务,必须先执行完主服务的重启后才能执行扩容服务
		// 同时 主服务重启时也需要重新执行所有流程
		fmt.Println("targetPort == fmt.Sprintf('v', sc.Server.Port)", fmt.Sprintf("%v", sc.Server.Port), " | target |", targetPort, " | ", targetPort == fmt.Sprintf("%v", sc.Server.Port))
		if targetPort == "" || targetPort == fmt.Sprintf("%v", sc.Server.Port) {
			err = utils2.IFExistThenRemove(storageStaticPath)
			if err != nil {
				fmt.Println("remove File Error storageStaticPath "+storageStaticPath, err.Error())
			}
			err = utils2.IFExistThenRemove(storageExEPath)
			if err != nil {
				fmt.Println("remove File Error storageExEPath "+storageExEPath, err.Error())
			}
			err = utils2.IFExistThenRemove(storageYmlEPath)
			if err != nil {
				fmt.Println("remove File Error storageYmlEPath "+storageYmlEPath, err.Error())
			}

			err = utils2.IFExistThenRemove(storageNodePath)
			if err != nil {
				fmt.Println("remove File Error storageNodePath "+storageNodePath, err.Error())
			}
			err = utils2.Unzip(storagePath, dest)
			if err != nil {
				fmt.Println("Error To Unzip", err.Error())
			}
			_, err = os.Stat(storageYmlProdPath)
			if err != nil {
				fmt.Println("os.Stat ", err.Error())
			}
			// 如果没有该文件，则将simp.yaml拷贝一份成simpProd.yaml
			if os.IsNotExist(err) {
				err = utils2.CopyFile(storageYmlEPath, storageYmlProdPath)
				if err != nil {
					fmt.Println("utils.CopyFile ", storageYmlEPath, err.Error())
				}
			}
		}

		var cmd *exec.Cmd
		fmt.Println("sc.Server.Type", sc.Server.Type)
		switch sc.Server.Type {
		case "node-http":
			{
				cmd = exec.Command("node", storageNodePath)
			}
		default:
			{
				cmd = exec.Command(storageExEPath)
			}
		}

		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println("Error Get StdoutPiper", err.Error())
		}
		stderrPipe, err := cmd.StderrPipe()
		if err != nil {
			fmt.Println("Error Get stderrPipe", err.Error())
		}
		// 设置环境变量
		cmd.Env = append(os.Environ(), "SIMP_PRODUCTION=Yes", "SIMP_SERVER_PATH="+dest)
		if targetPort != "" && targetPort != fmt.Sprintf("%v", sc.Server.Port) {
			cmd.Env = append(cmd.Env, "SIMP_TARGET_PORT="+targetPort, "SIMP_SERVER_INDEX="+targetPort)
		} else {
			cmd.Env = append(cmd.Env, "SIMP_SERVER_INDEX=1", "SIMP_TARGET_PORT="+fmt.Sprintf("%v", sc.Server.Port))
		}
		sm, err := utils2.NewSimpMonitor(serverName, "", targetPort)
		if err != nil {
			fmt.Println("Error To New Monitor", err.Error())
		}
		err = cmd.Start()
		if err != nil {
			fmt.Println("Error To EXEC Cmd Start", err.Error())
			fmt.Println("Cmd", cmd.Args)
		}
		var exit atomic.Bool
		cron := cron.New()
		RegistrhServicesCtx[ctxName] = ServerCtx{
			ExitSignal: &exit,
			Cron:       cron,
		}
		// 启动一个协程，用于读取并打印命令的输出
		go func() {
			spec := "0 0 0 * * *"
			// 添加定时任务
			err := cron.AddFunc(spec, func() {
				if !utils2.IsPidAlive(cmd.Process.Pid, ctxName) {
					cron.Stop()
				}
				newSM, err := utils2.NewSimpMonitor(serverName, "", targetPort)
				if err != nil {
					fmt.Println("Error To New Monitor", err.Error())
					return
				}
				sm = newSM
			})
			if err != nil {
				fmt.Println("AddFuncErr", err)
			}
			// 启动Cron调度器
			go cron.Start()
			go func() {
				for {
					if exit.Load() {
						fmt.Println(" serverName | ", ctxName, " |StopToReadOutput")
						return
					}
					// 读取输出
					buf := make([]byte, 1024)
					s := time.Now().Format(time.DateTime)
					n, err := stdoutPipe.Read(buf)
					if err != nil {
						break
					}
					// 打印输出
					content := s + " ServerName " + serverName + " || " + string(buf[:n]) + "\n"
					sm.AppendLogger(content)
				}
			}()
			go func() {
				for {
					if exit.Load() {
						fmt.Println("serverName | ", ctxName, " |StopToReadOutput")
						return
					}
					// 读取输出
					buf := make([]byte, 1024)
					s := time.Now().Format(time.DateTime)
					n, err := stderrPipe.Read(buf)
					if err != nil {
						break
					}
					// 打印输出
					content := s + " Error : ServerName " + serverName + " || " + string(buf[:n]) + "\n"
					fmt.Println(content)
					sm.AppendLogger(content)
				}
			}()

			go func() {
				for {
					time.Sleep(time.Second * 15)
					b := utils2.IsPidAlive(cmd.Process.Pid, ctxName)
					if !b {
						return
					}
					pInfo := utils2.GetProcessMemoryInfo(cmd.Process.Pid)
					cpuPercent, _ := pInfo.CPUPercent()
					cpuAffinity, _ := pInfo.CPUAffinity()
					createTime, _ := pInfo.CreateTime()
					Status, _ := pInfo.Status()
					pid := cmd.Process.Pid
					MemoryPercent, _ := pInfo.MemoryPercent()
					MemoryInfo, _ := pInfo.MemoryInfo()
					info := make(map[string]interface{})
					info["pid"] = pid
					info["MemoryInfo"] = MemoryInfo
					info["MemoryPercent"] = MemoryPercent
					info["CpuPercent"] = cpuPercent
					info["CpuAffinity"] = cpuAffinity
					info["CreateTime"] = createTime
					info["Status"] = Status
					info["ServerName"] = serverName
					pInfoContent, err := json.Marshal(info)
					if err != nil {
						break
					}
					s := time.Now().Format(time.DateTime)
					content := s + " ServerName " + serverName + " || " + string(pInfoContent) + "\n"
					sm.AppendLogger(content)
				}
			}()
		}()
		if err != nil {
			fmt.Println("Error To Err", err.Error())
		}
		v := make(map[string]interface{}, 10)
		fmt.Println("v", v)
		v["pid"] = cmd.Process.Pid
		v["status"] = true
		utils2.ServantAlives[ctxName] = cmd.Process.Pid

		c.JSON(http.StatusOK, handlers.Resp(0, "ok", v))
	})

	GROUP.POST("/testnode", func(c *gin.Context) {
		fileName := c.DefaultPostForm("fileName", "TestNodeServer.tar.gz")
		serverName := c.DefaultPostForm("serverName", "TestNodeServer")
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		storagePath := filepath.Join(cwd, utils2.PublishPath, serverName, fileName)
		dest := filepath.Join(cwd, utils2.PublishPath, serverName)
		err = utils2.Unzip(storagePath, dest)

		storageExEPath := filepath.Join(cwd, utils2.PublishPath, serverName, "app.js")
		cmd := exec.Command("node", storageExEPath)
		cmd.Env = append(os.Environ(), "SIMP_PRODUCTION=Yes", "SIMP_SERVER_PATH="+dest)
		stdoutPipe, _ := cmd.StdoutPipe()
		rc, _ := cmd.StderrPipe()
		go func() {

			for {
				// 读取输出
				buf := make([]byte, 1024)
				s := time.Now().Format(time.TimeOnly)
				n, err := stdoutPipe.Read(buf)
				if err != nil {
					break
				}
				// 打印输出
				content := s + "ServerName " + serverName + " || " + string(buf[:n]) + "\n"
				fmt.Println(content)
			}
			for {
				// 读取输出
				buf := make([]byte, 1024)
				s := time.Now().Format(time.TimeOnly)
				n, err := rc.Read(buf)
				if err != nil {
					break
				}
				// 打印输出
				content := s + "ServerName " + serverName + " || " + string(buf[:n]) + "\n"
				fmt.Println(content)
			}
		}()
		err = cmd.Start()

		if err != nil {
			fmt.Println("storageExEPath", storageExEPath)
			fmt.Println("err", err.Error())
		}

		c.AbortWithStatus(200)
	})

	GROUP.POST("/test/restart", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		filePath := filepath.Join(cwd, utils2.PublishPath, "SimpTestServer/childservice")
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

	GROUP.POST("/test/changeDoc", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		releaseDoc := c.PostForm("doc")
		storageDocPath := filepath.Join(cwd, utils2.PublishPath, "CalcServer", "doc.txt")
		E, err := utils2.IFNotExistThenCreate(storageDocPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "打开或创建文件失败"+err.Error(), nil))
		}
		defer E.Close()
		content := "\nCalcServer_2024_01_01_asdaasgjjhasioudh.tar.gz" + "\n" + releaseDoc + "\n"
		E.WriteString(content)
		c.JSON(http.StatusOK, handlers.Resp(0, "上传成功", nil))
	})
	GROUP.POST("/getServers", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverPath := filepath.Join(cwd, utils2.PublishPath)
		fmt.Println("serverPath", serverPath)
		subdirectories, err := utils2.GetSubdirectories(serverPath)
		if err != nil {
			fmt.Println("Error To GetSubdirectories")
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", subdirectories))
	})

	GROUP.POST("/createServer", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To GetWd", nil))
			return
		}
		value := c.PostForm("serverName")
		fmt.Println("createServer | serverName ", value)
		serverPath := filepath.Join(cwd, utils2.PublishPath, value)
		utils2.AutoCreateLoggerFile(value)
		err = os.Mkdir(serverPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error To Mkdir", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To Mkdir", nil))
			return
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	GROUP.POST("/getServerPackageList", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To GetWd", nil))
			return
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils2.PublishPath, serverName)
		fmt.Println("serverPath", serverPath)
		var packages []utils2.ReleasePackageVo
		err = filepath.Walk(serverPath, utils2.VisitTgzS(&packages, serverName))
		if err != nil {
			fmt.Printf("error walking the path %v: %v\n", serverPath, err)
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", packages))
	})

	GROUP.POST("/deletePackage", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		F := c.PostForm("fileName")
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		storagePath := filepath.Join(cwd, utils2.PublishPath, serverName, F)
		err = os.Remove(storagePath)
		if err != nil {
			fmt.Println("Error To RemoveFile", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To RemoveFile", nil))
			return
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	GROUP.POST("/checkServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		pid := c.DefaultPostForm("pid", fmt.Sprint(utils2.ServantAlives[serverName]))
		P, err := strconv.Atoi(pid)
		if err != nil {
			fmt.Println("Error to Atoi", err.Error())
		}
		b := utils2.IsPidAlive(P, serverName)
		v := make(map[string]interface{}, 10)
		v["status"] = false
		if b {
			v["pid"] = pid
			v["status"] = true
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", v))
	})

	GROUP.POST("/checkConfig", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		configPath := filepath.Join(utils2.PublishPath, serverName, "simp.yaml")
		configProdPath := filepath.Join(utils2.PublishPath, serverName, "simpProd.yaml")
		sc, err := config.NewConfig(configPath)
		prod, err := config.NewConfig(configProdPath)
		mergeConf := config.MergeYAML(prod, sc)
		if err != nil {
			fmt.Println("Error To Get NewConfig", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", mergeConf))
	})

	GROUP.POST("/coverConfig", func(c *gin.Context) {
		var reqVo config.CoverConfigVo
		if err := c.BindJSON(&reqVo); err != nil {
			c.JSON(http.StatusOK, handlers.Resp(0, "-1", err.Error()))
			return
		}
		serverName := reqVo.ServerName
		uploadConfig := reqVo.Conf
		if serverName == "" {
			fmt.Println("Server Name is Empty")
			c.JSON(http.StatusOK, handlers.Resp(0, "Server Name is Empty", nil))
			return
		}
		marshal, err := yaml.Marshal(uploadConfig)
		if err != nil {
			fmt.Println("Error To Stringify config", err.Error())
			c.JSON(http.StatusOK, handlers.Resp(0, "Error To Stringify config", nil))
			return
		}
		fmt.Println("serverName", serverName)
		fmt.Println("uploadConfig", string(marshal))
		if len(marshal) == 0 {
			fmt.Println("Error To Stringify config", err.Error())
			c.JSON(http.StatusOK, handlers.Resp(0, "Error To Stringify config", nil))
			return
		}
		configPath := filepath.Join(utils2.PublishPath, serverName, "simpProd.yaml")
		err = config.CoverConfig(string(marshal), configPath)
		if err != nil {
			fmt.Println("CoverConfig Error", err.Error())
			c.JSON(200, handlers.Resp(-1, "CoverConfig Error", nil))
		}
		c.JSON(200, handlers.Resp(0, "ok", nil))
	})

	GROUP.POST("/deleteServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		ErrorToRemoveAll := "Error To Remove All"
		ErrorToGetServerName := "Error To Get ServerName"
		if serverName == "" {
			c.AbortWithStatusJSON(200, handlers.Resp(-1, ErrorToGetServerName, nil))
			return
		}
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverPath := filepath.Join(cwd, utils2.PublishPath, serverName)
		fmt.Println("DeleteDirectory", serverPath)
		err = utils2.DeleteDirectory(serverPath)
		if err != nil {
			fmt.Println(ErrorToRemoveAll, err.Error())
			c.AbortWithStatusJSON(200, handlers.Resp(-1, ErrorToRemoveAll, nil))
			return
		}
		c.JSON(200, handlers.Resp(0, "ok", nil))
	})

	GROUP.POST("/getChildStats", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		if serverName == "" {
			c.AbortWithStatusJSON(200, handlers.Resp(-1, "missing params serverName", nil))
			return
		}
		m := make(map[string]map[string]interface{})
		for sName, ctx := range utils2.ServantAlives {
			if strings.HasPrefix(sName, serverName) {
				pid := ctx
				name := sName
				b := utils2.IsPidAlive(pid, sName)
				v := make(map[string]interface{}, 10)
				v["status"] = false
				if b {
					v["pid"] = pid
					v["status"] = true
				}

				m[name] = v
			}
		}
		c.AbortWithStatusJSON(200, handlers.Resp(0, "ok", m))
	})

	GROUP.POST("/shutdownServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		pid := utils2.ServantAlives[serverName]
		if pid == 0 {
			c.AbortWithStatusJSON(200, handlers.Resp(-1, "暂无PID", nil))
			return
		}
		if ctx, ok := RegistrhServicesCtx[serverName]; ok {
			fmt.Println("exit routine")
			fmt.Println("RegistrhServicesCtx[serverName]", serverName, "|", RegistrhServicesCtx[serverName])
			ctx.Cron.Stop()
			ctx.ExitSignal.Store(true)
		}

		fmt.Println("shoutDown server", serverName, "pid is ", pid)
		cmd := exec.Command("kill", "-9", strconv.Itoa(pid))
		// 执行命令
		err := cmd.Run()
		if err != nil {
			fmt.Println("x:", err)
			c.AbortWithStatusJSON(200, handlers.Resp(-1, "关闭服务异常", err.Error()))
			return
		}
		utils2.ServantAlives[serverName] = 0
		c.AbortWithStatusJSON(200, handlers.Resp(0, "ok", nil))
	})

	// tail -n rows log_file | grep "pattern"
	GROUP.POST("/getServerLog", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		fileName := c.PostForm("fileName")
		pattern := c.DefaultPostForm("pattern", "")
		rows := c.DefaultPostForm("rows", "100")
		sm, err := utils2.NewSearchLogMonitor(serverName, fileName)
		if err != nil {
			fmt.Println("Error To New SimMonitor", err.Error())
		}
		s, err := sm.GetLogger(pattern, rows)
		if err != nil {
			fmt.Println("Error To GetLogger", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", s))
	})

	GROUP.POST("/getApiJson", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils2.PublishPath, serverName, "API.json")
		Content, err := os.ReadFile(serverPath)
		if err != nil {
			fmt.Println("Error To ReadFile", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", string(Content)))
	})

	GROUP.POST("/getDoc", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils2.PublishPath, serverName, "doc.txt")
		Content, err := os.ReadFile(serverPath)
		if err != nil {
			fmt.Println("Error To ReadFile", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", string(Content)))
	})

	GROUP.POST("/getLogList", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils2.PublishPath, serverName)
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

	GROUP.POST("/main/getLogList", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverPath := filepath.Join(cwd, "static/main")
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

	GROUP.POST("/main/getServerLog", func(c *gin.Context) {
		logFile := c.PostForm("logFile")
		pattern := c.DefaultPostForm("pattern", "")
		rows := c.DefaultPostForm("rows", "100")
		sm, err := utils2.NewMainSearchLogMonitor(logFile)
		if err != nil {
			fmt.Println("Error To New SimMonitor", err.Error())
			c.JSON(200, handlers.Resp(-2, err.Error(), nil))
			return
		}
		s, err := sm.GetLogger(pattern, rows)
		if err != nil {
			fmt.Println("Error To GetLogger", err.Error())
			c.JSON(200, handlers.Resp(-1, err.Error(), nil))
			return
		}
		c.JSON(200, handlers.Resp(0, "ok", s))
		return

	})
	G.Use(GROUP.Handlers...)

}
