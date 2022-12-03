# Booking service

## Endpoints

```
GET  localhost:8080/v1/bookings/:id
GET  localhost:8080/v1/bookings?offset=0&limit=10
POST localhost:8080/v1/bookings
```

Booking JSON:
```
{
  "user_id": 1,
  "listing_id": 2,
  "number_of_adults": 2,
  "number_of_children": 0,
  "number_of_pets": 0
}
```