.PHONY: server client

setup:
	@docker-compose up -d

server: 
	@go run server/*.go

client:
	@go run client/*.go

gen-proto: 
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	userproto/user.proto