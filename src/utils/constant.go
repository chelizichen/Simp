package utils

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/robfig/cron"
)

const PublishPath = "static/publish/"

var ServantAlives = make(map[string]int, 1024)

func ServantMonitor() {
	go func() {
		c := cron.New()

		// 3分钟执行一次
		spec := "*/3 * * * * *"
		// 添加定时任务
		err := c.AddFunc(spec, func() {
			d := time.Now().Format(time.DateTime)
			for serverName, pid := range ServantAlives {
				b := IsPidAlive(pid, serverName)
				if !b {
					fmt.Println(d, "| IsNotAlive | ", serverName)
				}
				mis := GetProcessMemoryInfo(pid)
				fmt.Println(d, "| IsAlive | ", serverName)
				if mis != nil {
					data, err := json.Marshal(mis)
					if err != nil {
						fmt.Println("json Marshal Error ", data)
						return
					}
					fmt.Println(d, "| Memoryinfo |", data)
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
