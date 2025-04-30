package config

import (
	"log"
	"log/slog"
	"net"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/vovanwin/meetingsBot/pkg/validator"
)

var configPath = "config/config.yml"

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		slog.Debug("Нет конфиг :", "err", err)
	}
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
	}

	if err := validator.NewCustomValidator().Validate(cfg); err != nil {
		log.Fatalf(err.Error())
	}
	return cfg, nil
}

type (
	Config struct {
		Server `yaml:"server"`
		Log    `yaml:"log"`
	}

	Server struct {
		Host string `env-required:"true" yaml:"host" env:"HOST" validate:"required"`
		Port string `env-required:"true" yaml:"port" env:"PORT" validate:"required"`
		Env  string `env-required:"true" yaml:"env" env:"ENV" validate:"required,oneof=local dev prod"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"  validate:"required,oneof=debug info warn error"`
	}
)

func (c Config) Address() string {
	return net.JoinHostPort(c.Server.Host, c.Server.Port)
}

func (c Config) IsProduction() bool {
	return c.Server.Env == "prod"
}

func (c Config) IsLocal() bool {
	return c.Server.Env == "local"
}

func (c Config) IsTest() bool {
	return c.Server.Env == "test"
}

func (c Config) IsDebug() bool {
	return c.Log.Level == "debug"
}
