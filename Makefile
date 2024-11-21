go_out_path = "./proto/gen"
go_grpc_path = "./proto/gen"

.PHONY: generate

generate:
	@protoc --go_out=${go_out_path} --go_opt=paths=source_relative \
		--go-grpc_out=${go_grpc_path} \
		--go-grpc_opt=paths=source_relative \
		proto/*.proto
