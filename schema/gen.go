package schema

// Generates an embbedded http.FileSystem for all schema files
// using esc (https://github.com/mjibson/esc).

//go:generate esc -private -o fs.go -pkg=schema -ignore=".*go" .
