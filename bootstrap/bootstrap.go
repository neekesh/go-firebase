package bootstrap

import (
	"context"
	"letschat/api/controllers"
	middlewares "letschat/api/middlewares"
	"letschat/api/repository"
	"letschat/api/routes"
	"letschat/api/services"
	"letschat/api/validators"
	infrastructure "letschat/infrastructure"
	"letschat/socket"

	"go.uber.org/fx"
)

var Module = fx.Options(
	controllers.Module,
	infrastructure.Module,
	routes.Module,
	middlewares.Module,
	validators.Module,
	services.Module,
	repository.Module,
	socket.Module,

	fx.Invoke(bootstrap),
)

func bootstrap(
	lifecycle fx.Lifecycle,
	handler infrastructure.Router,
	routes routes.Routes,
	database infrastructure.Database,
	logger infrastructure.Logger,
	env infrastructure.Env,
	chatServer *socket.WsServer,

) {
	appStop := func(context.Context) error {
		logger.Zap.Info("Stopping Application")
		// conn.Close()
		return nil
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Zap.Info("Starting Application")
			logger.Zap.Info("------------------------")
			logger.Zap.Info("------  Chat Application  📺  ------")
			logger.Zap.Info("------------------------")

			go func() {
				//middlewares.Setup()
				routes.Setup()
				logger.Zap.Info("🌱 Seeding data...")
				//seeds.Run()
				if env.ServerPort == "" {
					handler.Gin.Run()
				} else {
					handler.Gin.Run(":" + env.ServerPort)
				}
			}()
			return nil
		},
		OnStop: appStop,
	})

}
