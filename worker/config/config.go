package config

import (
	"encoding/json"
	"io/ioutil"
)

type configData struct {
	EtcdEndPoints []string `json:"etcdEndPoints"`
	EtcdDialTimeOut int `json:"etcdDialTimeOut"`
}

var G_Config *configData

func InitConfig(fileName string)(err error)  {
	var conf configData
	file, err := ioutil.ReadFile(fileName)
	if err !=                                                                                                                                           nil {
		return err
	}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return err
	}
	G_Config = &conf
	return
}
