package user

import (
	mu "book-nest/internal/models/user"
	hh "book-nest/utils/handlerhelper"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService mu.UserService
}

func NewUserHandler(userService mu.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (hdl *UserHandler) Create(c *gin.Context) {
	userReq := mu.UserCreateRequest{}
	err := c.Bind(&userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	result, errCreate := hdl.userService.Create(userReq)

	if errCreate != nil {
		c.JSON(http.StatusInternalServerError, errCreate)
		return
	}

	c.JSON(http.StatusCreated, hh.ResponseData{
		Message: "success",
		Data:    result,
	})
}
