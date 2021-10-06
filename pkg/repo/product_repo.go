package repo

import (
	"context"

	"github.com/quyenphamkhac/gmd-productsrv/pkg/entity"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]entity.Product, error)
	FindById(ctx context.Context, id int) (*entity.Product, error)
}

type ProductPublisher interface {
	Publish(body []byte, contentType string) error
}

type ProductConsumer interface {
	Subscribe(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error
}
