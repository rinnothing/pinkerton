.PHONY: test
test:
	go test ./... -v

.PHONY: integration-test
integration-test:
	go test ./integration -v --tags integration

.PHONY: docker-integration-test
docker-integration-test:
	docker compose up -d
	EXTERNAL=TRUE go test ./integration -v --tags integration
	docker compose down

.PHONY: clean-testcache
clean-testcache:
	go clean --testcache

.PHONY: run
run:
	go run ./cmd/server

.PHONY: build
build:
	go build -o server ./cmd/server

.PHONY: build-cli
build-cli:
	go build -o pinkerton-cli ./cmd/cli

.PHONY: clean-dir
clean-dir:
	rm server pinkerton-cli
