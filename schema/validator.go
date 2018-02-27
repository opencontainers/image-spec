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
	"fmt"
	"io"
	"io/ioutil"
	"regexp"

	digest "github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

// Validator wraps a media type string identifier
// and implements validation against a JSON schema.
type Validator string

type validateFunc func(r io.Reader) ValidationError

var mapValidate = map[Validator]validateFunc{
	ValidatorMediaTypeImageConfig: validateConfig,
	ValidatorMediaTypeDescriptor:  validateDescriptor,
	ValidatorMediaTypeImageIndex:  validateIndex,
	ValidatorMediaTypeManifest:    validateManifest,
}

// ValidationError contains all the errors that happened during validation.
type ValidationError struct {
	Errs []error
	Warns []error
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%v", e.Errs)
}

// Validate validates the given reader against the schema of the wrapped media type.
func (v Validator) Validate(src io.Reader) ValidationError {
	var e ValidationError
	buf, err := ioutil.ReadAll(src)
	if err != nil {
		e.Errs = append(e.Errs, errors.Wrap(err, "unable to read the document file"))
		return e
	}

	if f, ok := mapValidate[v]; ok {
		if f == nil {
			e.Errs = append(e.Errs, fmt.Errorf("internal error: mapValidate[%q] is nil", v))
			return e
		}
		err = f(bytes.NewReader(buf))
		if err != nil {
			e.Errs = append(e.Errs, err)
			return e
		}
	}

	sl := newFSLoaderFactory(schemaNamespaces, fs).New(specs[v])
	ml := gojsonschema.NewStringLoader(string(buf))

	result, err := gojsonschema.Validate(sl, ml)
	if err != nil {
		e.Errs = append(e.Errs, errors.Wrapf(
			WrapSyntaxError(bytes.NewReader(buf), err),
			"schema %s: unable to validate", v))
		return e
	}

	if result.Valid() {
		return e
	}

	errs := make([]error, 0, len(result.Errors()))
	for _, desc := range result.Errors() {
		errs = append(errs, fmt.Errorf("%s", desc))
	}

	return ValidationError{
		Errs: errs,
	}
}

type unimplemented string

func (v unimplemented) Validate(src io.Reader) error {
	return fmt.Errorf("%s: unimplemented", v)
}

func validateManifest(r io.Reader) ValidationError {
	var e ValidationError
	header := v1.Manifest{}

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		e.Errs = append(e.Errs, errors.Wrapf(err, "error reading the io stream"))
		return e
	}

	err = json.Unmarshal(buf, &header)
	if err != nil {
		e.Errs = append(e.Errs, errors.Wrap(err, "manifest format mismatch"))
		return e
	}

	if header.Config.MediaType != string(v1.MediaTypeImageConfig) {
		e.Warns = append(e.Warns, errors.Errorf("config %s has an unknown media type: %s\n", header.Config.Digest, header.Config.MediaType))
	}

	for _, layer := range header.Layers {
		if layer.MediaType != string(v1.MediaTypeImageLayer) &&
			layer.MediaType != string(v1.MediaTypeImageLayerGzip) &&
			layer.MediaType != string(v1.MediaTypeImageLayerNonDistributable) &&
			layer.MediaType != string(v1.MediaTypeImageLayerNonDistributableGzip) {
			e.Warns = append(e.Warns, errors.Errorf("layer %s has an unknown media type: %s\n", layer.Digest, layer.MediaType))
		}
	}
	return e
}

func validateDescriptor(r io.Reader) ValidationError {
	var e ValidationError
	header := v1.Descriptor{}

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		e.Errs = append(e.Errs, errors.Wrapf(err, "error reading the io stream"))
		return e
	}

	err = json.Unmarshal(buf, &header)
	if err != nil {
		e.Errs = append(e.Errs, errors.Wrap(err, "descriptor format mismatch"))
		return e
	}

	err = header.Digest.Validate()
	if err == digest.ErrDigestUnsupported {
		// we ignore unsupported algorithms
		e.Warns = append(e.Warns, errors.Errorf("unsupported digest: %q: %v\n", header.Digest, err))
		return e
	}
	e.Errs = append(e.Errs, err)
	return e
}

func validateIndex(r io.Reader) ValidationError {
	var e ValidationError
	header := v1.Index{}

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		e.Errs = append(e.Errs, errors.Wrapf(err, "error reading the io stream"))
		return e
	}

	err = json.Unmarshal(buf, &header)
	if err != nil {
		e.Errs = append(e.Errs, errors.Wrap(err, "index format mismatch"))
		return e
	}

	for _, manifest := range header.Manifests {
		if manifest.MediaType != string(v1.MediaTypeImageManifest) {
			 e.Warns = append(e.Warns, errors.Errorf("manifest %s has an unknown media type: %s\n", manifest.Digest, manifest.MediaType))
		}
		if manifest.Platform != nil {
			warns := checkPlatform(manifest.Platform.OS, manifest.Platform.Architecture)
			e.Warns = append(e.Warns, warns...)
		}

	}

	return e
}

func validateConfig(r io.Reader) ValidationError {
	var e ValidationError
	header := v1.Image{}

	buf, err := ioutil.ReadAll(r)
	if err != nil {
		e.Errs = append(e.Errs, errors.Wrapf(err, "error reading the io stream"))
		return e
	}

	err = json.Unmarshal(buf, &header)
	if err != nil {
		e.Errs = append(e.Errs, errors.Wrap(err, "config format mismatch"))
		return e
	}

	warns := checkPlatform(header.OS, header.Architecture)
	e.Warns = append(e.Warns, warns...)

	envRegexp := regexp.MustCompile(`^[^=]+=.*$`)
	for _, env := range header.Config.Env {
		if !envRegexp.MatchString(env) {
			e.Errs = append(e.Errs, errors.Errorf("unexpected env: %q", env))
			return e
		}
	}

	return e
}

func checkPlatform(OS string, Architecture string) []error {
	warns := []error{}
	validCombins := map[string][]string{
		"android":   {"arm"},
		"darwin":    {"386", "amd64", "arm", "arm64"},
		"dragonfly": {"amd64"},
		"freebsd":   {"386", "amd64", "arm"},
		"linux":     {"386", "amd64", "arm", "arm64", "ppc64", "ppc64le", "mips64", "mips64le", "s390x"},
		"netbsd":    {"386", "amd64", "arm"},
		"openbsd":   {"386", "amd64", "arm"},
		"plan9":     {"386", "amd64"},
		"solaris":   {"amd64"},
		"windows":   {"386", "amd64"}}
	for os, archs := range validCombins {
		if os == OS {
			for _, arch := range archs {
				if arch == Architecture {
					return warns
				}
			}
			warns = append(warns, errors.Errorf("combination of %q and %q is invalid.\n", OS, Architecture))
		}
	}
	return append(warns, errors.Errorf("operating system %q of the bundle is not supported yet.\n", OS))
}
