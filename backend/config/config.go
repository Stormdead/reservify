package config

import (
	"log"
	"os"

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
	JWTExpiration string
	FrontendURL   string
}

var AppConfig *Config

func LoadConfig() {
	//Carga de archivos desde el .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontro el archivo .env usando las variables del sistema")
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
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
		JWTExpiration: getEnv("JWT_EXPIRATION_HOURS", "24"),
		FrontendURL:   getEnv("FRONTEND_URL", "http://localhost:4200"),
	}

	log.Println("Configuración cargada correctamente")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
