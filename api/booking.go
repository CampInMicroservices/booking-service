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

type BookingResponse struct {
	Booking db.Booking
	Payment *proto.PaymentResponse
}

// @BasePath /booking-service/v1

// Booking godoc
// @Summary Bookings by ID
// @Schemes
// @Description Returns booking by ID
// @Tags Bookings
// @Accept json
// @Produce json
// @Param id path int true "Booking ID"
// @Success 200 {array} db.Booking
// @Router /v1/bookings/{id} [get]
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

// @BasePath /booking-service/v1

// Booking godoc
// @Summary Bookings list
// @Schemes
// @Description Returns bookings by ID
// @Tags Bookings
// @Accept json
// @Produce json
// @Param limit query int true "Limit"
// @Param offset query int true "Offset"
// @Success 200 {array} db.Booking
// @Router /v1/bookings [get]
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

// @BasePath /booking-service/v1

// Booking godoc
// @Summary Bookings create
// @Schemes
// @Description Creates a booking
// @Tags Bookings
// @Accept json
// @Produce json
// @Param request body db.Booking true "Booking"
// @Success 200 {array} db.Booking
// @Router /v1/bookings [post]
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

	response := BookingResponse{Booking: result, Payment: paymentResponse}

	ctx.JSON(http.StatusCreated, response)
}
