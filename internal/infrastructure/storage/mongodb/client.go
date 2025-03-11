package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

type Config interface {
	ConnectionString() string
	DB() string
}

func NewDatabase(cfg Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client()
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(cfg.ConnectionString())
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	return client.Database(cfg.DB()), nil
}
