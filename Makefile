.PHONY: protos init_rabbitmq stop_rabbitmq \
	start_rabbitmq docker_compose start_product_svc \
	start_product_client stop_docker_compose \
	build_docker_image

protos:
	protoc --proto_path=api/v1 --go-grpc_out=pkg/api/v1 \
	--go_out=pkg/api/v1 --go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative product.proto

init_rabbitmq:
	docker run -d --hostname gmd-rabbitmq \
	--name rabbit-tank -p 5672:5672 -p 15672:15672 \
	rabbitmq:3.9.7-management

start_rabbitmq:
	docker start rabbit-tank

stop_rabbitmq:
	docker stop rabbit-tank

docker_compose:
	docker-compose -f docker-compose.local.yml up -d

stop_docker_compose:
	docker-compose -f docker-compose.local.yml down

start_product_svc:
	go run cmd/product_svc/main.go

start_product_client:
	go run cmd/product_client/main.go

build_docker_image:
	docker build -t gmd-product-svc:v1.0.0 -f deployments/docker/Dockerfile . 