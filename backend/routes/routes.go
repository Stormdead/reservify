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
	resourceRepo := repositories.NewResourceRepository(config.DB)
	availabilityRepo := repositories.NewAvailabilityRepository(config.DB)
	bookingRepo := repositories.NewBookingRepository(config.DB)

	// Inicializar servicios
	authService := services.NewAuthService(authRepo)
	userService := services.NewUserService(userRepo, authRepo)
	resourceService := services.NewResourceService(resourceRepo)
	availabilityService := services.NewAvailabilityService(availabilityRepo, resourceRepo)
	bookingService := services.NewBookingService(bookingRepo, resourceRepo, userRepo)

	// Inicializar controladores
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	resourceController := controllers.NewResourceController(resourceService)
	availabilityController := controllers.NewAvailabilityController(availabilityService)
	bookingController := controllers.NewBookingController(bookingService)

	// Grupo de API
	api := router.Group("/api")
	{
		// ==================== RUTAS PÚBLICAS ====================

		// Autenticación
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Recursos (públicos - solo activos)
		resources := api.Group("/resources")
		{
			resources.GET("", resourceController.GetAllResources)
			resources.GET("/categories", resourceController.GetCategories)
			resources.GET("/category/:category", resourceController.GetResourcesByCategory)
			resources.GET("/:id", resourceController.GetResourceByID)
			resources.GET("/:id/availability", availabilityController.GetAvailabilityByResource)
		}

		// ==================== RUTAS PROTEGIDAS ====================
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			// Perfil del usuario autenticado
			protected.GET("/auth/me", authController.GetMe)

			// Gestión de perfil
			protected.PUT("/users/me/password", userController.ChangePassword)
			protected.GET("/users/:id", userController.GetUserByID)
			protected.PUT("/users/:id", userController.UpdateUser)

			// ==================== RESERVAS (USUARIOS AUTENTICADOS) ====================
			bookings := protected.Group("/bookings")
			{
				bookings.GET("/my", bookingController.GetMyBookings)             // Mis reservas
				bookings.GET("/upcoming", bookingController.GetUpcomingBookings) // Próximas reservas
				bookings.GET("/:id", bookingController.GetBookingByID)           // Detalle de reserva
				bookings.POST("", bookingController.CreateBooking)               // Crear reserva
				bookings.PUT("/:id", bookingController.UpdateBooking)            // Actualizar reserva
				bookings.DELETE("/:id", bookingController.CancelBooking)         // Cancelar reserva
			}

			// ==================== RUTAS DE ADMIN ====================
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware())
			{
				// Gestión de usuarios
				admin.GET("/users", userController.GetAllUsers)
				admin.GET("/users/stats", userController.GetUserStats)
				admin.DELETE("/users/:id", userController.DeleteUser)

				// Gestión de recursos
				admin.POST("/resources", resourceController.CreateResource)
				admin.PUT("/resources/:id", resourceController.UpdateResource)
				admin.DELETE("/resources/:id", resourceController.DeleteResource)
				admin.GET("/resources/stats", resourceController.GetResourceStats)

				// Gestión de disponibilidad
				admin.POST("/resources/:id/availability", availabilityController.CreateAvailability)
				admin.PUT("/availability/:id", availabilityController.UpdateAvailability)
				admin.DELETE("/availability/:id", availabilityController.DeleteAvailability)

				// Gestión de reservas (admin)
				admin.GET("/bookings", bookingController.GetAllBookings)
				admin.GET("/bookings/stats", bookingController.GetBookingStats)
				admin.PATCH("/bookings/:id/status", bookingController.ChangeBookingStatus)
			}
		}
	}
}
