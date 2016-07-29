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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	caslayout "github.com/opencontainers/image-spec/image/cas/layout"
	imagelayout "github.com/opencontainers/image-spec/image/layout"
	"github.com/opencontainers/image-spec/image/refs"
	"github.com/opencontainers/image-spec/specs-go"
	"golang.org/x/net/context"
)

// TarEngine is a refs.Engine backed by a tar file.
type TarEngine struct {
	file caslayout.ReadWriteSeekCloser
}

// NewTarEngine returns a new TarEngine.
func NewTarEngine(ctx context.Context, file caslayout.ReadWriteSeekCloser) (eng refs.Engine, err error) {
	engine := &TarEngine{
		file: file,
	}

	err = imagelayout.CheckTarVersion(ctx, engine.file)
	if err != nil {
		return nil, err
	}

	return engine, nil
}

// Put adds a new reference to the store.
func (engine *TarEngine) Put(ctx context.Context, name string, descriptor *specs.Descriptor) (err error) {
	data, err := json.Marshal(descriptor)
	if err != nil {
		return err
	}

	size := int64(len(data))
	reader := bytes.NewReader(data)
	targetName := fmt.Sprintf("./refs/%s", name)
	return imagelayout.WriteTarEntryByName(ctx, engine.file, targetName, reader, &size)
}

// Get returns a reference from the store.
func (engine *TarEngine) Get(ctx context.Context, name string) (descriptor *specs.Descriptor, err error) {
	targetName := fmt.Sprintf("./refs/%s", name)

	_, tarReader, err := imagelayout.TarEntryByName(ctx, engine.file, targetName)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(tarReader)
	var desc specs.Descriptor
	err = decoder.Decode(&desc)
	if err != nil {
		return nil, err
	}
	return &desc, nil
}

// List returns available names from the store.
func (engine *TarEngine) List(ctx context.Context, prefix string, size int, from int, nameCallback refs.ListNameCallback) (err error) {
	var i = 0

	_, err = engine.file.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil
	}

	pathPrefix := fmt.Sprintf("./refs/%s", prefix)

	tarReader := tar.NewReader(engine.file)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var header *tar.Header
		header, err = tarReader.Next()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		if strings.HasPrefix(header.Name, pathPrefix) && len(header.Name) > 7 {
			i++
			if i > from {
				err = nameCallback(ctx, header.Name[7:])
				if err != nil {
					return err
				}
				if i-from == size {
					return nil
				}
			}
		}
	}
}

// Delete removes a reference from the store.
func (engine *TarEngine) Delete(ctx context.Context, name string) (err error) {
	// FIXME
	return errors.New("TarEngine.Delete is not supported yet")
}

// Close releases resources held by the engine.
func (engine *TarEngine) Close() (err error) {
	return engine.file.Close()
}
