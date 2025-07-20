package repository

import (
	"context"
	"database/sql"
	"fmt"

	"go-backend-service/internal/domain"
)

type postgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) domain.ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
        INSERT INTO products (id, name, description, price, stock, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := r.db.ExecContext(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

func (r *postgresProductRepository) GetByID(ctx context.Context, id string) (*domain.Product, error) {
	query := `
        SELECT id, name, description, price, stock, created_at, updated_at
        FROM products
        WHERE id = $1
    `

	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found")
	}

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *postgresProductRepository) List(ctx context.Context, page, limit int32) ([]*domain.Product, int32, error) {
	offset := (page - 1) * limit

	query := `
        SELECT id, name, description, price, stock, created_at, updated_at
        FROM products
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		products = append(products, product)
	}

	// Get total count
	var total int32
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM products").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
