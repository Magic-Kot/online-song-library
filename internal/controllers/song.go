package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Magic-Kot/effective-mobile/internal/models"
	"github.com/Magic-Kot/effective-mobile/internal/services/song"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ApiController struct {
	songService song.SongService
	logger      *zerolog.Logger
	validator   *validator.Validate
}

func NewApiController(songService *song.SongService, logger *zerolog.Logger, validator *validator.Validate) *ApiController {
	return &ApiController{
		songService: *songService,
		logger:      logger,
		validator:   validator,
	}
}

// @Summary Add Song
// @Tags songs
// @Description add a new song
// @ID add-song
// @Accept  json
// @Produce  json
// @Param input body models.CreateSong true "You need to specify the name of the band and the song in the request body"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /song/create [post]
func (ac *ApiController) AddSong(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	ac.logger.Debug().Msg("starting the handler 'AddSong'")

	req := new(models.CreateSong)
	if err := c.Bind(req); err != nil {
		ac.logger.Debug().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	err := ac.validator.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructField() == "group" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter the name of the group"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum length of the group name is 2 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum length of the group name is 20 characters"))
				}
			}

			if err.StructField() == "song" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter the name of the song"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum length of the song name is 2 characters"))
				}
			}
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := ac.songService.AddSong(ctx, *req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successfully created song, id: %d", id))
}

// @Summary Get All Song
// @Tags songs
// @Description get all saved songs
// @ID get-all-song
// @Accept  json
// @Produce  json
// @Param filter query string false "Enter the column name"
// @Param value query string false "Enter the required column value"
// @Success 200 {object} []models.SongsResponse
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /song/all [get]
func (ac *ApiController) GetAllSong(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	ac.logger.Debug().Msg("starting the handler 'GetAllSong'")

	var req models.RequestGetAll
	req.Filter = c.QueryParam("filter")
	req.Value = c.QueryParam("value")

	result, err := ac.songService.GetAllSong(ctx, req)
	if result == nil {
		return c.JSON(http.StatusNotFound, fmt.Sprintf("no songs found"))
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// @Summary Get Lyrics Song
// @Tags songs
// @Description get the lyrics by id
// @ID get-lyrics-song
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /song/get/{id} [get]
func (ac *ApiController) GetLyricsSong(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	ac.logger.Debug().Msg("starting the handler 'GetLyricsSong'")

	id := c.Param("id")

	result, err := ac.songService.GetLyricsSong(ctx, id)
	if err != nil {
		ac.logger.Debug().Msgf("error receiving song data: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// @Summary Update Song
// @Tags songs
// @Description update information about a saved song
// @ID update-song
// @Accept  json
// @Produce  json
// @Param id path int true "Enter the song ID"
// @Param input body models.SongsResponse true "You need to specify the name of the band and the song in the request body"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Router /song/update/{id} [patch]
func (ac *ApiController) UpdateSong(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	ac.logger.Debug().Msgf("starting the handler 'UpdateSong'")

	var req models.SongsResponse
	if err := c.Bind(&req); err != nil {
		ac.logger.Warn().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	id := c.Param("id")

	userIdInt, err := strconv.Atoi(id)
	if err != nil {
		ac.logger.Debug().Msgf("updateUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	req.Id = userIdInt

	err = ac.validator.Struct(&req)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructField() == "Id" && err.Value() != "" {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("incorrect id"))
			}

			if err.StructField() == "group" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter the name of the group"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum length of the group name is 2 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum length of the group name is 20 characters"))
				}
			}

			if err.StructField() == "song" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter the name of the song"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum length of the song name is 2 characters"))
				}
			}

			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	err = ac.songService.UpdateSong(ctx, req)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprint("successfully updated"))
}

// @Summary Delete Song
// @Tags songs
// @Description delete a song from the library
// @ID delete-song
// @Accept  json
// @Produce  json
// @Param id path int true "Enter the user ID"
// @Success 200 {string} string
// @Failure 400,404 {string} string
// @Router /song/delete/{id} [delete]
func (ac *ApiController) DeleteSong(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	ac.logger.Debug().Msgf("starting the handler 'DeleteSong'")

	id := c.Param("id")

	userIdInt, err := strconv.Atoi(id)
	if err != nil {
		ac.logger.Debug().Msgf("updateUser: invalid id: %s", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	err = ac.songService.DeleteSong(ctx, userIdInt)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successfully deleted song: %d", userIdInt))
}
