package domain

import (
	"context"
	"time"
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int32     `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id string) (*Product, error)
	List(ctx context.Context, page, limit int32) ([]*Product, int32, error)
}

type ProductUsecase interface {
	CreateProduct(ctx context.Context, product *Product) error
	GetProduct(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, page, limit int32) ([]*Product, int32, error)
}
