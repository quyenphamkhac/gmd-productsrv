.PHONY: protos

protos:
	protoc --proto_path=protos --go-grpc_out=protos --go_out=protos --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative product.proto