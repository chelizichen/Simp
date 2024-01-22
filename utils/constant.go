package utils

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

const PublishPath = "static/publish/"

var ServantAlives = make(map[string]int, 1024)

func ServantMonitor() {
	go func() {
		c := cron.New()

		// 4小时执行一次，更换日志文件指定目录
		spec := "1 * * * * *"
		// 添加定时任务
		err := c.AddFunc(spec, func() {
			d := time.Now().Format(time.DateTime)
			for serverName, pid := range ServantAlives {
				b := IsPidAlive(pid, serverName)
				if b {
					fmt.Println(d, "| IsAlive | ", serverName)
				} else {
					fmt.Println(d, "| IsNotAlive | ", serverName)
				}
			}
		})
		if err != nil {
			fmt.Println("AddFuncErr", err)
		}
		// 启动Cron调度器
		go c.Start()
	}()
}
