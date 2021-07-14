package worker

import (
	"context"
	"crontab/common"
	"crontab/worker/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type LogSink struct {
	Client        *mongo.Client
	LogCollecting *mongo.Collection
	LogChan       chan *common.JobLog
}

var G_logSink *LogSink

func (l *LogSink) InsertData(logs []interface{}) {
	_, _ = l.LogCollecting.InsertMany(context.TODO(), logs)
}

func (l *LogSink) WriteLog() {
	var (
		log *common.JobLog
		logs []interface{}
	)
	for {
		select {
		case log = <-l.LogChan:
			if len(logs) == 0 {
				logs = make([]interface{},0)
			}
			logs = append(logs,log)
			
		}
	}
}

func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)
	option := options.Client().
		ApplyURI(config.G_Config.MongoUrl).
		SetConnectTimeout(time.Duration(config.G_Config.MongoTime) * time.Millisecond)
	//链接mongoDB
	if client, err = mongo.Connect(context.TODO(), option); err != nil {
		return
	}
	G_logSink = &LogSink{
		Client:        client,
		LogCollecting: client.Database("cron").Collection("log"),
		LogChan:       make(chan *common.JobLog, 1000),
	}
	return
}
