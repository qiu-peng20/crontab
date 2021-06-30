package apiServer

import (
	"crontab/common"
	"crontab/master/config"
	"crontab/master/jobManager"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	HttpServer *ApiServer
)

func handleJobSave(rw http.ResponseWriter, rq *http.Request) {
	var (
		job common.Job
		err error
		body []byte
		oldJob *common.Job
	)
	err = rq.ParseForm()
	if err != nil {
		return
	}
	//获取post结构体
	defer rq.Body.Close()
	body, err = ioutil.ReadAll(rq.Body)
	if err != nil {
		goto ERR
	}
	err = json.Unmarshal(body, &job)
	if err != nil {
		goto ERR
	}
	//获取修改前的job
	oldJob, err = jobManager.G_jobMgr.SaveJob(job)
	if err != nil {
		goto ERR
	}
	fmt.Print(oldJob)
	//如果成功，返回正常应答

	//如果失败，返回失败应答

	ERR:
}

func InitApiServer() (err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	//启动http监听
	listen, err := net.Listen("tcp", ":"+strconv.Itoa(config.G_Config.ApiPort))
	if err != nil {
		return err
	}
	httpServer := &http.Server{
		ReadTimeout:  time.Duration(config.G_Config.ReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(config.G_Config.WriteTimeout) * time.Millisecond,
		Handler:      mux,
	}
	HttpServer = &ApiServer{
		httpServer: httpServer,
	}

	//启动服务器
	go func() {
		err := httpServer.Serve(listen)
		if err != nil {
			fmt.Print(err)
			return
		}
	}()
	return
}
