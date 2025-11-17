package controller

import (
	"net/http"
	"wishlister/internal/domain"
	"wishlister/internal/service"

	"github.com/gin-gonic/gin"
)

type wishHandler struct {
	wishService *service.WishService
}

func NewWishHandler(gr *gin.RouterGroup,
	wishService *service.WishService) *wishHandler {
	h := &wishHandler{
		wishService: wishService,
	}
	gr.GET("/", h.getList)
	gr.GET("/:id", h.getByID)
	gr.POST("/", h.create)
	gr.DELETE("/:id", h.delete)
	return h
}

func (h *wishHandler) getList(c *gin.Context) {
	ctx := c.Request.Context()
	wishesList, err := h.wishService.GetList(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, wishesList)
}

func (h *wishHandler) getByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	wish, err := h.wishService.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, wish)
}

func (h *wishHandler) create(c *gin.Context) {
	ctx := c.Request.Context()
	var req domain.Wish
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	wish, err := h.wishService.Create(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, wish)
}

func (h *wishHandler) delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	err := h.wishService.Delete(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
