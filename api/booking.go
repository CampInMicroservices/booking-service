package api

import (
	"booking-service/db"
	"booking-service/proto"
	"context"
	"log"
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
	UserID           int64  `json:"user_id" binding:"required"`
	ListingID        int64  `json:"listing_id" binding:"required"`
	NumberOfAdults   *int64 `json:"number_of_adults" binding:"required"`
	NumberOfChildren *int64 `json:"number_of_children" binding:"required"`
	NumberOfPets     *int64 `json:"number_of_pets" binding:"required"`
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
		UserID:           req.UserID,
		ListingID:        req.ListingID,
		NumberOfAdults:   *req.NumberOfAdults,
		NumberOfChildren: *req.NumberOfChildren,
		NumberOfPets:     *req.NumberOfPets,
	}

	// Execute query.
	result, err := server.store.CreateBooking(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err})
		ctx.Abort()
		return
	}

	paymentRequest := &proto.PaymentRequest{
		Payment: &proto.Payment{
			BookingId: result.ID,
			Price:     95.5,
			Paid:      false,
		},
	}

	// Create new payment via gRPC
	paymentResponse, err := server.grpcClient.CreatePaymentRequest(context.Background(), paymentRequest)
	if err != nil {
		log.Println("Payment service not reachable.")
	}

	if paymentResponse != nil {
		log.Println("New payment created: ", paymentResponse)
	}

	ctx.JSON(http.StatusCreated, result)
}
