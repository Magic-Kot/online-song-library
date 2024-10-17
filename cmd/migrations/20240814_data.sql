-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS music_group
(
    id            SERIAL         PRIMARY KEY,
    group_name    VARCHAR        NOT NULL
);

CREATE INDEX music_group_name ON music_group (group_name);

CREATE TABLE IF NOT EXISTS songs
(
    id            SERIAL         PRIMARY KEY,
    song_name     VARCHAR        NOT NULL,
    release_date  VARCHAR        NOT NULL,      --TIMESTAMP      DEFAULT(0000-00-00),
    text          VARCHAR        DEFAULT(''),
    link          VARCHAR        DEFAULT('')
);

CREATE TABLE IF NOT EXISTS mgs --music_group_songs
(
    id          SERIAL        PRIMARY KEY,
    group_id    INTEGER       references music_group (id) on delete cascade    NOT NULL,
    song_id     INTEGER       references songs (id) on delete cascade          NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS music_group_songs;
DROP TABLE IF EXISTS music_group;
DROP TABLE IF EXISTS songs;
-- +goose StatementEnd