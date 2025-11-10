package routes

import (
	"Reservify/config"
	"Reservify/controllers"
	"Reservify/middleware"
	"Reservify/repositories"
	"Reservify/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Aplicar middleware de CORS
	router.Use(middleware.CORSMiddleware())

	// Inicializar repositorios
	authRepo := repositories.NewAuthRepository(config.DB)
	userRepo := repositories.NewUserRepository(config.DB)

	// Inicializar servicios
	authService := services.NewAuthService(authRepo)
	userService := services.NewUserService(userRepo, authRepo)

	// Inicializar controladores
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	// Grupo de API
	api := router.Group("/api")
	{
		// Rutas públicas de autenticación
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Rutas protegidas (requieren autenticación)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Perfil del usuario autenticado
			protected.GET("/auth/me", authController.GetMe)

			// Cambiar contraseña del usuario autenticado
			protected.PUT("/users/me/password", userController.ChangePassword)

			// Obtener un usuario específico (cualquier usuario autenticado)
			protected.GET("/users/:id", userController.GetUserByID)

			// Actualizar usuario (propio perfil o admin)
			protected.PUT("/users/:id", userController.UpdateUser)

			// Rutas de admin (requieren rol de administrador)
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				// Gestión de usuarios
				admin.GET("/users", userController.GetAllUsers)
				admin.GET("/users/stats", userController.GetUserStats)
				admin.DELETE("/users/:id", userController.DeleteUser)
			}
		}
	}
}
