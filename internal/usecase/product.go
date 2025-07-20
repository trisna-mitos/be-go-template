package usecase

import (
	"context"
	"time"

	"go-backend-service/internal/domain"

	"github.com/google/uuid"
)

type productUsecase struct {
	productRepo domain.ProductRepository
}

func NewProductUsecase(productRepo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
	}
}

func (u *productUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	now := time.Now()
	product.ID = uuid.New().String()
	product.CreatedAt = now
	product.UpdatedAt = now

	return u.productRepo.Create(ctx, product)
}

func (u *productUsecase) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return u.productRepo.GetByID(ctx, id)
}

func (u *productUsecase) ListProducts(ctx context.Context, page, limit int32) ([]*domain.Product, int32, error) {
	return u.productRepo.List(ctx, page, limit)
}
