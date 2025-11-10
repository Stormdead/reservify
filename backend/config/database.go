package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Establece conexion a la base de datos
func ConnectDatabase() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		AppConfig.DBUser,
		AppConfig.DBPassword,
		AppConfig.DBHost,
		AppConfig.DBPort,
		AppConfig.DBName,
	)
	// Configurar GORM
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// Conectar a MySQL
	DB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		log.Fatal(" Error al conectar a la base de datos:", err)
	}

	// Configurar el pool de conexiones
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal(" Error al configurar el pool de conexiones:", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println(" Conexión a MySQL establecida correctamente")
}

// GetDB devuelve la instancia de la base de datos
func GetDB() *gorm.DB {
	return DB
}
