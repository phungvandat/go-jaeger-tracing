.PHONY: server1 client server2

setup:
	@docker-compose up -d

server1: 
	@go run server1/*.go

server2:
	@go run server2/*.go

client:
	@go run client/*.go

gen-proto: 
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	userproto/user.proto