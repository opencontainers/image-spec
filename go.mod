module github.com/opencontainers/image-spec

// The minimum Go release is only incremented when required by a feature.
// At least 3 Go releases will be supported by the spec.
// For example, updating this version to 1.19 first requires Go 1.21 to be released.
go 1.21

require (
	github.com/opencontainers/go-digest v1.0.0
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.1
)

require golang.org/x/text v0.14.0 // indirect
