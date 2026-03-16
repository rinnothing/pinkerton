.PHONY: test
test:
	go test ./... -v

.PHONY: integration-test
integration-test:
	go test ./integration -v --tags integration

.PHONY: clean-testcache
clean-testcache:
	go clean --testcache

.PHONY: run
run:
	go run ./cmd/server

# .PHONY: cli
# cli:
# 	go run ./cmd/cli
