LOCAL_BIN:=$(CURDIR)/backend/bin


install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

docker-build:
	docker compose up --build -d

docker-run:
	docker compose up -d