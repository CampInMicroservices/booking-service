CREATE TABLE "bookings" (
    "id"                    BIGSERIAL PRIMARY KEY,
    "user_id"               BIGINT NOT NULL,
    "listing_id"            BIGINT NOT NULL,
    "number_of_adults"      INT NOT NULL,
    "number_of_children"    INT NOT NULL,
    "number_of_pets"        INT NOT NULL,
    "created_at"            TIMESTAMP NOT NULL DEFAULT(now())
);