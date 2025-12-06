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
	wishIdParam = "wishID"
)

type WishService interface {
	GetWish(context.Context, service.GetWishCommand) (domain.Wish, error)
	CreateWish(context.Context, service.CreateWishCommand) (domain.Wish, error)
	UpdateWish(context.Context, service.UpdateWishCommand) (domain.Wish, error)
	DeleteWish(context.Context, service.DeleteWishCommand) error
}

// type WishesReservationService interface {
// 	ReserveWish(context.Context, domain.Reservation) error
// 	CancelWishReservation(context.Context, domain.Reservation) error
// }

type wishHandler struct {
	wishesService WishService
}

func NewWishesHandler(wishlistGr *gin.RouterGroup,
	wishesService WishService,
) *wishHandler {
	h := &wishHandler{
		wishesService: wishesService,
	}

	wishlistGr.POST("/", h.createWish)

	const wishPathPart = "/:" + wishIdParam
	wishGr := wishlistGr.Group(wishPathPart)

	wishGr.GET("", h.getWish)
	wishGr.PATCH("/", h.updateWish)
	wishGr.DELETE("/", h.deleteWish)

	return h
}

func (h *wishHandler) getWish(c *gin.Context) {
	ctx := c.Request.Context()
	au := auth.FromCtxOrEmpty(ctx)
	cmd := service.GetWishCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdParam),
		WishID:          c.Param(wishIdParam),
	}
	wish, err := h.wishesService.GetWish(ctx, cmd)
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
	cmd := service.CreateWishCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdParam),
		WishName:        req.WishName,
	}
	wishlist, err := h.wishesService.CreateWish(ctx, cmd)
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
	cmd := service.UpdateWishCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdParam),
		WishID:          c.Param(wishIdParam),
		WishName:        req.WishName,
	}
	wishlist, err := h.wishesService.UpdateWish(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusCreated, wishlist)
}

func (h *wishHandler) deleteWish(c *gin.Context) {
	ctx := c.Request.Context()
	au := auth.FromCtxOrEmpty(ctx)
	cmd := service.DeleteWishCommand{
		RequestorUserID: au.UserID,
		WishlistID:      c.Param(wishlistIdParam),
		WishID:          c.Param(wishIdParam),
	}
	err := h.wishesService.DeleteWish(ctx, cmd)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// func (h *wishesHandler) reserve(c *gin.Context) {
// 	ctx := c.Request.Context()
// 	reservation := domain.Reservation{
// 		WishID: c.Param("id"),
// 	}
// 	err := c.BindJSON(&reservation)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	err = h.wishesService.Reserve(ctx, reservation)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		return
// 	}
// 	c.JSON(http.StatusOK, reservation)
// }
