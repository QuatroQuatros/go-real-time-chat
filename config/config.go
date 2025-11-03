package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppEnv     string `envconfig:"APP_ENV" envDefault:"development"`
	Host       string `envconfig:"DB_SERVER" envDefault:"localhost"`
	Port       string `envconfig:"DB_PORT" envDefault:"5432"`
	Username   string `envconfig:"DB_USERNAME"`
	Password   string `envconfig:"DB_PASSWORD"`
	Database   string `envconfig:"DB_NAME"`
	ServerPort string `envconfig:"SERVER_PORT" envDefault:"8080"`
}

var Env Config

func LoadEnv() error {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Println("⚠️  Não foi possível carregar o .env, continuando com variáveis de ambiente")
		} else {
			log.Println("✅ Variáveis carregadas do .env")
		}
	} else {
		log.Println("⚠️  Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	if err := envconfig.Process("", &Env); err != nil {
		return err
	}

	return nil
}
