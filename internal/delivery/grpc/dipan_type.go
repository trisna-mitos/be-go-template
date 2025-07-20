package grpc

import (
	"context"
	"go-backend-service/internal/usecase"
	"go-backend-service/pkg/pb"
)

type DipanTypeHandler struct {
	pb.UnimplementedDipanTypeServiceServer
	uc usecase.DipanTypeUsecase
}

func NewDipanTypeHandler(uc usecase.DipanTypeUsecase) *DipanTypeHandler {
	return &DipanTypeHandler{uc: uc}
}

func (h *DipanTypeHandler) CreateDipanType(ctx context.Context, req *pb.CreateDipanTypeRequest) (*pb.DipanType, error) {
	d := &pb.DipanType{NamaType: req.NamaType}
	if err := h.uc.Create(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (h *DipanTypeHandler) GetDipanType(ctx context.Context, req *pb.GetDipanTypeRequest) (*pb.DipanType, error) {
	return h.uc.GetByID(ctx, req.Id)
}

func (h *DipanTypeHandler) ListDipanTypes(ctx context.Context, req *pb.ListDipanTypesRequest) (*pb.ListDipanTypesResponse, error) {
	items, total, err := h.uc.List(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, err
	}
	return &pb.ListDipanTypesResponse{
		DipanTypes: items,
		Total:      total,
	}, nil
}

func (h *DipanTypeHandler) UpdateDipanType(ctx context.Context, req *pb.UpdateDipanTypeRequest) (*pb.DipanType, error) {
	d := &pb.DipanType{Id: req.Id, NamaType: req.NamaType}
	if err := h.uc.Update(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (h *DipanTypeHandler) DeleteDipanType(ctx context.Context, req *pb.DeleteDipanTypeRequest) (*pb.DeleteDipanTypeResponse, error) {
	if err := h.uc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &pb.DeleteDipanTypeResponse{}, nil
}
