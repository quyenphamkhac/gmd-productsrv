package handler

import (
	"context"

	pb "github.com/quyenphamkhac/gmd-productsrv/pkg/api/v1"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/entity"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type productService struct {
	pb.UnimplementedProductSrvServer
	usecase usecase.ProductUsecase
}

func NewProductService(u usecase.ProductUsecase) *productService {
	return &productService{
		usecase: u,
	}
}

func (s *productService) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllReponse, error) {
	products, err := s.usecase.FindAll()
	if err != nil {
		return nil, err
	}
	res := &pb.GetAllReponse{
		Data:          marshalList(products),
		ResultPerPage: 10,
		CurrentPage:   1,
	}
	return res, nil
}

func (s *productService) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.GetByIdResponse, error) {
	product, err := s.usecase.FindById(int(req.GetId()))
	if err != nil {
		return nil, err
	}
	res := &pb.GetByIdResponse{
		Data: marshalItem(product),
	}
	return res, nil
}

func marshalList(products []entity.Product) []*pb.Product {
	res := make([]*pb.Product, len(products))
	for i, p := range products {
		res[i] = &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Sku:         p.Sku,
			CreatedAt:   timestamppb.New(p.CreatedAt),
			UpdatedAt:   timestamppb.New(p.UpdatedAt),
		}
	}
	return res
}

func marshalItem(p *entity.Product) *pb.Product {
	return &pb.Product{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Sku:         p.Sku,
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}
