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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// ValidateManifest validates the given CAS blob as
// application/vnd.oci.image.manifest.v1+json.  Calls
// ValidateJSONSchema as well.
func ValidateManifest(blob io.Reader, descriptor *v1.Descriptor, strict bool) (err error) {
	if descriptor.MediaType != v1.MediaTypeImageManifest {
		return fmt.Errorf("unexpected descriptor media type: %q", descriptor.MediaType)
	}

	buffer, err := ioutil.ReadAll(blob)
	if err != nil {
		return errors.Wrapf(err, "unable to read %s", descriptor.Digest)
	}

	err = ValidateJSONSchema(bytes.NewReader(buffer), descriptor, strict)
	if err != nil {
		return err
	}

	header := v1.Manifest{}
	err = json.Unmarshal(buffer, &header)
	if err != nil {
		return errors.Wrap(err, "manifest format mismatch")
	}

	if header.Config.MediaType != v1.MediaTypeImageConfig {
		error := fmt.Errorf("warning: config %s has an unknown media type: %s\n", header.Config.Digest, header.Config.MediaType)
		if strict {
			return error
		}
		fmt.Println(error)
	}

	for _, layer := range header.Layers {
		if layer.MediaType != v1.MediaTypeImageLayer &&
			layer.MediaType != v1.MediaTypeImageLayerNonDistributable {
			error := fmt.Errorf("warning: layer %s has an unknown media type: %s\n", layer.Digest, layer.MediaType)
			if strict {
				return error
			}
			fmt.Println(error)
		}
	}

	return nil
}
