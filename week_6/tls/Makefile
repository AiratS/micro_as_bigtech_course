generate:
	mkdir -p pkg/note_v1
	protoc \
    --proto_path api/note_v1 \
    --go_out=pkg/note_v1 \
    --go_opt=paths=source_relative \
    --go-grpc_out=pkg/note_v1 \
    --go-grpc_opt=paths=source_relative \
    api/note_v1/note.proto

gen-cert:
	openssl genrsa -out ca.key 4096
	openssl req -new -x509 -key ca.key -sha256 -subj "/C=US/ST=NJ/O=CA, Inc." -days 365 -out ca.cert
	openssl genrsa -out service.key 4096
	openssl req -new -key service.key -out service.csr -config certificate.conf
	openssl x509 -req -in service.csr -CA ca.cert -CAkey ca.key -CAcreateserial \
    		-out service.pem -days 365 -sha256 -extfile certificate.conf -extensions req_ext