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
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const (
	layoutStr = `{"imageLayoutVersion": "1.0.0"}`

	configStr = `{
    "created": "2015-10-31T22:22:56.015925234Z",
    "author": "Alyssa P. Hacker <alyspdev@example.com>",
    "architecture": "amd64",
    "os": "linux",
    "config": {
        "User": "alice",
        "Memory": 2048,
        "MemorySwap": 4096,
        "CpuShares": 8,
        "ExposedPorts": {
            "8080/tcp": {}
        },
        "Env": [
            "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "FOO=oci_is_a",
            "BAR=well_written_spec"
        ],
        "Entrypoint": [
            "/bin/my-app-binary"
        ],
        "Cmd": [
            "--foreground",
            "--config",
            "/etc/my-app.d/default.cfg"
        ],
        "Volumes": {
            "/var/job-result-data": {},
            "/var/log/my-app-logs": {}
        },
        "WorkingDir": "/home/alice"
    },
    "rootfs": {
      "diff_ids": [
        "sha256:c6f988f4874bb0add23a778f753c65efe992244e148a1d2ec2a8b664fb66bbd1",
        "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
      ],
      "type": "layers"
    },
    "history": [
      {
        "created": "2015-10-31T22:22:54.690851953Z",
        "created_by": "/bin/sh -c #(nop) ADD file:a3bc1e842b69636f9df5256c49c5374fb4eef1e281fe3f282c65fb853ee171c5 in /"
      },
      {
        "created": "2015-10-31T22:22:55.613815829Z",
        "created_by": "/bin/sh -c #(nop) CMD [\"sh\"]",
        "empty_layer": true
      }
    ]
}
`
)

var (
	refStr = `{"digest":"<manifest_digest>","mediaType":"application/vnd.oci.image.manifest.v1+json","size":<manifest_size>}`

	manifestStr = `{
    "annotations": null,
    "config": {
        "digest": "<config_digest>",
        "mediaType": "application/vnd.oci.image.config.v1+json",
        "size": <config_size>
    },
    "layers": [
        {
            "digest": "<layer_digest>",
            "mediaType": "application/vnd.oci.image.layer.tar+gzip",
            "size": <layer_size>
        }
    ],
    "mediaType": "application/vnd.oci.image.manifest.v1+json",
    "schemaVersion": 2
}
 `
)

func TestValidateLayout(t *testing.T) {
	root, err := ioutil.TempDir("", "oci-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(root)

	err = os.MkdirAll(filepath.Join(root, "blobs", "sha256"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(filepath.Join(root, "refs"), 0700)
	if err != nil {
		t.Fatal(err)
	}

	desc, err := createLayerFile(root)
	if err != nil {
		t.Fatal(err)
	}
	manifestStr = strings.Replace(manifestStr, "<layer_digest>", desc.Digest, 1)
	manifestStr = strings.Replace(manifestStr, "<layer_size>", strconv.FormatInt(desc.Size, 10), 1)

	desc, err = createConfigTestFile(root)
	if err != nil {
		t.Fatal(err)
	}
	manifestStr = strings.Replace(manifestStr, "<config_digest>", desc.Digest, 1)
	manifestStr = strings.Replace(manifestStr, "<config_size>", strconv.FormatInt(desc.Size, 10), 1)

	mft, err := createManifestFile(root, manifestStr)
	if err != nil {
		t.Fatal(err)
	}

	err = createRefFile(root, mft)
	if err != nil {
		t.Fatal(err)
	}

	err = createLayoutFile(root)
	if err != nil {
		t.Fatal(err)
	}

	err = ValidateLayout(root, []string{"latest"}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func createLayerFile(root string) (descriptor, error) {
	layerPath := filepath.Join(root, "blobs", "sha256", "test.tar")

	desc, err := createTarFile(layerPath, []tarContent{
		tarContent{&tar.Header{Name: "test", Size: 4, Mode: 0600}, []byte("test")},
	})
	if err != nil {
		return descriptor{}, err
	}

	err = os.Rename(layerPath, filepath.Join(root, "blobs", "sha256", desc.Digest))
	if err != nil {
		return descriptor{}, err
	}

	return descriptor{Digest: "sha256:" + desc.Digest, Size: desc.Size}, nil
}

func createConfigTestFile(root string) (descriptor, error) {
	oldpath := filepath.Join(root, "blobs", "sha256", "test-config")
	f, err := os.Create(oldpath)
	if err != nil {
		return descriptor{}, err
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewBuffer([]byte(configStr)))
	if err != nil {
		return descriptor{}, err
	}

	// generate sha256 hash
	h := sha256.New()
	size, err := io.Copy(h, bytes.NewBuffer([]byte(configStr)))
	if err != nil {
		return descriptor{}, err
	}
	digest := fmt.Sprintf("%x", h.Sum(nil))

	err = os.Rename(oldpath, filepath.Join(root, "blobs", "sha256", digest))
	if err != nil {
		return descriptor{}, err
	}
	return descriptor{Digest: "sha256:" + digest, Size: size}, nil
}

func createManifestFile(root, str string) (descriptor, error) {
	oldpath := filepath.Join(root, "blobs", "sha256", "test-manifest")
	f, err := os.Create(oldpath)
	if err != nil {
		return descriptor{}, err
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewBuffer([]byte(str)))
	if err != nil {
		return descriptor{}, err
	}

	// generate sha256 hash
	h := sha256.New()
	size, err := io.Copy(h, bytes.NewBuffer([]byte(str)))
	if err != nil {
		return descriptor{}, err
	}
	digest := fmt.Sprintf("%x", h.Sum(nil))

	err = os.Rename(oldpath, filepath.Join(root, "blobs", "sha256", digest))
	if err != nil {
		return descriptor{}, err
	}
	return descriptor{Digest: "sha256:" + digest, Size: size}, nil
}

func createRefFile(root string, mft descriptor) error {
	refpath := filepath.Join(root, "refs", "latest")
	f, err := os.Create(refpath)
	if err != nil {
		return err
	}
	defer f.Close()
	refStr = strings.Replace(refStr, "<manifest_digest>", mft.Digest, -1)
	refStr = strings.Replace(refStr, "<manifest_size>", strconv.FormatInt(mft.Size, 10), -1)
	_, err = io.Copy(f, bytes.NewBuffer([]byte(refStr)))
	return err
}

func createLayoutFile(root string) error {
	layoutPath := filepath.Join(root, "oci-layout")
	f, err := os.Create(layoutPath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewBuffer([]byte(layoutStr)))
	return err
}
