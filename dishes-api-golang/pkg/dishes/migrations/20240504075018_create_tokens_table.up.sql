CREATE TABLE IF NOT EXISTS tokens (
    hash bytea PRIMARY KEY,
    member_id bigint NOT NULL REFERENCES members ON DELETE CASCADE,
    expiry timestamp(0) with time zone NOT NULL,
    scope text NOT NULL
);