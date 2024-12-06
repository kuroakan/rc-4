package bootstrap

import (
	"errors"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASS"`
	DBName     string `env:"DB_NAME"`
	HTTPPort   string `env:"HTTP_PORT"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	config := Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASS"),
		DBName:     os.Getenv("DB_NAME"),
		HTTPPort:   os.Getenv("HTTP_PORT"),
	}

	return &config, nil
}

func (c *Config) Validate() error {
	var errorList []error

	if c.DBHost == "" {
		err := errors.New("invalid DB host field ")
		errorList = append(errorList, err)
	}

	if c.DBPort == "" {
		err := errors.New("invalid DB port field ")
		errorList = append(errorList, err)
	}

	if c.DBUser == "" {
		err := errors.New("invalid DB user field ")
		errorList = append(errorList, err)
	}

	if c.DBPassword == "" {
		err := errors.New("invalid DB password field ")
		errorList = append(errorList, err)
	}

	if c.DBName == "" {
		err := errors.New("invalid DB name field ")
		errorList = append(errorList, err)
	}

	if c.HTTPPort == "" {
		err := errors.New("invalid HTTP port field ")
		errorList = append(errorList, err)
	}

	if len(errorList) != 0 {
		return errors.Join(errorList...)
	}

	return nil
}
