package domain

import (
	"2024_2_FIGHT-CLUB/internal/service/type"
	"context"
)

type Ad struct {
	ID              string             `gorm:"primaryKey" json:"id"`
	LocationMain    string             `json:"location_main"`
	LocationStreet  string             `json:"location_street"`
	Position        ntype.Float64Array `gorm:"type:float[]" json:"position"`
	Images          ntype.StringArray  `gorm:"type:text[]"`
	AuthorUUID      string             `json:"author_uuid"`
	PublicationDate string             `json:"publication_date"`
	AvailableDates  ntype.StringArray  `gorm:"type:text[]" json:"available_dates"`
	Distance        float32            `json:"distance"`
	Requests        []Request          `gorm:"foreignKey:AdID" json:"requests"`
}

type AdFilter struct {
	Location    string
	Rating      string
	NewThisWeek string
	HostGender  string
	GuestCount  string
}

type AdRepository interface {
	GetAllPlaces(ctx context.Context, filter AdFilter) ([]Ad, error)
	GetPlaceById(ctx context.Context, adId string) (Ad, error)
	CreatePlace(ctx context.Context, ad *Ad) error
	SavePlace(ctx context.Context, ad *Ad) error
	UpdatePlace(ctx context.Context, ad *Ad, adId string, userId string) error
	DeletePlace(ctx context.Context, adId string, userId string) error
	GetPlacesPerCity(ctx context.Context, city string) ([]Ad, error)
}
