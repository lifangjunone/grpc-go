## proto go code generate command
```sh
> protoc -I=. --go_out=. --go_opt=module="grpc-go/simple" --go-grpc_out=. --go-grpc_opt=module="grpc-go/simple" hello.proto
```