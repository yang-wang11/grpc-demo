# 生成的go文件所在目录是../pb, 这个目录是相对于执行protoc命令的目录
# 最终的go文件所在路径，需要再次添加在protoc文件中的go_package配置
protoc  --go_out=.. --go-grpc_out=.. *.proto
