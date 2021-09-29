install:
	make build;
	go install -v

fmt:
	go fmt ./...

certs:
	openssl genrsa \
		-out ./certs/localhost.key \
		2048
	openssl req \
		-nodes	\
		-new -x509 \
		-sha256 \
		-config ./certs/cert.conf \
		-extensions 'req_ext' \
		-key ./certs/localhost.key \
		-out ./certs/localhost.cert \
		-days 3650 \
		-subj /CN=localhost

grpc:
	protoc --proto_path=proto proto/*.proto  --go_out=:. --go-grpc_out=:.

build:
	make grpc
	make certs
	go build -o grpc-upload main.go

run:
	make build
	docker build --tag image .
	docker run -it -p 8080:8080 image

.PHONY: fmt install grpc certs
