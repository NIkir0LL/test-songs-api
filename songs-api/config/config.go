package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	APIURL     string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("Не удалось загрузить .env файл, используются переменные окружения")
	}
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "yourpassword"),
		DBName:     getEnv("DB_NAME", "songs_db"),
		APIURL:     getEnv("API_URL", "http://localhost:8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
