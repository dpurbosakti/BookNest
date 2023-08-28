package user

import (
	mu "book-nest/internal/models/user"
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
		c.JSON(http.StatusInternalServerError, "tba")
		return
	}
	// dataUser := requestUser.ToCore(userReq)
	_, errCreate := hdl.userService.Create(userReq)

	if errCreate != nil {
		c.JSON(http.StatusInternalServerError, "tba")
		return
	}

	c.JSON(http.StatusOK, "tba")

	//c.JSON(http.StatusOK, gin.H{
	// "message": "success",
}
