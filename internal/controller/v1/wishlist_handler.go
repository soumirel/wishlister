package controller

import (
	"context"
	"net/http"
	"wishlister/internal/auth"
	"wishlister/internal/controller/v1/dto"
	"wishlister/internal/domain"
	service "wishlister/internal/service/wishlist"

	"github.com/gin-gonic/gin"
)

const (
	wishlistIdParam = "wishlistID"
)

type WishlistService interface {
	GetWishlists(context.Context, service.GetWishlistsCommand) ([]domain.Wishlist, error)
	GetWishlist(context.Context, service.GetWishlistCommand) (domain.Wishlist, error)
	CreateWishlist(context.Context, service.CreateWishlistCommand) (domain.Wishlist, error)
	UpdateWishlist(context.Context, service.UpdateWishlistCommand) (domain.Wishlist, error)
	DeleteWishlist(context.Context, service.DeleteWishlistCommand) error

	GrantWishlistPermission(context.Context, service.GrantWishlistPermissionCommand) error
	RevokeWishlistPermission(context.Context, service.RevokeWishlistPermissionCommand) error
}

type wishlistHandler struct {
	wishlistService WishlistService
}

func NewWishlistHandler(
	gr *gin.RouterGroup,
	wishlistService WishlistService,
	wishesService WishService,
) *wishlistHandler {
	h := &wishlistHandler{
		wishlistService: wishlistService,
	}

	gr.GET("/", h.getWishlits)
	gr.POST("/", h.createWishlist)

	const wishlistPathPart = "/:" + wishlistIdParam

	wishlistIdGr := gr.Group(wishlistPathPart)
	wishlistIdGr.GET("", h.getWishlist)
	wishlistIdGr.PATCH("", h.updateWishlist)
	wishlistIdGr.DELETE("", h.deleteWishlist)

	wishlistPermissiomGr := wishlistIdGr.Group("/permissions")
	wishlistPermissiomGr.POST("", h.grantWishlistPermission)
	wishlistPermissiomGr.DELETE("", h.revokeWishlistPermission)

	NewWishesHandler(wishlistIdGr, wishesService)

	return h
}

func (h *wishlistHandler) getWishlits(c *gin.Context) {
	ctx := c.Request.Context()
	au := auth.FromCtxOrEmpty(ctx)
	cmd := service.GetWishlistsCommand{
		RequestorUserID: au.UserID,
	}
	wishlistesList, err := h.wishlistService.GetWishlists(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, wishlistesList)
}

func (h *wishlistHandler) getWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	au := auth.FromCtxOrEmpty(ctx)
	cmd := service.GetWishlistCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdParam),
	}
	wishlist, err := h.wishlistService.GetWishlist(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, wishlist)
}

func (h *wishlistHandler) createWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.CreateWishlistRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}
	au := auth.FromCtxOrEmpty(ctx)
	cmd := service.CreateWishlistCommand{
		RequestorUserID: au.UserID,
		Name:            req.WishlistName,
	}
	wishlist, err := h.wishlistService.CreateWishlist(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, wishlist)
}

func (h *wishlistHandler) updateWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.UpdateWishlistRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}
	au := auth.FromCtxOrEmpty(ctx)
	cmd := service.UpdateWishlistCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdParam),
		Name:            req.WishlistName,
	}
	wishlist, err := h.wishlistService.UpdateWishlist(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, wishlist)
}

func (h *wishlistHandler) deleteWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	cmd := service.DeleteWishlistCommand{
		RequestorUserID: auth.FromCtxOrEmpty(ctx).UserID,
		WishlistID:      c.Param(wishlistIdParam),
	}
	err := h.wishlistService.DeleteWishlist(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func (h *wishlistHandler) grantWishlistPermission(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.GrantWishlistPermissionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}
	au := auth.FromCtxOrEmpty(ctx)
	cmd := service.GrantWishlistPermissionCommand{
		RequestorUserID:  au.UserID,
		WishlistID:       c.Param(wishlistIdParam),
		RequestingUserID: req.UserID,
		PersmissionLevel: req.PermissionLevel,
	}
	err = h.wishlistService.GrantWishlistPermission(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}

func (h *wishlistHandler) revokeWishlistPermission(c *gin.Context) {
	ctx := c.Request.Context()
	var req dto.RevokeWishlistPermissionRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.Error(err)
		return
	}
	au := auth.FromCtxOrEmpty(ctx)
	cmd := service.RevokeWishlistPermissionCommand{
		RequestorUserID:  au.UserID,
		WishlistID:       c.Param(wishlistIdParam),
		RequestingUserID: req.UserID,
	}
	err = h.wishlistService.RevokeWishlistPermission(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
