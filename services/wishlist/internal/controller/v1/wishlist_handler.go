package controller

import (
	"net/http"

	"github.com/soumirel/wishlister/wishlist/internal/auth"
	"github.com/soumirel/wishlister/wishlist/internal/controller/v1/dto"
	wishuc "github.com/soumirel/wishlister/wishlist/internal/usecase/wish"
	wishlistuc "github.com/soumirel/wishlister/wishlist/internal/usecase/wishlist"
	wishlistpermuc "github.com/soumirel/wishlister/wishlist/internal/usecase/wishlist_permission"

	"github.com/gin-gonic/gin"
)

const (
	wishlistIdPathParam = "wishlistID"
)

type wishlistHandler struct {
	wishlistUc           *wishlistuc.WishlistUsecase
	wishUc               *wishuc.WishUsecase
	wishlistPermissionUc *wishlistpermuc.WishlistPermissionUsecase
}

func NewWishlistHandler(
	gr *gin.RouterGroup,
	wishlistuc *wishlistuc.WishlistUsecase,
	wishUc *wishuc.WishUsecase,
	wishlistPermissionUc *wishlistpermuc.WishlistPermissionUsecase,
) *wishlistHandler {
	h := &wishlistHandler{
		wishlistUc:           wishlistuc,
		wishUc:               wishUc,
		wishlistPermissionUc: wishlistPermissionUc,
	}

	gr.GET("/", h.getWishlits)
	gr.POST("/", h.createWishlist)

	const wishlistPathPart = "/:" + wishlistIdPathParam
	wishlistIdGr := gr.Group(wishlistPathPart)
	wishlistIdGr.GET("", h.getWishlist)
	wishlistIdGr.PATCH("", h.updateWishlist)
	wishlistIdGr.DELETE("", h.deleteWishlist)

	wishlistPermissiomGr := wishlistIdGr.Group("/permissions")
	wishlistPermissiomGr.POST("", h.grantWishlistPermission)
	wishlistPermissiomGr.DELETE("", h.revokeWishlistPermission)

	wishlistWishesGr := wishlistIdGr.Group("/wishes")
	NewWishHandler(wishlistWishesGr, wishUc)

	return h
}

func (h *wishlistHandler) getWishlits(c *gin.Context) {
	ctx := c.Request.Context()
	au := auth.FromCtxOrEmpty(ctx)
	cmd := wishlistuc.GetWishlistsCommand{
		RequestorUserID: au.UserID,
	}
	wishlistesList, err := h.wishlistUc.GetWishlists(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, wishlistesList)
}

func (h *wishlistHandler) getWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	au := auth.FromCtxOrEmpty(ctx)
	cmd := wishlistuc.GetWishlistCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
	}
	wishlist, err := h.wishlistUc.GetWishlist(ctx, cmd)
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
	cmd := wishlistuc.CreateWishlistCommand{
		RequestorUserID: au.UserID,
		Name:            req.WishlistName,
	}
	wishlist, err := h.wishlistUc.CreateWishlist(ctx, cmd)
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
	cmd := wishlistuc.UpdateWishlistCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
		Name:            req.WishlistName,
	}
	wishlist, err := h.wishlistUc.UpdateWishlist(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, wishlist)
}

func (h *wishlistHandler) deleteWishlist(c *gin.Context) {
	ctx := c.Request.Context()
	cmd := wishlistuc.DeleteWishlistCommand{
		RequestorUserID: auth.FromCtxOrEmpty(ctx).UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
	}
	err := h.wishlistUc.DeleteWishlist(ctx, cmd)
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
	cmd := wishlistpermuc.GrantWishlistPermissionCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
		TargetUserID:    req.UserID,
		PermissionLevel: req.PermissionLevel,
	}
	err = h.wishlistPermissionUc.GrantWishlistPermission(ctx, cmd)
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
	cmd := wishlistpermuc.RevokeWishlistPermissionCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdPathParam),
		TargetUserID:    req.UserID,
	}
	err = h.wishlistPermissionUc.RevokeWishlistPermissionCommand(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
