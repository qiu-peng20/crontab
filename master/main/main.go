package main

import (
	apiServer2 "crontab/master/apiServer"
	"fmt"
	"runtime"
)

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main()  {
	//启动api服务
	err := apiServer2.InitApiServer()
	if err != nil {
		goto ERR
	}
	return
	ERR:
		fmt.Print(err)
}
