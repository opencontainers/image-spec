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

package image

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type descriptor struct {
	MediaType string `json:"mediaType"`
	Digest    string `json:"digest"`
	Size      int64  `json:"size"`
}

func (d *descriptor) normalizeDigest() string {
	return strings.Replace(d.Digest, ":", "-", -1)
}

func findDescriptor(w walker, name string) (*descriptor, error) {
	var d descriptor
	dpath := filepath.Join("refs", name)

	if err := w.walk(func(path string, info os.FileInfo, r io.Reader) error {
		if info.IsDir() || filepath.Clean(path) != dpath {
			return fmt.Errorf("%s: descriptor not found", dpath)
		}

		if err := json.NewDecoder(r).Decode(&d); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &d, nil
}

func (d *descriptor) validate(w walker) error {
	if err := w.walk(func(path string, info os.FileInfo, r io.Reader) error {
		if info.IsDir() {
			return fmt.Errorf("%s: not found", d.normalizeDigest())
		}

		digest, err := filepath.Rel("blobs", filepath.Clean(path))
		if err != nil || d.normalizeDigest() != digest {
			return fmt.Errorf("%s: not found", d.normalizeDigest())
		}

		if err := d.validateContent(r); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return errors.Wrapf(err, "%s: validation failed", d.normalizeDigest())
	}

	return nil
}

func (d *descriptor) validateContent(r io.Reader) error {
	h := sha256.New()
	n, err := io.Copy(h, r)
	if err != nil {
		return errors.Wrap(err, "error generating hash")
	}

	digest := "sha256:" + hex.EncodeToString(h.Sum(nil))

	if digest != d.Digest {
		return errors.New("digest mismatch")
	}

	if n != d.Size {
		return errors.New("size mismatch")
	}

	return nil
}
