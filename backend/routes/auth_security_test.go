package routes

import (
	"Reservify/config"
	"Reservify/middleware"
	"Reservify/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestProtectedRoutesRejectUnauthorizedRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)
	previous := config.AppConfig
	t.Cleanup(func() { config.AppConfig = previous })
	config.AppConfig = &config.Config{JWTSecret: "route-test-fixture-key-0123456789abcdef", JWTExpiration: time.Hour, FrontendURL: "http://localhost:4200"}
	router := gin.New()
	SetupRoutes(router)
	router.GET("/admin-check", middleware.AuthMiddleware(), middleware.AdminMiddleware(), func(c *gin.Context) { c.Status(http.StatusNoContent) })
	user, err := utils.GenerateToken(1, "user@example.invalid", "user")
	if err != nil {
		t.Fatal(err)
	}
	admin, err := utils.GenerateToken(2, "admin@example.invalid", "admin")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		method, path, token string
		status              int
	}{
		{"GET", "/api/users/2", "", 401},
		{"GET", "/api/users/2", "malformed-token", 401},
		{"GET", "/api/users/2", user, 403},
		{"PUT", "/api/users/2", user, 403},
		{"GET", "/api/admin/users", "", 401},
		{"GET", "/api/admin/users", user, 403},
		{"GET", "/admin-check", admin, 204},
	} {
		t.Run(tc.method+" "+tc.path+" "+http.StatusText(tc.status), func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)
			if res.Code != tc.status {
				t.Fatalf("status = %d, want %d", res.Code, tc.status)
			}
			if strings.Contains(res.Body.String(), `"error":`) {
				t.Fatal("response exposes internal error")
			}
		})
	}
}
