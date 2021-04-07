OS_NAME := $(shell uname -s | tr A-Z a-z)

all: generate fmt lint vet test build

generate:
		protoc -I=protobuf --go_out=plugins=grpc:src/ protobuf/*.proto
		go generate

fmt:
		go list ./... | grep -v /src/ | go fmt

lint:
		go list ./... | grep -v /src/ | xargs -L1 golint -set_exit_status

vet:
		go vet ./...

test:
		go test -v ./...

build: generate
		go mod tidy
ifeq ($(OS_NAME),darwin)
			./build-osx.sh
else
			./build-linux.sh
endif

clean:
		rm -r herald*
		rm assets.go