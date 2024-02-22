package engine

import (
	"fmt"

	"github.com/Kahono0/simple-go-api/models"
	"github.com/Kahono0/simple-go-api/utils"
	"github.com/google/uuid"
)

func getTotal(items []*models.Item) float64 {
	var total float64
	for _, item := range items {
		total += item.Price
	}
	return total
}

func CreateOrder(userID string, items []string, contact string) (*models.Order, error) {
	db := utils.DB
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	//get items
	var itemsDB []*models.Item
	result := db.Find(&itemsDB, "id IN ?", items)
	if result.Error != nil {
		return nil, result.Error
	}

	//create order
	order := models.Order{
		UserID: userID,
		Items:  itemsDB,
		Total:  getTotal(itemsDB),
		Status: models.OrderStatusPending,
	}

	result = db.Create(&order)
	if result.Error != nil {
		return nil, result.Error
	}

	return &order, nil
}

func UpdateOrder(id string, status string) (*models.Order, error) {
	db := utils.DB
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}

	order := models.Order{
		Base: models.Base{ID: uid},
	}

	result := db.First(&order)
	if result.Error != nil {
		return nil, result.Error
	}

	order.Status = status

	result = db.Save(&order)
	if result.Error != nil {
		return nil, result.Error
	}

	return &order, nil
}

func DeleteOrder(id string) error {
	db := utils.DB
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid id")
	}

	order := models.Order{
		Base: models.Base{ID: uid},
	}

	result := db.Delete(&order)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetOrder(id string) (*models.Order, error) {
	db := utils.DB
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id")
	}

	order := models.Order{
		Base: models.Base{ID: uid},
	}

	result := db.Preload("Items").First(&order)
	if result.Error != nil {
		return nil, result.Error
	}

	return &order, nil
}

func GetOrders() ([]*models.Order, error) {
	db := utils.DB
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var orders []*models.Order
	result := db.Preload("Items").Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}

func GetOrdersByUserID(userID string) ([]*models.Order, error) {
	db := utils.DB
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	var orders []*models.Order
	result := db.Preload("Items").Find(&orders, "user_id = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}

	return orders, nil
}
