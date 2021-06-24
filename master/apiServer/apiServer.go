package apiServer

import (
	"crontab/master/config"
	"fmt"
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

func handleJobSave(rw http.ResponseWriter, rp *http.Request) {

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
