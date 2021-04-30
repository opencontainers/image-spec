EPOCH_TEST_COMMIT ?= v0.2.0

DOCKER ?= $(shell command -v docker 2>/dev/null)
PANDOC ?= $(shell command -v pandoc 2>/dev/null)

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

TOOLS := esc gitvalidation glide glide-vc

default: check-license lint test

help:
	@echo "Usage: make <target>"
	@echo
	@echo " * 'docs' - produce document in the $(OUTPUT_DIRNAME) directory"
	@echo " * 'fmt' - format the json with indentation"
	@echo " * 'validate-examples' - validate the examples in the specification markdown files"
	@echo " * 'schema-fs' - regenerate the virtual schema http/FileSystem"
	@echo " * 'check-license' - check license headers in source files"
	@echo " * 'lint' - Execute the source code linter"
	@echo " * 'test' - Execute the unit tests"
	@echo " * 'img/*.png' - Generate PNG from dot file"

fmt:
	for i in schema/*.json ; do jq --indent 2 -M . "$${i}" > xx && cat xx > "$${i}" && rm xx ; done

docs: $(OUTPUT_DIRNAME)/$(DOC_FILENAME).pdf $(OUTPUT_DIRNAME)/$(DOC_FILENAME).html

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

validate-examples: schema/fs.go
	go test -run TestValidate ./schema

schema/fs.go: $(wildcard schema/*.json) schema/gen.go
	cd schema && printf "%s\n\n%s\n" "$$(cat ../.header)" "$$(go generate)" > fs.go

schema-fs: schema/fs.go
	@echo "generating schema fs"

check-license:
	@echo "checking license headers"
	@./.tool/check-license

lint:
	@echo "checking lint"
	@./.tool/lint

test: schema/fs.go
	go test -race -cover $(shell go list ./... | grep -v /vendor/)

img/%.png: img/%.dot
	dot -Tpng $^ > $@


# When this is running in GitHub, it will only check the commit range
.gitvalidation:
	@which git-validation > /dev/null 2>/dev/null || (echo "ERROR: git-validation not found. Consider 'make install.tools' target" && false)
ifdef GITHUB_SHA
	git-validation -q -run DCO,short-subject,dangling-whitespace -range $(GITHUB_SHA)..HEAD
else
	git-validation -v -run DCO,short-subject,dangling-whitespace -range $(EPOCH_TEST_COMMIT)..HEAD
endif

install.tools: $(TOOLS:%=.install.%)

.install.esc:
	go get -u github.com/mjibson/esc

.install.gitvalidation:
	go get -u github.com/vbatts/git-validation

.install.glide:
	go get -u github.com/Masterminds/glide

.install.glide-vc:
	go get -u github.com/sgotti/glide-vc

clean:
	rm -rf *~ $(OUTPUT_DIRNAME) header.html

.PHONY: \
	$(TOOLS:%=.install.%) \
	validate-examples \
	check-license \
	clean \
	lint \
	install.tools \
	docs \
	test \
	.gitvalidation \
	schema/fs.go \
	schema-fs
