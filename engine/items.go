package engine

import (
	"fmt"

	"github.com/Kahono0/simple-go-api/models"
	"github.com/Kahono0/simple-go-api/utils"
	"github.com/google/uuid"
)

func CreateItem(name string, description *string, price float64) (*models.Item, error) {
	db := utils.DB
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	item := models.Item{
		Name:        name,
		Description: description,
		Price:       price,
	}

	result := db.Create(&item)

	if result.Error != nil {
		return nil, result.Error
	}

	return &item, nil
}

func UpdateItem(id string, name *string, description *string, price *float64) (*models.Item, error) {
	db := utils.DB
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}

	item := models.Item{
		Base: models.Base{ID: uid},
	}

	result := db.First(&item)

	if result.Error != nil {
		return nil, result.Error
	}

	if name != nil {
		item.Name = *name
	}

	if description != nil {
		item.Description = description
	}

	if price != nil {
		item.Price = *price
	}

	result = db.Save(&item)

	if result.Error != nil {
		return nil, result.Error
	}

	return GetItem(id)
}

func DeleteItem(id string) error {
	db := utils.DB
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id")
	}

	item := models.Item{
		Base: models.Base{ID: uid},
	}

	result := db.Delete(&item)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetItem(id string) (*models.Item, error) {
	db := utils.DB
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}

	item := models.Item{
		Base: models.Base{ID: uid},
	}

	result := db.First(&item)

	if result.Error != nil {
		return nil, result.Error
	}

	return &item, nil
}

func GetItems() ([]*models.Item, error) {
	db := utils.DB
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var items []*models.Item

	result := db.Find(&items)

	if result.Error != nil {
		return nil, result.Error
	}

	return items, nil
}
