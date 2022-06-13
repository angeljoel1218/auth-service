package registry

import (
	"auth-service/config"
	routers "auth-service/src/infrastructure/router"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm"
)

type registry struct {
	db *gorm.DB
}

type Registry interface {
	StartHandlerInterface() (*gin.Engine, error)
}

func NewRegistry(db *gorm.DB) Registry {
	return &registry{db: db}
}

func (r *registry) StartHandlerInterface() (*gin.Engine, error) {
	logger := iniLogger()

	handler := routers.NewHandlerInterface(&routers.HandlerInterface{
		UserService: r.NewUserSerivice(),
		ConfigRouter: routers.Config{
			MaxBodyBytes: config.C.Route.MaxBodyBytes,
			ApiURI:       config.C.Route.ApiUri,
			Logger:       logger,
			GinMode:      config.C.Route.GinMode,
		},
	})

	return handler, nil
}

func iniLogger() *zerolog.Logger {

	logger := zerolog.New(&lumberjack.Logger{
		Filename:   config.C.Logger.FileName,
		MaxSize:    config.C.Logger.FileSize,
		MaxBackups: config.C.Logger.MaxBackup,
		MaxAge:     config.C.Logger.MaxAge,
		Compress:   config.C.Logger.FileCompress,
	})

	logger = logger.With().Caller().Timestamp().Logger()

	return &logger
}
