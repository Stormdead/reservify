package utils_test

import (
	"Reservify/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetPaginationParams(t *testing.T) {
	// Test: Parámetros por defecto
	t.Run("Parámetros por defecto", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test", nil)

		params := utils.GetPaginationParams(c)

		assert.Equal(t, 1, params.Page, "Page por defecto debería ser 1")
		assert.Equal(t, 10, params.PageSize, "PageSize por defecto debería ser 10")
		assert.Equal(t, "", params.Search, "Search debería estar vacío")
	})

	// Test: Parámetros personalizados
	t.Run("Parámetros personalizados", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test?page=2&page_size=20&search=test", nil)

		params := utils.GetPaginationParams(c)

		assert.Equal(t, 2, params.Page)
		assert.Equal(t, 20, params.PageSize)
		assert.Equal(t, "test", params.Search)
	})

	// Test: Valores inválidos
	t.Run("Valores inválidos", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/test?page=0&page_size=200", nil)

		params := utils.GetPaginationParams(c)

		assert.Equal(t, 1, params.Page, "Page inválido debería ser 1")
		assert.Equal(t, 10, params.PageSize, "PageSize > 100 debería ser 10")
	})
}

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		pageSize int
		expected int
	}{
		{"Primera página", 1, 10, 0},
		{"Segunda página", 2, 10, 10},
		{"Tercera página", 3, 10, 20},
		{"Página 5 con pageSize 20", 5, 20, 80},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := utils.PaginationParams{
				Page:     tt.page,
				PageSize: tt.pageSize,
			}
			offset := params.CalculateOffset()
			assert.Equal(t, tt.expected, offset)
		})
	}
}
