LOCAL_BIN=${CURDIR}/bin

install-tools:
	make get-deps
	GOBIN=${LOCAL_BIN} go install google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN=${LOCAL_BIN} go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

generate:
	make generate-access-api
	make generate-auth-api

generate-access-api:
	mkdir -p pkg/access_v1
	protoc \
    --proto_path api/access_v1 \
    --go_out=pkg/access_v1 \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
    --go_opt=paths=source_relative \
    --go-grpc_out=pkg/access_v1 \
    --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
    api/access_v1/access.proto

generate-auth-api:
	mkdir -p pkg/auth_v1
	protoc \
    --proto_path api/auth_v1 \
    --go_out=pkg/auth_v1 \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
    --go_opt=paths=source_relative \
    --go-grpc_out=pkg/auth_v1 \
    --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
    api/auth_v1/auth.proto