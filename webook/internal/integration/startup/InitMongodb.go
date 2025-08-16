package startup

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func InitMongoDB() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
	}
	option := options.Client().ApplyURI("mongodb://root:12345@localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, option)
	if err != nil {
		panic(err)
	}
	return client.Database("webook")
}
