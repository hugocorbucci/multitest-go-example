SERVER_PATH=cmd/server

.DEFAULT_GOAL := test

test:
	go test ./...
.PHONY: test

local-run:
	$(MAKE) -C ${SERVER_PATH} local-run
.PHONY: local-run
