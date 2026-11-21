package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName       string
	Port          string
	Env           string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	JWTExpiration time.Duration
	FrontendURL   string
}

var AppConfig *Config

func LoadConfig() error {
	//Carga de archivos desde el .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontro el archivo .env usando las variables del sistema")
	}

	secret := os.Getenv("JWT_SECRET")
	if len(strings.TrimSpace(secret)) < 32 || strings.TrimSpace(secret) == "your-super-secret-jwt-key-change-this-in-production" {
		return errors.New("JWT_SECRET requiere una clave aleatoria de al menos 32 bytes")
	}
	hours, err := strconv.ParseInt(getEnv("JWT_EXPIRATION_HOURS", "24"), 10, 32)
	if err != nil || hours <= 0 || hours > int64((1<<63-1)/time.Hour) {
		return errors.New("JWT_EXPIRATION_HOURS debe indicar una cantidad positiva de horas representable")
	}
	AppConfig = &Config{
		AppName:       getEnv("APP_NAME", "Reservify"),
		Port:          getEnv("PORT", "8080"),
		Env:           getEnv("ENV", "development"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "3306"),
		DBUser:        getEnv("DB_USER", "root"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "reservify_db"),
		JWTSecret:     secret,
		JWTExpiration: time.Duration(hours) * time.Hour,
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:4200"),
	}

	log.Println("Configuración cargada correctamente")
	return nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
