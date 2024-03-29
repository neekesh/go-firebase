package routes

import (
	"letschat/api/controllers"
	"letschat/infrastructure"
)

// ObtainJwtTokenRoutes -> struct
type ObtainJwtTokenRoutes struct {
	logger        infrastructure.Logger
	router        infrastructure.Router
	jwtController controllers.JwtAuthController
}

func NewObtainJwtTokenRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	jwtController controllers.JwtAuthController,

) ObtainJwtTokenRoutes {
	return ObtainJwtTokenRoutes{
		router:        router,
		logger:        logger,
		jwtController: jwtController,
	}
}

// Setup Obtain Jwt Token Routes
func (i ObtainJwtTokenRoutes) Setup() {
	i.logger.Zap.Info(" Setting up jwt routes")
	jwt := i.router.Gin.Group("/login")
	{
		jwt.POST("", i.jwtController.ObtainJwtToken)
		jwt.POST("/refresh", i.jwtController.RefreshJwtToken)
	}
}
