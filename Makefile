gen:
	protoc --proto_path=proto --go_out=collector/proto --go_opt=paths=source_relative --go-grpc_out=collector/proto --go-grpc_opt=paths=source_relative proto/metrics.proto