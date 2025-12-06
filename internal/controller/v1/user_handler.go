package controller

// import (
// 	"net/http"
// 	"wishlister/internal/domain"
// 	"wishlister/internal/service"

// 	"github.com/gin-gonic/gin"
// )

// type userHandler struct {
// 	userService *service.UserService
// }

// func NewUserHandler(gr *gin.RouterGroup,
// 	userService *service.UserService) *userHandler {
// 	h := &userHandler{
// 		userService: userService,
// 	}
// 	gr.GET("/", h.getList)
// 	gr.GET("/:id", h.getByID)
// 	gr.POST("/", h.create)
// 	gr.DELETE("/:id", h.delete)
// 	return h
// }

// func (h *userHandler) getList(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	usersList, err := h.userService.GetList(ctx)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, usersList)
// }

// func (h *userHandler) getByID(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	id := c.Param("id")
// 	user, err := h.userService.GetByID(ctx, id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, user)
// }

// func (h *userHandler) create(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	var req domain.User
// 	err := c.BindJSON(&req)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	user, err := h.userService.Create(ctx, req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusCreated, user)
// }

// func (h *userHandler) delete(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	id := c.Param("id")
// 	err := h.userService.Delete(ctx, id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{})
// }
