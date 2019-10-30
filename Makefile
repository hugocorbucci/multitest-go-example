SERVER_PATH=cmd/server

.DEFAULT_GOAL := local-run

target:
	mkdir -p target

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
