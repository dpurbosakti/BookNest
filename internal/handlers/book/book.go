package book

import (
	i "book-nest/internal/interfaces"
	mb "book-nest/internal/models/book"
	hh "book-nest/utils/handlerhelper"
	"book-nest/utils/pagination"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BookHandler struct {
	BookService i.BookService
}

func NewBookHandler(bookService i.BookService) i.BookHandler {
	return &BookHandler{BookService: bookService}
}

func (hdl *BookHandler) Create(c *gin.Context) {
	bookReq := mb.BookCreateRequest{}

	errBind := c.ShouldBindJSON(&bookReq)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "create",
		"scope": "book handler",
		"data":  bookReq,
	})
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind book")
		c.JSON(http.StatusInternalServerError, gin.H{"error": errBind.Error()})
		return
	}

	result, errCreate := hdl.BookService.Create(&bookReq)

	if errCreate != nil {
		logger.WithError(errCreate).Error("failed to create book")
		c.JSON(http.StatusInternalServerError, gin.H{"message": errCreate.Error()})
		return
	}

	c.JSON(http.StatusCreated, hh.ResponseData{
		Message: "success",
		Data:    result,
	})
}

func (hdl *BookHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_detail",
		"scope": "book handler",
		"id":    id,
	})
	cnvId, err := strconv.Atoi(id)
	if err != nil {
		logger.WithError(err).Error("failed to convert id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := hdl.BookService.GetDetail(uint(cnvId))
	if err != nil {
		logger.WithError(err).Error("failed to get detail")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "get data book success",
		Data:    result,
	})
}

func (hdl *BookHandler) GetList(c *gin.Context) {
	var page pagination.Pagination
	limitInt, _ := strconv.Atoi(c.Query("limit"))
	pageInt, _ := strconv.Atoi(c.Query("page"))
	page.Limit = limitInt
	page.Page = pageInt
	page.Sort = c.Query("sort")
	page.Search = c.Query("search")
	column := "title"
	page.Column = &column
	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_list",
		"scope": "book handler",
		"data":  page,
	})

	result, err := hdl.BookService.GetList(page)
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
		Message: "get data books success",
		Data:    result,
	})
}

func (hdl *BookHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	logger := logrus.WithFields(logrus.Fields{
		"func":  "get_detail",
		"scope": "book handler",
		"id":    id,
	})
	cnvId, err := strconv.Atoi(id)
	if err != nil {
		logger.WithError(err).Error("failed to convert id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = hdl.BookService.Delete(uint(cnvId))
	if err != nil {
		logger.WithError(err).Error("failed to delete data")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete data book success"})
}

func (hdl *BookHandler) Update(c *gin.Context) {
	bookReq := new(mb.BookUpdateRequest)
	logger := logrus.WithFields(logrus.Fields{
		"func":  "update",
		"scope": "book handler",
		"data":  bookReq,
	})

	errBind := c.ShouldBindJSON(&bookReq)
	if errBind != nil {
		logger.WithError(errBind).Error("failed to bind book")
		c.JSON(http.StatusBadRequest, gin.H{"message": errBind.Error()})
		return
	}

	id := c.Param("id")
	cnvId, err := strconv.Atoi(id)
	if err != nil {
		logger.WithError(err).Error("failed to convert id")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	result, err := hdl.BookService.Update(bookReq, uint(cnvId))
	if err != nil {
		logger.WithError(err).Error("failed to update data")
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hh.ResponseData{
		Message: "update data book success",
		Data:    result,
	})
}
