package models

import "time"

// Link represents a shortened URL.
type Link struct {
	ID        string    `json:"id" xml:"id" form:"id" gorm:"primaryKey" validate:"required"`
	URL       string    `json:"url" xml:"url" form:"url" gorm:"index" validate:"required"`
	ExpiredAt time.Time `json:"expired_at" xml:"expired_at" form:"expired_at"`
}

// LinkForm is used to create or update a link.
type LinkForm struct {
	URL       string    `json:"url" xml:"url" form:"url" validate:"required"`
	ExpiredAt time.Time `json:"expired_at" xml:"expired_at" form:"expired_at"`
}
