package services

import (
	"Reservify/dto"
	"Reservify/models"
	"Reservify/repositories"
	"errors"
	"fmt"
	"strings"
)

type AvailabilityService struct {
	availabilityRepo *repositories.AvailabilityRepository
	resourceRepo     *repositories.ResourceRepository
}

func NewAvailabilityService(availabilityRepo *repositories.AvailabilityRepository, resourceRepo *repositories.ResourceRepository) *AvailabilityService {
	return &AvailabilityService{
		availabilityRepo: availabilityRepo,
		resourceRepo:     resourceRepo,
	}
}

// GetAvailabilityByResource obtiene la disponibilidad de un recurso
func (s *AvailabilityService) GetAvailabilityByResource(resourceID uint) ([]dto.AvailabilityResponse, error) {
	// Verificar que el recurso existe
	if _, err := s.resourceRepo.FindByID(resourceID); err != nil {
		return nil, err
	}

	slots, err := s.availabilityRepo.FindByResourceID(resourceID)
	if err != nil {
		return nil, err
	}

	var response []dto.AvailabilityResponse
	for _, slot := range slots {
		response = append(response, dto.AvailabilityResponse{
			ID:         slot.ID,
			ResourceID: slot.ResourceID,
			DayOfWeek:  string(slot.DayOfWeek),
			StartTime:  slot.StartTime,
			EndTime:    slot.EndTime,
			CreatedAt:  slot.CreatedAt,
		})
	}

	return response, nil
}

// CreateAvailability crea un nuevo horario de disponibilidad
func (s *AvailabilityService) CreateAvailability(resourceID uint, req *dto.CreateAvailabilityRequest) (*dto.AvailabilityResponse, error) {
	// Verificar que el recurso existe
	if _, err := s.resourceRepo.FindByID(resourceID); err != nil {
		return nil, err
	}

	// Normalizar tiempos (agregar :00 si solo tiene HH:MM)
	startTime := normalizeTime(req.StartTime)
	endTime := normalizeTime(req.EndTime)

	// Validar que start_time < end_time
	if startTime >= endTime {
		return nil, errors.New("la hora de inicio debe ser menor que la hora de fin")
	}

	// Verificar solapamiento
	overlap, err := s.availabilityRepo.CheckOverlap(resourceID, req.DayOfWeek, startTime, endTime, nil)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errors.New("el horario se solapa con uno existente")
	}

	slot := &models.AvailabilitySlot{
		ResourceID: resourceID,
		DayOfWeek:  models.DayOfWeek(req.DayOfWeek),
		StartTime:  startTime,
		EndTime:    endTime,
	}

	if err := s.availabilityRepo.Create(slot); err != nil {
		return nil, errors.New("error al crear el horario")
	}

	response := &dto.AvailabilityResponse{
		ID:         slot.ID,
		ResourceID: slot.ResourceID,
		DayOfWeek:  string(slot.DayOfWeek),
		StartTime:  slot.StartTime,
		EndTime:    slot.EndTime,
		CreatedAt:  slot.CreatedAt,
	}

	return response, nil
}

// UpdateAvailability actualiza un horario
func (s *AvailabilityService) UpdateAvailability(id uint, req *dto.UpdateAvailabilityRequest) (*dto.AvailabilityResponse, error) {
	slot, err := s.availabilityRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Actualizar campos
	if req.DayOfWeek != "" {
		slot.DayOfWeek = models.DayOfWeek(req.DayOfWeek)
	}
	if req.StartTime != "" {
		slot.StartTime = normalizeTime(req.StartTime)
	}
	if req.EndTime != "" {
		slot.EndTime = normalizeTime(req.EndTime)
	}

	// Validar que start_time < end_time
	if slot.StartTime >= slot.EndTime {
		return nil, errors.New("la hora de inicio debe ser menor que la hora de fin")
	}

	// Verificar solapamiento (excluyendo el ID actual)
	overlap, err := s.availabilityRepo.CheckOverlap(slot.ResourceID, string(slot.DayOfWeek), slot.StartTime, slot.EndTime, &id)
	if err != nil {
		return nil, err
	}
	if overlap {
		return nil, errors.New("el horario se solapa con uno existente")
	}

	if err := s.availabilityRepo.Update(slot); err != nil {
		return nil, errors.New("error al actualizar el horario")
	}

	response := &dto.AvailabilityResponse{
		ID:         slot.ID,
		ResourceID: slot.ResourceID,
		DayOfWeek:  string(slot.DayOfWeek),
		StartTime:  slot.StartTime,
		EndTime:    slot.EndTime,
		CreatedAt:  slot.CreatedAt,
	}

	return response, nil
}

// DeleteAvailability elimina un horario
func (s *AvailabilityService) DeleteAvailability(id uint) error {
	return s.availabilityRepo.Delete(id)
}

// normalizeTime normaliza el formato de tiempo (agrega :00 si es necesario)
func normalizeTime(timeStr string) string {
	// Si ya tiene formato HH:MM:SS, devolverlo tal cual
	if len(timeStr) == 8 && strings.Count(timeStr, ":") == 2 {
		return timeStr
	}
	// Si tiene formato HH:MM, agregar :00
	if len(timeStr) == 5 && strings.Count(timeStr, ":") == 1 {
		return fmt.Sprintf("%s:00", timeStr)
	}
	return timeStr
}
