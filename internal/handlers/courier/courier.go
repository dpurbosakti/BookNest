package courier

import (
	mc "book-nest/internal/models/courier"
	hh "book-nest/utils/handlerhelper"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CourierHandler struct {
	CourierService mc.CourierService
}

func NewCourierHandler(courierService mc.CourierService) mc.CourierHandler {
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
