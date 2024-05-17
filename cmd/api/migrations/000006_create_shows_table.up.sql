CREATE TABLE IF NOT EXISTS shows (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone DEFAULT NOW(),
    updated_at timestamp(0) with time zone DEFAULT NOW(),
    title text NOT NULL,
    release_date timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    description text NOT NULL,
    genre text NOT NULL,
    duration integer NOT NULL
);
