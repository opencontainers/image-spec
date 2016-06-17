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

package layout

import (
	"archive/tar"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/opencontainers/image-spec/image/cas"
	"golang.org/x/net/context"
)

// TarEngine is a cas.Engine backed by a tar file.
type TarEngine struct {
	file ReadWriteSeekCloser
}

// NewTarEngine returns a new TarEngine.
func NewTarEngine(file ReadWriteSeekCloser) (engine cas.Engine, err error) {
	engine = &TarEngine{
		file: file,
	}

	return engine, nil
}

// Put adds a new blob to the store.
func (engine *TarEngine) Put(ctx context.Context, reader io.Reader) (digest string, err error) {
	// FIXME
	return "", errors.New("TarEngine.Put is not supported yet")
}

// Get returns a reader for retrieving a blob from the store.
func (engine *TarEngine) Get(ctx context.Context, digest string) (reader io.ReadCloser, err error) {
	fields := strings.SplitN(digest, ":", 2)
	if len(fields) != 2 {
		return nil, fmt.Errorf("invalid digest: %q, %v", digest, fields)
	}
	algorithm := fields[0]
	hash := fields[1]

	targetName := fmt.Sprintf("./blobs/%s/%s", algorithm, hash)

	_, err = engine.file.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, err
	}

	tarReader := tar.NewReader(engine.file)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		header, err := tarReader.Next()
		if err == io.EOF {
			return nil, os.ErrNotExist
		} else if err != nil {
			return nil, err
		}

		if header.Name == targetName {
			return ioutil.NopCloser(tarReader), nil
		}
	}
}

// Delete removes a blob from the store.
func (engine *TarEngine) Delete(ctx context.Context, digest string) (err error) {
	// FIXME
	return errors.New("TarEngine.Delete is not supported yet")
}

// Close releases resources held by the engine.
func (engine *TarEngine) Close() (err error) {
	return engine.file.Close()
}
