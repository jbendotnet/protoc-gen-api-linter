MAKEGO := make/go
MAKEGO_REMOTE := https://github.com/jbendotnet/makego.git
PROJECT := protoc-gen-api-linter
GO_MODULE := github.com/jbendotnet/protoc-gen-api-linter
DOCKER_ORG := jbendotnet
DOCKER_PROJECT := protoc-gen-api-linter

include make/protoc-gen-api-linter/all.mk
