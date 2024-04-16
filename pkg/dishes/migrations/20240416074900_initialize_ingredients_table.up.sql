CREATE TABLE IF NOT EXISTS ingredients
(
    id         bigserial PRIMARY KEY,
    createdAt  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt  timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name       text                        NOT NULL,
    quantity   integer                     NOT NULL,
    dish_id    bigserial REFERENCES dishes (id)
);
