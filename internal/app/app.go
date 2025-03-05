package app

import (
	"github.com/vkuksa/shortly/internal/infrastructure/config"
	"github.com/vkuksa/shortly/internal/infrastructure/encoder"
	"github.com/vkuksa/shortly/internal/infrastructure/http"
	"github.com/vkuksa/shortly/internal/infrastructure/storage/mongodb"
	"github.com/vkuksa/shortly/internal/interface/controller/rest"
	"github.com/vkuksa/shortly/internal/link"
	"github.com/vkuksa/shortly/internal/usecase"
	"go.uber.org/fx"
)

type App struct {
	*fx.App
}

func New(
	mongoConfig config.MongoDB,
	serverConfig config.Server,
) App {
	components := fx.Options(
		fx.Supply(
			fx.Annotate(mongoConfig, fx.As(new(mongodb.Config))),
			fx.Annotate(serverConfig, fx.As(new(http.ServerConfig))),
		),
		fx.Provide(
			mongodb.NewDatabase,
			fx.Annotate(mongodb.NewStorage, fx.As(new(link.Repository))),
			fx.Annotate(encoder.NewBase64, fx.As(new(link.Encoder))),
			link.NewFactory,
			usecase.New,
			rest.NewLinkController,
			http.NewServer,
		),
		fx.Invoke(
			RegisterServer,
			RegisterMongoDBClient,
		),
	)
	return App{
		App: fx.New(components),
	}
}
