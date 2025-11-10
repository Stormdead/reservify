package middleware

import (
	"net/http"

	"Reservify/utils"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware verifica que el usuario sea admin
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "No autenticado", nil)
			c.Abort()
			return
		}

		if role != "admin" {
			utils.ErrorResponse(c, http.StatusForbidden, "Acceso denegado: se requiere rol de administrador", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
