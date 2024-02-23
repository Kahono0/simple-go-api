package graph_test

import (
	"context"
	"testing"

	"github.com/Kahono0/simple-go-api/engine"
	"github.com/Kahono0/simple-go-api/graph/model"
	"github.com/Kahono0/simple-go-api/models"
	"github.com/Kahono0/simple-go-api/utils"
	"github.com/joho/godotenv"
)

type Resolver struct{}
type mutationResolver struct{ *Resolver }

var user = models.User{
	ID:    "1",
	Email: "test@test.com",
	Name:  "Test User",
}

// CreateItem is the resolver for the createItem field.
func (r *mutationResolver) CreateItem(ctx context.Context, name string, description *string, price float64) (*model.Item, error) {
	item, err := engine.CreateItem(name, description, price)
	if err != nil {
		return nil, err
	}

	return &model.Item{
		ID:          item.ID.String(),
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
	}, nil
}

// UpdateItem is the resolver for the updateItem field.
func (r *mutationResolver) UpdateItem(ctx context.Context, id string, name *string, description *string, price *float64) (*model.Item, error) {
	item, err := engine.UpdateItem(id, name, description, price)
	if err != nil {
		return nil, err
	}

	return &model.Item{
		ID:          item.ID.String(),
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
	}, nil
}

// DeleteItem is the resolver for the deleteItem field.
func (r *mutationResolver) DeleteItem(ctx context.Context, id string) (*model.Item, error) {
	return nil, engine.DeleteItem(id)
}

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, input model.OrderInput) (*model.Order, error) {
	user := ctx.Value("user").(*models.User)
	order, err := engine.CreateOrder(user.ID, input.Items, input.Contact)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		ID:        order.ID.String(),
		Items:     nil,
		Total:     order.Total,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.String(),
	}, nil
}

// UpdateOrder is the resolver for the updateOrder field.
func (r *mutationResolver) UpdateOrder(ctx context.Context, id string, status string) (*model.Order, error) {
	order, err := engine.UpdateOrder(id, status)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		ID:        order.ID.String(),
		Items:     nil,
		Total:     order.Total,
		Status:    order.Status,
		CreatedAt: order.CreatedAt.String(),
	}, nil
}

func TestMutationResolver_AllActions(t *testing.T) {
	// Setup
	// Create a test engine with a mock implementation or use a real one with a test database
	_ = godotenv.Load("../.env")
	utils.ConnectDB()

	resolver := &mutationResolver{&Resolver{}}

	ctx := context.WithValue(context.Background(), "user", &user)

	// Test CreateItem
	createdItem, err := resolver.CreateItem(ctx, "TestItem", nil, 10.0)
	if err != nil {
		t.Fatalf("CreateItem failed: %v", err)
	}
	if createdItem.Name != "TestItem" {
		t.Fatalf("CreateItem failed: %v", createdItem)
	}

	// Test UpdateItem
	updatedItem, err := resolver.UpdateItem(ctx, createdItem.ID, nil, nil, nil)
	if err != nil {
		t.Fatalf("UpdateItem failed: %v", err)
	}
	if updatedItem.Name != "TestItem" {
		t.Fatalf("UpdateItem failed: %v", updatedItem)
	}

	// Test DeleteItem
	_, err = resolver.DeleteItem(ctx, createdItem.ID)
	if err != nil {
		t.Fatalf("DeleteItem failed: %v", err)
	}

	// Test CreateOrder
	createdOrder, err := resolver.CreateOrder(ctx, model.OrderInput{
		Items:   nil,
		Contact: "",
	})

	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	if createdOrder.Status != models.OrderStatusPending {
		t.Fatalf("CreateOrder failed: %v", createdOrder)
	}

	// Test UpdateOrder
	updatedOrder, err := resolver.UpdateOrder(ctx, createdOrder.ID, models.OrderStatusCompleted)
	if err != nil {
		t.Fatalf("UpdateOrder failed: %v", err)
	}

	if updatedOrder.Status != models.OrderStatusCompleted {
		t.Fatalf("UpdateOrder failed: %v", updatedOrder)
	}

	// Teardown
	// Delete the created order and item
	_ = engine.DeleteOrder(createdOrder.ID)

	// Delete the created item
	_ = engine.DeleteItem(createdItem.ID)
}
