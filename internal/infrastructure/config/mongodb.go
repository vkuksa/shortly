package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type MongoDB struct {
	connectionString string
	db               string
}

func (c MongoDB) ConnectionString() string {
	return c.connectionString
}

func (c MongoDB) DB() string {
	return c.db
}

func NewMongo() MongoDB {
	var cfg struct {
		ConnectionString string `env:"MONGO_CONNECTION_STRING"`
		Db               string `env:"MONGO_DB"`
	}
	if err := env.Parse(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		log.Fatalln(err)
	}
	return MongoDB{
		connectionString: cfg.ConnectionString,
		db:               cfg.Db,
	}
}
