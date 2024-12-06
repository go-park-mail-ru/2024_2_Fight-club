package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	"2024_2_FIGHT-CLUB/microservices/ads_service/usecase"
	"context"
	"errors"
	"github.com/microcosm-cc/bluemonday"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type GrpcAdHandler struct {
	gen.AdsServer
	sessionService session.InterfaceSession
	usecase        usecase.AdUseCase
	jwtToken       middleware.JwtTokenService
}

func NewGrpcAdHandler(sessionService session.InterfaceSession, usecase usecase.AdUseCase, jwtToken middleware.JwtTokenService) *GrpcAdHandler {
	return &GrpcAdHandler{
		sessionService: sessionService,
		usecase:        usecase,
		jwtToken:       jwtToken,
	}
}

func (adh *GrpcAdHandler) GetAllPlaces(ctx context.Context, in *gen.AdFilterRequest) (*gen.GetAllAdsResponseList, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received GetAllPlaces request in microservice",
		zap.String("request_id", requestID),
	)

	layout := "2006-01-02"
	var dateTo time.Time
	var dateFrom time.Time

	location := sanitizer.Sanitize(in.Location)
	rating := sanitizer.Sanitize(in.Rating)
	newThisWeek := sanitizer.Sanitize(in.NewThisWeek)
	hostGender := sanitizer.Sanitize(in.HostGender)
	guestCounter := sanitizer.Sanitize(in.GuestCount)

	offset := sanitizer.Sanitize(in.Offset)
	var offsetInt int
	if offset != "" {
		var err error
		offsetInt, err = strconv.Atoi(offset)
		if err != nil {
			logger.AccessLogger.Error("Failed to parse offset as int", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("query offset not int")
		}
	}

	limit := sanitizer.Sanitize(in.Limit)
	var limitInt int
	if offset != "" {
		var err error
		limitInt, err = strconv.Atoi(limit)
		if err != nil {
			logger.AccessLogger.Error("Failed to parse limit as int", zap.String("request_id", requestID), zap.Error(err))
			return nil, errors.New("query limit not int")
		}
	}

	dateFromStr := sanitizer.Sanitize(in.DateFrom)
	if dateFromStr != "" {
		var err error
		dateFrom, err = time.Parse(layout, dateFromStr)
		if err != nil {
			logger.AccessLogger.Error("Failed to parse dateFrom",
				zap.Error(err),
				zap.String("request_id", requestID))
			return nil, errors.New("query dateFrom not int")
		}
	}

	dateToStr := sanitizer.Sanitize(in.DateTo)

	if dateToStr != "" {
		var err error
		dateTo, err = time.Parse(layout, dateToStr)
		if err != nil {
			logger.AccessLogger.Error("Failed to parse dateTo",
				zap.Error(err),
				zap.String("request_id", requestID))
			return nil, errors.New("query dateTo not int")
		}
	}

	filter := domain.AdFilter{
		Location:    location,
		Rating:      rating,
		NewThisWeek: newThisWeek,
		HostGender:  hostGender,
		GuestCount:  guestCounter,
		Limit:       limitInt,
		Offset:      offsetInt,
		DateFrom:    dateFrom,
		DateTo:      dateTo,
	}

	places, err := adh.usecase.GetAllPlaces(ctx, filter)
	if err != nil {
		logger.AccessLogger.Error("Failed to get places",
			zap.Error(err),
			zap.String("request_id", requestID))
		return nil, err
	}
	var responseList gen.GetAllAdsResponseList
	for _, place := range places {
		ad := &gen.GetAllAdsResponse{
			Id:              place.UUID,
			CityId:          int32(place.CityID),
			AuthorUUID:      place.AuthorUUID,
			Address:         place.Address,
			PublicationDate: place.PublicationDate.Format(layout),
			Description:     place.Description,
			RoomsNumber:     int32(place.RoomsNumber),
			ViewsCount:      int32Ptr(int32(place.ViewsCount)),
			CityName:        place.CityName,
			AdDateFrom:      place.AdDateFrom.Format(layout),
			AdDateTo:        place.AdDateTo.Format(layout),
			AdAuthor: &gen.UserResponse{
				Rating:     float32Ptr(float32(place.AdAuthor.Rating)),
				Avatar:     place.AdAuthor.Avatar,
				Name:       place.AdAuthor.Name,
				GuestCount: int32Ptr(int32(place.AdAuthor.GuestCount)),
				Sex:        place.AdAuthor.Sex,
				BirthDate:  place.AdAuthor.Birthdate.Format(layout),
			},
			Images: convertImagesToGRPC(place.Images),
		}
		responseList.Housing = append(responseList.Housing, ad)
	}

	logger.AccessLogger.Info("Successfully fetched all places", zap.String("request_id", requestID), zap.Int("count", len(places)))
	return &responseList, nil
}

func (adh *GrpcAdHandler) GetOnePlace(ctx context.Context, in *gen.GetPlaceByIdRequest) (*gen.GetAllAdsResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received GetAllPlaces request in microservice",
		zap.String("request_id", requestID),
	)

	layout := "2006-01-02"

	in.AdId = sanitizer.Sanitize(in.AdId)

	place, err := adh.usecase.GetOnePlace(ctx, in.AdId, in.IsAuthorized)
	if err != nil {
		logger.AccessLogger.Error("Failed to get places",
			zap.Error(err),
			zap.String("request_id", requestID))
		return nil, err
	}

	return &gen.GetAllAdsResponse{
		Id:              place.UUID,
		CityId:          int32(place.CityID),
		AuthorUUID:      place.AuthorUUID,
		Address:         place.Address,
		PublicationDate: place.PublicationDate.Format(layout),
		Description:     place.Description,
		RoomsNumber:     int32(place.RoomsNumber),
		ViewsCount:      int32Ptr(int32(place.ViewsCount)),
		CityName:        place.CityName,
		AdDateFrom:      place.AdDateFrom.Format(layout),
		AdDateTo:        place.AdDateTo.Format(layout),
		AdAuthor: &gen.UserResponse{
			Rating:     float32Ptr(float32(place.AdAuthor.Rating)),
			Avatar:     place.AdAuthor.Avatar,
			Name:       place.AdAuthor.Name,
			GuestCount: int32Ptr(int32(place.AdAuthor.GuestCount)),
			Sex:        place.AdAuthor.Sex,
			BirthDate:  place.AdAuthor.Birthdate.Format(layout),
		},
		Images: convertImagesToGRPC(place.Images),
	}, nil
}

func (adh *GrpcAdHandler) CreatePlace(ctx context.Context, in *gen.CreateAdRequest) (*gen.Ad, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	layout := "2006-01-02"
	logger.AccessLogger.Info("Received CreatePlace request in microservice",
		zap.String("request_id", requestID),
	)

	if in.AuthHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("Missing X-CSRF-Token header")),
		)
		return nil, errors.New("missing X-CSRF-Token header")
	}

	in.CityName = sanitizer.Sanitize(in.CityName)
	in.Description = sanitizer.Sanitize(in.Description)
	in.Address = sanitizer.Sanitize(in.Address)

	tokenString := in.AuthHeader[len("Bearer "):]
	_, err := adh.jwtToken.Validate(tokenString, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("invalid JWT token")
	}

	userID, err := adh.sessionService.GetUserID(ctx, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		return nil, errors.New("no active session")
	}

	var place domain.Ad
	newPlace := domain.CreateAdRequest{
		CityName:    in.CityName,
		Description: in.Description,
		Address:     in.Address,
		RoomsNumber: int(in.RoomsNumber),
		DateFrom:    (in.DateFrom).AsTime(),
		DateTo:      (in.DateTo).AsTime(),
	}
	place.AuthorUUID = userID

	err = adh.usecase.CreatePlace(ctx, &place, in.Images, newPlace, userID)
	if err != nil {
		logger.AccessLogger.Warn("Failed to create place", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}

	return &gen.Ad{
		Uuid:            place.UUID,
		CityId:          int32(place.CityID),
		AuthorUUID:      place.AuthorUUID,
		Address:         place.Address,
		Description:     place.Description,
		RoomsNumber:     int32(place.RoomsNumber),
		ViewsCount:      int32(place.ViewsCount),
		PublicationDate: (place.PublicationDate).Format(layout),
	}, nil
}

func (adh *GrpcAdHandler) UpdatePlace(ctx context.Context, in *gen.UpdateAdRequest) (*gen.AdResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received UpdatePlace request in microservice",
		zap.String("request_id", requestID),
	)

	in.AdId = sanitizer.Sanitize(in.AdId)
	in.Description = sanitizer.Sanitize(in.Description)
	in.Address = sanitizer.Sanitize(in.Address)
	in.CityName = sanitizer.Sanitize(in.CityName)

	if in.AuthHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		return nil, errors.New("missing X-CSRF-Token header")
	}

	tokenString := in.AuthHeader[len("Bearer "):]
	_, err := adh.jwtToken.Validate(tokenString, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("invalid JWT token")
	}

	userID, err := adh.sessionService.GetUserID(ctx, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("no active session")
	}
	updatedPlace := domain.UpdateAdRequest{
		CityName:    in.CityName,
		Description: in.Description,
		Address:     in.Address,
		RoomsNumber: int(in.RoomsNumber),
		DateFrom:    (in.DateFrom).AsTime(),
		DateTo:      (in.DateTo).AsTime(),
	}
	var place domain.Ad
	err = adh.usecase.UpdatePlace(ctx, &place, in.AdId, userID, in.Images, updatedPlace)
	if err != nil {
		logger.AccessLogger.Warn("Failed to update place", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}
	return &gen.AdResponse{Response: "Update successfully"}, nil
}

func (adh *GrpcAdHandler) DeletePlace(ctx context.Context, in *gen.DeletePlaceRequest) (*gen.DeleteResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received DeletePlace request in microservice",
		zap.String("request_id", requestID))

	in.AdId = sanitizer.Sanitize(in.AdId)

	if in.AuthHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		return nil, errors.New("missing X-CSRF-Token header")
	}

	tokenString := in.AuthHeader[len("Bearer "):]
	_, err := adh.jwtToken.Validate(tokenString, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("invalid JWT token")
	}

	userID, err := adh.sessionService.GetUserID(ctx, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		return nil, errors.New("no active session")
	}

	err = adh.usecase.DeletePlace(ctx, in.AdId, userID)
	if err != nil {
		logger.AccessLogger.Warn("Failed to delete place", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}
	return &gen.DeleteResponse{Response: "Delete successfully"}, nil
}

func (adh *GrpcAdHandler) GetPlacesPerCity(ctx context.Context, in *gen.GetPlacesPerCityRequest) (*gen.GetAllAdsResponseList, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	layout := "2006-01-02"
	logger.AccessLogger.Info("Received GetPlacesPerCity request in microservice",
		zap.String("request_id", requestID))

	in.CityName = sanitizer.Sanitize(in.CityName)

	places, err := adh.usecase.GetPlacesPerCity(ctx, in.CityName)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get places per city", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}
	var responseList gen.GetAllAdsResponseList
	for _, place := range places {
		ad := &gen.GetAllAdsResponse{
			Id:              place.UUID,
			CityId:          int32(place.CityID),
			AuthorUUID:      place.AuthorUUID,
			Address:         place.Address,
			PublicationDate: place.PublicationDate.Format(layout),
			Description:     place.Description,
			RoomsNumber:     int32(place.RoomsNumber),
			ViewsCount:      int32Ptr(int32(place.ViewsCount)),
			CityName:        place.CityName,
			AdDateFrom:      place.AdDateFrom.Format(layout),
			AdDateTo:        place.AdDateTo.Format(layout),
			AdAuthor: &gen.UserResponse{
				Rating:     float32Ptr(float32(place.AdAuthor.Rating)),
				Avatar:     place.AdAuthor.Avatar,
				Name:       place.AdAuthor.Name,
				GuestCount: int32Ptr(int32(place.AdAuthor.GuestCount)),
				Sex:        place.AdAuthor.Sex,
				BirthDate:  place.AdAuthor.Birthdate.Format(layout),
			},
			Images: convertImagesToGRPC(place.Images),
		}
		responseList.Housing = append(responseList.Housing, ad)
	}
	return &responseList, nil
}

func (adh *GrpcAdHandler) GetUserPlaces(ctx context.Context, in *gen.GetUserPlacesRequest) (*gen.GetAllAdsResponseList, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	layout := "2006-01-02"
	logger.AccessLogger.Info("Received GetUserPlaces request in microservice",
		zap.String("request_id", requestID))

	in.UserId = sanitizer.Sanitize(in.UserId)
	places, err := adh.usecase.GetUserPlaces(ctx, in.UserId)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user places", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}
	var responseList gen.GetAllAdsResponseList
	for _, place := range places {
		ad := &gen.GetAllAdsResponse{
			Id:              place.UUID,
			CityId:          int32(place.CityID),
			AuthorUUID:      place.AuthorUUID,
			Address:         place.Address,
			PublicationDate: place.PublicationDate.Format(layout),
			Description:     place.Description,
			RoomsNumber:     int32(place.RoomsNumber),
			ViewsCount:      int32Ptr(int32(place.ViewsCount)),
			CityName:        place.CityName,
			AdDateFrom:      place.AdDateFrom.Format(layout),
			AdDateTo:        place.AdDateTo.Format(layout),
			AdAuthor: &gen.UserResponse{
				Rating:     float32Ptr(float32(place.AdAuthor.Rating)),
				Avatar:     place.AdAuthor.Avatar,
				Name:       place.AdAuthor.Name,
				GuestCount: int32Ptr(int32(place.AdAuthor.GuestCount)),
				Sex:        place.AdAuthor.Sex,
				BirthDate:  place.AdAuthor.Birthdate.Format(layout),
			},
			Images: convertImagesToGRPC(place.Images),
		}
		responseList.Housing = append(responseList.Housing, ad)
	}
	return &responseList, nil
}

func (adh *GrpcAdHandler) DeleteAdImage(ctx context.Context, in *gen.DeleteAdImageRequest) (*gen.DeleteResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received DeleteAdImage request in microservice",
		zap.String("request_id", requestID),
	)

	in.AdId = sanitizer.Sanitize(in.AdId)
	in.ImageId = sanitizer.Sanitize(in.ImageId)

	if in.AuthHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		return nil, errors.New("missing X-CSRF-Token header")
	}

	tokenString := in.AuthHeader[len("Bearer "):]
	_, err := adh.jwtToken.Validate(tokenString, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("invalid JWT token")
	}

	userID, err := adh.sessionService.GetUserID(ctx, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		return nil, errors.New("no active session")
	}

	err = adh.usecase.DeleteAdImage(ctx, in.AdId, in.ImageId, userID)
	if err != nil {
		logger.AccessLogger.Warn("Failed to delete ad image", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}
	return &gen.DeleteResponse{Response: "Delete image successfully"}, nil
}

func (adh *GrpcAdHandler) AddToFavorites(ctx context.Context, in *gen.AddToFavoritesRequest) (*gen.AdResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received AddToFavorites request in microservice",
		zap.String("request_id", requestID),
	)

	in.AdId = sanitizer.Sanitize(in.AdId)

	if in.AuthHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		return nil, errors.New("missing X-CSRF-Token header")
	}

	tokenString := in.AuthHeader[len("Bearer "):]
	_, err := adh.jwtToken.Validate(tokenString, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("invalid JWT token")
	}

	userID, err := adh.sessionService.GetUserID(ctx, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		return nil, errors.New("no active session")
	}

	err = adh.usecase.AddToFavorites(ctx, in.AdId, userID)
	if err != nil {
		logger.AccessLogger.Warn("Failed to add ad to favorites", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}
	return &gen.AdResponse{Response: "Add to favorites successfully"}, nil
}

func (adh *GrpcAdHandler) DeleteFromFavorites(ctx context.Context, in *gen.DeleteFromFavoritesRequest) (*gen.AdResponse, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received DeleteFromFavorites request in microservice",
		zap.String("request_id", requestID),
	)

	in.AdId = sanitizer.Sanitize(in.AdId)

	if in.AuthHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		return nil, errors.New("missing X-CSRF-Token header")
	}

	tokenString := in.AuthHeader[len("Bearer "):]
	_, err := adh.jwtToken.Validate(tokenString, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("invalid JWT token")
	}

	userID, err := adh.sessionService.GetUserID(ctx, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		return nil, errors.New("no active session")
	}

	err = adh.usecase.DeleteFromFavorites(ctx, in.AdId, userID)
	if err != nil {
		logger.AccessLogger.Warn("Failed to delete ad from favorites", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}
	return &gen.AdResponse{Response: "Delete ad from favorites successfully"}, nil
}

func (adh *GrpcAdHandler) GetUserFavorites(ctx context.Context, in *gen.GetUserFavoritesRequest) (*gen.GetAllAdsResponseList, error) {
	requestID := middleware.GetRequestID(ctx)
	sanitizer := bluemonday.UGCPolicy()
	logger.AccessLogger.Info("Received DeleteFromFavorites request in microservice",
		zap.String("request_id", requestID),
	)
	layout := "2006-01-02"
	in.UserId = sanitizer.Sanitize(in.UserId)

	if in.AuthHeader == "" {
		logger.AccessLogger.Warn("Missing X-CSRF-Token header",
			zap.String("request_id", requestID),
			zap.Error(errors.New("missing X-CSRF-Token header")),
		)
		return nil, errors.New("missing X-CSRF-Token header")
	}

	tokenString := in.AuthHeader[len("Bearer "):]
	_, err := adh.jwtToken.Validate(tokenString, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("Invalid JWT token", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("invalid JWT token")
	}

	userID, err := adh.sessionService.GetUserID(ctx, in.SessionID)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		return nil, errors.New("no active session")
	}
	if userID != in.UserId {
		logger.AccessLogger.Warn("cant access other user favorites", zap.String("request_id", requestID), zap.Error(err))
		return nil, errors.New("cant access other user favorites")
	}
	places, err := adh.usecase.GetUserFavorites(ctx, in.UserId)
	if err != nil {
		logger.AccessLogger.Warn("Failed to get user favorites", zap.String("request_id", requestID), zap.Error(err))
		return nil, err
	}
	var responseList gen.GetAllAdsResponseList
	for _, place := range places {
		ad := &gen.GetAllAdsResponse{
			Id:              place.UUID,
			CityId:          int32(place.CityID),
			AuthorUUID:      place.AuthorUUID,
			Address:         place.Address,
			PublicationDate: place.PublicationDate.Format(layout),
			Description:     place.Description,
			RoomsNumber:     int32(place.RoomsNumber),
			ViewsCount:      int32Ptr(int32(place.ViewsCount)),
			CityName:        place.CityName,
			AdDateFrom:      place.AdDateFrom.Format(layout),
			AdDateTo:        place.AdDateTo.Format(layout),
			AdAuthor: &gen.UserResponse{
				Rating:     float32Ptr(float32(place.AdAuthor.Rating)),
				Avatar:     place.AdAuthor.Avatar,
				Name:       place.AdAuthor.Name,
				GuestCount: int32Ptr(int32(place.AdAuthor.GuestCount)),
				Sex:        place.AdAuthor.Sex,
				BirthDate:  place.AdAuthor.Birthdate.Format(layout),
			},
			Images: convertImagesToGRPC(place.Images),
		}
		responseList.Housing = append(responseList.Housing, ad)
	}

	logger.AccessLogger.Info("Successfully fetched all favorites", zap.String("request_id", requestID), zap.Int("count", len(places)))
	return &responseList, nil
}

func float32Ptr(f float32) *float32 {
	return &f
}

func int32Ptr(i int32) *int32 {
	return &i
}

func convertImagesToGRPC(images []domain.ImageResponse) []*gen.ImageResponse {
	var grpcImages []*gen.ImageResponse
	for _, img := range images {
		grpcImages = append(grpcImages, &gen.ImageResponse{
			Id:   int32(img.ID),
			Path: img.ImagePath,
		})
	}
	return grpcImages
}
