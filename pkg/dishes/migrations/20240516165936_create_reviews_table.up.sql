CREATE TABLE IF NOT EXISTS reviews (
    id                        bigserial PRIMARY KEY,
    dish_id                   INTEGER REFERENCES dishes(id) ON DELETE CASCADE,
    drink_id                  INTEGER REFERENCES drinks(id) ON DELETE CASCADE,
    rating                    INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment                   TEXT NOT NULL,
    createdAt                 timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updatedAt                 timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    CONSTRAINT check_entity CHECK (
        (dish_id IS NOT NULL AND drink_id IS NULL) OR 
        (drink_id IS NOT NULL AND dish_id IS NULL)
    )
);

CREATE INDEX IF NOT EXISTS idx_reviews_dish ON reviews (dish_id);
CREATE INDEX IF NOT EXISTS idx_reviews_drink ON reviews (drink_id);

