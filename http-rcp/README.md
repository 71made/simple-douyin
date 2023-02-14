# http-rpc
根据 proto IDL 生成 GRPC 代码，终端命令：
``` shell
cd http_rpc/ 
protoc --proto_path=. --go_out=../../  idl/*.proto
protoc --proto_path=. --go-grpc_out=../../ idl/*.proto
```

