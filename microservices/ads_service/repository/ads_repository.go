package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type adRepository struct {
	db *gorm.DB
}

func NewAdRepository(db *gorm.DB) domain.AdRepository {
	return &adRepository{
		db: db,
	}
}

func (r *adRepository) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.GetAllAdsResponse, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetAllPlaces called", zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetAllPlaces", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetAllPlaces", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetAllPlaces").Observe(duration)
	}()
	var ads []domain.GetAllAdsResponse

	query := r.db.Model(&domain.Ad{}).Joins("JOIN cities ON  ads.\"cityId\" = cities.id").
		Joins("JOIN users ON ads.\"authorUUID\" = users.uuid").
		Joins("JOIN ad_available_dates ON ad_available_dates.\"adId\" = ads.uuid").
		Select("ads.*, cities.title as \"CityName\", ad_available_dates.\"availableDateFrom\" as \"AdDateFrom\", ad_available_dates.\"availableDateTo\" as \"AdDateTo\"")

	if filter.Location != "" {
		query = query.Where("cities.\"enTitle\" = ?", filter.Location)
	}

	if filter.Rating != "" {
		rating, err := strconv.ParseFloat(filter.Rating, 64) // поменял с 1 на 64 из-за линтера
		if err != nil {
			logger.DBLogger.Error("Invalid rating value", zap.String("request_id", requestID))
			return nil, errors.New("invalid rating value")
		}
		query = query.Where("users.score >= ?", rating)
	}

	if filter.NewThisWeek == "true" {
		lastWeek := time.Now().AddDate(0, 0, -7)
		query = query.Where("\"publicationDate\" >= ?", lastWeek)
	}

	if filter.HostGender != "" && filter.HostGender != "any" {
		if filter.HostGender == "male" {
			query = query.Where("users.sex = ?", "M")
		} else if filter.HostGender == "female" {
			query = query.Where("users.sex = ?", "F")
		}
	}

	if filter.GuestCount != "" {
		switch filter.GuestCount {
		case "5":
			query = query.Where("users.\"guestCount\" > ?", 5)
		case "10":
			query = query.Where("users.\"guestCount\" > ?", 10)
		case "20":
			query = query.Where("users.\"guestCount\" > ?", 20)
		case "50":
			query = query.Where("users.\"guestCount\" > ?", 50)
		}
	}

	switch {
	case !filter.DateFrom.IsZero() && !filter.DateTo.IsZero():
		query = query.Where("ad_available_dates.\"availableDateFrom\" <= ? AND ad_available_dates.\"availableDateTo\" >= ?", filter.DateTo, filter.DateFrom)
	case !filter.DateFrom.IsZero():
		query = query.Where("ad_available_dates.\"availableDateTo\" >= ?", filter.DateFrom)
	case !filter.DateTo.IsZero():
		query = query.Where("ad_available_dates.\"availableDateFrom\" <= ?", filter.DateTo)
	}

	if filter.Offset != 0 {
		query = query.Offset(filter.Offset)
	}

	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if err := query.Find(&ads).Error; err != nil {
		logger.DBLogger.Error("Error fetching all places", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("error fetching all places")
	}

	for i, ad := range ads {
		var images []domain.Image
		var user domain.User
		err := r.db.Model(&domain.Image{}).Where("\"adId\" = ?", ad.UUID).Find(&images).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("error fetching images for ad")
		}

		err = r.db.Model(&domain.User{}).Where("uuid = ?", ad.AuthorUUID).Find(&user).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching user", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("error fetching user")
		}
		ads[i].AdAuthor.Name = user.Name
		ads[i].AdAuthor.Avatar = user.Avatar
		ads[i].AdAuthor.Rating = user.Score
		ads[i].AdAuthor.GuestCount = user.GuestCount
		ads[i].AdAuthor.Sex = user.Sex
		ads[i].AdAuthor.Birthdate = user.Birthdate
		for _, img := range images {
			ads[i].Images = append(ads[i].Images, domain.ImageResponse{
				ID:        img.ID,
				ImagePath: img.ImageUrl,
			})
		}
	}

	logger.DBLogger.Info("Successfully fetched all places", zap.String("request_id", requestID), zap.Int("count", len(ads)))
	return ads, nil
}

func (r *adRepository) GetPlaceById(ctx context.Context, adId string) (domain.GetAllAdsResponse, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetPlaceById called", zap.String("adId", adId), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetPlaceById", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetPlaceById", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetPlaceById").Observe(duration)
	}()
	var ad domain.GetAllAdsResponse

	query := r.db.Model(&domain.Ad{}).Joins("JOIN users ON ads.\"authorUUID\" = users.uuid").
		Joins("JOIN cities ON ads.\"cityId\" = cities.id").
		Joins("JOIN ad_available_dates ON ad_available_dates.\"adId\" = ads.uuid").
		Select("ads.*, cities.title as \"CityName\", ad_available_dates.\"availableDateFrom\" as \"AdDateFrom\", ad_available_dates.\"availableDateTo\" as \"AdDateTo\"").
		Where("\"adId\" = ?", adId)

	if err := query.Find(&ad).Error; err != nil {
		logger.DBLogger.Error("Error fetching place", zap.String("request_id", requestID), zap.Error(err))
		return ad, errors.New("error fetching place")
	}

	var images []domain.Image
	var user domain.User
	err = r.db.Model(&domain.Image{}).Where("\"adId\" = ?", ad.UUID).Find(&images).Error
	if err != nil {
		logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
		return ad, errors.New("error fetching images for ad")
	}

	err = r.db.Model(&domain.User{}).Where("uuid = ?", ad.AuthorUUID).Find(&user).Error
	if err != nil {
		logger.DBLogger.Error("Error fetching user", zap.String("request_id", requestID), zap.Error(err))
		return ad, errors.New("error fetching user")
	}

	ad.AdAuthor.Name = user.Name
	ad.AdAuthor.Avatar = user.Avatar
	ad.AdAuthor.Rating = user.Score
	ad.AdAuthor.GuestCount = user.GuestCount
	ad.AdAuthor.Sex = user.Sex
	ad.AdAuthor.Birthdate = user.Birthdate

	for _, img := range images {
		ad.Images = append(ad.Images, domain.ImageResponse{
			ID:        img.ID,
			ImagePath: img.ImageUrl,
		})
	}

	if err != nil {
		logger.DBLogger.Error("Error fetching place views", zap.String("request_id", requestID), zap.Error(err))
		return ad, errors.New("error fetching place views")
	}
	logger.DBLogger.Info("Successfully fetched place by ID", zap.String("adId", adId), zap.String("request_id", requestID))
	return ad, nil
}

func (r *adRepository) UpdateViewsCount(ctx context.Context, ad domain.GetAllAdsResponse) (domain.GetAllAdsResponse, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("UpdateViewsCount called", zap.String("adId", ad.UUID), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetPlaceById", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetPlaceById", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetPlaceById").Observe(duration)
	}()
	ad.ViewsCount += 1
	err = r.db.Model(&domain.Ad{}).Where("uuid = ?", ad.UUID).Updates(&ad).Error
	if err != nil {
		logger.DBLogger.Error("Error updating views count", zap.String("request_id", requestID), zap.Error(err))
		return ad, errors.New("error updating views count")
	}
	logger.DBLogger.Info("Successfully updated views count", zap.String("request_id", requestID))
	return ad, nil
}

func (r *adRepository) CreatePlace(ctx context.Context, ad *domain.Ad, newAd domain.CreateAdRequest, userId string) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("CreatePlace called", zap.String("adId", ad.UUID), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("CreatePlace", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("CreatePlace", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("CreatePlace").Observe(duration)
	}()
	var city domain.City
	var user domain.User
	var date domain.AdAvailableDate

	if err := r.db.Where("uuid = ?", userId).First(&user).Error; err != nil {
		logger.DBLogger.Error("Error finding user", zap.String("userId", userId), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error finding user")
	}

	if !user.IsHost {
		logger.DBLogger.Error("User is not host", zap.String("userId", userId), zap.String("request_id", requestID))
		return errors.New("user is not host")
	}

	if err := r.db.Where("title = ?", newAd.CityName).First(&city).Error; err != nil {
		logger.DBLogger.Error("Error creating place", zap.String("adId", ad.UUID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error finding city")
	}
	ad.CityID = city.ID
	ad.AuthorUUID = userId
	ad.PublicationDate = time.Now()
	if err := r.db.Create(ad).Error; err != nil {
		logger.DBLogger.Error("Error creating place", zap.String("adId", ad.UUID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error creating place")
	}

	date.AdID = ad.UUID
	date.AvailableDateFrom = newAd.DateFrom
	date.AvailableDateTo = newAd.DateTo

	if err := r.db.Create(&date).Error; err != nil {
		logger.DBLogger.Error("Error creating date", zap.String("adId", ad.UUID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error creating date")
	}
	logger.DBLogger.Info("Successfully place", zap.String("adId", ad.UUID), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) SavePlace(ctx context.Context, ad *domain.Ad) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("SavePlace called", zap.String("adId", ad.UUID), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("SavePlace", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("SavePlace", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("SavePlace").Observe(duration)
	}()
	if err := r.db.Save(ad).Error; err != nil {
		logger.DBLogger.Error("Error saving place", zap.String("adId", ad.UUID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error saving place")
	}
	logger.DBLogger.Info("Successfully saved place", zap.String("adId", ad.UUID), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) UpdatePlace(ctx context.Context, ad *domain.Ad, adId string, userId string, updatedPlace domain.UpdateAdRequest) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("UpdatePlace called", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("UpdatePlace", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("UpdatePlace", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("UpdatePlace").Observe(duration)
	}()
	var oldAd domain.Ad
	var oldDate domain.AdAvailableDate
	if err := r.db.Where("uuid = ?", adId).First(&oldAd).Error; err != nil {
		logger.DBLogger.Error("Ad not found", zap.String("adId", adId), zap.String("request_id", requestID))
		return errors.New("ad not found")
	}

	if err := r.db.Where("\"adId\" = ?", adId).First(&oldDate).Error; err != nil {
		logger.DBLogger.Error("Ad date not found", zap.String("adId", adId), zap.String("request_id", requestID))
		return errors.New("ad date not found")
	}

	if oldAd.AuthorUUID != userId {
		logger.DBLogger.Warn("User is not the owner of the ad", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))
		return errors.New("not owner of ad")
	}
	var city domain.City
	if err := r.db.Where("title = ?", updatedPlace.CityName).First(&city).Error; err != nil {
		logger.DBLogger.Error("Error creating place", zap.String("adId", ad.UUID), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error finding city")
	}
	ad.CityID = city.ID
	if err := r.db.Model(&oldAd).Updates(ad).Error; err != nil {
		logger.DBLogger.Error("Error updating place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error updating place")
	}

	oldDate.AvailableDateFrom = updatedPlace.DateFrom
	oldDate.AvailableDateTo = updatedPlace.DateTo

	if err := r.db.Model(&oldDate).Updates(oldDate).Error; err != nil {
		logger.DBLogger.Error("Error updating date", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error updating date")
	}

	logger.DBLogger.Info("Successfully updated place", zap.String("adId", adId), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) DeletePlace(ctx context.Context, adId string, userId string) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("DeletePlace called", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("DeletePlace", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("DeletePlace", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("DeletePlace").Observe(duration)
	}()
	var ad domain.Ad
	if err := r.db.Where("uuid = ?", adId).First(&ad).Error; err != nil {
		logger.DBLogger.Error("Ad not found", zap.String("adId", adId), zap.String("request_id", requestID))
		return errors.New("ad not found")
	}

	if ad.AuthorUUID != userId {
		logger.DBLogger.Warn("User is not the owner of the ad", zap.String("adId", adId), zap.String("userId", userId), zap.String("request_id", requestID))
		return errors.New("not owner of ad")
	}

	if err := r.db.Where("\"adId\" = ?", adId).Delete(&domain.Image{}).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error deleting place")
	}

	if err := r.db.Where("\"adId\" = ?", adId).Delete(&domain.AdPosition{}).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error deleting place")
	}

	if err := r.db.Where("\"adId\" = ?", adId).Delete(&domain.AdAvailableDate{}).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error deleting place")
	}

	if err := r.db.Where("\"adId\" = ?", adId).Delete(&domain.Request{}).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error deleting place")
	}

	if err := r.db.Delete(&ad).Error; err != nil {
		logger.DBLogger.Error("Error deleting place", zap.String("adId", adId), zap.String("request_id", requestID), zap.Error(err))
		return errors.New("error deleting place")
	}

	logger.DBLogger.Info("Successfully deleted place", zap.String("adId", adId), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) GetPlacesPerCity(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetPlacesPerCity called", zap.String("city", city), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetPlacesPerCity", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetPlacesPerCity", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetPlacesPerCity").Observe(duration)
	}()
	var ads []domain.GetAllAdsResponse
	query := r.db.Model(&domain.Ad{}).Joins("JOIN users ON ads.\"authorUUID\" = users.uuid").Joins("JOIN cities ON  ads.\"cityId\" = cities.id").
		Select("ads.*, cities.title as \"CityName\"").Where("cities.\"enTitle\" = ?", city)
	if err := query.Find(&ads).Error; err != nil {
		logger.DBLogger.Error("Error fetching places per city", zap.String("city", city), zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("error fetching places per city")
	}

	for i, ad := range ads {
		var images []domain.Image
		var user domain.User
		err := r.db.Model(&domain.Image{}).Where("\"adId\" = ?", ad.UUID).Find(&images).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("error fetching images for ad")
		}

		err = r.db.Model(&domain.User{}).Where("uuid = ?", ad.AuthorUUID).Find(&user).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching user", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("error fetching user")
		}
		ads[i].AdAuthor.Name = user.Name
		ads[i].AdAuthor.Avatar = user.Avatar
		ads[i].AdAuthor.Rating = user.Score
		ads[i].AdAuthor.GuestCount = user.GuestCount
		ads[i].AdAuthor.Sex = user.Sex
		ads[i].AdAuthor.Birthdate = user.Birthdate

		for _, img := range images {
			ads[i].Images = append(ads[i].Images, domain.ImageResponse{
				ID:        img.ID,
				ImagePath: img.ImageUrl,
			})
		}
	}

	logger.DBLogger.Info("Successfully fetched places per city", zap.String("city", city), zap.Int("count", len(ads)), zap.String("request_id", requestID))
	return ads, nil
}

func (r *adRepository) SaveImages(ctx context.Context, adUUID string, imagePaths []string) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetPlacesPerCity", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetPlacesPerCity", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetPlacesPerCity").Observe(duration)
	}()
	for _, path := range imagePaths {
		image := domain.Image{
			AdID:     adUUID,
			ImageUrl: path,
		}
		if err := r.db.Create(&image).Error; err != nil {
			logger.DBLogger.Error("Error creating images", zap.String("request_id", requestID), zap.Error(err))
			return errors.New("error creating image")
		}
	}
	return nil
}

func (r *adRepository) GetAdImages(ctx context.Context, adId string) ([]string, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetAdImages called", zap.String("request_id", requestID), zap.String("adId", adId))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetAdImages", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetAdImages", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetAdImages").Observe(duration)
	}()
	var imageUrls []string

	err = r.db.Model(&domain.Image{}).Where("\"adId\" = ?", adId).Pluck("imageUrl", &imageUrls).Error
	if err != nil {
		logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("error fetching images for ad")
	}

	logger.DBLogger.Info("Successfully fetched images for ad", zap.String("request_id", requestID), zap.Int("count", len(imageUrls)))
	return imageUrls, nil
}

func (r *adRepository) GetUserPlaces(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetUserPlaces called", zap.String("city", userId), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetUserPlaces", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetUserPlaces", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetUserPlaces").Observe(duration)
	}()
	var ads []domain.GetAllAdsResponse
	query := r.db.Model(&domain.Ad{}).Joins("JOIN users ON ads.\"authorUUID\" = users.uuid").Joins("JOIN cities ON  ads.\"cityId\" = cities.id").
		Select("ads.*, users.avatar, users.name, users.score as rating, cities.title as \"CityName\"").Where("users.uuid = ?", userId)
	if err := query.Find(&ads).Error; err != nil {
		logger.DBLogger.Error("Error fetching user places", zap.String("city", userId), zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("error fetching user places")
	}

	for i, ad := range ads {
		var images []domain.Image
		err := r.db.Model(&domain.Image{}).Where("\"adId\" = ?", ad.UUID).Find(&images).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("error fetching images for ad")
		}

		for _, img := range images {
			ads[i].Images = append(ads[i].Images, domain.ImageResponse{
				ID:        img.ID,
				ImagePath: img.ImageUrl,
			})
		}
	}

	logger.DBLogger.Info("Successfully fetched user places", zap.String("city", userId), zap.Int("count", len(ads)), zap.String("request_id", requestID))
	return ads, nil
}

func (r *adRepository) DeleteAdImage(ctx context.Context, adId string, imageId int, userId string) (string, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("DeleteAdImage called", zap.String("ad", adId), zap.Int("image", imageId), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("DeleteAdImage", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("DeleteAdImage", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("DeleteAdImage").Observe(duration)
	}()
	var ad domain.Ad
	if err := r.db.First(&ad, "uuid = ?", adId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("ad not found")
		}
		return "", errors.New("error fetching ad")
	}

	if ad.AuthorUUID != userId {
		return "", errors.New("not owner of ad")
	}

	var image domain.Image
	if err := r.db.First(&image, "id = ? AND \"adId\" = ?", imageId, adId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("image not found")
		}
		return "", errors.New("error finding image")
	}

	if err := r.db.Delete(&image).Error; err != nil {
		return "", errors.New("error deleting image from database")
	}

	logger.DBLogger.Info("Image deleted successfully", zap.Int("image_id", imageId), zap.String("ad_id", adId), zap.String("request_id", requestID))
	return image.ImageUrl, nil
}

func (r *adRepository) AddToFavorites(ctx context.Context, adId string, userId string) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("AddToFavorites called", zap.String("ad", adId), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("AddToFavorites", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("AddToFavorites", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("AddToFavorites").Observe(duration)
	}()
	var ad domain.Ad
	if err := r.db.First(&ad, "uuid = ?", adId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ad not found")
		}
		return errors.New("error fetching ad")
	}

	var favorite domain.Favorites
	favorite.AdId = adId
	favorite.UserId = userId
	if err := r.db.Create(&favorite).Error; err != nil {
		return errors.New("error create favorite")
	}

	logger.DBLogger.Info("Favorite create successfully", zap.String("ad_id", adId), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) DeleteFromFavorites(ctx context.Context, adId string, userId string) error {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("AddToFavorites called", zap.String("ad", adId), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("AddToFavorites", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("AddToFavorites", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("AddToFavorites").Observe(duration)
	}()
	var ad domain.Ad
	if err := r.db.First(&ad, "uuid = ?", adId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ad not found")
		}
		return errors.New("error fetching ad")
	}

	var favorite domain.Favorites
	favorite.AdId = adId
	favorite.UserId = userId
	if err := r.db.Delete(&favorite).Error; err != nil {
		return errors.New("error create favorite")
	}

	logger.DBLogger.Info("Favorite create successfully", zap.String("ad_id", adId), zap.String("request_id", requestID))
	return nil
}

func (r *adRepository) GetUserFavorites(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
	start := time.Now()
	requestID := middleware.GetRequestID(ctx)
	logger.DBLogger.Info("GetUserFavorites called", zap.String("user", userId), zap.String("request_id", requestID))
	var err error
	defer func() {
		if err != nil {
			metrics.RepoErrorsTotal.WithLabelValues("GetUserFavorites", "error", err.Error()).Inc()
		} else {
			metrics.RepoRequestTotal.WithLabelValues("GetUserFavorites", "success").Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.RepoRequestDuration.WithLabelValues("GetUserFavorites").Observe(duration)
	}()
	var ads []domain.GetAllAdsResponse

	query := r.db.Model(&domain.Ad{}).
		Joins("JOIN favorites ON favorites.\"adId\" = ads.uuid").
		Where("favorites.\"userId\" = ?", userId).
		Select("ads.*, favorites.\"userId\" AS \"FavoriteUserId\"")

	if err := query.Find(&ads).Error; err != nil {
		logger.DBLogger.Error("Error fetching user favorites", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("error fetching user favorites")
	}

	for i, ad := range ads {
		var images []domain.Image
		var user domain.User
		err := r.db.Model(&domain.Image{}).Where("\"adId\" = ?", ad.UUID).Find(&images).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching images for ad", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("error fetching images for ad")
		}

		err = r.db.Model(&domain.User{}).Where("uuid = ?", ad.AuthorUUID).Find(&user).Error
		if err != nil {
			logger.DBLogger.Error("Error fetching user", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("error fetching user")
		}
		ads[i].AdAuthor.Name = user.Name
		ads[i].AdAuthor.Avatar = user.Avatar
		ads[i].AdAuthor.Rating = user.Score
		ads[i].AdAuthor.GuestCount = user.GuestCount
		ads[i].AdAuthor.Sex = user.Sex
		ads[i].AdAuthor.Birthdate = user.Birthdate
		for _, img := range images {
			ads[i].Images = append(ads[i].Images, domain.ImageResponse{
				ID:        img.ID,
				ImagePath: img.ImageUrl,
			})
		}
	}

	logger.DBLogger.Info("Successfully fetched user favorites", zap.String("request_id", requestID), zap.Int("count", len(ads)))

	return ads, nil
}
