package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type adRepository struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) domain.AdRepository {
	return &adRepository{
		db: db,
	}
}

func (r *adRepository) GetAllPlaces(ctx context.Context) ([]domain.Ad, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetAllPlaces called", zap.String("request_id", requestID))
	var ads []domain.Ad
	if err := r.db.Find(&ads).Error; err != nil {
		logger.DBLogger.Error("Error fetching all places", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}

	logger.DBLogger.Info("Successfully fetched all places", zap.String("request_id", requestID), zap.Int("count", len(ads)))
	return ads, nil
}

func (r *adRepository) GetPlaceById(ctx context.Context, adId string) (domain.Ad, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetPlaceById called", zap.String("adId", adId), zap.String("request_id", requestID))

	var ad domain.Ad
	if err := r.db.Where("id = ?", adId).First(&ad).Error; err != nil {
		logger.DBLogger.Error("Error fetching place by ID", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return ad, err
	}

	logger.DBLogger.Info("Successfully fetched place by ID", zap.String("adId", adId), zap.String("request_id", requestID))
	return ad, nil
}

func (r *adRepository) CreatePlace(ctx context.Context, ad *domain.Ad) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("CreatePlace called", zap.String("adId", ad.ID), zap.String("request_id", requestID))

	if err := r.db.Create(ad).Error; err != nil {
		logger.DBLogger.Error("Error creating place", zap.String("adId", ad.ID), zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	if err := r.SavePlace(ctx, ad); err != nil {
		logger.DBLogger.Error("Error saving place after creation", zap.String("adId", ad.ID), zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	logger.DBLogger.Info("Successfully created and saved place", zap.String("adId", ad.ID), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) SavePlace(ctx context.Context, ad *domain.Ad) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("SavePlace called", zap.String("adId", ad.ID), zap.String("request_id", requestID))
	if err := r.db.Save(ad).Error; err != nil {
		logger.DBLogger.Error("Error saving place", zap.String("adId", ad.ID), zap.String("request_id", requestID), zap.Error(err))
		return err
	}
	logger.DBLogger.Info("Successfully saved place", zap.String("adId", ad.ID), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) UpdatePlace(ctx context.Context, ad *domain.Ad, adId string, userId string) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("UpdatePlace called", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))

	var oldAd domain.Ad
	if err := r.db.Where("id = ?", adId).First(&oldAd).Error; err != nil {
		logger.DBLogger.Error("Ad not found", zap.String("adId", adId), zap.String("request_id", requestID))
		return errors.New("ad not found")
	}

	if oldAd.AuthorUUID != userId {
		logger.DBLogger.Warn("User is not the owner of the ad", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))
		return errors.New("not owner of ad")
	}

	if err := r.db.Model(&oldAd).Updates(ad).Error; err != nil {
		logger.DBLogger.Error("Error updating place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	logger.DBLogger.Info("Successfully updated place", zap.String("adId", adId), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) DeletePlace(ctx context.Context, adId string, userId string) error {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("DeletePlace called", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))

	var ad domain.Ad
	if err := r.db.Where("id = ?", adId).First(&ad).Error; err != nil {
		logger.DBLogger.Error("Ad not found", zap.String("adId", adId), zap.String("request_id", requestID))
		return errors.New("ad not found")
	}

	if ad.AuthorUUID != userId {
		logger.DBLogger.Warn("User is not the owner of the ad", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))
		return errors.New("not owner of ad")
	}

	if err := r.db.Delete(&ad).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return err
	}

	logger.DBLogger.Info("Successfully deleted place", zap.String("adId", adId), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) GetPlacesPerCity(ctx context.Context, city string) ([]domain.Ad, error) {
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetPlacesPerCity called", zap.String("city", city), zap.String("request_id", requestID))

	var ads []domain.Ad
	if err := r.db.Where("location_main = ?", city).Find(&ads).Error; err != nil {
		logger.DBLogger.Error("Error fetching places per city", zap.String("city", city), zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}

	logger.DBLogger.Info("Successfully fetched places per city", zap.String("city", city), zap.Int("count", len(ads)), zap.String("request_id", requestID))
	return ads, nil
}
