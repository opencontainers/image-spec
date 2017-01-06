
default: help

help:
	@echo "Usage: make <target>"
	@echo
	@echo " * 'fmt' - format the json with indentation"

fmt:
	for i in *.json ; do jq -M . "$${i}" > xx && cat xx > "$${i}" && rm xx ; done
