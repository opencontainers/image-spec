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

package compatibility

const (
	// MediaTypeDockerImageManifestList specifies the docker image's mediaType for manifest lists.
	MediaTypeDockerImageManifestList = "application/vnd.docker.distribution.manifest.list.v2+json"

	// MediaTypeDockerImageManifest specifies the docker image's mediaType for the current version.
	MediaTypeDockerImageManifest = "application/vnd.docker.distribution.manifest.v2+json"

	// MediaTypeDockerImageLayer specifies the docker image's mediaType used for layers referenced by the manifest.
	MediaTypeDockerImageLayer = "application/vnd.docker.image.rootfs.diff.tar.gzip"

	// MediaTypeDockerImageConfig specifies the docker image's mediaType for the image configuration.
	MediaTypeDockerImageConfig = "application/vnd.docker.container.image.v1+json"
)
