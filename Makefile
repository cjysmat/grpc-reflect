descriptor:
	protoc --descriptor_set_out=./test.protoset --include_imports -I. test.proto
proto:
	protoc --go_out=plugins=grpc:. test.proto