# Server


## Tools

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

## GRPC

After changing the proto file, run the following command to generate the go code. This will trigger the commands in gen.go

```
go generate ./...
```