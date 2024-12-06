package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/images"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	ntype "2024_2_FIGHT-CLUB/internal/service/type"
	"2024_2_FIGHT-CLUB/internal/service/validation"
	"go.uber.org/zap"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"context"
	"errors"
)

type AdUseCase interface {
	GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error)
	GetOnePlace(ctx context.Context, adId string, isAuthorized bool) (domain.GetAllAdsResponse, error)
	CreatePlace(ctx context.Context, place *domain.Ad, fileHeader [][]byte, newPlace domain.CreateAdRequest, userId string) error
	UpdatePlace(ctx context.Context, place *domain.Ad, adId string, userId string, fileHeader [][]byte, updatedPlace domain.UpdateAdRequest) error
	DeletePlace(ctx context.Context, adId string, userId string) error
	GetPlacesPerCity(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error)
	GetUserPlaces(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error)
	DeleteAdImage(ctx context.Context, adId string, imageId string, userId string) error
	AddToFavorites(ctx context.Context, adId string, userId string) error
	DeleteFromFavorites(ctx context.Context, adId string, userId string) error
	GetUserFavorites(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error)
}

type adUseCase struct {
	adRepository domain.AdRepository
	minioService images.MinioServiceInterface
}

func NewAdUseCase(adRepository domain.AdRepository, minioService images.MinioServiceInterface) AdUseCase {
	return &adUseCase{
		adRepository: adRepository,
		minioService: minioService,
	}
}

func (uc *adUseCase) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
	ads, err := uc.adRepository.GetAllPlaces(ctx, filter)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

func (uc *adUseCase) GetOnePlace(ctx context.Context, adId string, isAuthorized bool) (domain.GetAllAdsResponse, error) {
	const maxLen = 255
	requestID := middleware.GetRequestID(ctx)
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(adId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return domain.GetAllAdsResponse{}, errors.New("input contains invalid characters")
	}

	if len(adId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return domain.GetAllAdsResponse{}, errors.New("input exceeds character limit")
	}

	ad, err := uc.adRepository.GetPlaceById(ctx, adId)
	if err != nil {
		return ad, err
	}

	if isAuthorized {
		ad, err = uc.adRepository.UpdateViewsCount(ctx, ad)
		if err != nil {
			return ad, err
		}
	}

	return ad, nil
}

func (uc *adUseCase) CreatePlace(ctx context.Context, place *domain.Ad, files [][]byte, newPlace domain.CreateAdRequest, userId string) error {
	const maxLen = 255
	requestID := middleware.GetRequestID(ctx)

	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s\-!?&;#()/$*^%+=|]*$`)
	if !validCharPattern.MatchString(newPlace.CityName) ||
		!validCharPattern.MatchString(newPlace.Description) ||
		!validCharPattern.MatchString(newPlace.Address) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if len(newPlace.CityName) > maxLen || len(newPlace.Description) > maxLen || len(newPlace.Address) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}

	const minRooms, maxRooms = 1, 100
	if newPlace.RoomsNumber < minRooms || newPlace.RoomsNumber > maxRooms {
		logger.AccessLogger.Warn("RoomsNumber out of range", zap.String("request_id", requestID))
		return errors.New("RoomsNumber out of range")
	}

	if err := validation.ValidateImages(files, 5<<20, []string{"image/jpeg", "image/png", "image/jpg"}, 2000, 2000); err != nil {
		logger.AccessLogger.Warn("Invalid image", zap.String("request_id", requestID), zap.Error(err))
		return errors.New("invalid size, type or resolution of image")
	}

	place.Description = newPlace.Description
	place.Address = newPlace.Address
	place.RoomsNumber = newPlace.RoomsNumber
	err := uc.adRepository.CreatePlace(ctx, place, newPlace, userId)
	if err != nil {
		return err
	}
	var uploadedPaths ntype.StringArray

	for _, file := range files {
		if file != nil {
			contentType := http.DetectContentType(file[:512])

			uploadedPath, err := uc.minioService.UploadFile(file, contentType, "ads/"+place.UUID)
			if err != nil {
				for _, path := range uploadedPaths {
					_ = uc.minioService.DeleteFile(path)
				}
				return err
			}
			uploadedPaths = append(uploadedPaths, "/images/"+uploadedPath)
		}
	}

	err = uc.adRepository.SaveImages(ctx, place.UUID, uploadedPaths)
	if err != nil {
		return err
	}

	return nil
}

func (uc *adUseCase) UpdatePlace(ctx context.Context, place *domain.Ad, adId string, userId string, files [][]byte, updatedPlace domain.UpdateAdRequest) error {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255

	if len(files) > 0 {
		if err := validation.ValidateImages(files, 5<<20, []string{"image/jpeg", "image/png", "image/jpg"}, 2000, 2000); err != nil {
			logger.AccessLogger.Warn("Invalid image", zap.String("request_id", requestID), zap.Error(err))
			return errors.New("invalid size, type or resolution of image")
		}
	}

	validCharPatternUrl := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPatternUrl.MatchString(adId) {
		logger.AccessLogger.Warn("URL contains invalid characters", zap.String("request_id", requestID))
		return errors.New("URL contains invalid characters")
	}

	if len(adId) > maxLen {
		logger.AccessLogger.Warn("URL exceeds character limit", zap.String("request_id", requestID))
		return errors.New("URL exceeds character limit")
	}

	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-Я0-9@.,\s\-!?&;#()/$*^%+=|]*$`)
	if !validCharPattern.MatchString(updatedPlace.CityName) ||
		!validCharPattern.MatchString(updatedPlace.Description) ||
		!validCharPattern.MatchString(updatedPlace.Address) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if len(updatedPlace.CityName) > maxLen || len(updatedPlace.Description) > maxLen || len(updatedPlace.Address) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}

	const minRooms, maxRooms = 1, 100
	if updatedPlace.RoomsNumber < minRooms || updatedPlace.RoomsNumber > maxRooms {
		logger.AccessLogger.Warn("RoomsNumber out of range", zap.String("request_id", requestID))
		return errors.New("RoomsNumber out of range")
	}

	_, err := uc.adRepository.GetPlaceById(ctx, adId)
	if err != nil {
		return err
	}
	place.Description = updatedPlace.Description
	place.Address = updatedPlace.Address
	place.RoomsNumber = updatedPlace.RoomsNumber
	var newUploadedPaths ntype.StringArray

	for _, file := range files {
		if file != nil {
			contentType := http.DetectContentType(file[:512])
			uploadedPath, err := uc.minioService.UploadFile(file, contentType, "ads/"+adId)
			if err != nil {
				for _, path := range newUploadedPaths {
					_ = uc.minioService.DeleteFile(path)
				}
				return err
			}
			newUploadedPaths = append(newUploadedPaths, "/images/"+uploadedPath)
		}
	}

	err = uc.adRepository.UpdatePlace(ctx, place, adId, userId, updatedPlace)
	if err != nil {
		return err
	}

	err = uc.adRepository.SaveImages(ctx, adId, newUploadedPaths)
	if err != nil {
		return err
	}
	return nil
}

func (uc *adUseCase) DeletePlace(ctx context.Context, adId string, userId string) error {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(adId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if len(adId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}

	_, err := uc.adRepository.GetPlaceById(ctx, adId)
	if err != nil {
		return err
	}
	imagesPath, err := uc.adRepository.GetAdImages(ctx, adId)
	if err != nil {
		return err
	}
	for _, imagePath := range imagesPath {
		_ = uc.minioService.DeleteFile(imagePath)
	}

	err = uc.adRepository.DeletePlace(ctx, adId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc *adUseCase) GetPlacesPerCity(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(city) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return []domain.GetAllAdsResponse{}, errors.New("input contains invalid characters")
	}

	if len(city) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return []domain.GetAllAdsResponse{}, errors.New("input exceeds character limit")
	}

	places, err := uc.adRepository.GetPlacesPerCity(ctx, city)
	if err != nil {
		return nil, err
	}
	return places, nil
}

func (uc *adUseCase) GetUserPlaces(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(userId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return []domain.GetAllAdsResponse{}, errors.New("input contains invalid characters")
	}

	if len(userId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return []domain.GetAllAdsResponse{}, errors.New("input exceeds character limit")
	}

	places, err := uc.adRepository.GetUserPlaces(ctx, userId)
	if err != nil {
		return nil, err
	}
	return places, nil
}

func (uc *adUseCase) DeleteAdImage(ctx context.Context, adId string, imageId string, userId string) error {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(adId) || !validCharPattern.MatchString(imageId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if len(adId) > maxLen || len(imageId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}

	imageIdInt, err2 := strconv.Atoi(imageId)
	if err2 != nil {
		logger.AccessLogger.Warn("Failed to ATOI image url", zap.String("request_id", requestID), zap.Error(err2))
		return errors.New("failed to ATOI image url")
	}

	imageURL, err := uc.adRepository.DeleteAdImage(ctx, adId, imageIdInt, userId)
	if err != nil {
		return err
	}

	if err := uc.minioService.DeleteFile(imageURL); err != nil {
		log.Printf("Warning: failed to delete file from MinIO: %v", err)
	}

	return nil
}

func (uc *adUseCase) AddToFavorites(ctx context.Context, adId string, userId string) error {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(adId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if len(adId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}

	err := uc.adRepository.AddToFavorites(ctx, adId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc *adUseCase) DeleteFromFavorites(ctx context.Context, adId string, userId string) error {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(adId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return errors.New("input contains invalid characters")
	}

	if len(adId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return errors.New("input exceeds character limit")
	}

	err := uc.adRepository.DeleteFromFavorites(ctx, adId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (uc *adUseCase) GetUserFavorites(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	const maxLen = 255
	validCharPattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ0-9\s\-_]*$`)
	if !validCharPattern.MatchString(userId) {
		logger.AccessLogger.Warn("Input contains invalid characters", zap.String("request_id", requestID))
		return nil, errors.New("input contains invalid characters")
	}

	if len(userId) > maxLen {
		logger.AccessLogger.Warn("Input exceeds character limit", zap.String("request_id", requestID))
		return nil, errors.New("input exceeds character limit")
	}

	places, err := uc.adRepository.GetUserFavorites(ctx, userId)
	if err != nil {
		return nil, err
	}

	return places, nil
}
