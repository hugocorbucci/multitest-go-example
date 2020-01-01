SERVER_PATH=cmd/server

.DEFAULT_GOAL := test

target:
	mkdir -p target

test:
	go test ./...
.PHONY: test

smoke_test:
	$(MAKE) -C ${SERVER_PATH} $@
.PHONY: smoke_test

local-run:
	$(MAKE) -C ${SERVER_PATH} $@
.PHONY: local-run

build: target
	$(MAKE) -C ${SERVER_PATH} $@
	cp ${SERVER_PATH}/target/* target
.PHONY: local-run

package:
	$(MAKE) -C ${SERVER_PATH} $@
.PHONY: package

push:
	$(MAKE) -C ${SERVER_PATH} $@
.PHONY: push

clean:
	-rm -Rf target
.PHONY: clean
