package repo

import "github.com/quyenphamkhac/gmd-productsrv/pkg/entity"

type ProductRepository interface {
	FindAll() ([]entity.Product, error)
	FindById(id int) (*entity.Product, error)
}
