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
	LogBash       chan *common.LogBash
}

var G_logSink *LogSink

func (l *LogSink) InsertData(logs []interface{}) {
	_, _ = l.LogCollecting.InsertMany(context.TODO(), logs)
}

func (l *LogSink) WriteLog() {
	var (
		log     *common.JobLog
		logBash *common.LogBash
	)
	for {
		select {
		case log = <-l.LogChan:
			if logBash == nil {
				logBash = &common.LogBash{}
			}
			logBash.Logs = append(logBash.Logs, log)
			timer := time.AfterFunc(time.Duration(5*time.Second), func() {
				l.InsertData(logBash.Logs)
			})
			if len(logBash.Logs) > 100 {
				l.InsertData(logBash.Logs)
				logBash.Logs = nil
				timer.Stop()
			}
		case timeOutBash := <-l.LogBash :
			l.InsertData(timeOutBash.Logs)
			logBash = nil
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
		LogBash:       make(chan *common.LogBash, 1000),
	}
	go G_logSink.WriteLog()
	return
}

func (l *LogSink) AppendData(log *common.JobLog)  {
	l.LogChan <- log
}
