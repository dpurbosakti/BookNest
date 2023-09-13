package rent

import (
	"book-nest/clients/midtrans"
	mr "book-nest/internal/models/rent"
	hh "book-nest/utils/handlerhelper"
	jh "book-nest/utils/jwthelper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RentHandler struct {
	RentService mr.RentService
}

func NewRentHandler(rentService mr.RentService) mr.RentHandler {
	return &RentHandler{RentService: rentService}
}

func (hdl *RentHandler) Create(c *gin.Context) {
	rentReq := new(mr.RentCreateRequest)
	errBind := c.ShouldBindJSON(&rentReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "rent handler",
		"data":  rentReq,
	})
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind rent")
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

	result, errCreate := hdl.RentService.Create(rentReq, userData.Id)

	if errCreate != nil {
		logger.WithError(errCreate).Error("failed to create rent data")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errCreate.Error()})
		return
	}

	c.JSON(http.StatusCreated, hh.ResponseData{
		Message: "success",
		Data:    result,
	})
}

func (hdl *RentHandler) MidtransCallback(c *gin.Context) {
	midtransReq := new(midtrans.MidtransRequest)
	errBind := c.ShouldBindJSON(&midtransReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "rent handler",
		"data":  midtransReq,
	})
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind midtrans request")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBind.Error()})
		return
	}
	updateReq := new(mr.RentUpdateRequest)
	updateReq.Copier(midtransReq.TransactionStatus, midtransReq.OrderId, midtransReq.TransactionTime, midtransReq.PaymentType)
	result, errUpdate := hdl.RentService.Update(updateReq)

	if errUpdate != nil {
		logger.WithError(errUpdate).Error("failed to update rent data")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errUpdate.Error()})
		return
	}

	c.JSON(http.StatusCreated, hh.ResponseData{
		Message: "success",
		Data:    result,
	})
}
