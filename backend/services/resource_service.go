package services

import (
	"Reservify/dto"
	"Reservify/models"
	"Reservify/repositories"
	"Reservify/utils"
	"errors"
)

type ResourceService struct {
	resourceRepo *repositories.ResourceRepository
}

func NewResourceService(resourceRepo *repositories.ResourceRepository) *ResourceService {
	return &ResourceService{resourceRepo: resourceRepo}
}

// GetAllResources obtiene todos los recursos con paginación
func (s *ResourceService) GetAllResources(params utils.PaginationParams, activeOnly bool) ([]dto.ResourceListResponse, int64, error) {
	resources, total, err := s.resourceRepo.FindAll(params, activeOnly)
	if err != nil {
		return nil, 0, err
	}

	var response []dto.ResourceListResponse
	for _, resource := range resources {
		response = append(response, dto.ResourceListResponse{
			ID:           resource.ID,
			Name:         resource.Name,
			Capacity:     resource.Capacity,
			PricePerHour: resource.PricePerHour,
			Category:     resource.Category,
			ImageURL:     resource.ImageURL,
			IsActive:     resource.IsActive,
		})
	}

	return response, total, nil
}

// GetResourcesByCategory obtiene recursos por categoría
func (s *ResourceService) GetResourcesByCategory(category string, params utils.PaginationParams) ([]dto.ResourceListResponse, int64, error) {
	resources, total, err := s.resourceRepo.FindByCategory(category, params)
	if err != nil {
		return nil, 0, err
	}

	var response []dto.ResourceListResponse
	for _, resource := range resources {
		response = append(response, dto.ResourceListResponse{
			ID:           resource.ID,
			Name:         resource.Name,
			Capacity:     resource.Capacity,
			PricePerHour: resource.PricePerHour,
			Category:     resource.Category,
			ImageURL:     resource.ImageURL,
			IsActive:     resource.IsActive,
		})
	}

	return response, total, nil
}

// GetResourceByID obtiene un recurso por ID
func (s *ResourceService) GetResourceByID(id uint) (*dto.ResourceResponse, error) {
	resource, err := s.resourceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := &dto.ResourceResponse{
		ID:           resource.ID,
		Name:         resource.Name,
		Description:  resource.Description,
		Capacity:     resource.Capacity,
		PricePerHour: resource.PricePerHour,
		Category:     resource.Category,
		ImageURL:     resource.ImageURL,
		IsActive:     resource.IsActive,
		CreatedAt:    resource.CreatedAt,
		UpdatedAt:    resource.UpdatedAt,
	}

	return response, nil
}

// CreateResource crea un nuevo recurso
func (s *ResourceService) CreateResource(req *dto.CreateResourceRequest) (*dto.ResourceResponse, error) {
	resource := &models.Resource{
		Name:         req.Name,
		Description:  req.Description,
		Capacity:     req.Capacity,
		PricePerHour: req.PricePerHour,
		Category:     req.Category,
		ImageURL:     req.ImageURL,
		IsActive:     true,
	}

	if err := s.resourceRepo.Create(resource); err != nil {
		return nil, errors.New("error al crear el recurso")
	}

	return s.GetResourceByID(resource.ID)
}

// UpdateResource actualiza un recurso
func (s *ResourceService) UpdateResource(id uint, req *dto.UpdateResourceRequest) (*dto.ResourceResponse, error) {
	resource, err := s.resourceRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.Name != "" {
		resource.Name = req.Name
	}
	if req.Description != "" {
		resource.Description = req.Description
	}
	if req.Capacity > 0 {
		resource.Capacity = req.Capacity
	}
	if req.PricePerHour >= 0 {
		resource.PricePerHour = req.PricePerHour
	}
	if req.Category != "" {
		resource.Category = req.Category
	}
	if req.ImageURL != "" {
		resource.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		resource.IsActive = *req.IsActive
	}

	if err := s.resourceRepo.Update(resource); err != nil {
		return nil, errors.New("error al actualizar el recurso")
	}

	return s.GetResourceByID(resource.ID)
}

// DeleteResource elimina un recurso
func (s *ResourceService) DeleteResource(id uint) error {
	return s.resourceRepo.Delete(id)
}

// GetCategories obtiene todas las categorías
func (s *ResourceService) GetCategories() ([]string, error) {
	return s.resourceRepo.GetCategories()
}

// GetResourceStats obtiene estadísticas de recursos
func (s *ResourceService) GetResourceStats() (map[string]interface{}, error) {
	total, err := s.resourceRepo.Count()
	if err != nil {
		return nil, err
	}

	categories, err := s.resourceRepo.GetCategories()
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_resources": total,
		"categories":      categories,
	}

	return stats, nil
}
