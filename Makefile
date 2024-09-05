EPOCH_TEST_COMMIT ?= v0.2.0

DOCKER ?= $(shell command -v docker 2>/dev/null)
PANDOC ?= $(shell command -v pandoc 2>/dev/null)

GOPATH:=$(shell go env GOPATH)

OUTPUT_DIRNAME	?= output
DOC_FILENAME	?= oci-image-spec

PANDOC_CONTAINER ?= ghcr.io/opencontainers/pandoc:2.9.2.1-8.fc33.x86_64@sha256:5d81ff930a043295a557be8b003ece2a33d14e91b28c50d368413b83372f8d28
ifeq "$(strip $(PANDOC))" ''
	ifneq "$(strip $(DOCKER))" ''
		PANDOC = $(DOCKER) run \
			--rm \
			-v $(shell pwd)/:/input/:ro \
			-v $(shell pwd)/$(OUTPUT_DIRNAME)/:/$(OUTPUT_DIRNAME)/ \
			-u $(shell id -u) \
			--workdir /input \
			$(PANDOC_CONTAINER)
		PANDOC_SRC := /input/
		PANDOC_DST := /
	endif
endif

# These docs are in an order that determines how they show up in the PDF/HTML docs.
DOC_FILES := \
	spec.md \
	media-types.md \
	descriptor.md \
	image-layout.md \
	manifest.md \
	image-index.md \
	layer.md \
	config.md \
	annotations.md \
	conversion.md \
	considerations.md \
	implementations.md

FIGURE_FILES := \
	img/media-types.png

MARKDOWN_LINT_VER?=v0.8.1

TOOLS := gitvalidation

default: check-license lint test

.PHONY: fmt
fmt: ## format the json with indentation
	for i in schema/*.json ; do jq --indent 2 -M . "$${i}" > xx && cat xx > "$${i}" && rm xx ; done

.PHONY: docs
docs: $(OUTPUT_DIRNAME)/$(DOC_FILENAME).pdf $(OUTPUT_DIRNAME)/$(DOC_FILENAME).html ## generate a PDF/HTML version of the OCI image specification

ifeq "$(strip $(PANDOC))" ''
$(OUTPUT_DIRNAME)/$(DOC_FILENAME).pdf: $(DOC_FILES) $(FIGURE_FILES)
	$(error cannot build $@ without either pandoc or docker)
else
$(OUTPUT_DIRNAME)/$(DOC_FILENAME).pdf: $(DOC_FILES) $(FIGURE_FILES)
	@mkdir -p $(OUTPUT_DIRNAME)/ && \
	$(PANDOC) -f gfm -t latex --pdf-engine=xelatex -V geometry:margin=0.5in,bottom=0.8in -V block-headings -o $(PANDOC_DST)$@ $(patsubst %,$(PANDOC_SRC)%,$(DOC_FILES))
	ls -sh $(realpath $@)

$(OUTPUT_DIRNAME)/$(DOC_FILENAME).html: header.html $(DOC_FILES) $(FIGURE_FILES)
	@mkdir -p $(OUTPUT_DIRNAME)/ && \
	cp -ap img/ $(shell pwd)/$(OUTPUT_DIRNAME)/&& \
	$(PANDOC) -f gfm -t html5 -H $(PANDOC_SRC)header.html --standalone -o $(PANDOC_DST)$@ $(patsubst %,$(PANDOC_SRC)%,$(DOC_FILES))
	ls -sh $(realpath $@)
endif

header.html: .tool/genheader.go specs-go/version.go
	go run .tool/genheader.go > $@

.PHONY: validate-examples
validate-examples: schema/schema.go ## validate examples in the specification markdown files
	go test -run TestValidate ./schema

.PHONY: check-license
check-license: ## check license headers in source files
	@echo "checking license headers"
	@./.tool/check-license

.PHONY: lint

.PHONY: lint
lint: lint-go lint-md ## Run all linters

.PHONY: lint-go
lint-go: .install.lint ## lint check of Go files using golangci-lint
	@echo "checking Go lint"
	@GO111MODULE=on $(GOPATH)/bin/golangci-lint run

.PHONY: lint-md
lint-md: ## Run linting for markdown
	docker run --rm -v "$(PWD):/workdir:ro" docker.io/davidanson/markdownlint-cli2:$(MARKDOWN_LINT_VER) \
	  "**/*.md" "#vendor"

.PHONY: test
test: ## run the unit tests
	go test -race -cover $(shell go list ./... | grep -v /vendor/)

img/%.png: img/%.dot ## generate PNG from dot file
	dot -Tpng $^ > $@

# When this is running in GitHub, it will only check the commit range
.PHONY: .gitvalidation
.gitvalidation:
	@which git-validation > /dev/null 2>/dev/null || (echo "ERROR: git-validation not found. Consider 'make install.tools' target" && false)
ifdef GITHUB_SHA
	$(GOPATH)/bin/git-validation -q -run DCO,short-subject,dangling-whitespace -range $(GITHUB_SHA)..HEAD
else
	$(GOPATH)/bin/git-validation -v -run DCO,short-subject,dangling-whitespace -range $(EPOCH_TEST_COMMIT)..HEAD
endif

.PHONY: .install.tools
install.tools: $(TOOLS:%=.install.%)

.PHONY: .install.lint
.install.lint:
	case "$$(go env GOVERSION)" in \
	go1.18.*)	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.47.3;; \
	go1.19.*)	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1;; \
	go1.20.*)	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2;; \
	go1.21.*)	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1;; \
	*) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest;; \
	esac

.PHONY: .install.gitvalidation
.install.gitvalidation:
	go install github.com/vbatts/git-validation@latest

.PHONY: clean
clean: ## clean all generated and compiled artifacts
	rm -rf *~ $(OUTPUT_DIRNAME) header.html

.PHONY: help
help: # Display help
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST)
