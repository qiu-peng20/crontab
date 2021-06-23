package apiServer

import (
	"fmt"
	"net"
	"net/http"
	"time"
)
type ApiServer struct {
	httpServer *http.Server
}

var (
	HttpServer *ApiServer
)

func handleJobSave(rw http.ResponseWriter, rp *http.Request)  {

}

func InitApiServer() (err error)  {
	mux :=  http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	
	//启动http监听
	listen, err := net.Listen("tcp", ":8888")
	if err != nil {
		return err
	}
	httpServer := &http.Server{
		ReadTimeout: 5*time.Second,
		WriteTimeout: 5*time.Second,
		Handler: mux,
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
