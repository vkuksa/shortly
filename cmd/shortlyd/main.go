package main

import (
	"github.com/vkuksa/shortly/internal/app"
	"github.com/vkuksa/shortly/internal/infrastructure/config"
)

func main() {
	app.New(
		config.NewMongo(),
		config.NewServer(),
	).Run()
}
