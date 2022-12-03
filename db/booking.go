package db

import (
	"context"
	"time"
)

type Booking struct {
	ID               int64     `json:"id" db:"id"`
	UserID           int64     `json:"user_id" db:"user_id"`
	ListingID        int64     `json:"listing_id" db:"listing_id"`
	NumberOfAdults   int64     `json:"number_of_adults" db:"number_of_adults"`
	NumberOfChildren int64     `json:"number_of_children" db:"number_of_children"`
	NumberOfPets     int64     `json:"number_of_pets" db:"number_of_pets"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type CreateBookingParam struct {
	UserID           int64
	ListingID        int64
	NumberOfAdults   int64
	NumberOfChildren int64
	NumberOfPets     int64
}

type UpdateBookingParam struct {
	UserID           int64
	ListingID        int64
	NumberOfAdults   int64
	NumberOfChildren int64
	NumberOfPets     int64
}

type ListBookingParam struct {
	Offset int32
	Limit  int32
}

func (store *Store) GetBookingByID(ctx context.Context, id int64) (booking Booking, err error) {

	const query = `SELECT * FROM "bookings" WHERE "id" = $1`
	err = store.db.GetContext(ctx, &booking, query, id)

	return
}

func (store *Store) GetAllBookings(ctx context.Context, arg ListBookingParam) (booking []Booking, err error) {

	const query = `SELECT * FROM "bookings" OFFSET $1 LIMIT $2`
	booking = []Booking{}
	err = store.db.SelectContext(ctx, &booking, query, arg.Offset, arg.Limit)

	return
}

func (store *Store) CreateBooking(ctx context.Context, arg CreateBookingParam) (Booking, error) {

	const query = `
	INSERT INTO "bookings" ("user_id", "listing_id", "number_of_adults", "number_of_children", "number_of_pets") 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING "id", "user_id", "listing_id", "number_of_adults", "number_of_children", "number_of_pets", "created_at"
	`
	row := store.db.QueryRowContext(ctx, query, arg.UserID, arg.ListingID, arg.NumberOfAdults, arg.NumberOfChildren, arg.NumberOfPets)

	var booking Booking
	err := row.Scan(
		&booking.ID,
		&booking.UserID,
		&booking.ListingID,
		&booking.NumberOfAdults,
		&booking.NumberOfChildren,
		&booking.NumberOfPets,
		&booking.CreatedAt,
	)

	return booking, err
}
