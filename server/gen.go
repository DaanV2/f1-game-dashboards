package main

//go:generate protoc --proto_path=../shared/proto/v1 ../shared/proto/v1/*.proto --go_out=./api/grpc --go-grpc_out=./api/grpc
