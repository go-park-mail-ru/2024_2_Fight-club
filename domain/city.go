package domain

import "context"

type City struct {
	ID          int    `gorm:"primary_key;auto_increment;column:id" json:"id"`
	Title       string `gorm:"type:varchar(100);column:title" json:"title"`
	EnTitle     string `gorm:"type:varchar(100);column:enTitle" json:"enTitle"`
	Description string `gorm:"type:text;size:3000;column:description" json:"description"`
	Image       string `gorm:"type:text;size:1000;column:image" json:"image"`
}

type CityRepository interface {
	GetCities(ctx context.Context) ([]City, error)
	GetCityByEnName(ctx context.Context, cityEnName string) (City, error)
}