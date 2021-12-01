package controllers

import (
	apierror "auth-service/src/domain/apierrors"
	fixture "auth-service/src/domain/fixtures"
	u "auth-service/src/domain/fixtures"
	"auth-service/src/domain/models"
	"auth-service/src/usecase/interactor"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserController struct {
	UserInteractor interactor.UserInteractor
	Logger         *zerolog.Logger
}

func (h *UserController) SignUpAction(c *gin.Context) {
	user := &models.User{}
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(apierror.HttpStatus(err), u.Message(1, err.Error()))
		return
	}

	user, err := h.UserInteractor.SignUp(user)

	if err != nil {
		c.JSON(apierror.HttpStatus(err), u.Message(1, err.Error()))
		return
	}

	c.JSON(http.StatusOK, u.DataResponse(0, "User created", user))
}

func (h *UserController) SignInAction(c *gin.Context) {
	login := &models.User{}

	if err := c.ShouldBindJSON(login); err != nil {
		c.JSON(apierror.HttpStatus(err), u.Message(1, err.Error()))
		return
	}

	res, err := h.UserInteractor.SignIn(login)

	if err != nil {
		c.JSON(apierror.HttpStatus(err), u.Message(1, err.Error()))
		return
	}

	c.JSON(http.StatusOK, u.DataResponse(0, "User logged in", res))
}

func (*UserController) AuthAction(c *gin.Context) {
	const prefix = "Bearer "

	header := c.GetHeader("Authorization")
	if header == "" {
		err := apierror.NewErrorApi(apierror.BadRequest, "Authorization header missing or empty")
		c.AbortWithStatusJSON(apierror.HttpStatus(err), u.MessageError(1, err))
		return
	}

	token := header[len(prefix):]
	claims, err := fixture.DecodeToken(token)
	if err != nil {
		c.AbortWithStatusJSON(apierror.HttpStatus(err), u.MessageError(1, err))
		return
	}

	c.JSON(http.StatusOK, u.DataResponse(0, "Authenticated user", claims))
}
