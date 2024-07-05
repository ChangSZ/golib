package mongodb

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestMongoDB(t *testing.T) {
	t.Skip("skip")
	// 结构化
	cfg := Config{
		Host:         "localhost",
		Port:         27017,
		Database:     "local",
		Monitor:      true,
		MonitorLevel: "all",
	}
	Init(cfg)
	filter := bson.M{}
	DB().Connection().Collection("startup_log").FindOne(context.TODO(), filter)
}
