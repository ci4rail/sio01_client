NAME = io4edge-devsim
BIN_DIR ?= ../../bin
VERSION ?= $(shell git describe --match=NeVeRmAtCh --always --abbrev=40 --dirty)
GO_LDFLAGS = -ldflags "-X github.com/ci4rail/sio01_host/devsim/pkg/version.Version=$(VERSION)"
PROTO_DIR=../proto
PROTO_ELOC_OUT=internal/eloc
all: proto test build-io4edge-devsim

build: build-io4edge-devsim

proto:
	mkdir -p ${PROTO_ELOC_OUT}/pb
	protoc -I=${PROTO_DIR} \
	  --go_opt=Mtracelet_location.proto=./pb \
	  --go_opt=Mtracelet_status.proto=./pb \
	  --go_out=${PROTO_ELOC_OUT} \
	  ${PROTO_DIR}/tracelet_location.proto \
	  ${PROTO_DIR}/tracelet_status.proto

build-io4edge-devsim:
	go mod tidy
	GOOS=linux go build $(GO_LDFLAGS) -o ${BIN_DIR}/${NAME} main.go

test:
	go test ./...

clean:
	rm -f ${BIN_DIR}/${NAME}

.PHONY: all proto build test clean clean build-io4edge-devsim
