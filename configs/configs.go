package configs

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Environment
		Postgres
		Logger
	}

	Environment struct {
		Status string `yaml:"status" env-default:"prod"`
	}

	Postgres struct {
		User     string `yaml:"user" env-default:"postgres"`
		Password string `yaml:"password" env-default:"postgres"`

		Host string `yaml:"host" env-default:"localhost"`
		Port string `yaml:"port" env-default:"5432"`

		Name string `yaml:"name" env-default:"template1"`
	}

	Logger struct {
		MaxFileSize  int  `yaml:"maxfilesize" env-default:"128"`
		MaxBackups   int  `yaml:"maxbackups" env-default:"8"`
		MaxAge       int  `yaml:"maxage" env-default:"28"`
		IsCompressed bool `yaml:"compress" env-default:"28"`
		IsLocalTime  bool `yaml:"localtime" env-default:"false"`
	}
)

const (
	// filePath defines the path when config.yaml resides
	filePath = `../../tools/config.yaml`
)

var (
	Cfg Config
)

// GenPostgresURI generates URI based on Cfg to be used in conns
func GenPostgresURI() string {
	const (
		PostgreConnFormatURI = "postgres://%s:%s@%s:%s/%s"
	)

	return fmt.Sprintf(PostgreConnFormatURI, Cfg.User, Cfg.Password, Cfg.Host, Cfg.Port, Cfg.Name)

}

// Load loads config from config.yaml and returns error if any
func Load() error {
	if err := cleanenv.ReadConfig(filePath, &Cfg); err != nil {
		return fmt.Errorf("reading config from %q: %w", filePath, err)
	}

	return nil
}
