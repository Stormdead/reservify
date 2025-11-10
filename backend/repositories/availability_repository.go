package repositories

import (
	"Reservify/models"
	"errors"

	"gorm.io/gorm"
)

type AvailabilityRepository struct {
	db *gorm.DB
}

func NewAvailabilityRepository(db *gorm.DB) *AvailabilityRepository {
	return &AvailabilityRepository{db: db}
}

// FindByResourceID obtiene todos los horarios de un recurso
func (r *AvailabilityRepository) FindByResourceID(resourceID uint) ([]models.AvailabilitySlot, error) {
	var slots []models.AvailabilitySlot
	err := r.db.Where("resource_id = ?", resourceID).Order("day_of_week, start_time").Find(&slots).Error
	return slots, err
}

// FindByID busca un horario por ID
func (r *AvailabilityRepository) FindByID(id uint) (*models.AvailabilitySlot, error) {
	var slot models.AvailabilitySlot
	err := r.db.First(&slot, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("horario no encontrado")
		}
		return nil, err
	}
	return &slot, nil
}

// Create crea un nuevo horario
func (r *AvailabilityRepository) Create(slot *models.AvailabilitySlot) error {
	return r.db.Create(slot).Error
}

// Update actualiza un horario
func (r *AvailabilityRepository) Update(slot *models.AvailabilitySlot) error {
	return r.db.Save(slot).Error
}

// Delete elimina un horario
func (r *AvailabilityRepository) Delete(id uint) error {
	result := r.db.Delete(&models.AvailabilitySlot{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("horario no encontrado")
	}
	return nil
}

// DeleteByResourceID elimina todos los horarios de un recurso
func (r *AvailabilityRepository) DeleteByResourceID(resourceID uint) error {
	return r.db.Where("resource_id = ?", resourceID).Delete(&models.AvailabilitySlot{}).Error
}

// CheckOverlap verifica si hay solapamiento de horarios
func (r *AvailabilityRepository) CheckOverlap(resourceID uint, dayOfWeek, startTime, endTime string, excludeID *uint) (bool, error) {
	query := r.db.Model(&models.AvailabilitySlot{}).
		Where("resource_id = ? AND day_of_week = ?", resourceID, dayOfWeek).
		Where("(start_time < ? AND end_time > ?) OR (start_time < ? AND end_time > ?) OR (start_time >= ? AND end_time <= ?)",
			endTime, startTime, // Caso 1: El nuevo horario está dentro de uno existente
			endTime, startTime, // Caso 2: El nuevo horario envuelve uno existente
			startTime, endTime) // Caso 3: El nuevo horario está completamente dentro

	// Excluir el ID actual si estamos actualizando
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
