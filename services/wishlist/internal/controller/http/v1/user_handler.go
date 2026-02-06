package controller

import (
	"net/http"

	"github.com/soumirel/wishlister/services/wishlist/internal/controller/http/v1/dto"
	useruc "github.com/soumirel/wishlister/services/wishlist/internal/usecase/user"

	"github.com/gin-gonic/gin"
)

const (
	userIdPathParam = "userID"
)

type userHandler struct {
	userUc *useruc.UserUsecase
}

func NewUserHandler(gr *gin.RouterGroup,
	userUc *useruc.UserUsecase) *userHandler {
	h := &userHandler{
		userUc: userUc,
	}
	gr.GET("/", h.getUsers)
	gr.POST("/", h.createUser)

	const userPathPart = "/:" + userIdPathParam
	userIdGr := gr.Group(userPathPart)
	userIdGr.GET("", h.getUser)
	userIdGr.DELETE("", h.deleteUser)

	return h
}

func (h *userHandler) getUsers(c *gin.Context) {
	ctx := c.Request.Context()
	usersList, err := h.userUc.GetUsers(ctx, useruc.GetUsersCommand{})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, usersList)
}

func (h *userHandler) getUser(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	user, err := h.userUc.GetUser(ctx, useruc.GetUserCommand{
		UserID: id,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *userHandler) createUser(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateUserRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}
	user, err := h.userUc.CreateUser(ctx, useruc.CreateUserCommand{
		Name: req.Name,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *userHandler) deleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	err := h.userUc.DeleteUser(ctx, useruc.DeleteUserCommand{
		UserID: id,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
