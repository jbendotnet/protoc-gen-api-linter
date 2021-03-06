GO_BINS := $(GO_BINS) cmd/protoc-gen-api-linter
DOCKER_BINS := $(DOCKER_BINS) protoc-gen-api-linter

LICENSE_HEADER_LICENSE_TYPE := apache
LICENSE_HEADER_COPYRIGHT_HOLDER := Jon Bennett
LICENSE_HEADER_YEAR_RANGE := 2021
LICENSE_HEADER_IGNORES := \/testdata

include make/go/bootstrap.mk
include make/go/go.mk
include make/go/docker.mk
include make/go/buf.mk
include make/go/license_header.mk
include make/go/dep_protoc_gen_go.mk

bufgeneratedeps:: $(BUF) $(PROTOC_GEN_GO)

.PHONY: bufgeneratecleango
bufgeneratecleango:
	rm -rf internal/gen/proto

bufgenerateclean:: bufgeneratecleango

.PHONY: bufgeneratego
bufgeneratego:
	buf generate

bufgeneratesteps:: bufgeneratego
