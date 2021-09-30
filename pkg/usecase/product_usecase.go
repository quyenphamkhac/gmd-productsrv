package usecase

import (
	"github.com/quyenphamkhac/gmd-productsrv/pkg/entity"
	"github.com/quyenphamkhac/gmd-productsrv/pkg/repo"
)

type ProductUsecase interface {
	FindAll() ([]entity.Product, error)
	FindById(id int) (*entity.Product, error)
}

type productUsecase struct {
	repo repo.ProductRepository
}

func NewProductUseCase(r repo.ProductRepository) *productUsecase {
	return &productUsecase{
		repo: r,
	}
}

func (p *productUsecase) FindAll() ([]entity.Product, error) {
	products, errror := p.repo.FindAll()
	return products, errror
}

func (p *productUsecase) FindById(id int) (*entity.Product, error) {
	product, error := p.repo.FindById(id)
	return product, error
}
