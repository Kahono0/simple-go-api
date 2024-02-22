package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key"`
}

type Item struct {
	Base
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Price       float64 `json:"price"`
}

type User struct {
	ID    string `gorm:"primary_key"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Order struct {
	Base
	UserID string  `json:"user_id"`
	Items  []*Item `gorm:"many2many:order_items;" json:"items"`
	Total  float64 `json:"total"`
	Status string  `json:"status"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
