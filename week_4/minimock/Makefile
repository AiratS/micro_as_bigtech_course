include .env

LOCAL_BIN=${CURDIR}/bin
MIGRATION_DSN="host=localhost port=$(PG_PORT) dbname=$(PG_DATABASE_NAME) user=$(PG_USER) password=$(PG_PASSWORD) sslmode=disable"

install-tools:
	GOBIN=${LOCAL_BIN} go install google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN=${LOCAL_BIN} go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	GOBIN=${LOCAL_BIN} go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=${LOCAL_BIN} go install github.com/gojuno/minimock/v3/cmd/minimock@latest


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

test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/AiratS/micro_as_bigtech_course/week_3/internal/service/...,github.com/AiratS/micro_as_bigtech_course/week_3/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "./coverage.out" .gitignore || echo "./coverage.out" .gitignore