PROJECT_DIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: start-db
start-db:
	docker run --rm \
		--name=postgresql \
		-v $(PROJECT_DIR)/deployments/db/:/docker-entrypoint-initdb.d \
		-v $(PROJECT_DIR)/deployments/db/data/:/var/lib/postgresql/data \
		-e POSTGRES_USER=vaultme \
		-e POSTGRES_PASSWORD=vaultme \
		-d \
		-p 5432:5432 \
		postgres:16.3

.PHONY: stop-db
stop-db:
	docker stop postgresql

.PHONY: clean-pgdata
clean-pgdata:
	rm -rf ./deployments/db/data/

.PHONY: build-server
build-server:
	go build -o cmd/server/server cmd/server/*.go

.PHONY: start-server
start-server: build-server
	env $$(cat .env | xargs) ./cmd/server/server

.PHONY: build-client
build-client:
	go build -o cmd/client/client cmd/client/*.go

.PHONY: start-client
start-client: build-client
	env $$(cat .env | xargs) ./cmd/client/client

.PHONY: proto-generate
proto-generate:
	mkdir -p pkg/vaultme_v1
	protoc --proto_path api/v1 \
        --go_out=./pkg/vaultme_v1 --go_opt=paths=source_relative \
        --go-grpc_out=./pkg/vaultme_v1 --go-grpc_opt=paths=source_relative \
        api/v1/user.proto
