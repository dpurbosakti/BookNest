package order

import (
	"book-nest/clients/midtrans"
	mo "book-nest/internal/models/order"
	hh "book-nest/utils/handlerhelper"
	jh "book-nest/utils/jwthelper"
	"book-nest/utils/pagination"
	"net/http"
	"strconv"

	i "book-nest/internal/interfaces"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	OrderService i.OrderService
}

func NewOrderHandler(orderService i.OrderService) i.OrderHandler {
	return &OrderHandler{OrderService: orderService}
}

func (hdl *OrderHandler) Create(c *gin.Context) {
	orderReq := new(mo.OrderCreateRequest)
	errBind := c.ShouldBindJSON(&orderReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "order handler",
		"data":  orderReq,
	})
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind order")
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

	result, errCreate := hdl.OrderService.Create(orderReq, userData.Id)

	if errCreate != nil {
		logger.WithError(errCreate).Error("failed to create order data")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errCreate.Error()})
		return
	}

	c.JSON(http.StatusCreated, hh.ResponseData{
		Message: "success",
		Data:    result,
	})
}

func (hdl *OrderHandler) MidtransCallback(c *gin.Context) {
	midtransReq := new(midtrans.MidtransRequest)
	errBind := c.ShouldBindJSON(&midtransReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "midtrans_callback",
		"scope": "order handler",
		"data":  midtransReq,
	})
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind midtrans request")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBind.Error()})
		return
	}
	updateReq := new(mo.OrderUpdateRequest)
	updateReq.Copier(midtransReq.TransactionStatus, midtransReq.OrderId, midtransReq.TransactionTime, midtransReq.PaymentType)
	result, errUpdate := hdl.OrderService.Update(updateReq)

	if errUpdate != nil {
		logger.WithError(errUpdate).Error("failed to update order data")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errUpdate.Error()})
		return
	}

	c.JSON(http.StatusCreated, hh.ResponseData{
		Message: "success",
		Data:    result,
	})
}

func (hdl *OrderHandler) Accept(c *gin.Context) {
	id := c.Param("reference_id")

	logger := logrus.WithFields(logrus.Fields{
		"func":         "accept",
		"scope":        "order handler",
		"reference_id": id,
	})

	if id == "" {
		logger.Warn("no reference id provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "no reference id provided"})
		return
	}
	errAccept := hdl.OrderService.Accept(c, id)

	if errAccept != nil {
		logger.WithError(errAccept).Error("failed to accept order")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errAccept.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func (hdl *OrderHandler) Reject(c *gin.Context) {
	id := c.Param("reference_id")

	logger := logrus.WithFields(logrus.Fields{
		"func":         "reject",
		"scope":        "order handler",
		"reference_id": id,
	})

	if id == "" {
		logger.Warn("no reference id provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "no reference id provided"})
		return
	}
	errReject := hdl.OrderService.Reject(c, id)

	if errReject != nil {
		logger.WithError(errReject).Error("failed to reject order")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errReject.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func (hdl *OrderHandler) GetDetail(c *gin.Context) {
	id := c.Param("reference_id")

	logger := logrus.WithFields(logrus.Fields{
		"func":         "reject",
		"scope":        "order handler",
		"reference_id": id,
	})

	if id == "" {
		logger.Warn("no reference id provided")
		c.JSON(http.StatusBadRequest, gin.H{"message": "no reference id provided"})
		return
	}

	data, errGet := hdl.OrderService.GetDetail(id)

	if errGet != nil {
		logger.WithError(errGet).Error("failed to get data order")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errGet.Error()})
		return
	}

	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "success",
		Data:    data,
	})
}

func (hdl *OrderHandler) GetList(c *gin.Context) {
	var page pagination.Pagination
	limitInt, _ := strconv.Atoi(c.Query("limit"))
	pageInt, _ := strconv.Atoi(c.Query("page"))
	page.Limit = limitInt
	page.Page = pageInt
	page.Sort = c.Query("sort")
	page.Search = c.Query("search")
	column := "user_id"
	page.Column = &column
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_list",
		"scope": "order handler",
		"data":  page,
	})

	result, err := hdl.OrderService.GetList(page)
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
		Message: "get data orders success",
		Data:    result,
	})
}
