package repositories

import (
	"Reservify/models"
	"errors"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser crea un nuevo usuario
func (r *AuthRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByEmail busca un usuario por email
func (r *AuthRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("usuario no encontrado")
		}
		return nil, err
	}
	return &user, nil
}

// FindByID busca un usuario por ID
func (r *AuthRepository) FindByID(id uint) (*models.User, error) {
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

// EmailExists verifica si un email ya está registrado
func (r *AuthRepository) EmailExists(email string) bool {
	var count int64
	r.db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count > 0
}
