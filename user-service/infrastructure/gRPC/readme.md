# Build Proto

```
cd /user-service/infrastructure/gRPC
protoc --go_out=./user --go_opt=paths=source_relative ./proto/*.proto
protoc --go-grpc_out=./user --go-grpc_opt=paths=source_relative ./proto/*.proto
```