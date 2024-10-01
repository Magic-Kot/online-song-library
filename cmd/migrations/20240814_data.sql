-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS songs
(
    id            SERIAL         PRIMARY KEY,
    group_song    VARCHAR        NOT NULL,
    song          VARCHAR        NOT NULL,
    release_date  VARCHAR,
    text          VARCHAR,
    link          VARCHAR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd