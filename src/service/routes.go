package service

import (
	"Simp/src/config"
	handlers "Simp/src/http"
	utils2 "Simp/src/utils"
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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

func Registry(ctx *handlers.SimpHttpServerCtx, pre string) {

	f := utils2.Join(pre)
	G := ctx.Engine

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

	GROUP.POST("/restartServer", func(c *gin.Context) {
		fileName := c.PostForm("fileName")                                         // 文件名称
		serverName := c.PostForm("serverName")                                     // 服务名称
		targetPort := c.PostForm("targetPort")                                     // 集群模式下需要指定端口
		releaseType := c.DefaultPostForm("releaseType", utils2.RELEASE_SINGLENODE) // 集群模式下需要指定Type 默认普通单节点发布
		isSame := utils2.ConfirmFileName(serverName, fileName)
		if !isSame {
			msg := "Error File!" + fileName + "  | " + serverName
			fmt.Println(msg)
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, msg, nil))
			return
		}

		svr := utils2.GetServant(serverName, targetPort)
		// svr.Language = Language
		err := svr.StopServant()
		if err != nil {
			msg := "Error StopServant!" + fileName + "  | " + serverName
			fmt.Println(msg)
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, msg, nil))
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			msg := "Error To GetWd :" + err.Error()
			fmt.Printf("msg: %v\n", msg)
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, msg, nil))
			return
		}
		storageYmlEPath := utils2.GetFilePath(cwd, serverName, utils2.DevConfEntry)
		storageYmlProdPath := utils2.GetFilePath(cwd, serverName, utils2.ProdConfEntry)
		storagePath := utils2.GetFilePath(cwd, serverName, fileName)
		dest := filepath.Join(cwd, utils2.PublishPath, serverName)

		isFirstRelease := !utils2.IsExist(storageYmlProdPath)
		var runScript func() *exec.Cmd
		var confPort int
		if isFirstRelease {
			err = utils2.Unzip(storagePath, dest) // 直接解压
			if err != nil {
				fmt.Println("Error To Unzip", err.Error())
			}
			err = utils2.CopyProdYml(storageYmlEPath, storageYmlProdPath) // 拷贝配置文件
			if err != nil {
				fmt.Println("Error To CopyProdYml", err.Error())
			}
			sc, err := config.NewConfig(storageYmlProdPath) // 引入配置文件
			if err != nil {
				fmt.Println("Error To NewConfig", err.Error())
			}
			confPort = sc.Server.Port
			svr.Language = sc.Server.Type
			if sc.Server.Type == utils2.RELEASE_TYPE_NODEJS {
				runScript = func() *exec.Cmd {
					storageNodePath := utils2.GetFilePath(cwd, serverName, utils2.NodeJsEntry)
					var cmd *exec.Cmd = exec.Command("node", storageNodePath)
					return cmd
				}
			} else if sc.Server.Type == utils2.RELEASE_TYPE_GO {
				runScript = func() *exec.Cmd {
					storageExEPath := utils2.GetFilePath(cwd, serverName, utils2.GoEntry)
					var cmd *exec.Cmd = exec.Command(storageExEPath)
					return cmd
				}
			} else if sc.Server.Type == utils2.RELEASE_TYPE_JAVA {
				// java -jar your-application.jar -Dspring.config.location=file:/path/to/application.yml,file:/path/to/another-config.yaml
				runScript = func() *exec.Cmd {
					storageJavaPath := utils2.GetFilePath(cwd, serverName, utils2.SpringEntry)
					var cmd *exec.Cmd = exec.Command("java", "-jar", storageJavaPath, "-Dspring.config.location="+"file:"+storageYmlProdPath)
					return cmd
				}
			}
		} else if releaseType == utils2.RELEASE_SINGLENODE {
			var clearScript func(sc config.SimpConfig)
			sc, err := config.NewConfig(storageYmlProdPath)
			svr.Language = sc.Server.Type

			if svr.Language == utils2.RELEASE_TYPE_NODEJS {
				storageNodePath := utils2.GetFilePath(cwd, serverName, utils2.NodeJsEntry)
				clearScript = func(sc config.SimpConfig) {
					confPort = sc.Server.Port
					storageStaticPath := utils2.GetFilePath(cwd, serverName, sc.Server.StaticPath)
					err = utils2.IFExistThenRemove(storageStaticPath, true)
					if err != nil {
						fmt.Println("remove File Error storageStaticPath "+storageStaticPath, err.Error())
					}
					err = utils2.IFExistThenRemove(storageNodePath, false)
					if err != nil {
						fmt.Println("remove File Error storageNodePath "+storageNodePath, err.Error())
					}
					err = utils2.Unzip(storagePath, dest) // 直接解压
					if err != nil {
						fmt.Println("Error To Unzip", err.Error())
					}
				}
				runScript = func() *exec.Cmd {
					var cmd *exec.Cmd = exec.Command("node", storageNodePath)
					return cmd
				}
			} else if svr.Language == utils2.RELEASE_TYPE_GO {
				storageExEPath := utils2.GetFilePath(cwd, serverName, utils2.GoEntry)
				clearScript = func(sc config.SimpConfig) {
					confPort = sc.Server.Port
					storageStaticPath := utils2.GetFilePath(cwd, serverName, sc.Server.StaticPath)
					err = utils2.IFExistThenRemove(storageStaticPath, true)
					if err != nil {
						fmt.Println("remove File Error storageStaticPath "+storageStaticPath, err.Error())
					}
					err = utils2.IFExistThenRemove(storageExEPath, false)
					if err != nil {
						fmt.Println("remove File Error storageExEPath "+storageExEPath, err.Error())
					}
					err = utils2.Unzip(storagePath, dest) // 直接解压
					if err != nil {
						fmt.Println("Error To Unzip", err.Error())
					}
				}
				runScript = func() *exec.Cmd {
					var cmd *exec.Cmd = exec.Command(storageExEPath)
					return cmd
				}
			} else if svr.Language == utils2.RELEASE_TYPE_JAVA {
				clearScript = func(sc config.SimpConfig) {
					confPort = sc.Server.Port
					storageStaticPath := utils2.GetFilePath(cwd, serverName, sc.Server.StaticPath)
					err = utils2.IFExistThenRemove(storageStaticPath, true)
					if err != nil {
						fmt.Println("remove File Error storageStaticPath "+storageStaticPath, err.Error())
					}
					err = utils2.Unzip(storagePath, dest) // 直接解压
					if err != nil {
						fmt.Println("Error To Unzip", err.Error())
					}
				}
				// java -jar your-application.jar -Dspring.config.location=file:/path/to/application.yml,file:/path/to/another-config.yaml
				runScript = func() *exec.Cmd {
					storageJavaPath := utils2.GetFilePath(cwd, serverName, utils2.SpringEntry)
					var cmd *exec.Cmd = exec.Command("java", "-jar", storageJavaPath, "-D", "spring.config.location="+"file:"+storageYmlProdPath)
					return cmd
				}
			}
			confPort = sc.Server.Port
			if err != nil {
				fmt.Println("Error To NewConfig", err.Error())
			}
			err = utils2.IFExistThenRemove(storageYmlEPath, false)
			if err != nil {
				fmt.Println("remove File Error storageYmlEPath "+storageYmlEPath, err.Error())
			}
			clearScript(sc) // 执行清除
		} else if releaseType == utils2.RELEASE_CLUSTER {
			sc, err := config.NewConfig(storageYmlProdPath)
			if err != nil {
				fmt.Println("err NewConfig", err.Error())
			}
			svr.Language = sc.Server.Type
			if svr.Language == utils2.RELEASE_TYPE_NODEJS {
				runScript = func() *exec.Cmd {
					storageNodePath := utils2.GetFilePath(cwd, serverName, utils2.NodeJsEntry)
					var cmd *exec.Cmd = exec.Command("node", storageNodePath)
					return cmd
				}
			} else if svr.Language == utils2.RELEASE_TYPE_GO {
				runScript = func() *exec.Cmd {
					storageExEPath := utils2.GetFilePath(cwd, serverName, utils2.GoEntry)
					var cmd *exec.Cmd = exec.Command(storageExEPath)
					return cmd
				}
			} else if svr.Language == utils2.RELEASE_TYPE_JAVA {
				runScript = func() *exec.Cmd {
					storageJavaPath := utils2.GetFilePath(cwd, serverName, utils2.SpringEntry)
					var cmd *exec.Cmd = exec.Command("java", "-jar", storageJavaPath, "-Dspring.config.location="+"file:"+storageYmlProdPath)
					return cmd
				}
			}
		}

		script := runScript()
		stdoutPipe, err := script.StdoutPipe()
		if err != nil {
			fmt.Println("Error Get StdoutPiper", err.Error())
		}
		stderrPipe, err := script.StderrPipe()
		if err != nil {
			fmt.Println("Error Get stderrPipe", err.Error())
		}
		// 设置环境变量
		script.Env = append(os.Environ(), "SIMP_PRODUCTION=Yes", "SIMP_SERVER_PATH="+dest)
		if svr.Language == utils2.RELEASE_TYPE_JAVA {
			script.Env = append(script.Env, "SIMP_PROD_YAML="+storageYmlProdPath)
		}
		if releaseType == utils2.RELEASE_CLUSTER {
			script.Env = append(script.Env, "SIMP_TARGET_PORT="+targetPort, "SIMP_SERVER_INDEX="+targetPort)
		} else {
			script.Env = append(script.Env, "SIMP_SERVER_INDEX=1", fmt.Sprintf("SIMP_TARGET_PORT=%v", confPort))
		}
		if err != nil {
			msg := "Error To New Monitor" + err.Error()
			fmt.Printf("msg: %v\n", msg)
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, msg, nil))
			return
		}
		err = script.Start()

		if err != nil {
			msg := "Error To EXEC Cmd Start ：" + err.Error()
			fmt.Println(msg)
			fmt.Println("Cmd", script.Args)
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, msg, nil))
		}
		var exit atomic.Bool
		cronSvr := cron.New()
		svr.Process = script
		svr.Pid = script.Process.Pid
		svr.ExitSignal = exit
		svr.Cron = cronSvr

		// 启动一个协程，用于读取并打印命令的输出
		go func() {
			spec := "0 0 0 * * *"
			sm, err := utils2.NewSimpMonitor(serverName, "", targetPort)
			if err != nil {
				fmt.Println("Error To New Monitor", err.Error())
			}
			// 添加定时任务
			err = cronSvr.AddFunc(spec, func() {

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
			go cronSvr.Start()

			go func() {
				for {
					if exit.Load() {
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
					content, isAlive := svr.ServantMonitor()
					if !isAlive {
						break
					}
					sm.AppendLogger(content)
				}
			}()
		}()
		v := make(map[string]interface{}, 2)
		v["pid"] = svr.Pid
		v["status"] = true
		utils2.SubServants[svr.GetContextName()] = svr
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
		port := c.PostForm("port")
		s := utils2.GetServant(serverName, port)
		b := utils2.IsPidAlive(s.Pid)
		v := make(map[string]interface{}, 10)
		v["status"] = false
		if b {
			v["pid"] = s.Pid
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
		for sName, ctx := range utils2.SubServants {
			context := &ctx
			if strings.HasPrefix(sName, serverName) {
				pid := context.Pid
				name := sName
				b := utils2.IsPidAlive(pid)
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
		port := c.PostForm("port")
		s := utils2.GetServant(serverName, port)

		if s.Pid == 0 {
			c.AbortWithStatusJSON(200, handlers.Resp(-1, "暂无PID", nil))
			return
		}
		err := s.StopServant()
		if err != nil {
			fmt.Println("x:", err)
			c.AbortWithStatusJSON(200, handlers.Resp(-1, "关闭服务异常", err.Error()))
			return
		}
		delete(utils2.SubServants, s.GetContextName())
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
