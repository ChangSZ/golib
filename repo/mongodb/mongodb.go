package mongodb

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/ChangSZ/golib/log"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var _ Repo = (*dbRepo)(nil)

type Repo interface {
	i()

	Connection() *mongo.Database
	Disconnect() error
}

type dbRepo struct {
	connection *mongo.Database
	client     *mongo.Client
}

type Config struct {
	Host         string `toml:"host"`
	Port         int    `toml:"port"`
	Database     string `toml:"database"`
	Monitor      bool   `toml:"monitor"`
	MonitorLevel string `toml:"monitorLevel"` // succeeded、failed、all
}

var db *dbRepo

func (d *dbRepo) i() {}

var (
	commandStartedTime sync.Map
	config             Config
)

var monitor = &event.CommandMonitor{
	Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
		if config.MonitorLevel != "all" && config.MonitorLevel != "" {
			return
		}
		if evt.CommandName == "aggregate" {
			return
		}
		_, file, line, _ := runtime.Caller(6) // 获取调用文件和行数
		commandStartedTime.Store(evt.RequestID, time.Now())
		log.WithTrace(ctx).Infof("%s:%d [Started]: %s", file, line, evt.Command)
	},
	Succeeded: func(ctx context.Context, evt *event.CommandSucceededEvent) {
		if config.MonitorLevel == "failed" {
			return
		}
		if evt.CommandName == "aggregate" {
			return
		}
		startTime, ok := commandStartedTime.Load(evt.RequestID)
		if !ok {
			return
		}
		_, file, line, _ := runtime.Caller(6) // 获取调用文件和行数
		var cnt int64
		// {\"n\": {\"$numberInt\":\"1\"},\"nModified\": {\"$numberInt\":\"1\"},\"ok\": {\"$numberDouble\":\"1.0\"}}
		switch evt.CommandName {
		case "find":
			cnt = int64(len(gjson.Get(evt.Reply.String(), "cursor.firstBatch").Array()))
		default:
			cnt = gjson.Get(evt.Reply.String(), "n.$numberInt").Int()
		}

		duration := time.Since(startTime.(time.Time))
		log.WithTrace(ctx).Infof("%s:%d[%v] [cnt:%v] [Succeeded]: %s", file, line, duration, cnt, evt.Reply)
		commandStartedTime.Delete(evt.RequestID)
	},
	Failed: func(ctx context.Context, evt *event.CommandFailedEvent) {
		if config.MonitorLevel == "succeeded" {
			return
		}
		startTime, ok := commandStartedTime.Load(evt.RequestID)
		if !ok {
			return
		}
		_, file, line, _ := runtime.Caller(6) // 获取调用文件和行数
		duration := time.Since(startTime.(time.Time))
		log.WithTrace(ctx).Errorf("%s:%d[%v] [Failed]: %s, err: %v", file, line, duration, evt.CommandName, evt.Failure)
		commandStartedTime.Delete(evt.RequestID)
	},
}

func Init(cfg Config) {
	config = cfg
	uri := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	opts := options.Client().ApplyURI(uri)
	if cfg.Monitor {
		opts.SetMonitor(monitor)
	}
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	connection := client.Database(cfg.Database)
	db = &dbRepo{client: client, connection: connection}
	log.Info("Connected to mongo db")
}

func DB() *dbRepo {
	if db == nil {
		log.Fatal("请先进行Init")
	}
	return db
}

func (d *dbRepo) Connection() *mongo.Database {
	return d.connection
}

func (d *dbRepo) Disconnect() error {
	return db.client.Disconnect(context.Background())
}
