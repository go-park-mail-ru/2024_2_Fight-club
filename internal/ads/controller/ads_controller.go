package controller

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/ads/usecase"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/internal/service/middleware"
	"2024_2_FIGHT-CLUB/internal/service/session"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type AdHandler struct {
	adUseCase      usecase.AdUseCase
	sessionService *session.ServiceSession
}

func NewAdHandler(adUseCase usecase.AdUseCase, sessionService *session.ServiceSession) *AdHandler {
	return &AdHandler{
		adUseCase:      adUseCase,
		sessionService: sessionService,
	}
}

func (h *AdHandler) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())

	logger.AccessLogger.Info("Received GetAllPlaces request",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("query", r.URL.Query().Encode()),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	queryParams := r.URL.Query()

	location := queryParams.Get("location")
	rating := queryParams.Get("rating")
	newThisWeek := queryParams.Get("new")
	hostGender := queryParams.Get("gender")
	guestCounter := queryParams.Get("guests")

	filter := domain.AdFilter{
		Location:    location,
		Rating:      rating,
		NewThisWeek: newThisWeek,
		HostGender:  hostGender,
		GuestCount:  guestCounter,
	}

	places, err := h.adUseCase.GetAllPlaces(r.Context(), filter)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"places": places,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	logger.AccessLogger.Info("Received GetOnePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	place, err := h.adUseCase.GetOnePlace(r.Context(), adId)
	if err != nil {
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"place": place,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	logger.AccessLogger.Info("Received CreatePlace request",
		zap.String("request_id", requestID),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var place domain.Ad
	if err := json.NewDecoder(r.Body).Decode(&place); err != nil {
		logger.AccessLogger.Error("Failed to decode request body", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.sessionService.GetUserID(r.Context(), r, w)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		h.handleError(w, errors.New("no active session"), requestID)
		return
	}
	place.AuthorUUID = userID

	err = h.adUseCase.CreatePlace(r.Context(), &place)
	if err != nil {
		logger.AccessLogger.Error("Failed to create place", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"place": place,
	}
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logger.AccessLogger.Error("Failed to encode response", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	// Логирование начала обработки запроса
	logger.AccessLogger.Info("Received UpdatePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var place domain.Ad
	if err := json.NewDecoder(r.Body).Decode(&place); err != nil {
		logger.AccessLogger.Error("Failed to decode request body", zap.String("request_id", requestID), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := h.sessionService.GetUserID(r.Context(), r, w)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		h.handleError(w, errors.New("no active session"), requestID)
		return
	}

	err = h.adUseCase.UpdatePlace(r.Context(), &place, adId, userID)
	if err != nil {
		logger.AccessLogger.Error("Failed to update place", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}
	w.WriteHeader(http.StatusOK)
	updateResponse := map[string]string{"response": "Update successfully"}
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

	logger.AccessLogger.Info("Received DeletePlace request",
		zap.String("request_id", requestID),
		zap.String("adId", adId),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userID, err := h.sessionService.GetUserID(r.Context(), r, w)
	if err != nil {
		logger.AccessLogger.Warn("No active session", zap.String("request_id", requestID))
		h.handleError(w, errors.New("no active session"), requestID)
		return
	}

	err = h.adUseCase.DeletePlace(r.Context(), adId, userID)
	if err != nil {
		logger.AccessLogger.Error("Failed to delete place", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}

	w.WriteHeader(http.StatusOK)
	duration := time.Since(start)
	logger.AccessLogger.Info("Completed DeletePlace request",
		zap.String("request_id", requestID),
		zap.Duration("duration", duration),
		zap.Int("status", http.StatusOK),
	)
}

func (h *AdHandler) GetPlacesPerCity(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := middleware.GetRequestID(r.Context())
	city := mux.Vars(r)["city"]

	// Логирование начала обработки запроса
	logger.AccessLogger.Info("Received GetPlacesPerCity request",
		zap.String("request_id", requestID),
		zap.String("city", city),
	)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ads, err := h.adUseCase.GetPlacesPerCity(r.Context(), city)
	if err != nil {
		logger.AccessLogger.Error("Failed to get places per city", zap.String("request_id", requestID), zap.Error(err))
		h.handleError(w, err, requestID)
		return
	}
	body := map[string]interface{}{
		"places": ads,
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

func (h *AdHandler) SearchPlace(w http.ResponseWriter, r *http.Request) {

}

func (h *AdHandler) handleError(w http.ResponseWriter, err error, requestID string) {
	logger.AccessLogger.Error("Handling error",
		zap.String("request_id", requestID),
		zap.Error(err),
	)

	w.Header().Set("Content-Type", "application/json")
	errorResponse := map[string]string{"error": err.Error()}

	switch err.Error() {
	case "ad not found":
		w.WriteHeader(http.StatusNotFound)
	case "ad already exists":
		w.WriteHeader(http.StatusConflict)
	case "not owner of ad", "no active session":
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	if jsonErr := json.NewEncoder(w).Encode(errorResponse); jsonErr != nil {
		logger.AccessLogger.Error("Failed to encode error response",
			zap.String("request_id", requestID),
			zap.Error(jsonErr),
		)
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}
}
