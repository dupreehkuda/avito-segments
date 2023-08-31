.PHONY: test
test:
	go test -covermode=atomic -v -race ./internal/...

.PHONY: gen
gen:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run --color=always

.PHONY: deploy
deploy:
	docker-compose -f docker-compose.prod.yml pull
	docker-compose -f docker-compose.prod.yml down
	docker-compose -f docker-compose.prod.yml up -d

.PHONY:
compose-up:
	docker-compose -f docker-compose.dev.yml up

.PHONY:
compose-down:
	docker-compose -f docker-compose.dev.yml down