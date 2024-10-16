package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Magic-Kot/effective-mobile/internal/models"
	"github.com/Magic-Kot/effective-mobile/pkg/client/postg"
	"github.com/Magic-Kot/effective-mobile/pkg/musicinfo"

	"github.com/rs/zerolog"
)

var (
	errTransaction  = errors.New("transaction error")
	errSongNotFound = errors.New("song not found")
	errCreateSong   = errors.New("failed to create song")
	errGetAllSong   = errors.New("error getting all songs")
	errGetSong      = errors.New("failed to get song")
	errUpdateSong   = errors.New("failed to update song")
	errDeleteSong   = errors.New("failed to delete song")
)

type SongRepository struct {
	client postg.Client
}

func NewSongRepository(client postg.Client) *SongRepository {
	return &SongRepository{
		client: client,
	}
}

// AddSong - add a new song
func (s *SongRepository) AddSong(ctx context.Context, req models.CreateSong, res musicinfo.SongDetail) (int, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'AddSong' method")

	tx, err := s.client.Begin()
	if err != nil {
		logger.Debug().Msgf("transaction creation error. err: %s", err)
		return 0, errTransaction
	}

	checkGroupQuery := fmt.Sprint(`SELECT id FROM music_group WHERE group_name = $1`)

	var idGroup, idSong int

	rowGroup := s.client.QueryRowx(checkGroupQuery, req.Group).Scan(&idGroup)

	if errors.Is(rowGroup, sql.ErrNoRows) {
		logger.Debug().Msgf("music group not found. err: %s", err)

		addGroupQuery := fmt.Sprint("INSERT INTO music_group (group_name) VALUES ($1) RETURNING id")

		row := tx.QueryRow(addGroupQuery, req.Group)
		if err = row.Scan(&idGroup); err != nil {
			logger.Debug().Msgf("error writing to the 'music_group' table. err: %s", err)

			tx.Rollback()
			return 0, errTransaction
		}

		return 0, errSongNotFound
	} else if rowGroup != nil {
		logger.Debug().Msgf("error getting a music group. err: %s", err)

		tx.Rollback()
		return 0, errTransaction
	}

	addSongQuery := fmt.Sprint("INSERT INTO songs (song_name, release_date, text, link) VALUES ($1, $2, $3, $4) RETURNING id")

	rowSong := tx.QueryRow(addSongQuery, req.Song, res.ReleaseData, res.Text, res.Link)
	if err = rowSong.Scan(&idSong); err != nil {
		logger.Debug().Msgf("error writing to the 'songs' table. err: %s", err)

		tx.Rollback()
		return 0, errTransaction
	}

	addMgsQuery := fmt.Sprint("INSERT INTO mgs (group_id, song_id) VALUES ($1, $2)")
	_, err = tx.Exec(addMgsQuery, idGroup, idSong)
	if err != nil {
		logger.Debug().Msgf("error writing to the 'mgs' table. err: %s", err)

		tx.Rollback()
		return 0, errTransaction
	}

	return idSong, tx.Commit()
}

// GetAllSong - get all the songs
func (s *SongRepository) GetAllSong(ctx context.Context, req models.RequestGetAll) ([]models.SongsResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'GetAllSong' method")
	logger.Debug().Msgf("postgres: get songs by id: %s, limit: %s", req.Id, req.Limit)

	var songs []models.SongsResponse

	query := fmt.Sprintf(`SELECT id, song_name, release_date, text, link FROM songs WHERE id > %s ORDER BY id LIMIT %s`, req.Id, req.Limit)

	err := s.client.Select(&songs, query)
	if err != nil {
		logger.Debug().Msgf("error getting all songs. err: %s", err)
		return nil, errGetAllSong
	}

	return songs, nil
}

// GetAllSongFilter - get all the songs using the filter
func (s *SongRepository) GetAllSongFilter(ctx context.Context, req models.RequestGetAll) ([]models.SongsResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'GetAllSongFilter' method")
	logger.Debug().Msgf("postgres: get songs by id: %s, limit: %s, filter: %s, value: %s", req.Id, req.Limit, req.Filter, req.Value)

	var songs []models.SongsResponse

	query := fmt.Sprintf(`SELECT id, song_name, release_date, text, link FROM songs WHERE %s = $1 AND id > %s ORDER BY id LIMIT %s`, req.Filter, req.Id, req.Limit)

	err := s.client.Select(&songs, query, req.Value)
	if err != nil {
		logger.Debug().Msgf("error getting all songs. err: %s", err)
		return nil, errGetAllSong
	}

	return songs, nil
}

// GetLyricsSong - get the lyrics by id
func (s *SongRepository) GetLyricsSong(ctx context.Context, id string) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'GetLyricsSong' method")
	logger.Debug().Msgf("postgres: get song by id: %s", id)

	query := fmt.Sprint(`SELECT text FROM songs WHERE id = $1`)

	var lyrics string

	err := s.client.QueryRowx(query, id).Scan(&lyrics)

	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)

		return "", errSongNotFound
	} else if err != nil {
		fmt.Println(err)

		return "", errGetSong
	}

	return lyrics, nil
}

// UpdateSong - update information about a saved song
func (s *SongRepository) UpdateSong(ctx context.Context, value string, arg []interface{}) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'UpdateSong' method")
	logger.Debug().Msgf("postgres: update table by value: %s, arg: %v", value, arg)

	q := fmt.Sprintf(`UPDATE songs SET %s WHERE id = $1`, value)

	commandTag, err := s.client.Exec(q, arg...)

	if err != nil {
		logger.Debug().Msgf("failed table updates: %s", err)
		return errUpdateSong
	}

	if str, _ := commandTag.RowsAffected(); str != 1 {
		logger.Debug().Msgf("song not found: %s", err)
		return errSongNotFound
	}

	return nil
}

// DeleteSong - delete a song from the library
func (s *SongRepository) DeleteSong(ctx context.Context, id int) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'DeleteSong' method")

	q := `
		DELETE FROM songs
		WHERE id = $1
	`

	commandTag, err := s.client.Exec(q, id)

	if err != nil {
		return errDeleteSong
	}

	if str, _ := commandTag.RowsAffected(); str != 1 {
		return errSongNotFound
	}

	return nil
}
