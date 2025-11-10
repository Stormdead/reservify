package utils

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Representa los parametros de la pagina
type PaginationParams struct {
	Page     int
	PageSize int
	Search   string
}

// Representa la metadata de pagina
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

// Representa la respuesta paginada
type PaginatedResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    interface{}    `json:"data"`
	Meta    PaginationMeta `json:"meta"`
}

// Obtiene los parametros de paginacion del request
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.DefaultQuery("search", "")

	// Validar valores
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Search:   search,
	}
}

// Calcula el offset para la consulta SQL
func (p *PaginationParams) CalculateOffset() int {
	return (p.Page - 1) * p.PageSize
}

// Envía una respuesta paginada exitosa
func PaginatedSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, total int64, params PaginationParams) {
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))

	c.JSON(statusCode, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: PaginationMeta{
			Page:       params.Page,
			PageSize:   params.PageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	})
}
