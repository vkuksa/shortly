package app

import (
	"context"

	"github.com/go-chi/chi/v5"
	"github.com/vkuksa/shortly/internal/infrastructure/http"
	"github.com/vkuksa/shortly/internal/interface/controller/rest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
)

func RegisterServer(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, server http.Server, controller rest.LinkController) {
	router, ok := server.Handler.(*chi.Mux)
	if !ok {
		panic("handler is not a chi.Router")
	}
	controller.Register(router)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := server.Run(); err != nil {
					shutdowner.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}

func RegisterMongoDBClient(lifecycle fx.Lifecycle, client *mongo.Client) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Ping(ctx, nil)
		},
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})
}
