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
	"io/ioutil"
	"os"
	"strings"
	"time"

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

// WriteTarEntryByName reads content from reader into an entry at name
// in the tarball at file, replacing a previous entry with that name
// (if any).  The current implementation avoids writing a temporary
// file to disk, but risks leaving a corrupted tarball if the program
// crashes mid-write.
//
// To add an entry to a tarball (with Go's interface) you need to know
// the size ahead of time.  If you set the size argument,
// WriteTarEntryByName will use that size in the entry header (and
// Go's implementation will check to make sure it matches the length
// of content read from reader).  If unset, WriteTarEntryByName will
// copy reader into a local buffer, measure its size, and then write
// the entry header and content.
func WriteTarEntryByName(ctx context.Context, file io.ReadWriteSeeker, name string, reader io.Reader, size *int64) (err error) {
	var buffer bytes.Buffer
	tarWriter := tar.NewWriter(&buffer)

	components := strings.Split(name, "/")
	if components[0] != "." {
		return fmt.Errorf("tar name entry does not start with './': %q", name)
	}

	var parents []string
	for i := 2; i < len(components); i++ {
		parents = append(parents, strings.Join(components[:i], "/"))
	}

	_, err = file.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(file)
	found := false
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var header *tar.Header
		header, err = tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		dirName := strings.TrimRight(header.Name, "/")
		for i, parent := range parents {
			if dirName == parent {
				parents = append(parents[:i], parents[i+1:]...)
				break
			}
		}

		if header.Name == name {
			found = true
			err = writeTarEntry(ctx, tarWriter, name, reader, size)
		} else {
			err = tarWriter.WriteHeader(header)
			if err != nil {
				return err
			}
			_, err = io.Copy(tarWriter, tarReader)
		}
		if err != nil {
			return err
		}
	}

	if !found {
		now := time.Now()
		for _, parent := range parents {
			header := &tar.Header{
				Name:     parent + "/",
				Mode:     0777,
				ModTime:  now,
				Typeflag: tar.TypeDir,
			}
			err = tarWriter.WriteHeader(header)
			if err != nil {
				return err
			}
		}
		err = writeTarEntry(ctx, tarWriter, name, reader, size)
		if err != nil {
			return err
		}
	}

	err = tarWriter.Close()
	if err != nil {
		return err
	}

	_, err = file.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}
	// FIXME: truncate file

	_, err = buffer.WriteTo(file)
	return err
}

func writeTarEntry(ctx context.Context, writer *tar.Writer, name string, reader io.Reader, size *int64) (err error) {
	if size == nil {
		var data []byte
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
		_size := int64(len(data))
		size = &_size
	}
	now := time.Now()
	header := &tar.Header{
		Name:     name,
		Mode:     0666,
		Size:     *size,
		ModTime:  now,
		Typeflag: tar.TypeReg,
	}
	err = writer.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, reader)
	return err
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

// CreateTarFile creates a new image-layout tar file at the given path.
func CreateTarFile(ctx context.Context, path string) (err error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	tarWriter := tar.NewWriter(file)
	defer tarWriter.Close()

	now := time.Now()
	for _, name := range []string{"./blobs/", "./refs/"} {
		header := &tar.Header{
			Name:     name,
			Mode:     0777,
			ModTime:  now,
			Typeflag: tar.TypeDir,
		}
		err = tarWriter.WriteHeader(header)
		if err != nil {
			return err
		}
	}

	imageLayoutVersion := specs.ImageLayoutVersion{
		Version: "1.0.0",
	}
	imageLayoutVersionBytes, err := json.Marshal(imageLayoutVersion)
	if err != nil {
		return err
	}
	header := &tar.Header{
		Name:     "./oci-layout",
		Mode:     0666,
		Size:     int64(len(imageLayoutVersionBytes)),
		ModTime:  now,
		Typeflag: tar.TypeReg,
	}
	err = tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}
	_, err = tarWriter.Write(imageLayoutVersionBytes)
	return err
}
