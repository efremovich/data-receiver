LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.2
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@v0.10.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.19.1
	GOBIN=$(LOCAL_BIN) go install github.com/rakyll/statik@v0.1.7
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.18.0
	GOBIN=$(LOCAL_BIN) go install github.com/gojuno/minimock/v3/cmd/minimock@latest

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

lint:
	golangci-lint run --config=./golangci-lint.yml --timeout=5m

gen-envs:
	conf2env -struct Config -file config/config.go -out config/local.env

create-migration: ### make create-migration name=init_db создает миграцию для бд маркировки
	cd migrations/package_receiver_db && goose create $(name) sql && cd ../..
name ?= default_name

migrate-up:
	cd migrations/package_receiver_db && goose postgres "user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DBNAME) host=$(PG_HOST) port=$(PG_PORT) sslmode=disable" up && cd ../..

generate-data-receiver-api:
	mkdir -p pkg/data-receiver-service
	protoc --proto_path proto/data-receiver-service --proto_path=vendor.protogen  \
		--go_out=pkg/data-receiver-service/  --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=bin/protoc-gen-go \
		--go-grpc_out=pkg/data-receiver-service --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
		--grpc-gateway_out=pkg/data-receiver-service --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=bin/protoc-gen-grpc-gateway \
		--openapiv2_out=allow_merge=true,merge_file_name=api:docs/swagger \
		--plugin=protoc-gen-openapiv2=bin/protoc-gen-openapiv2 \
		proto/data-receiver-service/receiver-service.proto

vendor-proto:
		@if [ ! -d vendor.protogen/validate ]; then \
			mkdir -p vendor.protogen/validate &&\
			git clone https://github.com/envoyproxy/protoc-gen-validate vendor.protogen/protoc-gen-validate &&\
			mv vendor.protogen/protoc-gen-validate/validate/*.proto vendor.protogen/validate &&\
			rm -rf vendor.protogen/protoc-gen-validate ;\
		fi
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi
		@if [ ! -d vendor.protogen/protoc-gen-openapiv2 ]; then \
			mkdir -p vendor.protogen/protoc-gen-openapiv2/options &&\
			git clone https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/openapiv2 &&\
			mv vendor.protogen/openapiv2/protoc-gen-openapiv2/options/*.proto vendor.protogen/protoc-gen-openapiv2/options &&\
			rm -rf vendor.protogen/openapiv2 ;\
		fi
