package interactor

import (
	"github.com/quyenphamkhac/gmd-productsrv/pkg/entity"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/repo"
)

type ProductUsecase interface {
	FindAll() ([]entity.Product, error)
	FindById() (*entity.Product, error)
}

type productUsecase struct {
	repo *repo.ProductRepository
}

func NewProductUseCase(r *repo.ProductRepository) (*productUsecase, error) {
	return &productUsecase{
		repo: r,
	}, nil
}

func (p *productUsecase) FindAll() ([]entity.Product, error) {
	return []entity.Product{}, nil
}

func (p *productUsecase) FindById() (*entity.Product, error) {
	return &entity.Product{}, nil
}
