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
	mpath := filepath.Join("blobs", d.normalizeDigest())

	if err := w.walk(func(path string, info os.FileInfo, r io.Reader) error {
		if info.IsDir() || filepath.Clean(path) != mpath {
			return fmt.Errorf("%s: manifest not found", mpath)
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

		return nil
	}); err != nil {
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

		if err := w.walk(func(path string, info os.FileInfo, r io.Reader) error {
			if info.IsDir() {
				return fmt.Errorf("%s: blob digest not found", d.normalizeDigest())
			}

			dd, err := filepath.Rel("blobs", filepath.Clean(path))
			if err != nil || d.normalizeDigest() != dd {
				return fmt.Errorf("%s: blob digest not found", d.normalizeDigest())
			}

			if err := unpackLayer(dest, r); err != nil {
				return errors.Wrap(err, "error extracting layer")
			}

			return nil
		}); err != nil {
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

	var dirs []*tar.Header
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

		hdr.Name = filepath.Clean(hdr.Name)
		if !strings.HasSuffix(hdr.Name, string(os.PathSeparator)) {
			// Not the root directory, ensure that the parent directory exists
			parent := filepath.Dir(hdr.Name)
			parentPath := filepath.Join(dest, parent)
			if _, err2 := os.Lstat(parentPath); err2 != nil && os.IsNotExist(err2) {
				if err3 := os.MkdirAll(parentPath, 0755); err3 != nil {
					return err3
				}
			}
		}
		path := filepath.Join(dest, hdr.Name)
		rel, err := filepath.Rel(dest, path)
		if err != nil {
			return err
		}
		info := hdr.FileInfo()
		if strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			return fmt.Errorf("%q is outside of %q", hdr.Name, dest)
		}

		if strings.HasPrefix(info.Name(), ".wh.") {
			path = strings.Replace(path, ".wh.", "", 1)

			if err := os.RemoveAll(path); err != nil {
				return errors.Wrap(err, "unable to delete whiteout path")
			}

			continue loop
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if fi, err := os.Lstat(path); !(err == nil && fi.IsDir()) {
				if err2 := os.MkdirAll(path, info.Mode()); err2 != nil {
					return errors.Wrap(err2, "error creating directory")
				}
			}

		case tar.TypeReg, tar.TypeRegA:
			f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, info.Mode())
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
		case tar.TypeXGlobalHeader:
			return nil
		}
		// Directory mtimes must be handled at the end to avoid further
		// file creation in them to modify the directory mtime
		if hdr.Typeflag == tar.TypeDir {
			dirs = append(dirs, hdr)
		}
	}
	for _, hdr := range dirs {
		path := filepath.Join(dest, hdr.Name)

		finfo := hdr.FileInfo()
		// I believe the old version was using time.Now().UTC() to overcome an
		// invalid error from chtimes.....but here we lose hdr.AccessTime like this...
		if err := os.Chtimes(path, time.Now().UTC(), finfo.ModTime()); err != nil {
			return errors.Wrap(err, "error changing time")
		}
	}
	return nil
}
