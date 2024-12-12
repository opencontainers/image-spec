# Hacking Guide

## Overview

This guide contains instructions for building artifacts contained in this repository.

### Go

This spec includes several Go packages, and a command line tool considered to be a reference implementation of the OCI image specification.

Prerequisites:

- Go - latest version is recommended, see the [go.mod](go.mod) file for the minimum requirement
- make

The following make targets are relevant for any work involving the Go packages.

### Linting

The included Go source code is being examined for any linting violations not included in the standard Go compiler.
Linting is done using [golangci-lint][golangci-lint].

Invocation:

```shell
make lint
```

### Tests

This target executes all Go based tests.

Invocation:

```shell
make test
make validate-examples
```

### JSON schema formatting

This target auto-formats all JSON files in the `schema` directory using the `jq` tool.

Prerequisites:

- [jq][jq] >=1.5

Invocation:

```shell
make fmt
```

### OCI image specification PDF/HTML documentation files

This target generates a PDF/HTML version of the OCI image specification.

Prerequisites:

- [Docker][docker]

Invocation:

```shell
make docs
```

### License header check

This target checks if the source code includes necessary headers.

Invocation:

```shell
make check-license
```

### Clean build artifacts

This target cleans all generated/compiled artifacts.

Invocation:

```shell
make clean
```

### Create PNG images from dot files

This target generates PNG image files from DOT source files in the `img` directory.

Prerequisites:

- [graphviz][graphviz]

Invocation:

```shell
make img/media-types.png
```

[docker]: https://www.docker.com/
[golangci-lint]: https://github.com/golangci/golangci-lint
[graphviz]: https://www.graphviz.org/
[jq]: https://stedolan.github.io/jq/
