package config

import "time"

type Config struct {
	ServerDeps
	PostgresDeps
	LoggerDeps
	MusicInfo
}

type ServerDeps struct {
	Host    string        `env:"HOST"     env-default:"localhost"`
	Port    string        `env:"PORT"     env-default:":8080"`
	Timeout time.Duration `env:"TIMEOUT"  env-default:"5s"`
}

type PostgresDeps struct {
	MaxAttempts int           `env:"MAX_ATTEMPTS"       env-default:"3"`
	Delay       time.Duration `env:"DELAY"              env-default:"10s"`
	Username    string        `env:"USERNAME_POSTGRES"  env-default:"postgres"`
	Password    string        `env:"PASSWORD_POSTGRES"  env-default:"postgres"`
	Host        string        `env:"HOST_POSTGRES"      env-default:"127.0.0.1"`
	Port        string        `env:"PORT_POSTGRES"      env-default:"5432"`
	Database    string        `env:"DATABASE"           env-default:"postgres"`
	SSLMode     string        `env:"MODELESS"           env-default:"disable"`
}

type LoggerDeps struct {
	LogLevel string `env:"LOG_LEVEL"  env-default:"info"`
}

type MusicInfo struct {
	Url string `env:"MUSIC_URL"`
}
