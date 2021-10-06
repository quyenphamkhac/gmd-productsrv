package adapter

import (
	"context"
	"time"

	"github.com/quyenphamkhac/gmd-productsrv/pkg/entity"
)

type mockAdapter struct{}

func NewMockAdaper() *mockAdapter {
	return &mockAdapter{}
}

func (m *mockAdapter) FindAll(ctx context.Context) ([]entity.Product, error) {
	var products []entity.Product
	products = append(products, entity.Product{
		Id:          1,
		Name:        "Glass",
		Description: "GlassOn",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	products = append(products, entity.Product{
		Id:          2,
		Name:        "Tshirt",
		Description: "Marvel Tshirt",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	return products, nil
}

func (m *mockAdapter) FindById(ctx context.Context, id int) (*entity.Product, error) {
	return &entity.Product{
		Id:          int32(id),
		Name:        "Just a product",
		Description: "Make with love",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
