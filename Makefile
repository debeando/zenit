SHELL := /bin/bash
.PHONY: help

help: ## Show this help.
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

deps: ## Install dependencies
	go get -u github.com/go-ini/ini
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/hpcloud/tail

deps-devel:
	brew install jq

tests: ## Run tests
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic ./...

build: ## Build binary
	go build -ldflags "-s -w" -o zenit main.go

release: ## Create release
	scripts/release.sh

docker-build: ## Build docker images
	docker-compose --file docker/docker-compose.yml build

docker-up: ## Run docker-compose
	docker-compose --file docker/docker-compose.yml --project-name=zenit up --detach
	docker cp assets/schema/clickhouse/zenit.sql zenit_clickhouse_1:/root
	docker exec -i -t -u root zenit_clickhouse_1 /bin/bash -c "cat /root/zenit.sql | /usr/bin/clickhouse-client --multiquery"

docker-down: ## Down docker-compose
	docker-compose --file docker/docker-compose.yml --project-name=zenit down

docker-clickhouse: ## Enter into ClickHouse Client
	docker exec -i -t -u root zenit_clickhouse_1 /usr/bin/clickhouse-client

docker-sandbox-bash: ## Enter into sandbox container
	docker exec -i -t -u root zenit_sandbox_1 /bin/bash

docker-sandbox-build: ## Build binary and copy to container
	GOOS=linux go build -ldflags "-s -w" -o zenit main.go
	docker cp sandbox zenit_sandbox_1:/root
