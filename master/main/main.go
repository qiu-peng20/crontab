package main

import (
	apiServer2 "crontab/master/apiServer"
	"crontab/master/config"
	"flag"
	"fmt"
	"runtime"
)
var confFile string //配置文件的路径

func initArgs()  {
	flag.StringVar(&confFile, "config", "./master/config/config.json","传入配置项的值")
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
	err := apiServer2.InitApiServer()
	if err != nil {
		goto ERR
	}
	fmt.Print("服务启动成功")
	return
	ERR:
		fmt.Print(err)
}
