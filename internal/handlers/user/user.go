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

	result, errCreate := hdl.UserService.Create(&userReq)

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

func (hdl *UserHandler) Verify(c *gin.Context) {
	var userReq mu.UserVerifyRequest

	errBind := c.ShouldBindJSON(&userReq)
	if errBind != nil {
		c.JSON(http.StatusBadRequest, errBind)
		return
	}

	err := hdl.UserService.Verify(&userReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, "account is verified")
}

func (hdl *UserHandler) RefreshVerCode(c *gin.Context) {
	userReq := new(mu.UserVerificationCodeRequest)

	errBind := c.ShouldBindJSON(&userReq)
	if errBind != nil {
		c.JSON(http.StatusBadRequest, errBind)
		return
	}

	err := hdl.UserService.RefreshVerificationCode(userReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, "new verification code has been sent to your registered email")
}
