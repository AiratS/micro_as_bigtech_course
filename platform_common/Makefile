LOCAL_BIN:=$(CURDIR)/bin

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.8.0

lint:
	${LOCAL_BIN}/golangci-lint run --config .golangci.pipeline.yaml