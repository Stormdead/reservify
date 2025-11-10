package middleware

import (
	"net/http"
	"strings"

	"Reservify/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware verifica el token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token no proporcionado", nil)
			c.Abort()
			return
		}

		// Verificar formato "Bearer {token}"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Formato de token inválido", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validar el token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token inválido o expirado", err)
			c.Abort()
			return
		}

		// Guardar información del usuario en el contexto
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}
