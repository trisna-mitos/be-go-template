package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-backend-service/pkg/pb"
)

type DipanTypeRepository interface {
	Create(ctx context.Context, dipanType *pb.DipanType) error
	GetByID(ctx context.Context, id int32) (*pb.DipanType, error)
	List(ctx context.Context, page, limit int32) ([]*pb.DipanType, int32, error)
	Update(ctx context.Context, dipanType *pb.DipanType) error
	Delete(ctx context.Context, id int32) error
}

type postgresDipanTypeRepository struct {
	db *sql.DB
}

func NewPostgresDipanTypeRepository(db *sql.DB) DipanTypeRepository {
	return &postgresDipanTypeRepository{db: db}
}

func (r *postgresDipanTypeRepository) Create(ctx context.Context, dipanType *pb.DipanType) error {
	query := `
        INSERT INTO dipan_types (nama_type)
        VALUES ($1)
    `
	_, err := r.db.ExecContext(ctx, query, dipanType.NamaType)
	return err
}

func (r *postgresDipanTypeRepository) GetByID(ctx context.Context, id int32) (*pb.DipanType, error) {
	query := `
        SELECT id, nama_type
        FROM dipan_types
        WHERE id = $1
    `
	dipanType := &pb.DipanType{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dipanType.Id,
		&dipanType.NamaType,
	)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return dipanType, nil
}

func (r *postgresDipanTypeRepository) List(ctx context.Context, page, limit int32) ([]*pb.DipanType, int32, error) {
	offset := (page - 1) * limit
	query := `
		SELECT id, nama_type
		FROM dipan_types
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var dipanTypes []*pb.DipanType
	for rows.Next() {
		dipanType := &pb.DipanType{}
		err := rows.Scan(
			&dipanType.Id,
			&dipanType.NamaType,
		)
		if err != nil {
			return nil, 0, err
		}
		dipanTypes = append(dipanTypes, dipanType)
	}

	// Get total count
	var total int32
	err = r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM dipan_types").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return dipanTypes, total, nil
}

func (r *postgresDipanTypeRepository) Update(ctx context.Context, dipanType *pb.DipanType) error {
	return fmt.Errorf("not implemented")
}
func (r *postgresDipanTypeRepository) Delete(ctx context.Context, id int32) error {
	return fmt.Errorf("not implemented")
}
