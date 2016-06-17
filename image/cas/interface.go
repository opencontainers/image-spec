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

// Package cas implements generic content-addressable storage.
package cas

import (
	"io"

	"golang.org/x/net/context"
)

// Engine represents a content-addressable storage engine.
type Engine interface {

	// Put adds a new blob to the store.  The action is idempotent; a
	// nil return means "that content is stored at DIGEST" without
	// implying "because of your Put()".
	Put(ctx context.Context, reader io.Reader) (digest string, err error)

	// Get returns a reader for retrieving a blob from the store.
	// Returns os.ErrNotExist if the digest is not found.
	Get(ctx context.Context, digest string) (reader io.ReadCloser, err error)

	// Delete removes a blob from the store. Returns os.ErrNotExist if
	// the digest is not found.
	Delete(ctx context.Context, digest string) (err error)

	// Close releases resources held by the engine.  Subsequent engine
	// method calls will fail.
	Close() (err error)
}
