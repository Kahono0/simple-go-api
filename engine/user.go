package engine

import (
	"fmt"

	"github.com/Kahono0/simple-go-api/models"
	"github.com/Kahono0/simple-go-api/utils"
)

func CreateUser(id, email, name string) error {
	db := utils.DB

	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	user := models.User{
		ID:    id,
		Email: email,
		Name:  name,
	}

	fmt.Println("=====>", user)

	result := db.FirstOrCreate(&user, "id = ?", id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetUserByID(id string) (*models.User, error) {
	db := utils.DB

	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	user := models.User{}

	result := db.First(&user, "id = ?", id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
