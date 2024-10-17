package usecase

import (
	"2024_2_FIGHT-CLUB/domain"
	"context"
	"errors"
)

type AdUseCase interface {
	GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error)
	GetOnePlace(ctx context.Context, adId string) (domain.Ad, error)
	CreatePlace(ctx context.Context, place *domain.Ad) error
	UpdatePlace(ctx context.Context, place *domain.Ad, adId string, userId string) error
	DeletePlace(ctx context.Context, adId string, userId string) error
	GetPlacesPerCity(ctx context.Context, city string) ([]domain.Ad, error)
}

type adUseCase struct {
	adRepository domain.AdRepository
}

func NewAdUseCase(adRepository domain.AdRepository) AdUseCase {
	return &adUseCase{
		adRepository: adRepository,
	}
}

func (uc *adUseCase) GetAllPlaces(ctx context.Context, filter domain.AdFilter) ([]domain.Ad, error) {
	ads, err := uc.adRepository.GetAllPlaces(ctx, filter)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

func (uc *adUseCase) GetOnePlace(ctx context.Context, adId string) (domain.Ad, error) {
	ad, err := uc.adRepository.GetPlaceById(ctx, adId)
	if err != nil {
		return ad, errors.New("ad not found")
	}
	return ad, nil
}

func (uc *adUseCase) CreatePlace(ctx context.Context, place *domain.Ad) error {
	err := uc.adRepository.CreatePlace(ctx, place)
	if err != nil {
		return err
	}
	return nil
}

func (uc *adUseCase) UpdatePlace(ctx context.Context, place *domain.Ad, adId string, userId string) error {
	err := uc.adRepository.UpdatePlace(ctx, place, adId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *adUseCase) DeletePlace(ctx context.Context, adId string, userId string) error {
	err := uc.adRepository.DeletePlace(ctx, adId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *adUseCase) GetPlacesPerCity(ctx context.Context, city string) ([]domain.Ad, error) {
	places, err := uc.adRepository.GetPlacesPerCity(ctx, city)
	if err != nil || len(places) == 0 {
		return nil, errors.New("ad not found")
	}
	return places, nil
}
