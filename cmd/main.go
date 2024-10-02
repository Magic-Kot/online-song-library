package main

import (
	"context"
	"embed"
	"time"

	"github.com/Magic-Kot/effective-mobile/internal/config"
	"github.com/Magic-Kot/effective-mobile/internal/controllers"
	"github.com/Magic-Kot/effective-mobile/internal/delivery/httpecho"
	"github.com/Magic-Kot/effective-mobile/internal/repository/postgres"
	"github.com/Magic-Kot/effective-mobile/internal/services/song"
	"github.com/Magic-Kot/effective-mobile/pkg/client/postg"
	"github.com/Magic-Kot/effective-mobile/pkg/httpserver"
	"github.com/Magic-Kot/effective-mobile/pkg/logging"
	"github.com/Magic-Kot/effective-mobile/pkg/musicinfo"
	"github.com/Magic-Kot/effective-mobile/pkg/ossignal"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/speakeasy-api/goose/v3"
	"golang.org/x/sync/errgroup"
)

// @title Online Song Library
// @version 1.0
// @description This project was developed as part of a test assignment from Effective Mobile

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	// read config
	var cfg config.Config

	err := cleanenv.ReadConfig("internal/config/.env", &cfg) // Local
	if err != nil {
		log.Fatal().Err(err).Msg("error initializing config")
	}

	// create logger
	logCfg := logging.LoggerDeps{
		LogLevel: cfg.LoggerDeps.LogLevel,
	}

	logger, err := logging.NewLogger(&logCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init logger")
	}

	logger.Info().Msg("init logger")

	ctx := context.Background()
	ctx = logger.WithContext(ctx)

	logger.Info().Msgf("config: %+v", cfg)

	// create server
	serv := httpserver.ConfigDeps{
		Host:    cfg.ServerDeps.Host,
		Port:    cfg.ServerDeps.Port,
		Timeout: cfg.ServerDeps.Timeout,
	}

	server := httpserver.NewServer(&serv)

	// create client Postgres
	repo := postg.ConfigDeps{
		MaxAttempts: cfg.PostgresDeps.MaxAttempts,
		Delay:       cfg.PostgresDeps.Delay,
		Username:    cfg.PostgresDeps.Username,
		Password:    cfg.PostgresDeps.Password,
		Host:        cfg.PostgresDeps.Host,
		Port:        cfg.PostgresDeps.Port,
		Database:    cfg.PostgresDeps.Database,
		SSLMode:     cfg.PostgresDeps.SSLMode,
	}

	pool, err := postg.NewClient(ctx, &repo)
	if err != nil {
		logger.Fatal().Err(err).Msgf("NewClient: %s", err)
	}

	// migrations
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(pool, "migrations"); err != nil {
		panic(err)
	}

	// create validator
	validate := validator.New()

	// Song
	songRepository := postgres.NewSongRepository(pool)
	musicInfo := musicinfo.NewMusicInfo(cfg.MusicInfo.Url)
	songService := song.NewSongService(songRepository, musicInfo)
	songController := controllers.NewApiController(songService, logger, validate)
	httpecho.SetSongRoutes(server.Server(), songController)

	runner, ctx := errgroup.WithContext(ctx)

	// start server
	logger.Info().Msg("starting server")
	runner.Go(func() error {
		if err := server.Start(); err != nil {
			logger.Fatal().Msgf("%v", err)
		}

		return nil
	})

	runner.Go(func() error {
		if err := ossignal.DefaultSignalWaiter(ctx); err != nil {
			return errors.Wrap(err, "waiting os signal")
		}

		return nil
	})

	runner.Go(func() error {
		<-ctx.Done()

		ctxSignal, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		if err := server.Shutdown(ctxSignal); err != nil {
			logger.Error().Err(err).Msg("shutdown http server")
		}

		return nil
	})

	if err := runner.Wait(); err != nil {
		switch {
		case ossignal.IsExitSignal(err):
			logger.Info().Msg("exited by exit signal")
		default:
			logger.Fatal().Msgf("exiting with error: %v", err)
		}
	}
}
