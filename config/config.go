package config

import (
	"log"
	"log/slog"
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

	if err := validator.Validator.Struct(cfg); err != nil {
		log.Fatalf(err.Error())
	}
	return cfg, nil
}

type (
	Config struct {
		Server   `yaml:"server"`
		Log      `yaml:"log"`
		Telegram `yaml:"telegram"`
	}

	Server struct {
		Env string `env-required:"true" yaml:"env" env:"ENV" validate:"required,oneof=local dev prod"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"  validate:"required,oneof=debug info warn error"`
	}

	Telegram struct {
		Token      string `env:"TG_TOKEN" env-required:"true" yaml:"token" validate:"required"`
		Webhook    string `env:"TG_WEBHOOK_URL" yaml:"webhook"`
		PublicURL  string `env:"TG_PUBLIC_URL" yaml:"public_url"` // адрес домена, если Webhook. Пример: https://mydomain.com
		Listen     string `env:"TG_LISTEN" yaml:"listen" env-default:":8080"`
		UseWebhook bool   `env:"TG_USE_WEBHOOK" yaml:"use_webhook" env-default:"false"` // true = webhook, false = polling
	}
)

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
