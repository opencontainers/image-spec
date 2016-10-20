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
	"fmt"
	"io"
	"io/ioutil"

	"github.com/opencontainers/go-digest"
	"github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
)

// ValidationError contains all the errors that happened during
// validation.
type ValidationError struct {
	Errs []error
}

// Validator is a template for validating a CAS blob.  The 'strict'
// parameter distinguishes between compliant blobs (which should only
// pass when strict is false) and blobs that only use features which
// the spec requires implementations to support (which should pass
// regardless of strict).
type Validator func(blob io.Reader, descriptor *v1.Descriptor, strict bool) (err error)

// Validators is a map from media types to an appropriate Validator
// function.
var Validators = map[string]Validator{
	v1.MediaTypeDescriptor:        ValidateJSONSchema,
	v1.MediaTypeImageManifestList: ValidateJSONSchema,
	v1.MediaTypeImageManifest:     ValidateManifest,
	v1.MediaTypeImageConfig:       ValidateJSONSchema,
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%v", e.Errs)
}

// Validate retrieves the appropriate Validator from Validators and
// uses it to validate the given CAS blob.  Validate uses the
// Validator template; see the Validator docs for usage information.
func Validate(blob io.Reader, descriptor *v1.Descriptor, strict bool) (err error) {
	validator, ok := Validators[descriptor.MediaType]
	if !ok {
		return fmt.Errorf("unrecognized media type %q", descriptor.MediaType)
	}
	return validator(blob, descriptor, strict)
}

// ValidateJSONSchema validates the given CAS blob against the schema
// for the descriptor's media type.  Calls ValidateByteSize and
// ValidateByteDigest as well.
func ValidateJSONSchema(blob io.Reader, descriptor *v1.Descriptor, strict bool) (err error) {
	buffer, err := ioutil.ReadAll(blob)
	if err != nil {
		return errors.Wrapf(err, "unable to read %s", descriptor.Digest)
	}

	err = ValidateByteSize(buffer, descriptor)
	if err != nil {
		return err
	}

	err = ValidateByteDigest(buffer, descriptor)
	if err != nil {
		return err
	}

	url := "file:///" + Schemas[descriptor.MediaType]
	schemaLoader := gojsonschema.NewReferenceLoaderFileSystem(url, fs)
	docLoader := gojsonschema.NewStringLoader(string(buffer))

	result, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return errors.Wrapf(
			WrapSyntaxError(bytes.NewReader(buffer), err),
			"unable to validate JSON Schema for %s", descriptor.Digest)
	}

	if result.Valid() {
		return nil
	}

	errs := make([]error, 0, len(result.Errors()))
	for _, description := range result.Errors() {
		errs = append(errs, fmt.Errorf("%s", description))
	}

	return ValidationError{
		Errs: errs,
	}
}

// ValidateByteDigest checks the digest of blob against the expected
// descriptor.Digest.
func ValidateByteDigest(blob []byte, descriptor *v1.Descriptor) (err error) {
	parsed, err := digest.Parse(descriptor.Digest)
	if err != nil {
		return err
	}
	algorithm := parsed.Algorithm()
	if !algorithm.Available() {
		return fmt.Errorf("unsupported digest algorithm for %s", descriptor.Digest)
	}
	actualDigest := algorithm.FromBytes(blob).String()
	if actualDigest != descriptor.Digest {
		return fmt.Errorf("unexpected digest for %s: %s", descriptor.Digest, actualDigest)
	}

	return nil
}

// ValidateByteSize checks the size of blob against the expected
// descriptor.Size.  This isn't very complicated; the function is
// mostly useful for generating consistent error messages.
func ValidateByteSize(blob []byte, descriptor *v1.Descriptor) (err error) {
	if descriptor.Size > 0 && int64(len(blob)) != descriptor.Size {
		return fmt.Errorf("unexpected size for %s: %d != %d", descriptor.Digest, len(blob), descriptor.Size)
	}

	return nil
}
