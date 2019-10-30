NAME=server
REV=$(shell git rev-parse --short=10 HEAD)
REGISTRY=docker.io
IMAGE=${REGISTRY}/hugocorbucci/multitest-go-example
GOPATH:=$(shell echo ${GOPATH})
PROJECT_ROOT=$(shell cd ../.. && pwd)
SERVER_RELATIVE=$(shell pwd | sed -e "s|^${GOPATH}/||g" )
SHMIG_CONFIG=docker
.DEFAULT_GOAL := test

target:
	mkdir -p target

test:
	$(MAKE) -C "../.." $@
.PHONY: test

clean:
	@echo "==> Cleaning releases"
	-rm -Rf target
	$(MAKE) -C "../.." $@
.PHONY: clean

local-run:
	bin/local-run
.PHONY: local-run

target/$(NAME)-darwin: target $(shell find ../../internal -type f) $(shell find . -path ./target -prune -o -type f)
	([ "$(shell uname)" = "Darwin" ] && GOOS='darwin' go build -o $@ .) || echo "Can't compile darwin executable"

target/$(NAME)-linux64: target $(shell find ../../internal -type f) $(shell find . -path ./target -prune -o -type f)
	([ "$(shell uname)" = "Darwin" ] && docker run --rm -v "${GOPATH}":/home/guest -w "/home/guest/${SERVER_RELATIVE}" -e "CGO_ENABLED=0" -e "GOPATH=/home/guest" golang:1.13 go build -o $@ .) || GOOS='linux' go build -o $@ .

build: target/$(NAME)-linux64 target/$(NAME)-darwin
.PHONY: build

package: target/$(NAME)-linux64 Dockerfile
	docker build -t "${IMAGE}:latest" -t "${IMAGE}:${REV}" .
.PHONY: package

push: package
	docker push "${IMAGE}:${REV}"
	docker push "${IMAGE}:latest"
.PHONY: push
