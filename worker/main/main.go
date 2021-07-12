package main

import (
	"crontab/worker/config"
	"crontab/worker/schedule"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var confFile string //配置文件的路径

func initArgs() {
	flag.StringVar(&confFile, "config", "./worker/config/config.json", "传入配置项的值")
	flag.Parse()
}

func init() {
	initArgs()
	err := config.InitConfig(confFile) // 初始化连接器
	if err != nil {
		goto ERR
	}
	schedule.InitExecutor() //初始化执行器

	runtime.GOMAXPROCS(runtime.NumCPU())

	schedule.InitSchedule() //初始化调度器

ERR:
	fmt.Print(err)
}

func main() {
	err := schedule.InitJobMgr() //初始化api服务
	if err != nil {
		fmt.Print(err)
		goto ERR
	}
	fmt.Print("服务启动成功")
	for {
		time.Sleep(time.Second)
	}
ERR:
	fmt.Print(err)
}
