package main

import (
	"crontab/worker/config"
	"crontab/worker/jobManager"
	"flag"
	"fmt"
	"runtime"
	"time"
)
var confFile string //配置文件的路径

func initArgs()  {
	flag.StringVar(&confFile, "config", "./worker/config/config.json","传入配置项的值")
	flag.Parse()
}


func init()  {
	initArgs()
	err := config.InitConfig(confFile)
	if err != nil {
		return
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main()  {
	//启动api服务
	err := jobManager.InitJobMgr()
	if err != nil {
		fmt.Print(err)
		goto ERR
	}
	fmt.Print("服务启动成功")
	for  {
		time.Sleep(time.Second)
	}
	ERR:
		fmt.Print(err)
}
