package song

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Magic-Kot/effective-mobile/internal/models"
	"github.com/Magic-Kot/effective-mobile/pkg/musicinfo"

	"github.com/rs/zerolog"
)

type SongRepository interface {
	AddSong(ctx context.Context, req models.CreateSong, res musicinfo.SongDetail) (int, error)
	GetAllSong(ctx context.Context, req models.RequestGetAll) ([]models.SongsResponse, error)
	GetAllSongFilter(ctx context.Context, req models.RequestGetAll) ([]models.SongsResponse, error)
	GetLyricsSong(ctx context.Context, id string) (string, error)
	UpdateSong(ctx context.Context, value string, arg []interface{}) error
	DeleteSong(ctx context.Context, id int) error
}

type SongService struct {
	SongRepository SongRepository
	MusicInfo      *musicinfo.MusicInfo
}

func NewSongService(songRepository SongRepository, musicInfo *musicinfo.MusicInfo) *SongService {
	return &SongService{
		SongRepository: songRepository,
		MusicInfo:      musicInfo,
	}
}

// AddSong - add a new song
func (s *SongService) AddSong(ctx context.Context, req models.CreateSong) (int, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'AddSong' service")

	res, _ := s.MusicInfo.Info(req.Group, req.Song)

	id, err := s.SongRepository.AddSong(ctx, req, res)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetAllSong - get all the songs
func (s *SongService) GetAllSong(ctx context.Context, req models.RequestGetAll) ([]models.SongsResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'GetAllSong' service")

	if req.Filter != "" && req.Value != "" {
		res, err := s.SongRepository.GetAllSongFilter(ctx, req)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	res, err := s.SongRepository.GetAllSong(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetLyricsSong - get the lyrics by id
func (s *SongService) GetLyricsSong(ctx context.Context, songId string, verse string) (string, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'GetLyricsSong' service")

	text, err := s.SongRepository.GetLyricsSong(ctx, songId)
	if err != nil {
		return "", err
	}

	textSplit := strings.Split(text, "\n\n")

	verseInt, err := strconv.Atoi(verse)
	if err != nil {
		return "", err
	}

	return textSplit[verseInt], nil
}

// UpdateSong - update information about a saved song
func (s *SongService) UpdateSong(ctx context.Context, song models.UpdateRequest) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'UpdateSong' service")

	value := make([]string, 0)
	arg := make([]interface{}, 0)
	argId := 2

	arg = append(arg, song.Id)

	values := reflect.ValueOf(song)
	types := values.Type()

	for i := 1; i < values.NumField(); i++ {
		if types.Field(i).Name == "Song" {
			value = append(value, fmt.Sprintf("song_name=$%d", argId))
		} else if types.Field(i).Name == "ReleaseDate" {
			value = append(value, fmt.Sprintf("release_date=$%d", argId))
		} else {
			value = append(value, fmt.Sprintf("%s=$%d", types.Field(i).Name, argId))
		}

		arg = append(arg, values.Field(i).String())

		argId++
	}

	valueQuery := strings.Join(value, ", ")

	err := s.SongRepository.UpdateSong(ctx, valueQuery, arg)
	if err != nil {
		return err
	}

	return nil
}

// DeleteSong - delete a song from the library
func (s *SongService) DeleteSong(ctx context.Context, id int) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'DeleteSong' service")

	err := s.SongRepository.DeleteSong(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
