# OCI Artifact Manifest Specification

The goal of the Artifact Manifest Specification is to define content addressable artifacts in order to store them along side container images in a registry.
Like [OCI Images](manifest.md), OCI Artifacts may be referenced by the hash of their manifest.
Unlike OCI Images, OCI Artifacts are not meant to be used by any container runtime.

Examples of artifacts that may be stored along with container images are Software Bill of Materials (SBOM), Digital Signatures, Provenance data, Supply Chain Attestations, scan results, and Helm charts.

This section defines the `application/vnd.oci.artifact.manifest.v1+json` [media type](media-types.md).
For the media type(s) that this is compatible with see the [matrix](media-types.md#compatibility-matrix).

# Artifact Manifest

## *Artifact Manifest* Property Descriptions

- **`mediaType`** *string*

  This property MUST be used and contain the media type `application/vnd.oci.artifact.manifest.v1+json`.

- **`artifactType`** *string*

  This property SHOULD be used and contain the mediaType of the referenced artifact.
  If defined, the value MUST comply with [RFC 6838][rfc6838], including the [naming requirements in its section 4.2][rfc6838-s4.2], and MAY be registered with [IANA][iana].

- **`blobs`** *array of objects*

  This OPTIONAL property is an array of objects and each item in the array MUST be a [descriptor](descriptor.md).
  Each descriptor represents an artifact of any IANA mediaType.
  The list MAY be ordered for certain artifact types like scan results.

- **`refers`** *[descriptor](descriptor.md)*

  This OPTIONAL property specifies a [descriptor](descriptor.md) of another manifest.
  This value, used by the [`referrers` API](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#listing-referrers), indicates a relationship to the specified manifest.

- **`annotations`** *string-string map*

  This OPTIONAL property contains additional metadata for the artifact manifest.
  This OPTIONAL property MUST use the [annotation rules](annotations.md#rules).

  See [Pre-Defined Annotation Keys](annotations.md#pre-defined-annotation-keys).

  Annotations MAY be used to filter the response from the [`referrers` API](https://github.com/opencontainers/distribution-spec/blob/main/spec.md#listing-referrers).

## Examples

*Example showing an artifact manifest for an example SBOM referencing an image:*

```jsonc,title=Manifest&mediatype=application/vnd.oci.artifact.manifest.v1%2Bjson
{
  "mediaType": "application/vnd.oci.artifact.manifest.v1+json",
  "artifactType": "application/vnd.example.sbom.v1"
  "blobs": [
    {
      "mediaType": "application/gzip",
      "size": 123,
      "digest": "sha256:87923725d74f4bfb94c9e86d64170f7521aad8221a5de834851470ca142da630"
    }
  ],
  "refers": {
    "mediaType": "application/vnd.oci.image.manifest.v1+json",
    "size": 1234,
    "digest": "sha256:cc06a2839488b8bd2a2b99dcdc03d5cfd818eed72ad08ef3cc197aac64c0d0a0"
  },
  "annotations": {
    "org.opencontainers.artifact.created": "2022-01-01T14:42:55Z",
    "org.example.sbom.format": "json"
  }
}
```

[iana]:         https://www.iana.org/assignments/media-types/media-types.xhtml
[rfc6838]:      https://tools.ietf.org/html/rfc6838
[rfc6838-s4.2]: https://tools.ietf.org/html/rfc6838#section-4.2
