package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.44

import (
	"context"
	"fmt"

	"github.com/Kahono0/simple-go-api/graph/model"
)

// CreateItem is the resolver for the createItem field.
func (r *mutationResolver) CreateItem(ctx context.Context, name string, description *string, price float64) (*model.Item, error) {
	panic(fmt.Errorf("not implemented: CreateItem - createItem"))
}

// UpdateItem is the resolver for the updateItem field.
func (r *mutationResolver) UpdateItem(ctx context.Context, id string, name *string, description *string, price *float64) (*model.Item, error) {
	panic(fmt.Errorf("not implemented: UpdateItem - updateItem"))
}

// DeleteItem is the resolver for the deleteItem field.
func (r *mutationResolver) DeleteItem(ctx context.Context, id string) (*model.Item, error) {
	panic(fmt.Errorf("not implemented: DeleteItem - deleteItem"))
}

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, itemIds []string) (*model.Order, error) {
	panic(fmt.Errorf("not implemented: CreateOrder - createOrder"))
}

// UpdateOrder is the resolver for the updateOrder field.
func (r *mutationResolver) UpdateOrder(ctx context.Context, id string, status string) (*model.Order, error) {
	panic(fmt.Errorf("not implemented: UpdateOrder - updateOrder"))
}

// DeleteOrder is the resolver for the deleteOrder field.
func (r *mutationResolver) DeleteOrder(ctx context.Context, id string) (*model.Order, error) {
	panic(fmt.Errorf("not implemented: DeleteOrder - deleteOrder"))
}

// Items is the resolver for the items field.
func (r *queryResolver) Items(ctx context.Context) ([]*model.Item, error) {
	panic(fmt.Errorf("not implemented: Items - items"))
}

// Item is the resolver for the item field.
func (r *queryResolver) Item(ctx context.Context, id string) (*model.Item, error) {
	panic(fmt.Errorf("not implemented: Item - item"))
}

// Orders is the resolver for the orders field.
func (r *queryResolver) Orders(ctx context.Context) ([]*model.Order, error) {
	panic(fmt.Errorf("not implemented: Orders - orders"))
}

// Order is the resolver for the order field.
func (r *queryResolver) Order(ctx context.Context, id string) (*model.Order, error) {
	panic(fmt.Errorf("not implemented: Order - order"))
}

// Me is the resolver for the me field.
func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	panic(fmt.Errorf("not implemented: Me - me"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
