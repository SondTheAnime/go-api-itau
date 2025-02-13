.PHONY: docs
docs:
	scalar generate --input docs/scalar/scalar.yaml --output docs/scalar/dist

.PHONY: docs-serve
docs-serve:
	scalar serve --input docs/scalar/scalar.yaml --port 8088

SHELL := /bin/bash

.PHONY: run
run:
	set -a && source .env && set +a && go run cmd/api/main.go

.PHONY: build
build:
	go build -o bin/api cmd/api/main.go
