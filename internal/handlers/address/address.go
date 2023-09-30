package address

import (
	i "book-nest/internal/interfaces"
	mad "book-nest/internal/models/address"
	hh "book-nest/utils/handlerhelper"
	jh "book-nest/utils/jwthelper"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddressHandler struct {
	AddressService i.AddressService
}

func NewAddressHandler(addressService i.AddressService) i.AddressHandler {
	return &AddressHandler{AddressService: addressService}
}

func (hdl *AddressHandler) Create(c *gin.Context) {
	addressReq := mad.AddressCreateRequest{}

	errBind := c.ShouldBindJSON(&addressReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "address handler",
		"data":  addressReq,
	})
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind address")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBind.Error()})
		return
	}

	userData, err := jh.ParseToken(c)
	if err != nil {
		logger.WithError(err).Error("failed to parse token")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	logger.WithField("user_id", userData.Id).Info()

	addressReq.UserId = userData.Id

	result, errCreate := hdl.AddressService.Create(&addressReq)

	if errCreate != nil {
		logger.WithError(errCreate).Error("failed to create address")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errCreate.Error()})
		return
	}

	c.JSON(http.StatusCreated, hh.ResponseData{
		Message: "success",
		Data:    result,
	})
}

func (hdl *AddressHandler) Update(c *gin.Context) {
	addressReq := new(mad.AddressUpdateRequest)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "update",
		"scope": "address handler",
		"data":  addressReq,
	})

	errBind := c.ShouldBindJSON(&addressReq)
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind address")
		c.JSON(http.StatusBadRequest, gin.H{"message": errBind.Error()})
		return
	}
	userData, err := jh.ParseToken(c)
	if err != nil {
		logger.WithError(err).Error("failed to parse token")
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id := c.Param("id")
	cnvId, err := strconv.Atoi(id)
	if err != nil {
		logger.WithError(err).Error("failed to convert id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := hdl.AddressService.Update(addressReq, uint(cnvId), userData.Id)
	if err != nil && errors.Is(errors.New("Unauthorized"), err) {
		logger.WithError(err).Error("Unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
	} else if err != nil {
		logger.WithError(err).Error("failed to update data")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "update data address success",
		Data:    result,
	})
}

func (hdl *AddressHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_detail",
		"scope": "address handler",
		"id":    id,
	})
	cnvId, err := strconv.Atoi(id)
	if err != nil {
		logger.WithError(err).Error("failed to convert id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := hdl.AddressService.GetDetail(uint(cnvId))
	if err != nil {
		logger.WithError(err).Error("failed to get detail")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "get data address success",
		Data:    result,
	})
}
