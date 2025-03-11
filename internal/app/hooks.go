package app

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/vkuksa/shortly/internal/infrastructure/http"
	"github.com/vkuksa/shortly/internal/interface/controller/rest"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/sdk/trace"
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
					slog.Info("Server exited with error", slog.Any("error", err))
					_ = shutdowner.Shutdown()
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Shutting down server...")
			return server.Shutdown(ctx)
		},
	})
}

func RegisterMongoDBClient(lifecycle fx.Lifecycle, db *mongo.Database) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return db.Client().Ping(ctx, nil)
		},
		OnStop: func(ctx context.Context) error {
			slog.Info("Shutting down mongo connection...")
			return db.Client().Disconnect(ctx)
		},
	})
}

func RegisterTracer(lifecycle fx.Lifecycle, tracer *trace.TracerProvider) {
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			slog.Info("Shutting down tracer...")
			return tracer.Shutdown(ctx)
		},
	})
}
