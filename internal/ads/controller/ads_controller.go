package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/metrics"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"2024_2_FIGHT-CLUB/microservices/ads_service/controller/gen"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"net/http"
	"time"
)

type AdHandler struct {
	client         gen.AdsClient
	sessionService session.InterfaceSession
	jwtToken       middleware.JwtTokenService
}

func NewAdHandler(client gen.AdsClient, sessionService session.InterfaceSession, jwtToken middleware.JwtTokenService) *AdHandler {
	return &AdHandler{
		client:         client,
		sessionService: sessionService,
		jwtToken:       jwtToken,
	}
}

func (h *AdHandler) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error

	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received GetAllPlaces request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	queryParams := r.URL.Query()

	response, err := h.client.GetAllPlaces(ctx, &gen.AdFilterRequest{
		Location:    queryParams.Get("location"),
		Rating:      queryParams.Get("rating"),
		NewThisWeek: queryParams.Get("new"),
		HostGender:  queryParams.Get("gender"),
		GuestCount:  queryParams.Get("guests"),
		Limit:       queryParams.Get("limit"),
		Offset:      queryParams.Get("offset"),
		DateFrom:    queryParams.Get("dateFrom"),
		DateTo:      queryParams.Get("dateTo"),
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to GetAllPlaces",
			zap.Error(err),
			zap.String("request_id", requestID),
			zap.String("method", r.Method))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetAllPlaces request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) GetOnePlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()
	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received GetOnePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	var isAuthorized bool

	sessionID, err := session.GetSessionId(r)
	if err != nil || sessionID == "" {
		logger.AccessLogger.Warn("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		isAuthorized = false
	}

	if _, err := h.sessionService.GetUserID(ctx, sessionID); err != nil {
		isAuthorized = false
	} else {
		isAuthorized = true
	}

	place, err := h.client.GetOnePlace(ctx, &gen.GetPlaceByIdRequest{
		AdId:         adId,
		IsAuthorized: isAuthorized,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to GetOnePlace",
			zap.String("request_id", requestID),
			zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	body := map[string]interface{}{
		"place": place,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetOnePlace request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) CreatePlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusCreated
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusCreated {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received CreatePlace request",
		zap.String("request_id", requestID),
	)

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	authHeader := r.Header.Get("X-CSRF-Token")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		logger.AccessLogger.Error("Failed to parse multipart form", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	metadata := r.FormValue("metadata")
	var newPlace domain.CreateAdRequest
	if err := json.Unmarshal([]byte(metadata), &newPlace); err != nil {
		logger.AccessLogger.Error("Failed to decode metadata", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, errors.New("failed to decode metadata"), requestID)
		return
	}

	fileHeaders := r.MultipartForm.File["images"]
	if len(fileHeaders) == 0 {
		logger.AccessLogger.Warn("No images", zap.String("request_id", requestID))
		statusCode = h.handleError(w, errors.New("no images provided"), requestID)
		return
	}

	// Преобразование файлов в [][]byte
	files := make([][]byte, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			logger.AccessLogger.Error("Failed to open file", zap.String("request_id", requestID), zap.Error(err))
			statusCode = h.handleError(w, errors.New("failed to open file"), requestID)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			logger.AccessLogger.Error("Failed to read file", zap.String("request_id", requestID), zap.Error(err))
			statusCode = h.handleError(w, errors.New("failed to read file"), requestID)
			return
		}

		files = append(files, data)
	}

	response, err := h.client.CreatePlace(ctx, &gen.CreateAdRequest{
		CityName:    newPlace.CityName,
		Description: newPlace.Description,
		Address:     newPlace.Address,
		RoomsNumber: int32(newPlace.RoomsNumber),
		DateFrom:    timestamppb.New(newPlace.DateFrom),
		DateTo:      timestamppb.New(newPlace.DateTo),
		Images:      files,
		AuthHeader:  authHeader,
		SessionID:   sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to create place", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	body := map[string]interface{}{
		"place": response,
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, errors.New("failed to encode response"), requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed CreatePlace request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) UpdatePlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received UpdatePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	authHeader := r.Header.Get("X-CSRF-Token")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err = r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		logger.AccessLogger.Error("Failed to parse multipart form", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, errors.New("invalid multipart form"), requestID)
		return
	}

	metadata := r.FormValue("metadata")
	var updatedPlace domain.UpdateAdRequest
	if err := json.Unmarshal([]byte(metadata), &updatedPlace); err != nil {
		logger.AccessLogger.Error("Failed to decode metadata", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, errors.New("invalid metadata JSON"), requestID)
		return
	}

	fileHeaders := r.MultipartForm.File["images"]

	// Преобразование `[]*multipart.FileHeader` в `[][]byte`
	files := make([][]byte, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			logger.AccessLogger.Error("Failed to open file", zap.String("request_id", requestID), zap.Error(err))
			statusCode = h.handleError(w, errors.New("failed to open file"), requestID)
			return
		}
		defer file.Close()

		// Чтение содержимого файла в []byte
		data, err := io.ReadAll(file)
		if err != nil {
			logger.AccessLogger.Error("Failed to read file", zap.String("request_id", requestID), zap.Error(err))
			statusCode = h.handleError(w, errors.New("failed to read file"), requestID)
			return
		}
		files = append(files, data)
	}

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	response, err := h.client.UpdatePlace(ctx, &gen.UpdateAdRequest{
		AdId:        adId,
		CityName:    updatedPlace.CityName,
		Address:     updatedPlace.Address,
		Description: updatedPlace.Description,
		RoomsNumber: int32(updatedPlace.RoomsNumber),
		SessionID:   sessionID,
		AuthHeader:  authHeader,
		Images:      files,
		DateFrom:    timestamppb.New(updatedPlace.DateFrom),
		DateTo:      timestamppb.New(updatedPlace.DateTo),
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to update place", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": response.Response}
	if err := json.NewEncoder(w).Encode(updateResponse); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed UpdatePlace request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) DeletePlace(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received DeletePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	authHeader := r.Header.Get("X-CSRF-Token")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	response, err := h.client.DeletePlace(ctx, &gen.DeletePlaceRequest{
		AdId:       adId,
		SessionID:  sessionID,
		AuthHeader: authHeader,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to delete place", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": response.Response}
	if err := json.NewEncoder(w).Encode(updateResponse); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *AdHandler) GetPlacesPerCity(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	city := mux.Vars(r)["city"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received GetPlacesPerCity request",
		zap.String("request_id", requestID),
		zap.String("city", city),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response, err := h.client.GetPlacesPerCity(ctx, &gen.GetPlacesPerCityRequest{
		CityName: city,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to get places per city", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}
	body := map[string]interface{}{
		"places": response,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetPlacesPerCity request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) GetUserPlaces(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	userId := mux.Vars(r)["userId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received GetUserPlaces request",
		zap.String("request_id", requestID),
		zap.String("userId", userId),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response, err := h.client.GetUserPlaces(ctx, &gen.GetUserPlacesRequest{
		UserId: userId,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to get places per user",
			zap.String("request_id", requestID),
			zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetUserPlaces request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) DeleteAdImage(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	imageId := mux.Vars(r)["imageId"]
	adId := mux.Vars(r)["adId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received DeleteAdImage request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
		zap.String("imageId", imageId))

	authHeader := r.Header.Get("X-CSRF-Token")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	response, err := h.client.DeleteAdImage(ctx, &gen.DeleteAdImageRequest{
		AdId:       adId,
		ImageId:    imageId,
		AuthHeader: authHeader,
		SessionID:  sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to delete ad image", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": response.Response}
	if err := json.NewEncoder(w).Encode(updateResponse); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed DeleteAdImage request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
		zap.String("imageId", imageId),
		zap.Duration("duration", duration),
	)
}

func (h *AdHandler) AddToFavorites(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received AddToFavorites request",
		zap.String("request_id", requestID),
		zap.String("adId", adId))

	authHeader := r.Header.Get("X-CSRF-Token")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	response, err := h.client.AddToFavorites(ctx, &gen.AddToFavoritesRequest{
		AdId:       adId,
		AuthHeader: authHeader,
		SessionID:  sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to add ad to favorites", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": response.Response}
	if err := json.NewEncoder(w).Encode(updateResponse); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed AddToFavorites request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
		zap.Duration("duration", duration),
	)

}

func (h *AdHandler) DeleteFromFavorites(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	adId := mux.Vars(r)["adId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received DeleteFromFavorites request",
		zap.String("request_id", requestID),
		zap.String("adId", adId))

	authHeader := r.Header.Get("X-CSRF-Token")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	response, err := h.client.DeleteFromFavorites(ctx, &gen.DeleteFromFavoritesRequest{
		AdId:       adId,
		AuthHeader: authHeader,
		SessionID:  sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to delete ad from favorites", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": response.Response}
	if err := json.NewEncoder(w).Encode(updateResponse); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed DeleteFromFavorites request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
		zap.Duration("duration", duration),
	)

}

func (h *AdHandler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	userId := mux.Vars(r)["userId"]
	ctx, cancel := middleware.WithTimeout(r.Context())
	defer cancel()

	statusCode := http.StatusOK
	var err error
	clientIP := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		clientIP = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	defer func() {
		if statusCode == http.StatusOK {
			metrics.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), clientIP).Inc()
		} else {
			metrics.HttpErrorsTotal.WithLabelValues(r.Method, r.URL.Path, http.StatusText(statusCode), err.Error(), clientIP).Inc()
		}
		duration := time.Since(start).Seconds()
		metrics.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, clientIP).Observe(duration)
	}()

	ctx = middleware.WithLogger(ctx, logger.AccessLogger)

	logger.AccessLogger.Info("Received GetUserFavorites request",
		zap.String("request_id", requestID),
		zap.String("userId", userId))

	authHeader := r.Header.Get("X-CSRF-Token")

	sessionID, err := session.GetSessionId(r)
	if err != nil {
		logger.AccessLogger.Error("Failed to get session ID",
			zap.String("request_id", requestID),
			zap.Error(err))
		statusCode = h.handleError(w, err, requestID)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	response, err := h.client.GetUserFavorites(ctx, &gen.GetUserFavoritesRequest{
		UserId:     userId,
		AuthHeader: authHeader,
		SessionID:  sessionID,
	})
	if err != nil {
		logger.AccessLogger.Error("Failed to delete ad from favorites", zap.String("request_id", requestID), zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			statusCode = h.handleError(w, errors.New(st.Message()), requestID)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start)
	logger.AccessLogger.Info("Completed GetUserFavorites request",
		zap.String("request_id", requestID),
		zap.String("userId", userId),
		zap.Duration("duration", duration),
	)

}

func (h *AdHandler) handleError(w http.ResponseWriter, err error, requestID string) int {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)
	var status int
	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"error": err.Error()}

	switch err.Error() {
	case "ad not found", "ad date not found", "image not found":
		w.WriteHeader(http.StatusNotFound)
		status = http.StatusNotFound
	case "ad already exists", "RoomsNumber out of range", "not owner of ad":
		w.WriteHeader(http.StatusConflict)
		status = http.StatusConflict
	case "no active session", "missing X-CSRF-Token header",
		"invalid JWT token", "user is not host", "session not found", "user ID not found in session":
		w.WriteHeader(http.StatusUnauthorized)
		status = http.StatusUnauthorized
	case "invalid metadata JSON", "invalid multipart form", "input contains invalid characters",
		"input exceeds character limit", "invalid size, type or resolution of image",
		"query offset not int", "query limit not int", "query dateFrom not int",
		"query dateTo not int", "URL contains invalid characters", "URL exceeds character limit",
		"token parse error", "token invalid", "token expired", "bad sign method",
		"failed to decode metadata", "no images provided", "failed to open file",
		"failed to read file", "failed to encode response", "invalid rating value",
		"cant access other user favorites":
		w.WriteHeader(http.StatusBadRequest)
		status = http.StatusBadRequest
	case "error fetching all places", "error fetching images for ad", "error fetching user",
		"error finding user", "error finding city", "error creating place", "error creating date",
		"error saving place", "error updating place", "error updating date",
		"error updating views count", "error deleting place", "get places error",
		"get places per city error", "get user places error", "error creating image",
		"delete ad image error", "failed to generate session id", "failed to save session",
		"failed to delete session", "error generating random bytes for session ID",
		"failed to get session id from request cookie":
		w.WriteHeader(http.StatusInternalServerError)
		status = http.StatusInternalServerError
	default:
		w.WriteHeader(http.StatusInternalServerError)
		status = http.StatusInternalServerError
	}

	if jsonErr := json.NewEncoder(w).Encode(errorResponse); jsonErr != nil {
		logger.AccessLogger.Error("Failed to encode error response",
			zap.String("request_id", requestID),
			zap.Error(jsonErr),
		)
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}

	return status
}
