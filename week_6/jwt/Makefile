LOCAL_BIN=${CURDIR}/bin

install-tools:
	make get-deps
	GOBIN=${LOCAL_BIN} go install google.golang.org/protobuf/cmd/protoc-gen-go
	GOBIN=${LOCAL_BIN} go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc