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
	Base
	Email string `json:"email"`
}

type Order struct {
	Base
	UserID *uuid.UUID `json:"user"`
	Items  []*Item    `gorm:"many2many:order_items;" json:"items"`
	Total  float64    `json:"total"`
	Status string     `json:"status"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
