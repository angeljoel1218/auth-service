package routers

import (
	"auth-service/src/interface/middleware"
	"auth-service/src/usecase/interactor"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Config struct {
	MaxBodyBytes int64
	ApiURI       string
	Logger       *zerolog.Logger
	GinMode      string
}

type HandlerInterface struct {
	UserInteractor interactor.UserInteractor
	ConfigRouter   Config
}

func NewHandlerInterface(h *HandlerInterface) *gin.Engine {
	gin.SetMode(h.ConfigRouter.GinMode)

	r := gin.Default()

	r.Use(middleware.Logger(middleware.LoggerConfig{
		SkipPaths: []string{},
		Logger:    h.ConfigRouter.Logger,
	}))

	r.Use(middleware.Recovery(h.ConfigRouter.Logger))

	g := r.Group(h.ConfigRouter.ApiURI)

	h.userRoutes(g)

	return r
}
