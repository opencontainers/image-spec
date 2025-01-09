module github.com/opencontainers/image-spec

// The minimum Go release is only incremented when required by a feature.
// At least 3 Go releases will be supported by the spec.
// For example, updating this version to 1.19 first requires Go 1.21 to be released.
go 1.18

require (
	github.com/opencontainers/go-digest v1.0.0
	github.com/russross/blackfriday v1.6.0
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1
)
