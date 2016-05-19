// Copyright 2016 The Linux Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
