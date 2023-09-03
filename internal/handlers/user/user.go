package user

import (
	mu "book-nest/internal/models/user"
	hh "book-nest/utils/handlerhelper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	UserService mu.UserService
}

func NewUserHandler(userService mu.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (hdl *UserHandler) Create(c *gin.Context) {
	userReq := mu.UserCreateRequest{}

	err := c.ShouldBindJSON(&userReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "user handler",
		"data":  userReq,
	})
	if err != nil {
		logger.WithError(err).Error("failed to bind user")
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	result, errCreate := hdl.UserService.Create(userReq)

	if errCreate != nil {
		logger.WithError(errCreate).Error("failed to create user")
		c.JSON(http.StatusInternalServerError, errCreate.Error())
		return
	}

	c.JSON(http.StatusCreated, hh.ResponseData{
		Message: "success",
		Data:    result,
	})
}
