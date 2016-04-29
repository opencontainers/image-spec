
default: help

help:
	@echo "Usage: make <target>"
	@echo
	@echo " * 'fmt' - format the json with indentation"
	@echo " * 'validate' - build the validation tool"

fmt:
	for i in *.json ; do jq --indent 2 -M . "$${i}" > xx && cat xx > "$${i}" && rm xx ; done

.PHONY: validate-examples
validate-examples: oci-validate-examples
	./oci-validate-examples < manifest.md

oci-validate-json: validate.go
	go build ./cmd/oci-validate-json

oci-validate-examples: cmd/oci-validate-examples/main.go
	go build ./cmd/oci-validate-examples

