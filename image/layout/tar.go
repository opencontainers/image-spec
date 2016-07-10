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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/opencontainers/image-spec/specs-go"
	"golang.org/x/net/context"
)

// TarEntryByName walks a tarball pointed to by reader, finds an
// entry matching the given name, and returns the header and reader
// for that entry.  Returns os.ErrNotExist if the path is not found.
func TarEntryByName(ctx context.Context, reader io.ReadSeeker, name string) (header *tar.Header, tarReader *tar.Reader, err error) {
	_, err = reader.Seek(0, os.SEEK_SET)
	if err != nil {
		return nil, nil, err
	}

	tarReader = tar.NewReader(reader)
	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
		}

		header, err := tarReader.Next()
		if err == io.EOF {
			return nil, nil, os.ErrNotExist
		}
		if err != nil {
			return nil, nil, err
		}

		if header.Name == name {
			return header, tarReader, nil
		}
	}
}

// CheckTarVersion walks a tarball pointed to by reader and returns an
// error if oci-layout is missing or has unrecognized content.
func CheckTarVersion(ctx context.Context, reader io.ReadSeeker) (err error) {
	_, tarReader, err := TarEntryByName(ctx, reader, "./oci-layout")
	if err == os.ErrNotExist {
		return errors.New("oci-layout not found")
	}
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(tarReader)
	var version specs.ImageLayoutVersion
	err = decoder.Decode(&version)
	if err != nil {
		return err
	}
	if version.Version != "1.0.0" {
		return fmt.Errorf("unrecognized imageLayoutVersion: %q", version.Version)
	}

	return nil
}
