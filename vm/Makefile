# Makefile for generating Go code from protobuf files

PROTOC = protoc
OUT_DIR = .

.PHONY: all nodeapi controllerapi clean

all: controllerapi

	# Generate Go files from the .proto files
controllerapi: $(PROTO_FILES)
		$(PROTOC) --go_out=. --go_opt=paths=source_relative \
		           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
				   controllerapi/controllerapi.proto

# Clean generated files
clean_nodeapi:
	rm -f nodeapi/*.pb.go nodeapi/*.pb.gw.go

clean: clean_nodeapi

build_nodeapi: controllerapi
	GOOS=linux GOARCH=amd64 go build -o node_exe cmd/node/main.go

node: build_nodeapi
	ssh root@10.100.0.101 rm -rf /root/node_exe
	scp config.yml root@10.100.0.101:/root/config.yml
	scp node_exe root@10.100.0.101:/root/node
	rm -rf node_exe
	ssh -t root@10.100.0.101 /root/node
