all:
	protoc -I/usr/local/include -I. \
		-I${HOME}/googleapis \
		--go-grpc_out=. \
		--grpc-gateway_out=logtostderr=true:${GOBIN} \
		echoserver.proto

	protoc -I . \
   --go_out . --go_opt paths=source_relative \
   --go-grpc_out . --go-grpc_opt paths=source_relative \
  --grpc-gateway_out . --grpc-gateway_opt paths=source_relative \
   ./echoserver.proto