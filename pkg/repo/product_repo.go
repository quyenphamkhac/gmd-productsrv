package repo

import "github.com/quyenphamkhac/gmd-productsrv/pkg/entity"

type ProductRepository interface {
	FindAll() ([]entity.Product, error)
	FindById(id int) (*entity.Product, error)
}

type ProductPublisher interface {
	Publish(body []byte, contentType string) error
}

type ProductConsumer interface {
	Subscribe(workerPoolSize int, exchange, queueName, bindingKey, consumerTag string) error
}
