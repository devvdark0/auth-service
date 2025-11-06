.PHONY: generate
generate:
	protoc -I api/proto api/proto/*.proto --go_out=api/gen --go_opt=paths=source_relative --go-grpc_out=api/gen --go-grpc_opt=paths=source_relative
.PHONY: run
run:
	go run cmd/auth/main.go