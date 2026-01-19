include .env

LOCAL_BIN=${CURDIR}/bin
MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-tools:
	GOBIN=${LOCAL_BIN} go install google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN=${LOCAL_BIN} go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	GOBIN=${LOCAL_BIN} go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.20.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate-grpc:
	mkdir -p pkg/note_v1
	protoc \
    --proto_path api/note_v1 \
    --go_out=pkg/note_v1 \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
    --go_opt=paths=source_relative \
    --go-grpc_out=pkg/note_v1 \
    --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
    api/note_v1/note.proto

goose-up:
	${LOCAL_BIN}/goose -dir ${GOOSE_MIGRATION_DIR} postgres ${MIGRATION_DSN} up -v

goose-down:
	${LOCAL_BIN}/goose -dir ${GOOSE_MIGRATION_DIR} postgres ${MIGRATION_DSN} down -v

vendor-proto:
	@if [ ! -d vendor.protogen/google ]; then \
		git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
		mkdir -p vendor.protogen/google/ && \
		mv  vendor.protogen/googleapis/google/api vendor.protogen/google && \
		rm -rf vendor.protogen/googleapis ;\
	fi