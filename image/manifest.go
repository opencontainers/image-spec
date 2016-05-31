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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opencontainers/image-spec/schema"
	"github.com/pkg/errors"
)

type manifest struct {
	Config descriptor   `json:"config"`
	Layers []descriptor `json:"layers"`
}

func findManifest(w walker, d *descriptor) (*manifest, error) {
	var m manifest
	mpath := filepath.Join("blobs", d.Digest)

	f := func(path string, info os.FileInfo, r io.Reader) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Clean(path) != mpath {
			return nil
		}

		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return errors.Wrapf(err, "%s: error reading manifest", path)
		}

		if err := schema.MediaTypeManifest.Validate(bytes.NewReader(buf)); err != nil {
			return errors.Wrapf(err, "%s: manifest validation failed", path)
		}

		if err := json.Unmarshal(buf, &m); err != nil {
			return err
		}

		if len(m.Layers) == 0 {
			return fmt.Errorf("%s: no layers found", path)
		}

		return errEOW
	}

	switch err := w.walk(f); err {
	case nil:
		return nil, fmt.Errorf("%s: manifest not found", mpath)
	case errEOW:
		// found, continue below
	default:
		return nil, err
	}

	return &m, nil
}

func (m *manifest) validate(w walker) error {
	if err := m.Config.validate(w); err != nil {
		return errors.Wrap(err, "config validation failed")
	}

	for _, d := range m.Layers {
		if err := d.validate(w); err != nil {
			return errors.Wrap(err, "layer validation failed")
		}
	}

	return nil
}

func (m *manifest) unpack(w walker, dest string) error {
	for _, d := range m.Layers {
		if d.MediaType != string(schema.MediaTypeImageSerialization) {
			continue
		}

		f := func(path string, info os.FileInfo, r io.Reader) error {
			if info.IsDir() {
				return nil
			}

			dd, err := filepath.Rel("blobs", filepath.Clean(path))
			if err != nil || d.Digest != dd {
				return nil // ignore
			}

			if err := unpackLayer(dest, r); err != nil {
				return errors.Wrap(err, "error extracting layer")
			}

			return errEOW
		}

		err := w.walk(f)
		if err != nil && err != errEOW {
			return err
		}
	}

	return nil
}

func unpackLayer(dest string, r io.Reader) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return errors.Wrap(err, "error creating gzip reader")
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

loop:
	for {
		hdr, err := tr.Next()
		switch err {
		case io.EOF:
			break loop
		case nil:
			// success, continue below
		default:
			return errors.Wrapf(err, "error advancing tar stream")
		}

		path := filepath.Join(dest, filepath.Clean(hdr.Name))
		info := hdr.FileInfo()

		if strings.HasPrefix(info.Name(), ".wh.") {
			path = strings.Replace(path, ".wh.", "", 1)

			if err := os.RemoveAll(path); err != nil {
				return errors.Wrap(err, "unable to delete whiteout path")
			}

			continue loop
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, info.Mode()); err != nil {
				return errors.Wrap(err, "error creating directory")
			}

		case tar.TypeReg, tar.TypeRegA:
			f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
			if err != nil {
				return errors.Wrap(err, "unable to open file")
			}

			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return errors.Wrap(err, "unable to copy")
			}
			f.Close()

		case tar.TypeLink:
			target := filepath.Join(dest, hdr.Linkname)

			if !strings.HasPrefix(target, dest) {
				return fmt.Errorf("invalid hardlink %q -> %q", target, hdr.Linkname)
			}

			if err := os.Link(target, path); err != nil {
				return err
			}

		case tar.TypeSymlink:
			target := filepath.Join(filepath.Dir(path), hdr.Linkname)

			if !strings.HasPrefix(target, dest) {
				return fmt.Errorf("invalid symlink %q -> %q", path, hdr.Linkname)
			}

			if err := os.Symlink(hdr.Linkname, path); err != nil {
				return err
			}

		}

		if err := os.Chtimes(path, time.Now().UTC(), info.ModTime()); err != nil {
			return errors.Wrap(err, "error changing time")
		}
	}

	return nil
}
