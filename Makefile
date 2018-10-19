SHELL := /bin/bash
.PHONY: help

help: ## Show this help.
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

deps: ## Install dependencies
	go get -u github.com/go-sql-driver/mysql
	go get -u github.com/go-yaml/yaml
	go get -u github.com/hpcloud/tail
	go get -u github.com/kardianos/service
	go get -u github.com/shirou/gopsutil
	go get -u golang.org/x/sys/unix

deps-devel:
	brew install jq

tests: ## Run tests
	go generate ./...
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic ./...

build: ## Build binary
	go build -ldflags "-s -w" -o zenit main.go

release: ## Create release
	scripts/release.sh

docker-build: ## Build docker images
	docker-compose --file docker/docker-compose.yml build

docker-up: ## Run docker-compose
	docker-compose --file docker/docker-compose.yml --project-name=zenit up

docker-ps: ## Show status for all containers
	docker-compose --file docker/docker-compose.yml --project-name=zenit ps

docker-down: ## Down docker-compose
	docker-compose --file docker/docker-compose.yml --project-name=zenit down

docker-clickhouse: ## Enter into ClickHouse Client
	docker exec -i -t -u root zenit_clickhouse /usr/bin/clickhouse-client

docker-mysql-primary: ## Enter in MySQL Primary Console
	docker exec -i -t -u root zenit_percona_server_primary /usr/bin/mysql

docker-mysql-primary-bash: ## Enter in MySQL Primary bash console
	docker exec -i -t -u root zenit_percona_server_primary /bin/bash

docker-mysql-secondary: ## Enter in MySQL Secondary Console
	docker exec -i -t -u root zenit_percona_server_secondary /usr/bin/mysql

docker-mysql-secondary-bash: ## Enter in MySQL Secondary bash console
	docker exec -i -t -u root zenit_percona_server_secondary /bin/bash

docker-proxysql: ## Enter in ProxySQL Console
	docker exec -i -t -u root zenit_proxysql /usr/bin/mysql --socket=/tmp/proxysql_admin.sock -u proxysql -padmin  --prompt='ProxySQLAdmin> '

docker-proxysql-bash: ## Enter in ProxySQL bash console
	docker exec -i -t -u root zenit_proxysql /bin/bash

docker-zenit-build: ## Build binary and copy to container
	GOOS=linux go build -ldflags "-s -w" -o zenit main.go
	docker cp zenit zenit_percona_server_primary:/usr/bin/
	docker cp zenit zenit_percona_server_secondary:/usr/bin/
	docker cp zenit zenit_proxysql:/usr/bin/
