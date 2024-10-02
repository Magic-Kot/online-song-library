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

	q := `
		INSERT INTO songs 
		    (group_song, song, release_date, text, link) 
		VALUES 
		       ($1, $2, $3, $4, $5) 
		RETURNING id
	`

	var id int

	if err := s.client.QueryRowx(q, req.Group, req.Song, res.ReleaseData, res.Text, res.Link).Scan(&id); err != nil {
		logger.Debug().Msgf("failed to create song. %s", err)
		return 0, errCreateSong
	}

	return id, nil
}

// GetAllSong - get all the songs
func (s *SongRepository) GetAllSong(ctx context.Context, req models.RequestGetAll) ([]models.SongsResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'GetAllSong' method")
	logger.Debug().Msgf("postgres: get songs by id: %s, limit: %s", req.Id, req.Limit)

	var songs []models.SongsResponse

	query := fmt.Sprintf(`SELECT * FROM songs WHERE id > %s ORDER BY id LIMIT %s`, req.Id, req.Limit)

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

	query := fmt.Sprintf(`SELECT * FROM songs WHERE %s = $1 AND id > %s ORDER BY id LIMIT %s`, req.Filter, req.Id, req.Limit)

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
