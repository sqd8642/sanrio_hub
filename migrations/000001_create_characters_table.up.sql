CREATE TABLE IF NOT EXISTS movies (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    debut timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    description text NOT NULL,
    hobbies text NOT NULL,
    affiliations text[],
    version integer NOT NULL DEFAULT 1
);
