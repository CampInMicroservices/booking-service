package api

import (
	"booking-service/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getBookingRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getBookingListRequest struct {
	Offset int32 `form:"offset"`
	Limit  int32 `form:"limit" binding:"required,min=1,max=20"`
}

type createBookingRequest struct {
	Name string `json:"name" binding:"required"`
}

func (server *Server) GetBookingByID(ctx *gin.Context) {

	// Check if request has ID field in URI.
	var req getBookingRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	// Execute query.
	result, err := server.store.GetBookingByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) GetAllBookings(ctx *gin.Context) {

	// Check if request has parameters offset and limit for pagination.
	var req getBookingListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	arg := db.ListBookingParam{
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	// Execute query.
	result, err := server.store.GetAllBookings(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) CreateBooking(ctx *gin.Context) {

	// Check if request has all required fields in json body.
	var req createBookingRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
		ctx.Abort()
		return
	}

	arg := db.CreateBookingParam{
		Name: req.Name,
	}

	// Execute query.
	result, err := server.store.CreateBooking(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusCreated, result)
}
