package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env"
)

type Config struct {
	DBName     string `env:"DB_NAME,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBUser     string `env:"DB_USER,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBPort     string `env:"DB_PORT,required"`

	CacheMinutes int    `env:"CACHE_MINUTES,required"`
	CacheHost    string `env:"CACHE_HOST,required"`
	CachePort    string `env:"CACHE_PORT,required"`

	Port string `env:"PORT,required"`
}

func (c Config) String() string {
	var sb strings.Builder

	sb.WriteString("Product Service Settings:\n")

	sb.WriteString("DB\n")
	sb.WriteString(fmt.Sprintf("DB_NAME: %s\n", c.DBName))
	sb.WriteString(fmt.Sprintf("DB_PASSWORD: %s\n", c.DBPassword))
	sb.WriteString(fmt.Sprintf("DB_NAME: %s\n", c.DBName))
	sb.WriteString(fmt.Sprintf("DB_HOST %s\n", c.DBHost))
	sb.WriteString(fmt.Sprintf("DB_PORT: %s\n", c.DBPort))

	sb.WriteString("CACHE\n")
	sb.WriteString(fmt.Sprintf("CACHE_MINUTES: %d\n", c.CacheMinutes))
	sb.WriteString(fmt.Sprintf("CACHE_HOST %s\n", c.CacheHost))
	sb.WriteString(fmt.Sprintf("CACHE_PORT: %s\n", c.CachePort))

	sb.WriteString("APP\n")
	sb.WriteString(fmt.Sprintf("PORT: %s\n", c.Port))

	return sb.String()
}

func MustLoadConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(fmt.Errorf("parsing config error: %s", err.Error()))
	}
	return cfg
}
