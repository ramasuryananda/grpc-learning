build-proto:
	 protoc --proto_path=proto proto/*.proto --go_out=./pb --go_opt=paths=source_relative --go-grpc_out=./pb --go-grpc_opt=paths=source_relative

run-server:
	go build .
	./grpc-learning

run-client:
	go run ./client/client.go