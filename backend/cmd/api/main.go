package main

import (
	"Reservify/config"
	"Reservify/models"
	"Reservify/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Iniciando Reservify API")

	config.LoadConfig()
	config.ConnectDatabase()

	//Ejecutar migraciones automaticas
	if err := models.AutoMigrate(config.GetDB()); err != nil {
		log.Fatal(" Error en migraciones:", err)
	}
	// Configurar modo de Gin
	if config.AppConfig.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Crear router
	router := gin.Default()

	// Ruta de health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "ok",
			"database": "connected",
		})
	})

	// Configurar rutas
	routes.SetupRoutes(router)

	// Iniciar servidor
	port := ":" + config.AppConfig.Port
	log.Printf(" Servidor corriendo en http://localhost%s", port)
	log.Printf(" Ambiente: %s", config.AppConfig.Env)
	log.Println(" Presiona Ctrl+C para detener")

	if err := router.Run(port); err != nil {
		log.Fatal(" Error al iniciar servidor:", err)
	}
}
