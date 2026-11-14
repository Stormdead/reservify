package config

import (
	"strings"
	"testing"
	"time"
)

func TestLoadConfigJWT(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("ENV", "production")
	for _, key := range []string{"", "secret", strings.Repeat(" ", 32), "your-super-secret-jwt-key-change-this-in-production"} {
		t.Setenv("JWT_SECRET", key)
		if LoadConfig() == nil {
			t.Fatal("accepted insecure JWT configuration")
		}
	}
	t.Setenv("JWT_SECRET", "test-only-config-fixture-0123456789")
	for _, hours := range []string{"0", "-1", "invalid", "2147483647"} {
		t.Setenv("JWT_EXPIRATION_HOURS", hours)
		if LoadConfig() == nil {
			t.Fatal("accepted invalid expiration")
		}
	}
	t.Setenv("JWT_EXPIRATION_HOURS", "2")
	if err := LoadConfig(); err != nil {
		t.Fatal(err)
	}
	if AppConfig.JWTExpiration != 2*time.Hour {
		t.Fatal("incorrect JWT expiration")
	}
}
