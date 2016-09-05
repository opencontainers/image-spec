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
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/opencontainers/image-spec/image/cas"
	caslayout "github.com/opencontainers/image-spec/image/cas/layout"
	"github.com/opencontainers/image-spec/image/refs"
	refslayout "github.com/opencontainers/image-spec/image/refs/layout"
	"golang.org/x/net/context"
)

// Validate validates the given reference.
func Validate(ctx context.Context, path, ref string) error {
	refEngine, err := refslayout.NewEngine(ctx, path)
	if err != nil {
		return err
	}
	defer refEngine.Close()

	casEngine, err := caslayout.NewEngine(ctx, path)
	if err != nil {
		return err
	}
	defer casEngine.Close()

	return validate(ctx, refEngine, casEngine, ref)
}

func validate(ctx context.Context, refEngine refs.Engine, casEngine cas.Engine, ref string) error {
	descriptor, err := refEngine.Get(ctx, ref)
	if err != nil {
		return err
	}

	err = validateDescriptor(ctx, casEngine, descriptor)
	if err != nil {
		return err
	}

	m, err := findManifest(ctx, casEngine, descriptor)
	if err != nil {
		return err
	}

	return m.validate(ctx, casEngine)
}

// Unpack unpacks the given reference to a destination directory.
func Unpack(ctx context.Context, path, dest, ref string) error {
	refEngine, err := refslayout.NewEngine(ctx, path)
	if err != nil {
		return err
	}
	defer refEngine.Close()

	casEngine, err := caslayout.NewEngine(ctx, path)
	if err != nil {
		return err
	}
	defer casEngine.Close()

	return unpack(ctx, refEngine, casEngine, dest, ref)
}

func unpack(ctx context.Context, refEngine refs.Engine, casEngine cas.Engine, dest, ref string) error {
	descriptor, err := refEngine.Get(ctx, ref)
	if err != nil {
		return err
	}

	err = validateDescriptor(ctx, casEngine, descriptor)
	if err != nil {
		return err
	}

	m, err := findManifest(ctx, casEngine, descriptor)
	if err != nil {
		return err
	}

	if err = m.validate(ctx, casEngine); err != nil {
		return err
	}

	return m.unpack(ctx, casEngine, dest)
}

// CreateRuntimeBundle creates an OCI runtime bundle in the given
// destination.
func CreateRuntimeBundle(ctx context.Context, path, dest, ref, rootfs string) error {
	refEngine, err := refslayout.NewEngine(ctx, path)
	if err != nil {
		return err
	}
	defer refEngine.Close()

	casEngine, err := caslayout.NewEngine(ctx, path)
	if err != nil {
		return err
	}
	defer casEngine.Close()

	return createRuntimeBundle(ctx, refEngine, casEngine, dest, ref, rootfs)
}

func createRuntimeBundle(ctx context.Context, refEngine refs.Engine, casEngine cas.Engine, dest, ref, rootfs string) error {
	descriptor, err := refEngine.Get(ctx, ref)
	if err != nil {
		return err
	}

	err = validateDescriptor(ctx, casEngine, descriptor)
	if err != nil {
		return err
	}

	m, err := findManifest(ctx, casEngine, descriptor)
	if err != nil {
		return err
	}

	if err = m.validate(ctx, casEngine); err != nil {
		return err
	}

	c, err := findConfig(ctx, casEngine, m.Config)
	if err != nil {
		return err
	}

	err = m.unpack(ctx, casEngine, filepath.Join(dest, rootfs))
	if err != nil {
		return err
	}

	spec, err := c.runtimeSpec(rootfs)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(dest, "config.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(spec)
}
