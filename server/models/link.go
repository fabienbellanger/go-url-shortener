package models

import (
	"time"

	"gorm.io/gorm"
)

// Link represents a shortened URL.
type Link struct {
	ID        string         `json:"id" xml:"id" form:"id" gorm:"primaryKey;size:8" validate:"required"`
	URL       string         `json:"url" xml:"url" form:"url" gorm:"index;size:191" validate:"required"`
	Name      *string        `json:"name" xml:"name" form:"name" gorm:"index;size:127"`
	CreatedAt time.Time      `json:"created_at" xml:"created_at" form:"created_at" gorm:"autoCreateTime"`
	ExpiredAt time.Time      `json:"expired_at" xml:"expired_at" form:"expired_at"`
	DeletedAt gorm.DeletedAt `json:"-" xml:"-" form:"deleted_at" gorm:"index"`
}

// LinkForm is used to create or update a link.
type LinkForm struct {
	URL       string    `json:"url" xml:"url" form:"url" validate:"required"`
	Name      *string   `json:"name" xml:"name" form:"name"`
	ExpiredAt time.Time `json:"expired_at" xml:"expired_at" form:"expired_at"`
}
