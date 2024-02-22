package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	OrderStatusPending   = "pending"
	OrderStatusCompleted = "completed"
)

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
}

type Item struct {
	Base
	Name        string
	Description *string
	Price       float64
}

type User struct {
	ID    string `gorm:"primary_key"`
	Email string
	Name  string
}

type Order struct {
	Base
	UserID string
	Items  []*Item `gorm:"many2many:order_items;"`
	Total  float64
	Status string
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
