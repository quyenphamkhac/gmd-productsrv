# ==============================================================================
# Generate gprc stubs

protos:
	protoc --proto_path=api/v1 --go-grpc_out=pkg/api/v1 \
	--go_out=pkg/api/v1 --go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative product.proto

# ==============================================================================
# Docker compose environment
docker_compose_up:
	docker-compose -f docker-compose.local.yml up -d

docker_compose_down:
	docker-compose -f docker-compose.local.yml down

# ==============================================================================
# Product service cmd
run_svc:
	go run cmd/product_svc/main.go

run_client:
	go run cmd/product_client/main.go

# ==============================================================================
# Build docker image
build_image:
	docker build -t gmd-product-svc:v1.0.0 -f deployments/docker/Dockerfile . 

# ==============================================================================
# Go migrate postgresql
DB_NAME = product_svc_db

force:
	migrate -database postgres://postgres:postgres@localhost:5432/$(DB_NAME)?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://postgres:postgres@localhost:5432/$(DB_NAME)?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://postgres:postgres@localhost:5432/$(DB_NAME)?sslmode=disable -path migrations up 1

migrate_down:
	migrate -database postgres://postgres:postgres@localhost:5432/$(DB_NAME)?sslmode=disable -path migrations down 1