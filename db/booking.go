package db

import (
	"context"
	"time"
)

type Booking struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateBookingParam struct {
	Name string
}

type UpdateBookingParam struct {
	Name string
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
	INSERT INTO "bookings" ("name") 
	VALUES ($1)
	RETURNING "id", "name", "created_at"
	`
	row := store.db.QueryRowContext(ctx, query, arg.Name)

	var user Booking
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.CreatedAt,
	)

	return user, err
}
