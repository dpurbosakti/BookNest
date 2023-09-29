package courier

import (
	i "book-nest/internal/interfaces"
	mc "book-nest/internal/models/courier"
	hh "book-nest/utils/handlerhelper"
	jh "book-nest/utils/jwthelper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CourierHandler struct {
	CourierService i.CourierService
}

func NewCourierHandler(courierService i.CourierService) i.CourierHandler {
	return &CourierHandler{CourierService: courierService}
}

func (hdl *CourierHandler) GetBiteshipCourier(c *gin.Context) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_biteship_courier",
		"scope": "courier handler",
	})

	err := hdl.CourierService.GetBiteshipCourier()
	if err != nil {
		logger.WithError(err).Error("failed to get biteship courier")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "get data couriers from biteship success",
	})
}

func (hdl *CourierHandler) GetList(c *gin.Context) {
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_list",
		"scope": "courier handler",
	})

	result, err := hdl.CourierService.GetList()
	if err != nil {
		logger.WithError(err).Error("failed to get list")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if result == nil {
		logger.Info("data is not found")
		c.JSON(http.StatusNotFound, gin.H{"message": "data not found"})
		return
	}

	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "get data couriers success",
		Data:    result,
	})
}

func (hdl *CourierHandler) CheckRates(c *gin.Context) {
	checkRatesReq := mc.CheckRatesRequest{}

	errBind := c.ShouldBindJSON(&checkRatesReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "check_rates",
		"scope": "courier handler",
		"data":  checkRatesReq,
	})
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind check rates request")
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

	result, err := hdl.CourierService.CheckRates(userData.Id, checkRatesReq.BookId)

	if err != nil {
		logger.WithError(err).Error("failed to check rates courier")
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hh.ResponseData{
		Message: "success",
		Data:    result,
	})
}
