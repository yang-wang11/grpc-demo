echo "生成 rpc server 代码"

OUT=../server
protoc --go_out=${OUT} \
--go-grpc_out=${OUT} \
--go-grpc_opt=require_unimplemented_servers=false \
server.proto


echo "生成 rpc client 代码"

OUT=../client
protoc --go_out=${OUT} \
--go-grpc_out=${OUT} \
--go-grpc_opt=require_unimplemented_servers=false \
server.proto



