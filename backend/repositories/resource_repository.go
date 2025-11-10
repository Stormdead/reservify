package repositories

import (
	"Reservify/models"
	"Reservify/utils"
	"errors"

	"gorm.io/gorm"
)

type ResourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

// FindAll obtiene todos los recursos con paginación
func (r *ResourceRepository) FindAll(params utils.PaginationParams, activeOnly bool) ([]models.Resource, int64, error) {
	var resources []models.Resource
	var total int64

	query := r.db.Model(&models.Resource{})

	// Filtrar solo activos si se especifica
	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	// Búsqueda por nombre
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("name LIKE ?", searchPattern)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener datos con paginación
	offset := params.CalculateOffset()
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&resources).Error; err != nil {
		return nil, 0, err
	}

	return resources, total, nil
}

// FindByCategory obtiene recursos por categoría
func (r *ResourceRepository) FindByCategory(category string, params utils.PaginationParams) ([]models.Resource, int64, error) {
	var resources []models.Resource
	var total int64

	query := r.db.Model(&models.Resource{}).Where("category = ? AND is_active = ?", category, true)

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener datos con paginación
	offset := params.CalculateOffset()
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&resources).Error; err != nil {
		return nil, 0, err
	}

	return resources, total, nil
}

// FindByID busca un recurso por ID
func (r *ResourceRepository) FindByID(id uint) (*models.Resource, error) {
	var resource models.Resource
	err := r.db.First(&resource, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("recurso no encontrado")
		}
		return nil, err
	}
	return &resource, nil
}

// Create crea un nuevo recurso
func (r *ResourceRepository) Create(resource *models.Resource) error {
	return r.db.Create(resource).Error
}

// Update actualiza un recurso
func (r *ResourceRepository) Update(resource *models.Resource) error {
	return r.db.Save(resource).Error
}

// Delete elimina un recurso (soft delete)
func (r *ResourceRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Resource{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("recurso no encontrado")
	}
	return nil
}

// GetCategories obtiene todas las categorías únicas
func (r *ResourceRepository) GetCategories() ([]string, error) {
	var categories []string
	err := r.db.Model(&models.Resource{}).
		Where("is_active = ? AND category != ''", true).
		Distinct("category").
		Pluck("category", &categories).Error
	return categories, err
}

// Count cuenta el total de recursos
func (r *ResourceRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Resource{}).Where("is_active = ?", true).Count(&count).Error
	return count, err
}
