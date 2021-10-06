package handler

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/quyenphamkhac/gmd-productsrv/internal/logger"
	pb "github.com/quyenphamkhac/gmd-productsrv/pkg/api/v1"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/entity"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/usecase"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type productService struct {
	pb.UnimplementedProductSrvServer
	usecase usecase.ProductUsecase
	logger  logger.Logger
}

func NewProductService(u usecase.ProductUsecase, logger logger.Logger) *productService {
	return &productService{
		usecase: u,
		logger:  logger,
	}
}

func (s *productService) GetAll(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllReponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUsecase.FindAll")
	defer span.Finish()

	s.logger.Info("start find all api", logger.LogFields{
		"request_id": "123",
		"user_id":    "1234",
		"user_ip":    "1234",
	})
	products, err := s.usecase.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	res := &pb.GetAllReponse{
		Data:          marshalList(products),
		ResultPerPage: 10,
		CurrentPage:   1,
	}
	s.logger.Info("log products data", logger.LogFields{
		"request_id":    "123",
		"user_id":       "1234",
		"user_ip":       "1234",
		"products_data": products,
	})
	return res, nil
}

func (s *productService) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.GetByIdResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUsecase.FindById")
	defer span.Finish()

	s.logger.Info("start find product by id api", logger.LogFields{
		"request_id": "123",
		"user_id":    "1234",
		"user_ip":    "1234",
	})
	product, err := s.usecase.FindById(ctx, int(req.GetId()))
	if err != nil {
		return nil, err
	}
	res := &pb.GetByIdResponse{
		Data: marshalItem(product),
	}
	s.logger.Info("log product data", logger.LogFields{
		"request_id":   "123",
		"user_id":      "1234",
		"user_ip":      "1234",
		"product_data": product,
	})
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
