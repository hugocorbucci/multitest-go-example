NAME=server
REV=$(shell git rev-parse --short=10 HEAD)
REGISTRY=docker.io
IMAGE=${REGISTRY}/hugocorbucci/multitest-go-example
GOPATH:=$(shell echo ${GOPATH})
PROJECT_ROOT=$(shell cd ../.. && pwd)
SERVER_RELATIVE=$(shell pwd | sed -e "s|^${GOPATH}/||g" )
SHMIG_CONFIG=docker
DOCKER_IP:=$(shell if (docker-machine env 2>/dev/null >/dev/null); then (docker-machine ip); else echo "127.0.0.1"; fi)
DB_PORT:=$(shell if (docker-compose ps mysql 2>/dev/null >/dev/null) && (docker-compose port mysql 3306 2>/dev/null >/dev/null); then (docker-compose port mysql 3306 | sed 's\#.*:\([0-9]*\)\#\1\#'); else echo "3306"; fi)
TARGET_URL:=$(shell if [ -z "${TARGET_URL}" ]; then echo "http://localhost:8080"; else echo "${TARGET_URL}"; fi)
DB_USER:=root
DB_PASSWORD:=$(shell if [ -z "${DB_PASSWORD}" ]; then if (docker-compose ps mysql 2>/dev/null >/dev/null); then grep MYSQL_ROOT_PASSWORD docker-compose.yml | sed "s\#.*MYSQL_ROOT_PASSWORD: '\([^']*\)'.*\#\1\#"; fi; else echo "${DB_PASSWORD}"; fi)
DB_CONN:=$(shell if [ -z "${DB_CONN}" ]; then echo "${DB_USER}:${DB_PASSWORD}@tcp(${DOCKER_IP}:${DB_PORT})/multitest"; else echo "${DB_CONN}"; fi)
.DEFAULT_GOAL := test

target:
	mkdir -p target

test:
	$(MAKE) -C "../.." $@
.PHONY: test

smoke_test:
	export TARGET_URL="${TARGET_URL}"; \
	export DB_CONN="${DB_CONN}"; \
	cd "../.." && go test ./...
.PHONY: smoke_test

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

db-up: docker-compose.yml
	[ "docker" = "${SHMIG_CONFIG}" ] && (([ -n "$(shell docker-compose ps -q mysql)" ] && echo "db is up") || (docker-compose up -d mysql && sleep 5)) || echo "using local mysql"
	bin/wait-for "$(shell ([ "docker" = "${SHMIG_CONFIG}" ] && docker-machine ip default) || echo "localhost"):$(shell ([ "docker" = "${SHMIG_CONFIG}" ] && docker-compose port mysql 3306 | sed -e 's|[\.0-9]*:||g') || echo '3306')" -t 60
.PHONY: db-up

migrations:
	mkdir -p migrations

migration: migrations bin/shmig
	([ -n "${MIGRATION}" ] && bin/shmig -c "shmig.${SHMIG_CONFIG}.conf" create "${MIGRATION}") || echo "Set env var MIGRATION to create the migration"

migrate: migrations db-up
	bin/shmig -c "shmig.${SHMIG_CONFIG}.conf" migrate
