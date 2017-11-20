all:	check

clean:	
		rm -rf cache_expt.o

check:	
		@go get -u github.com/golang/protobuf/proto
		@go get -u github.com/golang/protobuf/protoc-gen-go
		@go get stathat.com/c/consistent
		pip install grpcio grpcio-tools
		@go build kvstore