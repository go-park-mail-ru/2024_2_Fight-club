package domain

import (
	"context"
	"time"
)

type Review struct {
	ID        int       `gorm:"primary_key;auto_increment;column:id" json:"id"`
	UserID    string    `gorm:"column:userId;not null" json:"userId"`
	HostID    string    `gorm:"column:hostId;not null" json:"hostId"`
	Title     string    `gorm:"type:text;size:250;column:title;not null" json:"title"`
	Text      string    `gorm:"type:text;size:1000;column:text;not null" json:"text"`
	Rating    int       `gorm:"column:rating" json:"rating"`
	CreatedAt time.Time `gorm:"type:timestamp;column:createdAt" json:"createdAt"`
	User      User      `gorm:"foreignkey:UserID;references:UUID" json:"-"`
	Host      User      `gorm:"foreignkey:HostID;references:UUID" json:"-"`
}

type UserReviews struct {
	ID         int       `gorm:"primary_key;auto_increment;column:id" json:"id"`
	UserID     string    `gorm:"column:userId;not null" json:"userId"`
	HostID     string    `gorm:"column:hostId;not null" json:"hostId"`
	Title      string    `gorm:"type:text;size:250;column:title;not null" json:"title"`
	Text       string    `gorm:"type:text;size:1000;column:text;not null" json:"text"`
	Rating     int       `gorm:"column:rating" json:"rating"`
	CreatedAt  time.Time `gorm:"type:timestamp;column:createdAt" json:"createdAt"`
	User       User      `gorm:"foreignkey:UserID;references:UUID" json:"-"`
	Host       User      `gorm:"foreignkey:HostID;references:UUID" json:"-"`
	UserAvatar string    `json:"userAvatar"`
	UserName   string    `json:"userName"`
}

type ReviewRepository interface {
	CreateReview(ctx context.Context, review *Review) error
	GetUserReviews(ctx context.Context, userID string) ([]UserReviews, error)
	DeleteReview(ctx context.Context, userID, hostID string) error
	UpdateReview(ctx context.Context, userID, hostID string, updatedReview *Review) error
}
