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
		job      common.Job
		err      error
		body     []byte
		oldJob   *common.Job
		response []byte
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
	//如果成功，返回正常应答
	response, err = common.NewResponse(1, "success", oldJob)
	if err != nil {
		fmt.Print(err)
		return
	}
	_, _ = rw.Write(response)
	return
	//如果失败，返回失败应答

ERR:
	response, _ = common.NewResponse(0, "file", nil)
	_, _ = rw.Write(response)
}

func handleJobDelete(rw http.ResponseWriter, rq *http.Request) {
	var (
		job       common.Job
		err       error
		all       []byte
		deleteJob *common.Job
		response  []byte
	)
	err = rq.ParseForm()
	if err != nil {
		goto ERR
	}
	defer rq.Body.Close()
	all, err = ioutil.ReadAll(rq.Body)
	if err != nil {
		goto ERR
	}
	err = json.Unmarshal(all, &job)
	if err != nil {
		goto ERR
	}
	deleteJob, err = jobManager.G_jobMgr.DeleteJob(job.Name)
	if err != nil {
		goto ERR
	}
	response, err = common.NewResponse(1, "success", deleteJob)
	if err != nil {
		goto ERR
	}
	_, _ = rw.Write(response)
	return
ERR:
	response, _ = common.NewResponse(1, "error", nil)
	_, _ = rw.Write(response)
}

func handleJobList(rw http.ResponseWriter, rq *http.Request) {
	var (
		newResponse []byte
	)
	list, err := jobManager.G_jobMgr.FindList()
	if err != nil {
		goto ERR
	}
	newResponse, err = common.NewResponse(1, "success", list)
	if err != nil {
		goto ERR
	}
	rw.Write(newResponse)
ERR:
	newResponse, _ = common.NewResponse(0, "error", nil)
	rw.Write(newResponse)
}

func handleJobKill(rw http.ResponseWriter, rq *http.Request) {
	var (
		readAll  []byte
		job      common.Job
		response []byte
	)
	defer rq.Body.Close()

	err := rq.ParseForm()
	if err != nil {
		goto ERR
	}
	readAll, err = ioutil.ReadAll(rq.Body)
	if err != nil {
		goto ERR
	}
	err = json.Unmarshal(readAll, &job)
	if err != nil {
		goto ERR
	}
	err = jobManager.G_jobMgr.KillJob(job.Name)
	if err != nil {
		goto ERR
	}
	response, err = common.NewResponse(1, "success", "")
	if err != nil {
		goto ERR
	}
	_, _ = rw.Write(response)
	return
ERR:
	response, _ = common.NewResponse(1, "error", nil)
	_, _ = rw.Write(response)
}

func InitApiServer() (err error) {
	//配置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

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
