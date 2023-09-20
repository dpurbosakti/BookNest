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
		"func":  "midtrans_callback",
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

func (hdl *RentHandler) Accept(c *gin.Context) {
	id := c.Param("reference_id")

	logger := logrus.WithFields(logrus.Fields{
		"func":         "accept",
		"scope":        "rent handler",
		"reference_id": id,
	})

	if id == "" {
		logger.Warn("no reference id provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "no reference id provided"})
		return
	}
	errAccept := hdl.RentService.Accept(c, id)

	if errAccept != nil {
		logger.WithError(errAccept).Error("failed to accept rent")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errAccept.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func (hdl *RentHandler) Reject(c *gin.Context) {
	id := c.Param("reference_id")

	logger := logrus.WithFields(logrus.Fields{
		"func":         "reject",
		"scope":        "rent handler",
		"reference_id": id,
	})

	if id == "" {
		logger.Warn("no reference id provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "no reference id provided"})
		return
	}
	errReject := hdl.RentService.Reject(c, id)

	if errReject != nil {
		logger.WithError(errReject).Error("failed to reject rent")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errReject.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}
