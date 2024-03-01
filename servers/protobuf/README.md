## Protocol buffer

Install compiler:
```bash
apt install -y protobuf-compiler
```

Install plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Generate server gRPC (repeat when .proto file is modified):
```bash
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative servers/protobuf/name_of_the_project/server.proto
```