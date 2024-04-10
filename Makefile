CWD := $(strip $(CURDIR))#win
#CWD = $(shell pwd)#lin
BUILDER = docker run -it --rm \
	-v $(CWD):/mnt \
	--add-host git.astralnalog.ru:192.168.1.137 \
	harbor.infra.yandex.astral-dev.ru/astral-edo/go/edo-golang-builder:v2.0.6
SERVICE:=package-receiver-service

lint:
	$(BUILDER) golangci-lint run --config=./golangci-lint.yml --timeout=5m

gen-envs:
	$(BUILDER) conf2env -struct Config -file config/config.go -out config/local.env

create-migration: ### make create-migration name=init_db создает миграцию для бд маркировки
	cd migrations/package_receiver_db && goose create $(name) sql && cd ../..
name ?= default_name

## make migrate-up PG_USER=edo-user PG_PASSWORD=PazzW0rD PG_DBNAME=package_receiver PG_HOST=localhost PG_PORT=5437
# TODO: миграции добавить в docker-compose.yaml, здесь использовать goose из BUILDER
migrate-up:
	cd migrations/package_receiver_db && goose postgres "user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DBNAME) host=$(PG_HOST) port=$(PG_PORT) sslmode=disable" up && cd ../..

generate-package-receiver-api:
	$(BUILDER)  protoc --proto_path proto/$(SERVICE) --proto_path=proto/vendor  \
		--go_out=pkg/$(SERVICE)/  --go_opt=paths=source_relative \
		--plugin=protoc-gen-go=/go/bin/protoc-gen-go \
		--go-grpc_out=pkg/$(SERVICE) --go-grpc_opt=paths=source_relative \
		--plugin=protoc-gen-go-grpc=/go/bin/protoc-gen-go-grpc \
		--grpc-gateway_out=pkg/$(SERVICE) --grpc-gateway_opt=paths=source_relative \
		--plugin=protoc-gen-grpc-gateway=/go/bin/protoc-gen-grpc-gateway \
		--openapiv2_out=allow_merge=true,merge_file_name=api:docs/swagger \
		--plugin=protoc-gen-openapiv2=/go/bin/protoc-gen-openapiv2 \
		proto/$(SERVICE)/receiver-service.proto
	$(BUILDER)  python3 scripts/swagger.py docs/swagger/api.swagger.json
