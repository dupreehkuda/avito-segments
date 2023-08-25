.PHONY: test
test:
	go test -covermode=atomic -v -race ./internal/...

.PHONY: gen
gen:
	go generate ./...