package address

import (
	mad "book-nest/internal/models/address"
	hh "book-nest/utils/handlerhelper"
	jh "book-nest/utils/jwthelper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AddressHandler struct {
	AddressService mad.AddressService
}

func NewAddressHandler(addressService mad.AddressService) mad.AddressHandler {
	return &AddressHandler{AddressService: addressService}
}

func (hdl *AddressHandler) Create(c *gin.Context) {
	addressReq := mad.AddressCreateRequest{}

	errBind := c.ShouldBindJSON(&addressReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "book handler",
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
