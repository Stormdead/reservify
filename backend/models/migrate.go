package models

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	log.Println("Ejecutando migraciones automaticas")

	//Lista de modelos a migrar
	err := db.AutoMigrate(
		&User{},
	)
	if err != nil {
		return fmt.Errorf("Error en auto-migrate: %v", err)
	}

	log.Println("Migraciones completadas")
	return nil
}
