.PHONY: protos

protos:
	protoc --proto_path=protos --go_out=protos --go_opt=paths=source_relative product.proto