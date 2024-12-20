package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	"2024_2_FIGHT-CLUB/microservices/ads_service/mocks"
	"bytes"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"testing"
	"time"
)

func TestAdUseCase_GetAllPlaces(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	expectedAds := []domain.GetAllAdsResponse{
		{UUID: "1234", CityID: 1, AuthorUUID: "user123"},
	}
	mockRepo.MockGetAllPlaces = func(ctx context.Context, filter domain.AdFilter, userId string) ([]domain.GetAllAdsResponse, error) {
		return expectedAds, nil
	}

	ctx := context.Background()
	filter := domain.AdFilter{Location: "New York"}
	userId := "user123"
	ads, err := useCase.GetAllPlaces(ctx, filter, userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedAds, ads)
}

func TestAdUseCase_GetOnePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	adID := "ad123"
	isAutrhorized := true
	expectedAd := domain.GetAllAdsResponse{UUID: adID, CityID: 2, AuthorUUID: "user567"}
	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return expectedAd, nil
	}
	mockRepo.MockUpdateViewsCount = func(ctx context.Context, ad domain.GetAllAdsResponse) (domain.GetAllAdsResponse, error) {
		return expectedAd, nil
	}
	ctx := context.Background()
	ad, err := useCase.GetOnePlace(ctx, adID, isAutrhorized)

	assert.NoError(t, err)
	assert.Equal(t, expectedAd, ad)
}

func TestAdUseCase_CreatePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	newAd := domain.Ad{}
	fileHeaders := [][]byte{}
	userId := "user123"
	createRequest := domain.CreateAdRequest{
		CityName: "Los Angeles", Address: "123 Main St", Description: "Nice place", RoomsNumber: 2,
	}

	mockRepo.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, newAd domain.CreateAdRequest, userId string) error {
		return nil
	}

	mockMinioService.UploadFileFunc = func(file []byte, contentType, id string) (string, error) {
		return "uploadedPath", nil
	}

	mockRepo.MockSaveImages = func(ctx context.Context, adUUID string, imagePaths []string) error {
		return nil
	}

	ctx := context.Background()
	err := useCase.CreatePlace(ctx, &newAd, fileHeaders, createRequest, userId)

	assert.NoError(t, err)
}

func TestAdUseCase_UpdatePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	adID := "ad123"
	userID := "user456"
	existingAd := domain.Ad{}
	updateRequest := domain.UpdateAdRequest{
		CityName: "New City", Address: "456 New St", Description: "Updated description", RoomsNumber: 3,
	}
	fileHeaders := [][]byte{}

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{UUID: adID}, nil
	}

	mockRepo.MockUpdatePlace = func(ctx context.Context, ad *domain.Ad, aID, uID string, req domain.UpdateAdRequest) error {
		return nil
	}

	mockMinioService.UploadFileFunc = func(file []byte, contentType, id string) (string, error) {
		return "uploadedPath", nil
	}

	mockRepo.MockSaveImages = func(ctx context.Context, adUUID string, imagePaths []string) error {
		return nil
	}

	ctx := context.Background()
	err := useCase.UpdatePlace(ctx, &existingAd, adID, userID, fileHeaders, updateRequest)

	assert.NoError(t, err)
}

func TestAdUseCase_DeletePlace(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	adID := "ad123"
	userID := "user456"
	imagePaths := []string{"images/path1", "images/path2"}

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{UUID: adID}, nil
	}

	mockRepo.MockGetAdImages = func(ctx context.Context, aID string) ([]string, error) {
		return imagePaths, nil
	}

	mockMinioService.DeleteFileFunc = func(filePath string) error {
		return nil
	}

	mockRepo.MockDeletePlace = func(ctx context.Context, aID, uID string) error {
		return nil
	}

	ctx := context.Background()
	err := useCase.DeletePlace(ctx, adID, userID)

	assert.NoError(t, err)
}

func TestAdUseCase_GetPlacesPerCity(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	city := "New York"
	expectedPlaces := []domain.GetAllAdsResponse{
		{UUID: "1234", CityID: 1, AuthorUUID: "user123"},
	}

	ctx := context.Background()

	mockRepo.MockGetPlacesPerCity = func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
		return expectedPlaces, nil
	}

	places, err := useCase.GetPlacesPerCity(ctx, city)

	assert.NoError(t, err)
	assert.Equal(t, expectedPlaces, places)

	mockRepo.MockGetPlacesPerCity = func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
		return []domain.GetAllAdsResponse{}, nil
	}

	places, err = useCase.GetPlacesPerCity(ctx, city)

	assert.NoError(t, err)
	assert.Equal(t, []domain.GetAllAdsResponse{}, places)

	mockRepo.MockGetPlacesPerCity = func(ctx context.Context, city string) ([]domain.GetAllAdsResponse, error) {
		return nil, errors.New("database error")
	}

	places, err = useCase.GetPlacesPerCity(ctx, city)

	assert.Error(t, err)
	assert.Nil(t, places)
	if err != nil {
		assert.Equal(t, "database error", err.Error())
	}
}

func TestAdUseCase_GetUserPlaces(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	userID := "user123"
	expectedPlaces := []domain.GetAllAdsResponse{
		{UUID: "ad123", CityID: 2, AuthorUUID: userID},
	}

	// Успешный случай - объявления пользователя найдены
	mockRepo.MockGetUserPlaces = func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
		return expectedPlaces, nil
	}

	ctx := context.Background()
	places, err := useCase.GetUserPlaces(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedPlaces, places)

	// Случай, когда у пользователя нет объявлений
	mockRepo.MockGetUserPlaces = func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
		return []domain.GetAllAdsResponse{}, nil
	}

	places, err = useCase.GetUserPlaces(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(places))

	// Случай, когда произошла ошибка при запросе
	mockRepo.MockGetUserPlaces = func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
		return nil, errors.New("database error")
	}

	places, err = useCase.GetUserPlaces(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, places)
	assert.Equal(t, "database error", err.Error())
}

func TestAdUseCase_GetAllPlaces_Error(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	mockRepo.MockGetAllPlaces = func(ctx context.Context, filter domain.AdFilter, userId string) ([]domain.GetAllAdsResponse, error) {
		return nil, errors.New("database error")
	}

	ctx := context.Background()
	filter := domain.AdFilter{Location: "New York"}
	userId := "user123"
	ads, err := useCase.GetAllPlaces(ctx, filter, userId)

	assert.Error(t, err)
	assert.Nil(t, ads)
	assert.Equal(t, "database error", err.Error())
}

func TestAdUseCase_GetOnePlace_Error(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{}, errors.New("ad not found")
	}

	ctx := context.Background()
	adID := "invalid_ad_id"
	isAuthorized := true
	ad, err := useCase.GetOnePlace(ctx, adID, isAuthorized)

	assert.Error(t, err)
	assert.Equal(t, "ad not found", err.Error())
	assert.Equal(t, domain.GetAllAdsResponse{}, ad)
}

func TestAdUseCase_CreatePlace_ErrorOnCreate(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	newAd := domain.Ad{}
	userId := "user123"
	createRequest := domain.CreateAdRequest{
		CityName: "Los Angeles", Address: "123 Main St", Description: "Nice place", RoomsNumber: 2,
	}

	mockRepo.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, newAd domain.CreateAdRequest, userId string) error {
		return errors.New("creation failed")
	}

	ctx := context.Background()
	err := useCase.CreatePlace(ctx, &newAd, nil, createRequest, userId)

	assert.Error(t, err)
	assert.Equal(t, "creation failed", err.Error())
}

func TestAdUseCase_CreatePlace_ErrorOnUploadImage(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	newAd := domain.Ad{}
	userId := "user123"
	fileHeaders, err := createValidFileHeaders(3)
	if err != nil {
		t.Fatalf("Failed to create valid files: %v", err)
	}

	createRequest := domain.CreateAdRequest{
		CityName: "Los Angeles", Address: "123 Main St", Description: "Nice place", RoomsNumber: 2,
	}

	mockRepo.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, newAd domain.CreateAdRequest, userId string) error {
		return nil
	}

	mockMinioService.UploadFileFunc = func(file []byte, contentType, id string) (string, error) {
		return "", errors.New("upload failed")
	}
	mockRepo.MockSaveImages = func(ctx context.Context, adUUID string, imagePaths []string) error {
		return errors.New("save image failed")
	}

	ctx := context.Background()
	err = useCase.CreatePlace(ctx, &newAd, fileHeaders, createRequest, userId)

	assert.Error(t, err)
	assert.Equal(t, "upload failed", err.Error())
}

func TestAdUseCase_CreatePlace_ErrorOnSaveImage(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	newAd := domain.Ad{}
	userId := "user123"
	fileHeaders, err := createValidFileHeaders(3)
	if err != nil {
		t.Fatalf("Failed to create valid files: %v", err)
	}

	createRequest := domain.CreateAdRequest{
		CityName: "Los Angeles", Address: "123 Main St", Description: "Nice place", RoomsNumber: 2,
	}

	mockRepo.MockCreatePlace = func(ctx context.Context, ad *domain.Ad, newAd domain.CreateAdRequest, userId string) error {
		return nil
	}

	mockMinioService.UploadFileFunc = func(file []byte, contentType, id string) (string, error) {
		return "", nil
	}
	mockRepo.MockSaveImages = func(ctx context.Context, adUUID string, imagePaths []string) error {
		return errors.New("save image failed")
	}

	ctx := context.Background()
	err = useCase.CreatePlace(ctx, &newAd, fileHeaders, createRequest, userId)

	assert.Error(t, err)
	assert.Equal(t, "save image failed", err.Error())
}

func TestAdUseCase_UpdatePlace_ErrorOnUploadImage(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)
	fileHeaders, err := createValidFileHeaders(3)
	if err != nil {
		return
	}
	adID := "invalid_ad_id"
	userID := "user456"
	newAd := domain.Ad{}
	updateRequest := domain.UpdateAdRequest{
		CityName: "New City", Address: "456 New St", Description: "Updated description", RoomsNumber: 3,
	}

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{}, nil
	}
	mockMinioService.UploadFileFunc = func(file []byte, contentType, id string) (string, error) {
		return "", errors.New("upload failed")
	}
	ctx := context.Background()
	err = useCase.UpdatePlace(ctx, &newAd, adID, userID, fileHeaders, updateRequest)

	assert.Error(t, err)
	assert.Equal(t, "upload failed", err.Error())
}

func TestAdUseCase_UpdatePlace_ErrorOnGet(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	adID := "invalid_ad_id"
	userID := "user456"
	updateRequest := domain.UpdateAdRequest{
		CityName: "New City", Address: "456 New St", Description: "Updated description", RoomsNumber: 3,
	}

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{}, errors.New("ad not found")
	}

	ctx := context.Background()
	err := useCase.UpdatePlace(ctx, nil, adID, userID, nil, updateRequest)

	assert.Error(t, err)
	assert.Equal(t, "ad not found", err.Error())
}

func TestAdUseCase_DeletePlace_ErrorOnGet(t *testing.T) {
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	adID := "invalid_ad_id"
	userID := "user456"

	mockRepo.MockGetPlaceById = func(ctx context.Context, id string) (domain.GetAllAdsResponse, error) {
		return domain.GetAllAdsResponse{}, errors.New("ad not found")
	}

	ctx := context.Background()
	err := useCase.DeletePlace(ctx, adID, userID)

	assert.Error(t, err)
	assert.Equal(t, "ad not found", err.Error())
}

func TestAdUseCase_AddToFavorites(t *testing.T) {
	setupLogger()
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	ctx := context.Background()
	validAdID := "ad123"
	invalidAdID := "ad!@#"
	overMaxLenID := string(make([]rune, 256))
	userID := "user123"

	mockRepo.MockAddToFavorites = func(ctx context.Context, adId, userId string) error {
		return nil
	}
	mockRepo.MockUpdateFavoritesCount = func(ctx context.Context, adId string) error {
		return nil
	}

	err := useCase.AddToFavorites(ctx, validAdID, userID)
	assert.NoError(t, err)

	err = useCase.AddToFavorites(ctx, invalidAdID, userID)
	assert.Error(t, err)
	assert.Equal(t, "input contains invalid characters", err.Error())

	err = useCase.AddToFavorites(ctx, overMaxLenID, userID)
	assert.Error(t, err)
	assert.Equal(t, "input exceeds character limit", err.Error())

	mockRepo.MockAddToFavorites = func(ctx context.Context, adId, userId string) error {
		return errors.New("repository error")
	}
	err = useCase.AddToFavorites(ctx, validAdID, userID)
	assert.Error(t, err)
	assert.Equal(t, "repository error", err.Error())

	mockRepo.MockAddToFavorites = func(ctx context.Context, adId, userId string) error {
		return nil
	}
	mockRepo.MockUpdateFavoritesCount = func(ctx context.Context, adId string) error {
		return errors.New("update count error")
	}

	err = useCase.AddToFavorites(ctx, validAdID, userID)
	assert.Error(t, err)
	assert.Equal(t, "update count error", err.Error())
}

func TestAdUseCase_DeleteFromFavorites(t *testing.T) {
	setupLogger()
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	ctx := context.Background()
	validAdID := "ad123"
	invalidAdID := "ad!@#"
	overMaxLenID := string(make([]rune, 256))
	userID := "user123"

	mockRepo.MockDeleteFromFavorites = func(ctx context.Context, adId, userId string) error {
		return nil
	}
	mockRepo.MockUpdateFavoritesCount = func(ctx context.Context, adId string) error {
		return nil
	}

	err := useCase.DeleteFromFavorites(ctx, validAdID, userID)
	assert.NoError(t, err)

	err = useCase.DeleteFromFavorites(ctx, invalidAdID, userID)
	assert.Error(t, err)
	assert.Equal(t, "input contains invalid characters", err.Error())

	err = useCase.DeleteFromFavorites(ctx, overMaxLenID, userID)
	assert.Error(t, err)
	assert.Equal(t, "input exceeds character limit", err.Error())

	mockRepo.MockDeleteFromFavorites = func(ctx context.Context, adId, userId string) error {
		return errors.New("repository error")
	}
	err = useCase.DeleteFromFavorites(ctx, validAdID, userID)
	assert.Error(t, err)
	assert.Equal(t, "repository error", err.Error())

	mockRepo.MockDeleteFromFavorites = func(ctx context.Context, adId, userId string) error {
		return nil
	}
	mockRepo.MockUpdateFavoritesCount = func(ctx context.Context, adId string) error {
		return errors.New("update count error")
	}

	err = useCase.DeleteFromFavorites(ctx, validAdID, userID)
	assert.Error(t, err)
	assert.Equal(t, "update count error", err.Error())
}

func TestAdUseCase_GetUserFavorites(t *testing.T) {
	setupLogger()
	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	ctx := context.Background()
	validUserID := "user123"
	invalidUserID := "user!@#"
	overMaxLenUserID := string(make([]rune, 256))

	expectedAds := []domain.GetAllAdsResponse{
		{UUID: "ad123", AuthorUUID: "user123"},
	}

	mockRepo.MockGetUserFavorites = func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
		return expectedAds, nil
	}

	ads, err := useCase.GetUserFavorites(ctx, validUserID)
	assert.NoError(t, err)
	assert.Equal(t, expectedAds, ads)

	ads, err = useCase.GetUserFavorites(ctx, invalidUserID)
	assert.Error(t, err)
	assert.Nil(t, ads)
	assert.Equal(t, "input contains invalid characters", err.Error())

	ads, err = useCase.GetUserFavorites(ctx, overMaxLenUserID)
	assert.Error(t, err)
	assert.Nil(t, ads)
	assert.Equal(t, "input exceeds character limit", err.Error())

	mockRepo.MockGetUserFavorites = func(ctx context.Context, userId string) ([]domain.GetAllAdsResponse, error) {
		return nil, errors.New("repository error")
	}

	ads, err = useCase.GetUserFavorites(ctx, validUserID)
	assert.Error(t, err)
	assert.Nil(t, ads)
	assert.Equal(t, "repository error", err.Error())
}

func TestAdUseCase_UpdatePriority(t *testing.T) {
	setupLogger()

	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	ctx := context.Background()

	t.Run("Error: invalid characters in adId", func(t *testing.T) {
		invalidAdID := "invalid#ID" // Недопустимый символ #
		userID := "user123"
		amount := 5

		err := useCase.UpdatePriority(ctx, invalidAdID, userID, amount)
		assert.Error(t, err)
		assert.Equal(t, "input contains invalid characters", err.Error())
	})

	t.Run("Error: adId exceeds max length", func(t *testing.T) {
		overMaxLenID := string(make([]rune, 256)) // Строка длиной 256 символов
		userID := "user123"
		amount := 5

		err := useCase.UpdatePriority(ctx, overMaxLenID, userID, amount)
		assert.Error(t, err)
		assert.Equal(t, "input exceeds character limit", err.Error())
	})

	t.Run("Success: valid input", func(t *testing.T) {
		validAdID := "ad123"
		userID := "user123"
		amount := 10

		mockRepo.MockUpdatePriority = func(ctx context.Context, adId, userId string, amount int) error {
			return nil
		}

		err := useCase.UpdatePriority(ctx, validAdID, userID, amount)
		assert.NoError(t, err)
	})
}

func TestAdUseCase_StartPriorityResetWorker(t *testing.T) {
	setupLogger()

	mockRepo := &mocks.MockAdRepository{}
	mockMinioService := &mocks.MockMinioService{}
	useCase := NewAdUseCase(mockRepo, mockMinioService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resetCalled := false

	mockRepo.MockResetExpiredPriorities = func(ctx context.Context) error {
		resetCalled = true
		return nil
	}

	useCase.StartPriorityResetWorker(ctx, 10*time.Millisecond)

	time.Sleep(50 * time.Millisecond)
	cancel()

	assert.True(t, resetCalled, "ResetExpiredPriorities should be called")
}

func TestDeleteAdImage(t *testing.T) {
	ctx := context.Background()
	adId := "ad-uuid"
	imageId := "123"
	userId := "user-uuid"
	imageURL := "/images/image.jpg"

	// Тест успешного удаления изображения
	t.Run("success", func(t *testing.T) {
		adRepoMock := &mocks.MockAdRepository{
			MockDeleteAdImage: func(ctx context.Context, adId string, imageId int, userId string) (string, error) {
				return imageURL, nil
			},
		}
		minioServiceMock := &mocks.MockMinioService{
			DeleteFileFunc: func(imageURL string) error {
				return nil
			},
		}

		adUseCase := NewAdUseCase(adRepoMock, minioServiceMock)

		// Вызываем функцию
		err := adUseCase.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.NoError(t, err)
	})

	// Тест ошибки при удалении изображения в репозитории
	t.Run("repository delete error", func(t *testing.T) {
		expectedErr := errors.New("repository delete error")

		adRepoMock := &mocks.MockAdRepository{
			MockDeleteAdImage: func(ctx context.Context, adId string, imageId int, userId string) (string, error) {
				return "", expectedErr
			},
		}
		minioServiceMock := &mocks.MockMinioService{
			DeleteFileFunc: func(imageURL string) error {
				return nil
			},
		}

		adUseCase := NewAdUseCase(adRepoMock, minioServiceMock)

		// Вызываем функцию
		err := adUseCase.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	// Тест ошибки при удалении файла в MinIO
	t.Run("minio delete error", func(t *testing.T) {
		adRepoMock := &mocks.MockAdRepository{
			MockDeleteAdImage: func(ctx context.Context, adId string, imageId int, userId string) (string, error) {
				return imageURL, nil
			},
		}
		minioErr := errors.New("failed to delete file from MinIO")
		minioServiceMock := &mocks.MockMinioService{
			DeleteFileFunc: func(imageURL string) error {
				return minioErr
			},
		}

		adUseCase := NewAdUseCase(adRepoMock, minioServiceMock)

		// Вызываем функцию
		err := adUseCase.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат: основная операция успешна, но возникает ошибка в логировании MinIO
		assert.NoError(t, err)
	})
}

func generateValidImageBytes() ([]byte, error) {
	width, height := 100, 100
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, randColor())
		}
	}

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func randColor() color.Color {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

func createValidFileHeaders(numFiles int) ([][]byte, error) {
	var fileHeaders [][]byte

	for i := 0; i < numFiles; i++ {
		file, err := generateValidImageBytes()
		if err != nil {
			return nil, err
		}
		fileHeaders = append(fileHeaders, file)
	}

	return fileHeaders, nil
}

func setupLogger() {
	logger.AccessLogger = zap.NewNop()
}
