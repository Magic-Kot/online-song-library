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
	errUserNotFound = errors.New("song not found")
	errCreateSong   = errors.New("failed to create song")
	errGetAllSong   = errors.New("error getting all songs")
	errGetUser      = errors.New("failed to get song")
	errUpdateUser   = errors.New("failed to update song")
	errDeleteUser   = errors.New("failed to delete song")
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
func (s *SongRepository) GetAllSong(ctx context.Context) ([]models.SongsResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'GetAllSong' method")

	var songs []models.SongsResponse

	query := fmt.Sprint(`SELECT id, group_song, song FROM songs ORDER  BY group_song LIMIT 3 OFFSET 6`)

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
	logger.Debug().Msg("accessing Postgres using the 'GetAllSong' method")
	logger.Debug().Msgf("postgres: get songs by filter: %s, value: %s", req.Filter, req.Value)

	var songs []models.SongsResponse

	query := fmt.Sprintf(`SELECT id, group_song, song FROM songs WHERE %s = $1`, req.Filter)

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
		return "", errUserNotFound
	} else if err != nil {
		return "", errGetUser
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
		return errUpdateUser
	}

	if str, _ := commandTag.RowsAffected(); str != 1 {
		logger.Debug().Msgf("song not found: %s", err)
		return errUserNotFound
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
		return errDeleteUser
	}

	if str, _ := commandTag.RowsAffected(); str != 1 {
		return errUserNotFound
	}

	return nil
}
