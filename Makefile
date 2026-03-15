.PHONY: test
test:
	go test ./... -v

.PHONY: run
run:
	go run ./cmd/server
