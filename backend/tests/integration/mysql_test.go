//go:build integration

package integration

import (
	"Reservify/config"
	"Reservify/dto"
	"Reservify/models"
	"Reservify/routes"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	driver "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestMySQLAPI(t *testing.T) {
	db := openTestDatabase(t)
	t.Run("migrations", func(t *testing.T) {
		if err := models.AutoMigrate(db); err != nil {
			t.Fatal(err)
		}
		for _, model := range []any{&models.User{}, &models.Resource{}, &models.AvailabilitySlot{}, &models.Booking{}, &models.Notification{}} {
			if !db.Migrator().HasTable(model) {
				t.Errorf("missing table for %T", model)
			}
		}
	})
	if t.Failed() {
		return
	}

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	t.Cleanup(func() {
		if err := tx.Rollback().Error; err != nil {
			t.Errorf("rollback test data: %v", err)
		}
	})

	previousDB, previousConfig := config.DB, config.AppConfig
	t.Cleanup(func() { config.DB, config.AppConfig = previousDB, previousConfig })
	config.DB = tx
	config.AppConfig = &config.Config{
		JWTSecret:     "integration-test-signing-fixture-only",
		JWTExpiration: time.Hour,
		FrontendURL:   "http://localhost:4200",
	}
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router)

	const password = "integration-fixture"
	var owner, other dto.AuthResponse
	if !t.Run("registration persists a regular user", func(t *testing.T) {
		for index, account := range []*dto.AuthResponse{&owner, &other} {
			email := fmt.Sprintf("fixture-%d@example.invalid", index)
			request(t, router, http.MethodPost, "/api/auth/register", "", dto.RegisterRequest{
				Email: email, Password: password, FullName: "Integration Fixture",
			}, http.StatusCreated, account)

			var stored models.User
			if err := tx.First(&stored, account.User.ID).Error; err != nil {
				t.Fatal(err)
			}
			if account.User.ID == 0 || account.Token == "" || stored.ID != account.User.ID || stored.Email != email || stored.Role != models.RoleUser {
				t.Fatal("registration response or persisted user has unexpected values")
			}
			if err := bcrypt.CompareHashAndPassword([]byte(stored.PasswordHash), []byte(password)); err != nil {
				t.Fatal("persisted password hash does not match the fixture")
			}
		}
	}) {
		return
	}

	t.Run("own profile", func(t *testing.T) {
		for _, path := range []string{"/api/auth/me", fmt.Sprintf("/api/users/%d", owner.User.ID)} {
			var profile dto.UserResponse
			request(t, router, http.MethodGet, path, owner.Token, nil, http.StatusOK, &profile)
			if profile.ID != owner.User.ID || profile.Email != owner.User.Email {
				t.Fatalf("%s returned an unexpected profile", path)
			}
		}
	})
	t.Run("another user's profile is forbidden", func(t *testing.T) {
		request(t, router, http.MethodGet, fmt.Sprintf("/api/users/%d", other.User.ID),
			owner.Token, nil, http.StatusForbidden, nil)
	})
	t.Run("administration is forbidden for regular users", func(t *testing.T) {
		request(t, router, http.MethodGet, "/api/admin/users", owner.Token, nil, http.StatusForbidden, nil)
	})
	t.Run("valid booking is persisted", func(t *testing.T) {
		resource := models.Resource{Name: "Integration Room", Capacity: 4, PricePerHour: 10000, IsActive: true}
		if err := tx.Create(&resource).Error; err != nil {
			t.Fatal(err)
		}
		start := time.Now().UTC().Truncate(time.Second).Add(24 * time.Hour)
		end := start.Add(time.Hour)
		var response dto.BookingResponse
		request(t, router, http.MethodPost, "/api/bookings", owner.Token, dto.CreateBookingRequest{
			ResourceID: resource.ID, StartDatetime: start, EndDatetime: end,
		}, http.StatusCreated, &response)
		if response.ID == 0 {
			t.Fatal("booking response has no ID")
		}
		var stored models.Booking
		if err := tx.First(&stored, response.ID).Error; err != nil {
			t.Fatal(err)
		}
		if stored.UserID != owner.User.ID || stored.ResourceID != resource.ID ||
			stored.Status != models.StatusPending || stored.TotalPrice != resource.PricePerHour ||
			!stored.StartDatetime.Equal(start) || !stored.EndDatetime.Equal(end) {
			t.Fatalf("unexpected booking persisted for ID %d", stored.ID)
		}
	})
}

// The schema and account are reserved for these tests. Only the local service is allowed.
func openTestDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	port, password := os.Getenv("MYSQL_TEST_PORT"), os.Getenv("MYSQL_TEST_PASSWORD")
	if port == "" || password == "" {
		t.Fatal("set MYSQL_TEST_PORT and MYSQL_TEST_PASSWORD for the isolated local MySQL service")
	}
	dsn := driver.NewConfig()
	dsn.Net = "tcp"
	dsn.Addr = "127.0.0.1:" + port
	dsn.User = "reservify_test"
	dsn.Passwd = password
	dsn.DBName = "reservify_integration_test"
	dsn.ParseTime = true
	dsn.Timeout = 5 * time.Second
	db, err := gorm.Open(mysql.Open(dsn.FormatDSN()), &gorm.Config{})
	if err != nil {
		t.Fatalf("open isolated MySQL: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		t.Fatalf("ping isolated MySQL: %v", err)
	}
	var version string
	if err := sqlDB.QueryRowContext(ctx, "SELECT VERSION()").Scan(&version); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(version, "8.") || strings.Contains(strings.ToLower(version), "mariadb") {
		t.Fatalf("expected MySQL 8, got %q", version)
	}
	return db
}

func request(t *testing.T, router http.Handler, method, path, token string, body any, status int, data any) {
	t.Helper()
	var payload bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&payload).Encode(body); err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, &payload)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	if response.Code != status {
		t.Fatalf("%s %s returned %d, want %d: %s", method, path, response.Code, status, response.Body.String())
	}
	if data == nil {
		return
	}
	var envelope struct {
		Success bool            `json:"success"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil || !envelope.Success {
		t.Fatalf("%s %s returned an invalid success envelope: %v", method, path, err)
	}
	if err := json.Unmarshal(envelope.Data, data); err != nil {
		t.Fatal(err)
	}
}
