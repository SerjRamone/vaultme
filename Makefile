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

.PHONY: clean-data
clean-data:
	rm -rf ./deployments/db/data/

.PHONY: build-server
build-server:
	go build -o cmd/server/server cmd/server/*.go

.PHONY: start-server
start-server: build-server
	env $$(cat .env | xargs) ./cmd/server/server
