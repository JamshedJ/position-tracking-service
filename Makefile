generate:
	protoc -I protos \
    	protos/proto/pts/pts.proto \
    	--go_out=protos/gen/pts \
    	--go_opt=paths=source_relative \
    	--go-grpc_out=protos/gen/pts \
    	--go-grpc_opt=paths=source_relative \
    	--plugin=/Users/jamshed/go/bin/protoc-gen-go-grpc


run-local:
	go run cmd/pts/main.go --config=config/local.yaml