package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

type config struct {
	Database struct {
		Host     string `env:"DB_HOST,required"`
		Port     int    `env:"DB_PORT,required"`
		User     string `env:"DB_USER,required"`
		Password string `env:"DB_PASSWORD,required"`
		DBName   string `env:"DB_NAME,required"`
		SSLMode  string `env:"DB_SSLMODE,required"`
		TimeZone string `env:"DB_TIMEZONE,required"`
	}
	Server struct {
		Address string `env:"SV_ADDRESS,required"`
	}
	Token struct {
		Secret      string `env:"TK_SECRET,required"`
		ExpiredHour int64  `env:"TK_EXPIREDHOUR,required"`
	}
	Route struct {
		GinMode      string `env:"RT_GINMODE,required"`
		MaxBodyBytes int64  `env:"RT_MAXBODYBYTES,required"`
		ApiUri       string `env:"RT_APIURI,required"`
	}
	Logger struct {
		FileName     string `env:"LOG_FILENAME,required"`
		FileSize     int    `env:"LOG_FILESIZE,required"`
		MaxBackup    int    `env:"LOG_MAXBACKUP,required"`
		MaxAge       int    `env:"LOG_MAXAGE,required"`
		FileCompress bool   `env:"LOG_FILECOMPRESS,required"`
	}
}

var C *config

func LoadConfig() error {
	C = &config{}
	ctx := context.Background()

	if err := envconfig.Process(ctx, C); err != nil {
		return err
	}

	return nil
}
