package repositories

import (
	"Reservify/models"
	"Reservify/utils"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindAll obtiene todos los usuarios con paginación
func (r *UserRepository) FindAll(params utils.PaginationParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.db.Model(&models.User{})

	// Búsqueda por nombre o email
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("full_name LIKE ? OR email LIKE ?", searchPattern, searchPattern)
	}

	// Contar total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener datos con paginación
	offset := params.CalculateOffset()
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindByID busca un usuario por ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return &user, nil
}

// Update actualiza un usuario
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete elimina un usuario (soft delete)
func (r *UserRepository) Delete(id uint) error {
	result := r.db.Delete(&models.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("usuario no encontrado")
	}
	return nil
}

// UpdatePassword actualiza solo la contraseña
func (r *UserRepository) UpdatePassword(userID uint, newPasswordHash string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("password_hash", newPasswordHash).Error
}

// Count cuenta el total de usuarios
func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.User{}).Count(&count).Error
	return count, err
}
