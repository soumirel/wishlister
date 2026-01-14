package controller

import (
	"net/http"

	"github.com/soumirel/wishlister/wishlist/internal/auth"
	"github.com/soumirel/wishlister/wishlist/internal/controller/http/v1/dto"
	wishuc "github.com/soumirel/wishlister/wishlist/internal/usecase/wish"

	"github.com/gin-gonic/gin"
)

const (
	wishIdPathParam = "wishID"
)

type wishHandler struct {
	wishUc *wishuc.WishUsecase
}

func NewWishHandler(wishesGr *gin.RouterGroup,
	wishesService *wishuc.WishUsecase,
) *wishHandler {
	h := &wishHandler{
		wishUc: wishesService,
	}

	wishesGr.GET("/", h.getWishes)
	wishesGr.POST("/", h.createWish)

	const wishPathPart = "/:" + wishIdPathParam
	wishIdGr := wishesGr.Group(wishPathPart)

	wishIdGr.GET("", h.getWish)
	wishIdGr.PATCH("/", h.updateWish)
	wishIdGr.DELETE("/", h.deleteWish)

	wishIdGr.POST("/reserve", h.reserveWish)

	return h
}

func (h *wishHandler) getWishes(c *gin.Context) {
	ctx := c.Request.Context()
	au := auth.FromCtxOrEmpty(ctx)
	cmd := wishuc.GetWishesFromWishlistCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
	}
	wishes, err := h.wishUc.GetWishesFromWishlist(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, wishes)
}

func (h *wishHandler) getWish(c *gin.Context) {
	ctx := c.Request.Context()
	au := auth.FromCtxOrEmpty(ctx)
	cmd := wishuc.GetWishCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
		WishID:          c.Param(wishIdPathParam),
	}
	wish, err := h.wishUc.GetWish(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, wish)
}

func (h *wishHandler) createWish(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateWishRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}
	au := auth.FromCtxOrEmpty(ctx)
	cmd := wishuc.CreateWishCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
		WishName:        req.WishName,
	}
	wishlist, err := h.wishUc.CreateWish(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, wishlist)
}

func (h *wishHandler) updateWish(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.UpdateWishRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}
	au := auth.FromCtxOrEmpty(ctx)
	cmd := wishuc.UpdateWishCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
		WishID:          c.Param(wishIdPathParam),
		WishName:        req.WishName,
	}
	wishlist, err := h.wishUc.UpdateWish(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, wishlist)
}

func (h *wishHandler) deleteWish(c *gin.Context) {
	ctx := c.Request.Context()
	au := auth.FromCtxOrEmpty(ctx)
	cmd := wishuc.DeleteWishCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
		WishID:          c.Param(wishIdPathParam),
	}
	err := h.wishUc.DeleteWish(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *wishHandler) reserveWish(c *gin.Context) {
	ctx := c.Request.Context()
	au := auth.FromCtxOrEmpty(ctx)
	cmd := wishuc.ReserveWishCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
		WishID:          c.Param(wishIdPathParam),
	}
	err := h.wishUc.ReserveWish(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
