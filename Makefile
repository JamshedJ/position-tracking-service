generate:
	protoc -I protos \
    	protos/proto/pts/pts.proto \
    	--go_out=protos/gen/pts \
    	--go_opt=paths=source_relative \
    	--go-grpc_out=protos/gen/pts \
    	--go-grpc_opt=paths=source_relative \
    	--plugin=/Users/jamshed/go/bin/protoc-gen-go-grpc

mongo:
	mongod --dbpath /Users/jamshed/Desktop/position-tracking-service/data/db

run-local:
	go run cmd/pts/main.go --config=config/local.yaml

grpc-curl:
	grpcurl -plaintext -d '{"client_id":"1234","latitude":37.7749,"longitude":-122.4194,"radius":1000}' localhost:50051 position.PositionTracker.GetNearbyPositions