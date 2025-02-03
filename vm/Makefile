# Makefile for generating Go code from protobuf files

PROTOC = protoc
OUT_DIR = .

.PHONY: all nodeapi controllerapi clean

all: nodeapi

	# Generate Go files from the .proto files
controllerapi: $(PROTO_FILES)
		$(PROTOC) --go_out=. --go_opt=paths=source_relative \
		           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
				   controllerapi/controllerapi.proto

# Generate Go files from the .proto files
nodeapi: $(PROTO_FILES)
	$(PROTOC) --go_out=. --go_opt=paths=source_relative \
	           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
			   nodeapi/nodeapi.proto

# Clean generated files
clean_nodeapi:
	rm -f nodeapi/*.pb.go nodeapi/*.pb.gw.go

clean: clean_nodeapi

build_nodeapi: nodeapi
	GOOS=linux GOARCH=amd64 go build -o node cmd/node/main.go

node: build_nodeapi
	ssh root@10.100.0.101 rm -rf /root/node
	scp config.yml root@10.100.0.101:/root/config.yml
	scp node root@10.100.0.101:/root/node
	ssh -t root@10.100.0.101 /root/node
