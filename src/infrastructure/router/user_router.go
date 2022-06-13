package routers

import (
	controllers "auth-service/src/interface/controller"

	"github.com/gin-gonic/gin"
)

func (h *HandlerInterface) userRoutes(group *gin.RouterGroup) {
	uc := &controllers.UserController{
		UserService: h.UserService,
		Logger:      h.ConfigRouter.Logger,
	}
	group.POST("/auth/login", uc.SignInAction)
	group.POST("/auth/signup", uc.SignUpAction)
	group.POST("/auth/authenticated", uc.AuthAction)
}
