package config

import (
	"log"

	"github.com/caarlos0/env/v6"
)

type Server struct {
	port int
}

func (c Server) Port() int {
	return c.port
}

func NewServer() Server {
	var cfg struct {
		Port int `env:"APP_PORT" envDefault:"80"`
	}
	if err := env.Parse(&cfg, env.Options{RequiredIfNoDef: true}); err != nil {
		log.Fatalln(err)
	}
	return Server{
		port: cfg.Port,
	}
}
