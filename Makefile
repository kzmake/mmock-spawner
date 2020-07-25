SOURCES := $(shell find . \
	-type f -name '*.go' \
	-not -name '*.pb*.go' \
	-not -path "./vendor/*" \
	-not -path "./hack/tools/vendor/*")

.DEFAULT_GOAL := help

.PHONY: tools
tools: ## 開発に必要なツールをインストールします
	@echo "\033[31m"
	@echo "$$ brew install protobuf"
	@echo "\033[0m"
	make -C hack/tools install

.PHONY: lint
lint: ## コードを検証します
	golangci-lint run

.PHONY: fmt
fmt: ## コードをフォーマットします
	@goimports -l -w $(SOURCES)

.PHONY: mod
mod: ## mod.go / mod.sum を整理します
	go mod tidy
	go mod vendor


.PHONY: proto
proto: ## protoファイルからgoファイルを生成します
	@for f in proto/*.proto; do \
		protoc \
		--proto_path=.:. \
		--proto_path=.:${GOPATH}/src \
		--proto_path=.:${GOPATH}/src/github.com/envoyproxy/protoc-gen-validate \
		--go_out=paths=source_relative:. \
		--micro_out=paths=source_relative:. \
		--validate_out=lang=go,paths=source_relative:. \
		$$f; \
		echo "generating $$f"; \
	done

.PHONY: proto/lint
proto/lint: ## protoファイルを検証します
	@echo "linting protos"
	buf check lint

.PHONY: proto/fmt
proto/fmt: proto/lint ## protoファイルのフォーマットを行います
	@echo "formating protos"
	prototool format -d . || true
	prototool format -w .

.PHONY: __
__:
	@echo "\033[33m"
	@echo "kzmake/mmock-spawner"
	@echo "\033[0m"

.PHONY: help
help: __ ## ヘルプを表示します
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@cat $(MAKEFILE_LIST) \
	| grep -e "^[a-zA-Z_/\-]*: *.*## *" \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-24s\033[0m %s\n", $$1, $$2}'
