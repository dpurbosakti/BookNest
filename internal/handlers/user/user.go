package user

import (
	mu "book-nest/internal/models/user"
	hh "book-nest/utils/handlerhelper"
	jh "book-nest/utils/jwthelper"
	"book-nest/utils/pagination"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	UserService mu.UserService
}

func NewUserHandler(userService mu.UserService) mu.UserHandler {
	return &UserHandler{UserService: userService}
}

func (hdl *UserHandler) Create(c *gin.Context) {
	userReq := mu.UserCreateRequest{}

	errBind := c.ShouldBindJSON(&userReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "user handler",
		"data":  userReq,
	})
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBind.Error()})
		return
	}

	result, errCreate := hdl.UserService.Create(&userReq)

	if errCreate != nil {
		logger.WithError(errCreate).Error("failed to create user")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errCreate.Error()})
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
	logger := logrus.WithFields(logrus.Fields{
		"func":  "verify",
		"scope": "user handler",
		"data":  userReq,
	})
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind user")
		c.JSON(http.StatusBadRequest, gin.H{"message": errBind.Error()})
		return
	}

	err := hdl.UserService.Verify(&userReq)
	if err != nil {
		logger.WithError(err).Error("failed to verify user")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "account is verified"})
}

func (hdl *UserHandler) RefreshVerificationCode(c *gin.Context) {
	userReq := new(mu.UserVerificationCodeRequest)

	errBind := c.ShouldBindJSON(&userReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "refresh_verfication_code",
		"scope": "user handler",
		"data":  userReq,
	})
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind user")
		c.JSON(http.StatusBadRequest, gin.H{"message": errBind.Error()})
		return
	}

	err := hdl.UserService.RefreshVerificationCode(userReq)
	if err != nil {
		logger.WithError(err).Error("failed to refresh verification code")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "new verification code has been sent to your registered email"})
}

func (hdl *UserHandler) GetDetail(c *gin.Context) {
	userData, err := jh.ParseToken(c)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_detail",
		"scope": "user handler",
		"data":  userData,
	})
	if err != nil {
		logger.WithError(err).Error("failed to parse token")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	result, err := hdl.UserService.GetDetail(userData.Id)
	if err != nil {
		logger.WithError(err).Error("failed to get detail")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "get data user success",
		Data:    result,
	})
}

func (hdl *UserHandler) GetList(c *gin.Context) {
	var page pagination.Pagination
	limitInt, _ := strconv.Atoi(c.Query("limit"))
	pageInt, _ := strconv.Atoi(c.Query("page"))
	page.Limit = limitInt
	page.Page = pageInt
	page.Sort = c.Query("sort")
	page.Search = c.Query("search")
	column := "name"
	page.Column = &column
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_list",
		"scope": "user handler",
		"data":  page,
	})

	result, err := hdl.UserService.GetList(page)
	if err != nil {
		logger.WithError(err).Error("failed to get list")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if result.TotalRows == 0 {
		logger.Info("data is not found")
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "get data users success",
		Data:    result,
	})
}

func (hdl *UserHandler) Delete(c *gin.Context) {
	userData, err := jh.ParseToken(c)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "delete",
		"scope": "user handler",
		"data":  userData,
	})
	if err != nil {
		logger.WithError(err).Error("failed to parse token")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	err = hdl.UserService.Delete(userData.Id)
	if err != nil {
		logger.WithError(err).Error("failed to delete data")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete data user success"})
}

func (hdl *UserHandler) Update(c *gin.Context) {
	userReq := new(mu.UserUpdateRequest)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "update",
		"scope": "user handler",
		"data":  userReq,
	})

	errBind := c.ShouldBindJSON(&userReq)
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind user")
		c.JSON(http.StatusBadRequest, gin.H{"message": errBind.Error()})
		return
	}

	if userReq.Phone != nil && !hh.PhoneValidator(*userReq.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "phone number format is incorrect"})
		return
	}

	userData, err := jh.ParseToken(c)
	logger.WithField("user id", userData.Id)
	if err != nil {
		logger.WithError(err).Error("failed to parse token")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	result, err := hdl.UserService.Update(userReq, userData.Id)
	if err != nil {
		logger.WithError(err).Error("failed to update data")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "update data user success",
		Data:    result,
	})
}
