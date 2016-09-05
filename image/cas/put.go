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

package cas

import (
	"bytes"
	"encoding/json"

	"github.com/opencontainers/image-spec/specs-go"
	"golang.org/x/net/context"
)

// PutJSON writes a generic JSON object to content-addressable storage
// and returns a Descriptor referencing it.
func PutJSON(ctx context.Context, engine Engine, data interface{}, mediaType string) (descriptor *specs.Descriptor, err error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	size := len(jsonBytes)
	size64 := int64(size) // panics on overflow

	reader := bytes.NewReader(jsonBytes)
	digest, err := engine.Put(ctx, reader)
	if err != nil {
		return nil, err
	}

	descriptor = &specs.Descriptor{
		MediaType: mediaType,
		Digest:    digest,
		Size:      size64,
	}
	return descriptor, nil
}
