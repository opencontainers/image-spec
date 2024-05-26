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

package schema

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"

	digest "github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Validator wraps a media type string identifier and implements validation against a JSON schema.
type Validator string

// ValidationError contains all the errors that happened during validation.
//
// Deprecated: this is no longer used by [Validator].
type ValidationError struct {
	Errs []error
}

// Error returns the error message.
//
// Deprecated: this is no longer used by [Validator].
func (e ValidationError) Error() string {
	return fmt.Sprintf("%v", e.Errs)
}

// Validate validates the given reader against the schema of the wrapped media type.
func (v Validator) Validate(src io.Reader) error {
	// run the media type specific validation
	if fn, ok := validateByMediaType[v]; ok {
		if fn == nil {
			return fmt.Errorf("internal error: mapValidate is nil for %s", string(v))
		}
		// buffer the src so the media type validation and the schema validation can both read it
		buf, err := io.ReadAll(src)
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		src = bytes.NewReader(buf)
		err = fn(buf)
		if err != nil {
			return err
		}
	}

	// json schema validation
	return v.validateSchema(src)
}

func (v Validator) validateSchema(src io.Reader) error {
	if _, ok := specs[v]; !ok {
		return fmt.Errorf("no validator available for %s", string(v))
	}

	c := jsonschema.NewCompiler()

	// load the schema files from the embedded FS
	dir, err := specFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("spec embedded directory could not be loaded: %w", err)
	}
	for _, file := range dir {
		if file.IsDir() {
			continue
		}
		specBuf, err := specFS.ReadFile(file.Name())
		if err != nil {
			return fmt.Errorf("could not read spec file %s: %w", file.Name(), err)
		}
		err = c.AddResource(file.Name(), bytes.NewReader(specBuf))
		if err != nil {
			return fmt.Errorf("failed to add spec file %s: %w", file.Name(), err)
		}
		if len(specURLs[file.Name()]) == 0 {
			fmt.Fprintf(os.Stderr, "warning: spec file has no aliases: %s", file.Name())
		}
		for _, specURL := range specURLs[file.Name()] {
			err = c.AddResource(specURL, bytes.NewReader(specBuf))
			if err != nil {
				return fmt.Errorf("failed to add spec file %s as url %s: %w", file.Name(), specURL, err)
			}
		}
	}

	// compile based on the type of validator
	schema, err := c.Compile(specs[v])
	if err != nil {
		return fmt.Errorf("failed to compile schema %s: %w", string(v), err)
	}

	// read in the user input and validate
	var input interface{}
	err = json.NewDecoder(src).Decode(&input)
	if err != nil {
		return fmt.Errorf("unable to parse json to validate: %w", err)
	}
	err = schema.Validate(input)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

type validateFunc func([]byte) error

var validateByMediaType = map[Validator]validateFunc{
	ValidatorMediaTypeImageConfig: validateConfig,
	ValidatorMediaTypeDescriptor:  validateDescriptor,
	ValidatorMediaTypeImageIndex:  validateIndex,
	ValidatorMediaTypeManifest:    validateManifest,
}

func validateManifest(buf []byte) error {
	header := v1.Manifest{}

	err := json.Unmarshal(buf, &header)
	if err != nil {
		return fmt.Errorf("manifest format mismatch: %w", err)
	}

	if header.Config.MediaType != string(v1.MediaTypeImageConfig) {
		fmt.Printf("warning: config %s has an unknown media type: %s\n", header.Config.Digest, header.Config.MediaType)
	}

	for _, layer := range header.Layers {
		if layer.MediaType != string(v1.MediaTypeImageLayer) &&
			layer.MediaType != string(v1.MediaTypeImageLayerGzip) &&
			layer.MediaType != string(v1.MediaTypeImageLayerZstd) &&
			layer.MediaType != string(v1.MediaTypeImageLayerNonDistributable) && //nolint:staticcheck
			layer.MediaType != string(v1.MediaTypeImageLayerNonDistributableGzip) && //nolint:staticcheck
			layer.MediaType != string(v1.MediaTypeImageLayerNonDistributableZstd) { //nolint:staticcheck
			fmt.Printf("warning: layer %s has an unknown media type: %s\n", layer.Digest, layer.MediaType)
		}
	}
	return nil
}

func validateDescriptor(buf []byte) error {
	header := v1.Descriptor{}

	err := json.Unmarshal(buf, &header)
	if err != nil {
		return fmt.Errorf("descriptor format mismatch: %w", err)
	}

	err = header.Digest.Validate()
	if errors.Is(err, digest.ErrDigestUnsupported) {
		// we ignore unsupported algorithms
		fmt.Printf("warning: unsupported digest: %q: %v\n", header.Digest, err)
		return nil
	}
	return err
}

func validateIndex(buf []byte) error {
	header := v1.Index{}

	err := json.Unmarshal(buf, &header)
	if err != nil {
		return fmt.Errorf("index format mismatch: %w", err)
	}

	for _, manifest := range header.Manifests {
		if manifest.MediaType != string(v1.MediaTypeImageManifest) {
			fmt.Printf("warning: manifest %s has an unknown media type: %s\n", manifest.Digest, manifest.MediaType)
		}
		if manifest.Platform != nil {
			checkPlatform(manifest.Platform.OS, manifest.Platform.Architecture)
			checkArchitecture(manifest.Platform.Architecture, manifest.Platform.Variant)
		}

	}

	return nil
}

func validateConfig(buf []byte) error {
	header := v1.Image{}

	err := json.Unmarshal(buf, &header)
	if err != nil {
		return fmt.Errorf("config format mismatch: %w", err)
	}

	checkPlatform(header.OS, header.Architecture)
	checkArchitecture(header.Architecture, header.Variant)

	envRegexp := regexp.MustCompile(`^[^=]+=.*$`)
	for _, e := range header.Config.Env {
		if !envRegexp.MatchString(e) {
			return fmt.Errorf("unexpected env: %q", e)
		}
	}

	return nil
}

func checkArchitecture(Architecture string, Variant string) {
	validCombins := map[string][]string{
		"arm":      {"", "v6", "v7", "v8"},
		"arm64":    {"", "v8"},
		"386":      {""},
		"amd64":    {""},
		"ppc64":    {""},
		"ppc64le":  {""},
		"mips64":   {""},
		"mips64le": {""},
		"s390x":    {""},
		"riscv64":  {""},
	}
	for arch, variants := range validCombins {
		if arch == Architecture {
			for _, variant := range variants {
				if variant == Variant {
					return
				}
			}
			fmt.Printf("warning: combination of architecture %q and variant %q is not valid.\n", Architecture, Variant)
		}
	}
	fmt.Printf("warning: architecture %q is not supported yet.\n", Architecture)
}

func checkPlatform(OS string, Architecture string) {
	validCombins := map[string][]string{
		"android":   {"arm"},
		"darwin":    {"386", "amd64", "arm", "arm64"},
		"dragonfly": {"amd64"},
		"freebsd":   {"386", "amd64", "arm"},
		"linux":     {"386", "amd64", "arm", "arm64", "ppc64", "ppc64le", "mips64", "mips64le", "s390x", "riscv64"},
		"netbsd":    {"386", "amd64", "arm"},
		"openbsd":   {"386", "amd64", "arm"},
		"plan9":     {"386", "amd64"},
		"solaris":   {"amd64"},
		"windows":   {"386", "amd64"}}
	for os, archs := range validCombins {
		if os == OS {
			for _, arch := range archs {
				if arch == Architecture {
					return
				}
			}
			fmt.Printf("warning: combination of os %q and architecture %q is invalid.\n", OS, Architecture)
		}
	}
	fmt.Printf("warning: operating system %q of the bundle is not supported yet.\n", OS)
}
