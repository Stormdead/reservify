package repositories

import (
	"Reservify/models"
	"Reservify/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

// FindAll obtiene todas las reservas con paginación
func (r *BookingRepository) FindAll(params utils.PaginationParams) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).
		Preload("User").
		Preload("Resource")

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener datos con paginación
	offset := params.CalculateOffset()
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&bookings).Error; err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

// FindByUserID obtiene las reservas de un usuario
func (r *BookingRepository) FindByUserID(userID uint, params utils.PaginationParams) ([]models.Booking, int64, error) {
	var bookings []models.Booking
	var total int64

	query := r.db.Model(&models.Booking{}).
		Where("user_id = ?", userID).
		Preload("Resource")

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener datos con paginación
	offset := params.CalculateOffset()
	if err := query.Offset(offset).Limit(params.PageSize).Order("start_datetime DESC").Find(&bookings).Error; err != nil {
		return nil, 0, err
	}

	return bookings, total, nil
}

// FindByID busca una reserva por ID
func (r *BookingRepository) FindByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.Preload("User").Preload("Resource").First(&booking, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("reserva no encontrada")
		}
		return nil, err
	}
	return &booking, nil
}

// Create crea una nueva reserva
func (r *BookingRepository) Create(booking *models.Booking) error {
	return r.db.Create(booking).Error
}

// Update actualiza una reserva
func (r *BookingRepository) Update(booking *models.Booking) error {
	return r.db.Save(booking).Error
}

// Delete elimina una reserva (soft delete)
func (r *BookingRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Booking{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("reserva no encontrada")
	}
	return nil
}

// CheckOverlap verifica si hay solapamiento de reservas para un recurso
func (r *BookingRepository) CheckOverlap(resourceID uint, startDatetime, endDatetime time.Time, excludeID *uint) (bool, error) {
	query := r.db.Model(&models.Booking{}).
		Where("resource_id = ?", resourceID).
		Where("status IN ?", []string{"pending", "confirmed"}). // Solo considerar reservas activas
		Where(
			r.db.Where("start_datetime < ? AND end_datetime > ?", endDatetime, startDatetime), // Solapamiento
		)

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

// CountByStatus cuenta reservas por estado
func (r *BookingRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

// GetTotalRevenue calcula el total de ingresos
func (r *BookingRepository) GetTotalRevenue() (float64, error) {
	var total float64
	err := r.db.Model(&models.Booking{}).
		Where("status IN ?", []string{"confirmed", "completed"}).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&total).Error
	return total, err
}

// GetRevenueByStatus calcula ingresos por estado
func (r *BookingRepository) GetRevenueByStatus(status string) (float64, error) {
	var total float64
	err := r.db.Model(&models.Booking{}).
		Where("status = ?", status).
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&total).Error
	return total, err
}

// GetUpcomingBookings obtiene reservas próximas
func (r *BookingRepository) GetUpcomingBookings(userID uint, limit int) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("user_id = ? AND start_datetime > ? AND status IN ?",
		userID, time.Now(), []string{"pending", "confirmed"}).
		Preload("Resource").
		Order("start_datetime ASC").
		Limit(limit).
		Find(&bookings).Error
	return bookings, err
}

// GetBookingsByDateRange obtiene reservas en un rango de fechas
func (r *BookingRepository) GetBookingsByDateRange(resourceID uint, startDate, endDate time.Time) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("resource_id = ? AND start_datetime >= ? AND end_datetime <= ?",
		resourceID, startDate, endDate).
		Where("status IN ?", []string{"pending", "confirmed"}).
		Order("start_datetime ASC").
		Find(&bookings).Error
	return bookings, err
}
