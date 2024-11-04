package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/images"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/validation"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
)

type AuthUseCase interface {
	RegisterUser(ctx context.Context, creds *domain.User) error
	LoginUser(ctx context.Context, creds *domain.User) (*domain.User, error)
	PutUser(ctx context.Context, creds *domain.User, userID string, avatar *multipart.FileHeader) error
	GetAllUser(ctx context.Context) ([]domain.User, error)
	GetUserById(ctx context.Context, userID string) (*domain.User, error)
}

type authUseCase struct {
	authRepository domain.AuthRepository
	minioService   images.MinioServiceInterface
}

func NewAuthUseCase(authRepository domain.AuthRepository, minioService images.MinioServiceInterface) AuthUseCase {
	return &authUseCase{
		authRepository: authRepository,
		minioService:   minioService,
	}
}

func (uc *authUseCase) RegisterUser(ctx context.Context, creds *domain.User) error {
	if creds.Username == "" || creds.Password == "" || creds.Email == "" {
		return errors.New("username, password, and email are required")
	}
	errorResponse := map[string]interface{}{
		"error":       "Incorrect data forms",
		"wrongFields": []string{},
	}
	var wrongFields []string
	if !validation.ValidateLogin(creds.Username) {
		wrongFields = append(wrongFields, "username")
	}
	if !validation.ValidateEmail(creds.Email) {
		wrongFields = append(wrongFields, "email")
	}
	if !validation.ValidatePassword(creds.Password) {
		wrongFields = append(wrongFields, "password")
	}
	if !validation.ValidateName(creds.Name) {
		wrongFields = append(wrongFields, "name")
	}
	if len(wrongFields) > 0 {
		errorResponse["wrongFields"] = wrongFields
		errorResponseJSON, err := json.Marshal(errorResponse)
		if err != nil {
			return errors.New("failed to generate error response")
		}
		return errors.New(string(errorResponseJSON))
	}

	// Хэширование пароля
	hashedPassword, ok := middleware.HashPassword(creds.Password)
	if ok != nil {
		return errors.New("failed to hash password")
	}
	creds.Password = hashedPassword

	existingUser, _ := uc.authRepository.GetUserByName(ctx, creds.Username)
	if existingUser != nil {
		return errors.New("user already exists")
	}
	err := uc.authRepository.CreateUser(ctx, creds)
	if err != nil {
		return err
	}

	return uc.authRepository.SaveUser(ctx, creds)
}

func (uc *authUseCase) LoginUser(ctx context.Context, creds *domain.User) (*domain.User, error) {
	if creds.Username == "" || creds.Password == "" {
		return nil, errors.New("username and password are required")
	}
	errorResponse := map[string]interface{}{
		"error":       "Incorrect data forms",
		"wrongFields": []string{},
	}
	var wrongFields []string
	if !validation.ValidateLogin(creds.Username) {
		wrongFields = append(wrongFields, "username")
	}
	if !validation.ValidatePassword(creds.Password) {
		wrongFields = append(wrongFields, "password")
	}
	if len(wrongFields) > 0 {
		errorResponse["wrongFields"] = wrongFields
		errorResponseJSON, err := json.Marshal(errorResponse)
		if err != nil {
			return nil, errors.New("failed to generate error response")
		}
		return nil, errors.New(string(errorResponseJSON))
	}

	requestedUser, err := uc.authRepository.GetUserByName(ctx, creds.Username)
	if err != nil || requestedUser == nil {
		return nil, errors.New("user not found")
	}

	if !middleware.CheckPassword(requestedUser.Password, creds.Password) {
		return nil, errors.New("invalid credentials")
	}

	return requestedUser, nil
}

func (uc *authUseCase) PutUser(ctx context.Context, creds *domain.User, userID string, avatar *multipart.FileHeader) error {
	var wrongFields []string
	errorResponse := map[string]interface{}{
		"error":       "Incorrect data forms",
		"wrongFields": []string{},
	}
	if !validation.ValidateLogin(creds.Username) && len(creds.Username) > 0 {
		wrongFields = append(wrongFields, "username")
	}
	if !validation.ValidateEmail(creds.Email) && len(creds.Email) > 0 {
		wrongFields = append(wrongFields, "email")
	}
	if !validation.ValidatePassword(creds.Password) && len(creds.Password) > 0 {
		wrongFields = append(wrongFields, "password")
	}
	if !validation.ValidateName(creds.Name) && len(creds.Name) > 0 {
		wrongFields = append(wrongFields, "name")
	}
	if len(wrongFields) > 0 {
		errorResponse["wrongFields"] = wrongFields
		errorResponseJSON, err := json.Marshal(errorResponse)
		if err != nil {
			return errors.New("failed to generate error response")
		}
		return errors.New(string(errorResponseJSON))
	}

	if avatar != nil {
		uploadedPath, err := uc.minioService.UploadFile(avatar, "user/"+userID)
		if err != nil {
			return err
		}

		creds.Avatar = "/images/" + uploadedPath
	}

	err := uc.authRepository.PutUser(ctx, creds, userID)
	if err != nil {
		_ = uc.minioService.DeleteFile(creds.Avatar)
	}
	return nil
}

func (uc *authUseCase) GetAllUser(ctx context.Context) ([]domain.User, error) {
	users, err := uc.authRepository.GetAllUser(ctx)
	if err != nil {
		return nil, errors.New("there is none user in db")
	}
	return users, nil
}

func (uc *authUseCase) GetUserById(ctx context.Context, userID string) (*domain.User, error) {
	user, err := uc.authRepository.GetUserById(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}
