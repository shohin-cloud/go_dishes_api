CREATE TABLE IF NOT EXISTS dishes
(
    id          bigserial PRIMARY KEY,
    createdAt   timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt   timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name        text                        NOT NULL,
    description text                        NOT NULL,
    price       numeric(10, 2)              NOT NULL
);
