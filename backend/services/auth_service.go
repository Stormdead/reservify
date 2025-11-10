package services

import (
	"Reservify/dto"
	"Reservify/models"
	"Reservify/repositories"
	"Reservify/utils"
	"errors"
)

type AuthService struct {
	authRepo *repositories.AuthRepository
}

func NewAuthService(authRepo *repositories.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

// Register registra un nuevo usuario
func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Verificar si el email ya existe
	if s.authRepo.EmailExists(req.Email) {
		return nil, errors.New("el email ya está registrado")
	}

	// Hash de la contraseña
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}

	// Crear el usuario
	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Phone:        req.Phone,
		Role:         models.RoleUser, // Por defecto es user
	}

	if err := s.authRepo.CreateUser(user); err != nil {
		return nil, errors.New("error al crear el usuario")
	}

	// Generar token
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, errors.New("error al generar el token")
	}

	// Preparar respuesta
	response := &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		},
	}

	return response, nil
}

// Login autentica a un usuario
func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	// Buscar usuario por email
	user, err := s.authRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	// Verificar contraseña
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("credenciales inválidas")
	}

	// Generar token
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, errors.New("error al generar el token")
	}

	// Preparar respuesta
	response := &dto.AuthResponse{
		Token: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FullName:  user.FullName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
		},
	}

	return response, nil
}

// Obtiene un usuario por ID
func (s *AuthService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.authRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Phone:     user.Phone,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
	}

	return response, nil
}
