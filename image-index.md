# OCI Image Index Specification

The image index is a higher-level manifest which points to specific [image manifests](manifest.md), ideal for one or more platforms.
While the use of an image index is OPTIONAL for image providers, image consumers SHOULD be prepared to process them.

This section defines the `application/vnd.oci.image.index.v1+json` [media type](media-types.md).

For the media type(s) that this document is compatible with, see the [matrix][matrix].

## _Image Index_ Property Descriptions

- **`schemaVersion`** *int*

  This REQUIRED property specifies the image manifest schema version.
  For this version of the specification, this MUST be `2` to ensure backward compatibility with older versions of Docker.
  The value of this field will not change.
  This field MAY be removed in a future version of the specification.

- **`mediaType`** *string*

  This property SHOULD be used and [remain compatible][matrix] with earlier versions of this specification and with other similar external formats.
  When used, this field MUST contain the media type `application/vnd.oci.image.index.v1+json`.
  This field usage differs from the [descriptor](descriptor.md#properties) use of `mediaType`.

- **`artifactType`** *string*

  This OPTIONAL property contains the type of an artifact when the manifest is used for an artifact.
  If defined, the value MUST comply with [RFC 6838][rfc6838], including the [naming requirements in its section 4.2][rfc6838-s4.2], and MAY be registered with [IANA][iana].

- **`manifests`** *array of objects*

  This REQUIRED property contains a list of [manifests](manifest.md) for specific platforms.
  While this property MUST be present, the size of the array MAY be zero.

  Each object in `manifests` includes a set of [descriptor properties](descriptor.md#properties) with the following additional properties and restrictions:

  - **`mediaType`** *string*

    This [descriptor property](descriptor.md#properties) has additional restrictions for `manifests`.
    Implementations MUST support at least the following media types:

    - [`application/vnd.oci.image.manifest.v1+json`](manifest.md)

    Also, implementations SHOULD support the following media types:

    - `application/vnd.oci.image.index.v1+json` (nested index)

    Image indexes concerned with portability SHOULD use one of the above media types.
    Future versions of the spec MAY use a different mediatype (i.e. a new versioned format).
    An encountered `mediaType` that is unknown to the implementation MUST NOT generate an error.

  - **`platform`** *object*

    This OPTIONAL property describes the minimum runtime requirements of the image.
    This property SHOULD be present if its target is platform-specific.

    - **`architecture`** *string*

      This REQUIRED property specifies the CPU architecture.
      Image indexes SHOULD use, and implementations SHOULD understand, values listed in the Go Language document for [`GOARCH`][go-environment2].

    - **`os`** *string*

      This REQUIRED property specifies the operating system.
      Image indexes SHOULD use, and implementations SHOULD understand, values listed in the Go Language document for [`GOOS`][go-environment2].

    - **`os.version`** *string*

      This OPTIONAL property specifies the version of the operating system targeted by the referenced blob.
      Implementations MAY refuse to use manifests where `os.version` is not known to work with the host OS version.
      Valid values are implementation-defined. e.g. `10.0.14393.1066` on `windows`.

    - **`os.features`** *array of strings*

      This OPTIONAL property specifies an array of strings, each specifying a mandatory OS feature.
      When `os` is `windows`, image indexes SHOULD use, and implementations SHOULD understand the following values:

      - `win32k`: image requires `win32k.sys` on the host (Note: `win32k.sys` is missing on Nano Server)

      When `os` is not `windows`, values are implementation-defined and SHOULD be submitted to this specification for standardization.

    - **`variant`** *string*

      This OPTIONAL property specifies the variant of the CPU.
      Image indexes SHOULD use, and implementations SHOULD understand, `variant` values listed in the [Platform Variants](#platform-variants) table.

    - **`features`** *array of strings*

        This property is RESERVED for future versions of the specification.

  If multiple manifests match a client or runtime's requirements, the first matching entry SHOULD be used.

- **`subject`** *[descriptor](descriptor.md)*

    This OPTIONAL property specifies a [descriptor](descriptor.md) of another manifest.
    This value defines a weak association to a separate [Merkle Directed Acyclic Graph (DAG)][dag] structure, and is used by the [`referrers` API][referrers-api] to include this manifest in the list of responses for the subject digest.

- **`annotations`** *string-string map*

    This OPTIONAL property contains arbitrary metadata for the image index.
    This OPTIONAL property MUST use the [annotation rules](annotations.md#rules).

    See [Pre-Defined Annotation Keys](annotations.md#pre-defined-annotation-keys).

## Platform Variants

When the variant of the CPU is not listed in the table, values are implementation-defined and SHOULD be submitted to this specification for standardization.
These values SHOULD match (or be similar to) their analog listed in [the Go Language document][go-environment2].

| ISA/ABI    | `architecture` | `variant`             | Go analog   |
|------------|----------------|-----------------------|-------------|
| ARM 32-bit | `arm`          | `v6`, `v7`, `v8`      | `GOARM`     |
| ARM 64-bit | `arm64`        | `v8`, `v8.1`, …       | `GOARM64`   |
| POWER8+    | `ppc64le`      | `power8`, `power9`, … | `GOPPC64`   |
| RISC-V     | `riscv64`      | `rva20u64`, …         | `GORISCV64` |
| x86-64     | `amd64`        | `v1`, `v2`, `v3`, …   | `GOAMD64`   |

## Example Image Index

*Example showing a simple image index pointing to image manifests for two platforms:*

```json,title=Image%20Index&mediatype=application/vnd.oci.image.index.v1%2Bjson
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 7143,
      "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
      "platform": {
        "architecture": "ppc64le",
        "os": "linux"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 7682,
      "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
      "platform": {
        "architecture": "amd64",
        "os": "linux"
      }
    }
  ],
  "annotations": {
    "com.example.key1": "value1",
    "com.example.key2": "value2"
  }
}
```

## Example Image Index with multiple media types

*Example showing an image index pointing to manifests with multiple media types:*

```json,title=Image%20Index&mediatype=application/vnd.oci.image.index.v1%2Bjson
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 7143,
      "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
      "platform": {
        "architecture": "ppc64le",
        "os": "linux"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.index.v1+json",
      "size": 7682,
      "digest": "sha256:601570aaff1b68a61eb9c85b8beca1644e698003e0cdb5bce960f193d265a8b7"
    }
  ],
  "annotations": {
    "com.example.key1": "value1",
    "com.example.key2": "value2"
  }
}
```

[dag]:             https://en.wikipedia.org/wiki/Merkle_tree
[go-environment2]: https://golang.org/doc/install/source#environment
[iana]:            https://www.iana.org/assignments/media-types/media-types.xhtml
[matrix]:          media-types.md#compatibility-matrix
[referrers-api]:   https://github.com/opencontainers/distribution-spec/blob/main/spec.md#listing-referrers
[rfc6838]:         https://tools.ietf.org/html/rfc6838
[rfc6838-s4.2]:    https://tools.ietf.org/html/rfc6838#section-4.2
