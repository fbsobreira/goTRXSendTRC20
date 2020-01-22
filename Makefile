# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=sendTRC20
BINARY_UNIX=$(BINARY_NAME)_unix


all: deps build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN) ./...
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./
	./$(BINARY_NAME)

deps:
	$(GOGET) -u google.golang.org/grpc
	$(GOGET) -u github.com/golang/protobuf/proto
	$(GOGET) -u github.com/golang/protobuf/protoc-gen-go
	$(GOGET) -u github.com/sasaxie/go-client-api
	$(GOGET) -u golang.org/x/crypto/sha3

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
