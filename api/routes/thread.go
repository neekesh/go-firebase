package routes

import (
	"letschat/api/controllers"
	"letschat/api/middlewares"
	"letschat/infrastructure"
)

// ThreadRoutes -> struct
type ThreadRoutes struct {
	logger           infrastructure.Logger
	router           infrastructure.Router
	threadController controllers.ThreadController
	jwtMiddleware    middlewares.JWTAuthMiddleWare
}

// NewThreadRoutes -> creates new user controller
func NewThreadRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	threadController controllers.ThreadController,
	jwtMiddleware middlewares.JWTAuthMiddleWare,

) ThreadRoutes {
	return ThreadRoutes{
		router:           router,
		logger:           logger,
		threadController: threadController,
		jwtMiddleware:    jwtMiddleware,
	}
}

func (m ThreadRoutes) Setup() {
	m.logger.Zap.Info(" Setting up thread routes")
	threads := m.router.Gin.Group("websocket").Use()
	{
		threads.GET("", m.jwtMiddleware.Handle(), m.threadController.ServeWs)
	}
}
