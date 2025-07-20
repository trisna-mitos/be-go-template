package usecase

import (
	"context"

	"go-backend-service/internal/repository"
	"go-backend-service/pkg/pb"
)

type DipanTypeUsecase interface {
	Create(ctx context.Context, dipanType *pb.DipanType) error
	GetByID(ctx context.Context, id int32) (*pb.DipanType, error)
	List(ctx context.Context, page, limit int32) ([]*pb.DipanType, int32, error)
	Update(ctx context.Context, dipanType *pb.DipanType) error
	Delete(ctx context.Context, id int32) error
}

type dipanTypeUsecase struct {
	dipanTypeRepo repository.DipanTypeRepository
}

func NewDipanTypeUsecase(dipanTypeRepo repository.DipanTypeRepository) DipanTypeUsecase {
	return &dipanTypeUsecase{
		dipanTypeRepo: dipanTypeRepo,
	}
}

func (u *dipanTypeUsecase) Create(ctx context.Context, dipanType *pb.DipanType) error {
	return u.dipanTypeRepo.Create(ctx, dipanType)
}

func (u *dipanTypeUsecase) GetByID(ctx context.Context, id int32) (*pb.DipanType, error) {
	return u.dipanTypeRepo.GetByID(ctx, id)
}

func (u *dipanTypeUsecase) List(ctx context.Context, page, limit int32) ([]*pb.DipanType, int32, error) {
	return u.dipanTypeRepo.List(ctx, page, limit)
}

func (u *dipanTypeUsecase) Update(ctx context.Context, dipanType *pb.DipanType) error {
	return u.dipanTypeRepo.Update(ctx, dipanType)
}

func (u *dipanTypeUsecase) Delete(ctx context.Context, id int32) error {
	return u.dipanTypeRepo.Delete(ctx, id)
}
