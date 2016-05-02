package schema

import "net/http"

// Media types for the OCI image formats
const (
	MediaTypeManifest                   Validator     = `application/vnd.oci.image.manifest.v1+json`
	MediaTypeManifestList               Validator     = `application/vnd.oci.image.manifest.list.v1+json`
	MediaTypeImageSerialization         unimplemented = `application/vnd.oci.image.serialization.rootfs.tar.gzip`
	MediaTypeImageSerializationConfig   unimplemented = `application/vnd.oci.image.serialization.config.v1+json`
	MediaTypeImageSerializationCombined unimplemented = `application/vnd.oci.image.serialization.combined.v1+json`
)

var (
	// fs stores the embedded http.FileSystem
	// having the OCI JSON schema files in root "/".
	fs = _escFS(false)

	// specs maps OCI schema media types to schema files.
	specs = map[Validator]string{
		MediaTypeManifest:     "image-manifest-schema.json",
		MediaTypeManifestList: "manifest-list-schema.json",
	}
)

// FileSystem returns an in-memory file system including the schema files.
// The schema files are located at the root directory.
func FileSystem() http.FileSystem {
	return fs
}
