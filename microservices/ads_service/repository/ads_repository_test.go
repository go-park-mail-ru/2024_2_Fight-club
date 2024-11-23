package repository

import (
	"2024_2_FIGHT-CLUB/domain"
	"2024_2_FIGHT-CLUB/internal/service/logger"
	ntype "2024_2_FIGHT-CLUB/internal/service/type"
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"regexp"
	"testing"
	"time"
)

func setupDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	return gormDB, mock, err
}

func TestGetAllPlaces(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	filter := domain.AdFilter{}

	query := `SELECT ads.*, cities.title as cityName FROM "ads" JOIN cities ON ads."cityId" = cities.id JOIN users ON ads."authorUUID" = users.uuid`
	rows := sqlmock.NewRows([]string{"uuid", "cityId", "authorUUID", "address", "publicationDate", "description", "roomsNumber", "avatar", "name", "rating", "cityname"}).
		AddRow("some-uuid", 1, "author-uuid", "Some Address", time.Now(), "Some Description", 3, "avatar_url", "Username", 4.5, "City Name")
	mock.ExpectQuery(query).WillReturnRows(rows)

	imagesQuery := "SELECT * FROM \"images\" WHERE \"adId\" = $1"
	imageRows := sqlmock.NewRows(ntype.StringArray{"imageUrl"}).AddRow("images/image1.jpg").AddRow("images/image2.jpg")
	mock.ExpectQuery(regexp.QuoteMeta(imagesQuery)).WithArgs("some-uuid").WillReturnRows(imageRows)

	query2 := "SELECT * FROM \"users\" WHERE uuid = $1"
	rows2 := sqlmock.NewRows([]string{"uuid", "username", "password", "email", "username"}).
		AddRow("some-uuid", "test_username", "some_password", "test@example.com", "test_username")
	mock.ExpectQuery(regexp.QuoteMeta(query2)).WillReturnRows(rows2)

	filter = domain.AdFilter{Location: "", Rating: "", NewThisWeek: "", HostGender: "", GuestCount: ""}
	ads, err := repo.GetAllPlaces(context.Background(), filter)

	require.NoError(t, err)
	assert.Len(t, ads, 1)
	assert.Equal(t, "some-uuid", ads[0].UUID)
	assert.Equal(t, "Some Address", ads[0].Address)
	assert.Equal(t, 3, ads[0].RoomsNumber)
	assert.Equal(t, "City Name", ads[0].Cityname)
	assert.ElementsMatch(t, []domain.ImageResponse{
		{
			ID:        0,
			ImagePath: "images/image1.jpg",
		},
		{
			ID:        0,
			ImagePath: "images/image2.jpg",
		},
	}, ads[0].Images)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetAllPlaces_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)
	filter := domain.AdFilter{}
	query := `SELECT ads.*, users.avatar, users.name, users.score as rating , cities.title as cityName FROM "ads" JOIN users ON ads."authorUUID" = users.uuid JOIN cities ON ads."cityId" = cities.id`
	mock.ExpectQuery(query).
		WillReturnError(errors.New("db error"))

	ads, err := repo.GetAllPlaces(context.Background(), filter)
	assert.Error(t, err)
	assert.Nil(t, ads)
}

func TestGetPlaceById(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	expectedAd := domain.GetAllAdsResponse{
		UUID:       "some-uuid",
		CityID:     1,
		AuthorUUID: "author-uuid",
		Address:    "Some Address",
	}

	// Step 2: Define Mock Database Expectations
	query := "SELECT ads.*, cities.title as cityName FROM \"ads\" JOIN users ON ads.\"authorUUID\" = users.uuid JOIN cities ON ads.\"cityId\" = cities.id WHERE ads.uuid = $1"
	rows := sqlmock.NewRows([]string{"uuid", "cityId", "authorUUID", "address", "avatar", "name", "rating", "cityName"}).
		AddRow("some-uuid", 1, "author-uuid", "Some Address", "avatar_url", "UserName", 4.5, "CityName")
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("some-uuid").WillReturnRows(rows)

	imagesQuery := "SELECT * FROM \"images\" WHERE \"adId\" = $1"
	imageRows := sqlmock.NewRows(ntype.StringArray{"imageUrl"}).AddRow("images/image1.jpg").AddRow("images/image2.jpg")
	mock.ExpectQuery(regexp.QuoteMeta(imagesQuery)).WithArgs("some-uuid").WillReturnRows(imageRows)

	query2 := "SELECT * FROM \"users\" WHERE uuid = $1"
	rows2 := sqlmock.NewRows([]string{"uuid", "username", "password", "email", "username"}).
		AddRow("some-uuid", "test_username", "some_password", "test@example.com", "test_username")
	mock.ExpectQuery(regexp.QuoteMeta(query2)).WillReturnRows(rows2)

	ad, err := repo.GetPlaceById(context.Background(), "some-uuid")

	require.NoError(t, err)
	assert.Equal(t, expectedAd.UUID, ad.UUID)
	assert.Equal(t, expectedAd.CityID, ad.CityID)
	assert.Equal(t, expectedAd.AuthorUUID, ad.AuthorUUID)
	assert.Equal(t, expectedAd.Address, ad.Address)
	assert.ElementsMatch(t, []domain.ImageResponse{
		{
			ID:        0,
			ImagePath: "images/image1.jpg",
		},
		{
			ID:        0,
			ImagePath: "images/image2.jpg",
		},
	}, ad.Images)
	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestGetPlaceById_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()
	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	mock.ExpectQuery(`SELECT ads\.\*, users.avatar, users.name, users.score as rating , cities.title as cityName FROM "ads" JOIN users ON ads."authorUUID" = users.uuid [\s\S]* WHERE ads\.uuid = \$1`).
		WithArgs("ad1").
		WillReturnError(errors.New("db error"))

	ad, err := repo.GetPlaceById(context.Background(), "ad1")
	assert.Error(t, err)
	assert.Empty(t, ad)
}

func TestUpdatePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	adId := "existing-ad-id"
	userId := "author-id"
	updatedRequest := domain.UpdateAdRequest{
		CityName:    "New City",
		Address:     "New Address",
		Description: "Updated Description",
		RoomsNumber: 3,
	}

	adRows := sqlmock.NewRows([]string{"uuid", "authorUUID", "cityId", "address", "roomsNumber"}).
		AddRow("existing-ad-id", "author-id", 1, "Old Address", 2)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).WithArgs(adId, 1).WillReturnRows(adRows)

	cityRows := sqlmock.NewRows([]string{"id"}).AddRow(2)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE title = $1 ORDER BY "cities"."id" LIMIT $2`)).WithArgs("New City", 1).WillReturnRows(cityRows)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "cityId"=$1 WHERE "uuid" = $2`)).
		WithArgs(2, adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = repo.UpdatePlace(context.Background(), &domain.Ad{}, adId, userId, updatedRequest)

	assert.NoError(t, err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdatePlace_AdNotFound(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	adId := "non-existing-ad-id"
	userId := "author-id"
	updatedRequest := domain.UpdateAdRequest{
		CityName:    "New City",
		Address:     "New Address",
		Description: "Updated Description",
		RoomsNumber: 3,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adId, 1).WillReturnRows(sqlmock.NewRows([]string{"uuid", "authorUUID", "cityId", "address", "roomsNumber"}))

	err = repo.UpdatePlace(context.Background(), &domain.Ad{}, adId, userId, updatedRequest)

	assert.Error(t, err)
	assert.Equal(t, errors.New("ad not found"), err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdatePlace_UserNotAuthorized(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	repo := NewAdRepository(db)

	adId := "existing-ad-id"
	userId := "unauthorized-user-id"
	updatedRequest := domain.UpdateAdRequest{
		CityName:    "New City",
		Address:     "New Address",
		Description: "Updated Description",
		RoomsNumber: 3,
	}

	adRows := sqlmock.NewRows([]string{"uuid", "authorUUID", "cityId", "address", "roomsNumber"}).
		AddRow("existing-ad-id", "different-author-id", 1, "Old Address", 2)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adId, 1).WillReturnRows(adRows)

	err = repo.UpdatePlace(context.Background(), &domain.Ad{}, adId, userId, updatedRequest)

	assert.Error(t, err)
	assert.Equal(t, errors.New("not owner of ad"), err)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeletePlace(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adId := "1"
	userId := "123"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "address", "authorUUID"}).
			AddRow("1", "test_address", "123"))

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "images" WHERE "adId" = $1`)).
		WithArgs(adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ad_positions" WHERE "adId" = $1`)).
		WithArgs(adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ad_available_dates" WHERE "adId" = $1`)).
		WithArgs(adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "requests" WHERE "adId" = $1`)).
		WithArgs(adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "ads" WHERE "ads"."uuid" = $1`)).
		WithArgs(adId).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = adRepo.DeletePlace(ctx, adId, userId)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeletePlace_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adId := "1"
	userId := "not-owner"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "ads" WHERE uuid = $1 ORDER BY "ads"."uuid" LIMIT $2`)).
		WithArgs(adId, 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "location_main", "author_uuid"}).
			AddRow("1", "New York", "123"))

	err = adRepo.DeletePlace(ctx, adId, userId)

	assert.Error(t, err)
	assert.Equal(t, "not owner of ad", err.Error())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePlace_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	newAdReq := domain.CreateAdRequest{
		CityName:    "New York",
		Address:     "123 Main St",
		Description: "A lovely place",
		RoomsNumber: 3,
	}

	city := domain.City{
		ID:    1,
		Title: "New York",
	}

	ad := &domain.Ad{
		UUID:            "ad-uuid-123",
		CityID:          city.ID,
		AuthorUUID:      "author-uuid-456",
		Address:         "123 Main St",
		Description:     "A lovely place",
		RoomsNumber:     3,
		PublicationDate: time.Now(),
	}

	// Mock query to find the city
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE title = $1 ORDER BY "cities"."id" LIMIT $2`)).
		WithArgs(newAdReq.CityName, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(city.ID, city.Title))

	// Mock insert ad
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "ads" ("cityId","authorUUID","address","publicationDate","description","roomsNumber","uuid") VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING "uuid"`)).
		WithArgs(ad.CityID, ad.AuthorUUID, ad.Address, sqlmock.AnyArg(), ad.Description, ad.RoomsNumber, ad.UUID).
		WillReturnRows(sqlmock.NewRows([]string{"uuid"}).AddRow(ad.UUID))
	mock.ExpectCommit()

	err = adRepo.CreatePlace(ctx, ad, newAdReq)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestCreatePlace_CityNotFound(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	newAdReq := domain.CreateAdRequest{
		CityName:    "Unknown City",
		Address:     "123 Main St",
		Description: "A lovely place",
		RoomsNumber: 3,
	}

	ad := &domain.Ad{
		UUID:            "ad-uuid-123",
		AuthorUUID:      "author-uuid-456",
		Address:         "123 Main St",
		Description:     "A lovely place",
		RoomsNumber:     3,
		PublicationDate: time.Now(),
	}

	// Mock query to find the city - no rows
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "cities" WHERE title = $1 ORDER BY "cities"."id" LIMIT $2`)).
		WithArgs(newAdReq.CityName, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}))

	err = adRepo.CreatePlace(ctx, ad, newAdReq)

	if err == nil || err.Error() != "Error finding city" {
		t.Errorf("Expected 'Error finding city', got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestSavePlace_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	ad := &domain.Ad{
		UUID:            "ad-uuid-123",
		CityID:          1,
		AuthorUUID:      "author-uuid-456",
		Address:         "123 Main St",
		Description:     "Updated description",
		RoomsNumber:     4,
		PublicationDate: time.Now(),
	}

	// Mock update ad
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "cityId"=$1,"authorUUID"=$2,"address"=$3,"publicationDate"=$4,"description"=$5,"roomsNumber"=$6 WHERE "uuid" = $7`)).
		WithArgs(ad.CityID, ad.AuthorUUID, ad.Address, time.Now(), ad.Description, ad.RoomsNumber, ad.UUID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = adRepo.SavePlace(ctx, ad)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestSavePlace_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	ad := &domain.Ad{
		UUID:            "ad-uuid-123",
		CityID:          1,
		AuthorUUID:      "author-uuid-456",
		Address:         "123 Main St",
		Description:     "Updated description",
		RoomsNumber:     4,
		PublicationDate: time.Now(),
	}

	// Mock update ad with error
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "ads" SET "cityId"=$1,"authorUUID"=$2,"address"=$3,"publicationDate"=$4,"description"=$5,"roomsNumber"=$6 WHERE "uuid" = $7`)).
		WithArgs(ad.CityID, ad.AuthorUUID, ad.Address, time.Now(), ad.Description, ad.RoomsNumber, ad.UUID).
		WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	err = adRepo.SavePlace(ctx, ad)

	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetPlacesPerCity_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	city := "New York"

	rows := sqlmock.NewRows([]string{
		"uuid", "cityId", "authorUUID", "address",
		"publicationDate", "description", "roomsNumber",
		"avatar", "name", "rating", "cityName",
	}).
		AddRow("ad-uuid-123", 1, "author-uuid-456", "123 Main St",
			time.Now(), "A lovely place", 3, "avatar.png", "John Doe", 4.5, "New York")

	// Mock select ads
	mock.ExpectQuery(regexp.QuoteMeta("SELECT ads.*, cities.title as cityName FROM \"ads\" JOIN users ON ads.\"authorUUID\" = users.uuid JOIN cities ON ads.\"cityId\" = cities.id WHERE cities.\"enTitle\" = $1")).
		WithArgs(city).
		WillReturnRows(rows)

	imagesQuery := "SELECT * FROM \"images\" WHERE \"adId\" = $1"
	imageRows := sqlmock.NewRows(ntype.StringArray{"imageUrl"}).AddRow("images/image1.jpg").AddRow("images/image2.jpg")
	mock.ExpectQuery(regexp.QuoteMeta(imagesQuery)).WithArgs("ad-uuid-123").WillReturnRows(imageRows)

	query2 := "SELECT * FROM \"users\" WHERE uuid = $1"
	rows2 := sqlmock.NewRows([]string{"uuid", "username", "password", "email", "username"}).
		AddRow("some-uuid", "test_username", "some_password", "test@example.com", "test_username")
	mock.ExpectQuery(regexp.QuoteMeta(query2)).WillReturnRows(rows2)

	ads, err := adRepo.GetPlacesPerCity(ctx, city)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ads) != 1 {
		t.Errorf("Expected 1 ad, got %d", len(ads))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetPlacesPerCity_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	city := "Unknown City"

	// Mock select ads with error
	mock.ExpectQuery(regexp.QuoteMeta("SELECT ads.*, cities.title as cityName FROM \"ads\" JOIN users ON ads.\"authorUUID\" = users.uuid JOIN cities ON ads.\"cityId\" = cities.id WHERE cities.\"enTitle\" = $1")).
		WithArgs(city).
		WillReturnError(gorm.ErrInvalidData)

	ads, err := adRepo.GetPlacesPerCity(ctx, city)

	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if ads != nil {
		t.Errorf("Expected nil ads, got %v", ads)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestSaveImages_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adUUID := "ad-uuid-123"
	imagePaths := []string{"img1.png", "img2.png"}

	for _, path := range imagePaths {
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "images" ("adId","imageUrl") VALUES ($1,$2) RETURNING "id"`)).
			WithArgs(adUUID, path).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()
	}

	err = adRepo.SaveImages(ctx, adUUID, imagePaths)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestSaveImages_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adUUID := "ad-uuid-123"
	imagePaths := []string{"img1.png", "img2.png"}

	// Mock insert first image
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "images" ("adId","imageUrl") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(adUUID, imagePaths[0]).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	// Mock insert second image with error
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "images" ("adId","imageUrl") VALUES ($1,$2) RETURNING "id"`)).
		WithArgs(adUUID, imagePaths[1]).
		WillReturnError(gorm.ErrInvalidData)
	mock.ExpectRollback()

	err = adRepo.SaveImages(ctx, adUUID, imagePaths)

	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetAdImages_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adId := "ad-uuid-123"

	rows := sqlmock.NewRows([]string{"imageUrl"}).
		AddRow("img1.png").
		AddRow("img2.png")

	// Mock select images
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "imageUrl" FROM "images" WHERE "adId" = $1`)).
		WithArgs(adId).
		WillReturnRows(rows)

	images, err := adRepo.GetAdImages(ctx, adId)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(images) != 2 {
		t.Errorf("Expected 2 images, got %d", len(images))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetAdImages_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	adId := "ad-uuid-123"

	// Mock select images with error
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "imageUrl" FROM "images" WHERE "adId" = $1`)).
		WithArgs(adId).
		WillReturnError(gorm.ErrInvalidData)

	images, err := adRepo.GetAdImages(ctx, adId)

	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if images != nil {
		t.Errorf("Expected nil images, got %v", images)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetUserPlaces_Success(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	userId := "author-uuid-456"

	rows := sqlmock.NewRows([]string{
		"uuid", "cityId", "authorUUID", "address",
		"publicationDate", "description", "roomsNumber",
		"avatar", "name", "rating", "cityName",
	}).
		AddRow("ad-uuid-123", 1, userId, "123 Main St",
			time.Now(), "A lovely place", 3, "avatar.png", "John Doe", 4.5, "New York")

	// Mock select ads
	mock.ExpectQuery(`SELECT ads\.\*, users.avatar, users.name, users.score as rating , cities.title as cityName FROM "ads" JOIN users ON ads\."authorUUID" = users\.uuid JOIN cities ON  ads\."cityId" = cities\.id WHERE users\.uuid = \$1`).
		WithArgs(userId).
		WillReturnRows(rows)

	// Mock select images for each ad
	imageRows := sqlmock.NewRows([]string{"uuid", "adId", "imageUrl"}).
		AddRow("img-uuid-1", "ad-uuid-123", "img1.png").
		AddRow("img-uuid-2", "ad-uuid-123", "img2.png")

	mock.ExpectQuery(`SELECT \* FROM "images" WHERE "adId" = \$1`).
		WithArgs("ad-uuid-123").
		WillReturnRows(imageRows)

	ads, err := adRepo.GetUserPlaces(ctx, userId)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ads) != 1 {
		t.Errorf("Expected 1 ad, got %d", len(ads))
	}

	if len(ads[0].Images) != 2 {
		t.Errorf("Expected 2 images, got %d", len(ads[0].Images))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestGetUserPlaces_Failure(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	userId := "author-uuid-456"

	rows := sqlmock.NewRows([]string{
		"uuid", "cityId", "authorUUID", "address",
		"publicationDate", "description", "roomsNumber",
		"avatar", "name", "rating", "cityName",
	}).
		AddRow("ad-uuid-123", 1, userId, "123 Main St",
			time.Now(), "A lovely place", 3, "avatar.png", "John Doe", 4.5, "New York")

	// Mock select ads
	mock.ExpectQuery(`SELECT ads\.\*, users.avatar, users.name, users.score as rating , cities.title as cityName FROM "ads" JOIN users ON ads\."authorUUID" = users\.uuid JOIN cities ON  ads\."cityId" = cities\.id WHERE users\.uuid = \$1`).
		WithArgs(userId).
		WillReturnRows(rows)

	// Mock select images with error
	mock.ExpectQuery(`SELECT \* FROM "images" WHERE "adId" = \$1`).
		WithArgs("ad-uuid-123").
		WillReturnError(gorm.ErrInvalidData)

	ads, err := adRepo.GetUserPlaces(ctx, userId)

	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if ads != nil {
		t.Errorf("Expected nil ads, got %v", ads)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet expectations: %v", err)
	}
}

func TestDeleteAdImage(t *testing.T) {
	if err := logger.InitLoggers(); err != nil {
		log.Fatalf("Failed to initialize loggers: %v", err)
	}
	defer logger.SyncLoggers()

	db, mock, err := setupDBMock()
	assert.Nil(t, err)

	adRepo := NewAdRepository(db)
	ctx := context.Background()

	// Тест успешного удаления изображения
	t.Run("success", func(t *testing.T) {
		adId := "ad-uuid"
		imageId := 1
		userId := "user-uuid"
		imageUrl := "/images/image.jpg"

		// Настраиваем ожидания
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"ads\" WHERE uuid = $1 ORDER BY \"ads\".\"uuid\" LIMIT $2")).
			WithArgs(adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "authorUUID"}).AddRow(adId, userId))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"images\" WHERE id = $1 AND \"adId\" = $2 ORDER BY \"images\".\"id\" LIMIT $3")).
			WithArgs(imageId, adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "adId", "imageUrl"}).AddRow(imageId, adId, imageUrl))

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM \"images\"").
			WithArgs(imageId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Вызываем функцию
		result, err := adRepo.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.NoError(t, err)
		assert.Equal(t, imageUrl, result)

		// Проверяем все ожидания выполнены
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	// Тест случая, когда объявление не найдено
	t.Run("ad not found", func(t *testing.T) {
		adId := "ad-uuid"
		imageId := 1
		userId := "user-uuid"

		// Настраиваем ожидания
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"ads\" WHERE uuid = $1 ORDER BY \"ads\".\"uuid\" LIMIT $2")).
			WithArgs(adId, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Вызываем функцию
		result, err := adRepo.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "ad not found", err.Error())
	})

	// Тест случая, когда пользователь не является владельцем объявления
	t.Run("not owner of ad", func(t *testing.T) {
		adId := "ad-uuid"
		imageId := 1
		userId := "user-uuid"
		wrongUUID := "another-user-uuid"

		// Настраиваем ожидания
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"ads\" WHERE uuid = $1 ORDER BY \"ads\".\"uuid\" LIMIT $2")).
			WithArgs(adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "authorUUID"}).AddRow(adId, wrongUUID))

		// Вызываем функцию
		result, err := adRepo.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "not owner of ad", err.Error())
	})

	// Тест случая, когда изображение не найдено
	t.Run("image not found", func(t *testing.T) {
		adId := "ad-uuid"
		imageId := 1
		userId := "user-uuid"

		// Настраиваем ожидания
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"ads\" WHERE uuid = $1 ORDER BY \"ads\".\"uuid\" LIMIT $2")).
			WithArgs(adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "authorUUID"}).AddRow(adId, userId))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"images\" WHERE id = $1 AND \"adId\" = $2 ORDER BY \"images\".\"id\" LIMIT $3")).
			WithArgs(imageId, adId, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		// Вызываем функцию
		result, err := adRepo.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "image not found", err.Error())
	})

	// Тест случая ошибки при удалении изображения
	t.Run("error deleting image", func(t *testing.T) {
		adId := "ad-uuid"
		imageId := 1
		userId := "user-uuid"
		imageUrl := "/images/image.jpg"

		// Настраиваем ожидания
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"ads\" WHERE uuid = $1 ORDER BY \"ads\".\"uuid\" LIMIT $2")).
			WithArgs(adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"uuid", "authorUUID"}).AddRow(adId, userId))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM \"images\" WHERE id = $1 AND \"adId\" = $2 ORDER BY \"images\".\"id\" LIMIT $3")).
			WithArgs(imageId, adId, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "adId", "imageUrl"}).AddRow(imageId, adId, imageUrl))

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM \"images\"").
			WithArgs(imageId).
			WillReturnError(errors.New("delete error"))
		mock.ExpectRollback()

		// Вызываем функцию
		result, err := adRepo.DeleteAdImage(ctx, adId, imageId, userId)

		// Проверяем результат
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Equal(t, "error deleting image from database", err.Error())
	})
}