package usecase

import (
	"context"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/entity"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/repo"
)

type ProductUsecase interface {
	FindAll(ctx context.Context) ([]entity.Product, error)
	FindById(ctx context.Context, id int) (*entity.Product, error)
}

type productUsecase struct {
	repo repo.ProductRepository
}

func NewProductUseCase(r repo.ProductRepository) *productUsecase {
	return &productUsecase{
		repo: r,
	}
}

func (p *productUsecase) FindAll(ctx context.Context) ([]entity.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUsecase.FindAll")
	defer span.Finish()

	products, errror := p.repo.FindAll(ctx)
	return products, errror
}

func (p *productUsecase) FindById(ctx context.Context, id int) (*entity.Product, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductUsecase.FindById")
	defer span.Finish()

	product, error := p.repo.FindById(ctx, id)

	span.LogFields(log.String("product_id", strconv.Itoa(id)))
	return product, error
}
