package bootstrap

import (
	"errors"
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	DBHost     string `env:"DB_HOST"`
	DBPort     string `env:"DB_PORT"`
	DBUser     string `env:"DB_USER"`
	DBPassword string `env:"DB_PASS"`
	DBName     string `env:"DB_NAME"`
	HTTPPort   string `env:"HTTP_PORT"`
	Mail       MailConfig
}

type MailConfig struct {
	Host     string `env:"MAIL_HOST"`
	Port     int    `env:"MAIL_PORT"`
	Username string `env:"MAIL_USERNAME"`
	Password string `env:"MAIL_PASSWORD"`
	From     string `env:"MAIL_FROM"`
}

func NewConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	var config Config

	_, err = env.UnmarshalFromEnviron(&config)
	if err != nil {
		log.Fatal(err)
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) validate() error {
	var errorList []error

	if c.DBHost == "" {
		errorList = append(errorList, errors.New("empty DB host field "))
	}

	if c.DBPort == "" {
		errorList = append(errorList, errors.New("empty DB port field "))
	}

	if c.DBUser == "" {
		errorList = append(errorList, errors.New("empty DB user field "))
	}

	if c.DBPassword == "" {
		errorList = append(errorList, errors.New("empty DB password field "))
	}

	if c.DBName == "" {
		errorList = append(errorList, errors.New("empty DB name field "))
	}

	if c.HTTPPort == "" {
		errorList = append(errorList, errors.New("empty HTTP port field "))
	}

	if c.Mail.Host == "" {
		errorList = append(errorList, errors.New("empty mail host field "))
	}

	if c.Mail.Port == 0 {
		errorList = append(errorList, errors.New("empty mail port field "))
	}

	if c.Mail.Username == "" {
		errorList = append(errorList, errors.New("empty mail username field "))
	}

	if c.Mail.Password == "" {
		errorList = append(errorList, errors.New("empty mail password field "))
	}

	if c.Mail.From == "" {
		errorList = append(errorList, errors.New("empty mail from field "))
	}

	return errors.Join(errorList...)
}
