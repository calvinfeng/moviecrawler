CREATE TABLE IF NOT EXISTS movies (
    id SERIAL,
    created timestamp,
    modified timestamp,
    title varchar(255),
    description text,
    link varchar(255) UNIQUE,
    imdb_rating double precision,
    release_year integer
)
