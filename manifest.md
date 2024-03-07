# OCI Image Manifest Specification

There are three main goals of the Image Manifest Specification.
The first goal is content-addressable images, by supporting an image model where the image's configuration can be hashed to generate a unique ID for the image and its components.
The second goal is to allow multi-architecture images, through a "fat manifest" which references image manifests for platform-specific versions of an image.
In OCI, this is codified in an [image index](image-index.md).
The third goal is to be [translatable](conversion.md) to the [OCI Runtime Specification](https://github.com/opencontainers/runtime-spec).

This section defines the `application/vnd.oci.image.manifest.v1+json` [media type](media-types.md).
For the media type(s) that this is compatible with see the [matrix](media-types.md#compatibility-matrix).

## Image Manifest

Unlike the [image index](image-index.md), which contains information about a set of images that can span a variety of architectures and operating systems, an image manifest provides a configuration and set of layers for a single container image for a specific architecture and operating system.

## _Image Manifest_ Property Descriptions

- **`schemaVersion`** *int*

  This REQUIRED property specifies the image manifest schema version.
  For this version of the specification, this MUST be `2` to ensure backward compatibility with older versions of Docker. The value of this field will not change. This field MAY be removed in a future version of the specification.

- **`mediaType`** *string*

  This property SHOULD be used and [remain compatible](media-types.md#compatibility-matrix) with earlier versions of this specification and with other similar external formats.
  When used, this field MUST contain the media type `application/vnd.oci.image.manifest.v1+json`.
  This field usage differs from the [descriptor](descriptor.md#properties) use of `mediaType`.

- **`artifactType`** *string*

  This OPTIONAL property contains the type of an artifact when the manifest is used for an artifact.
  This MUST be set when `config.mediaType` is set to the [empty value](#guidance-for-an-empty-descriptor).
  If defined, the value MUST comply with [RFC 6838][rfc6838], including the [naming requirements in its section 4.2][rfc6838-s4.2], and MAY be registered with [IANA][iana].
  Implementations storing or copying image manifests MUST NOT error on encountering an `artifactType` that is unknown to the implementation.

- **`config`** *[descriptor](descriptor.md)*

  This REQUIRED property references a configuration object for a container, by digest.
  Beyond the [descriptor requirements](descriptor.md#properties), the value has the following additional restrictions:

  - **`mediaType`** *string*

    This [descriptor property](descriptor.md#properties) has additional restrictions for `config`.

    Implementations MUST NOT attempt to parse the referenced content if this media type is unknown and instead consider the referenced content as arbitrary binary data (e.g.: as `application/octet-stream`).

    Implementations storing or copying image manifests MUST NOT error on encountering a value that is unknown to the implementation.

    Implementations MUST support at least the following media types:

    - [`application/vnd.oci.image.config.v1+json`](config.md)

    Manifests for container images concerned with portability SHOULD use one of the above media types.
    Manifests for artifacts concerned with portability SHOULD use `config.mediaType` as described in [Guidelines for Artifact Usage](#guidelines-for-artifact-usage).

    If the manifest uses a different media type than the above, it MUST comply with [RFC 6838][rfc6838], including the [naming requirements in its section 4.2][rfc6838-s4.2], and MAY be registered with [IANA][iana].

  To set an effectively null or empty config and maintain portability see the [guidance for an empty descriptor](#guidance-for-an-empty-descriptor) below, and `DescriptorEmptyJSON` of the reference code.

  If this image manifest will be "runnable" by a runtime of some kind, it is strongly recommended to ensure it includes enough data to be unique (such as the `rootfs` and `diff_ids` included in `application/vnd.oci.image.config.v1+json`) so that it has a unique [`ImageID`](config.md#imageid).

- **`layers`** *array of objects*

  Each item in the array MUST be a [descriptor](descriptor.md).
  For portability, `layers` SHOULD have at least one entry.
  See the [guidance for an empty descriptor](#guidance-for-an-empty-descriptor) below, and `DescriptorEmptyJSON` of the reference code.

  When the `config.mediaType` is set to `application/vnd.oci.image.config.v1+json`, the following additional restrictions apply:

  - The array MUST have the base layer at index 0.
  - Subsequent layers MUST then follow in stack order (i.e. from `layers[0]` to `layers[len(layers)-1]`).
  - The final filesystem layout MUST match the result of [applying](layer.md#applying-changesets) the layers to an empty directory.
  - The [ownership, mode, and other attributes](layer.md#file-attributes) of the initial empty directory are unspecified.

  Beyond the [descriptor requirements](descriptor.md#properties), the value has the following additional restrictions:

  - **`mediaType`** *string*

    This [descriptor property](descriptor.md#properties) has additional restrictions for `layers[]`.
    Implementations MUST support at least the following media types:

    - [`application/vnd.oci.image.layer.v1.tar`](layer.md)
    - [`application/vnd.oci.image.layer.v1.tar+gzip`](layer.md#gzip-media-types)
    - [`application/vnd.oci.image.layer.nondistributable.v1.tar`](layer.md#non-distributable-layers)
    - [`application/vnd.oci.image.layer.nondistributable.v1.tar+gzip`](layer.md#gzip-media-types)

    Manifests concerned with portability SHOULD use one of the above media types.
    Implementations storing or copying image manifests MUST NOT error on encountering a `mediaType` that is unknown to the implementation.

    Entries in this field will frequently use the `+gzip` types.

    If the manifest uses a different media type than the above, it MUST comply with [RFC 6838][rfc6838], including the [naming requirements in its section 4.2][rfc6838-s4.2], and MAY be registered with [IANA][iana].

  See [Guidelines for Artifact Usage](#guidelines-for-artifact-usage) for other uses of the `layers`.

- **`subject`** *[descriptor](descriptor.md)*

  This OPTIONAL property specifies a [descriptor](descriptor.md) of another manifest.
  This value defines a weak association to a separate [Merkle Directed Acyclic Graph (DAG)][dag] structure, and is used by the [`referrers` API][referrers-api] to include this manifest in the list of responses for the subject digest.

- **`annotations`** *string-string map*

  This OPTIONAL property contains arbitrary metadata for the image manifest.
  This OPTIONAL property MUST use the [annotation rules](annotations.md#rules).

  See [Pre-Defined Annotation Keys](annotations.md#pre-defined-annotation-keys).

## Example Image Manifest

*Example showing an image manifest:*

```json,title=Manifest&mediatype=application/vnd.oci.image.manifest.v1%2Bjson
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "digest": "sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7",
    "size": 7023
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:9834876dcfb05cb167a5c24953eba58c4ac89b1adf57f28f2f9d09af107ee8f0",
      "size": 32654
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:3c3a4604a545cdc127456d94e421cd355bca5b528f4a9c1905b15da2eb4a4c6b",
      "size": 16724
    },
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "digest": "sha256:ec4b8955958665577945c89419d1af06b5f7636b4ac3da7f12184802ad867736",
      "size": 73109
    }
  ],
  "subject": {
    "mediaType": "application/vnd.oci.image.manifest.v1+json",
    "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
    "size": 7682
  },
  "annotations": {
    "com.example.key1": "value1",
    "com.example.key2": "value2"
  }
}
```

## Guidance for an Empty Descriptor

*Implementers note*: The following is considered GUIDANCE for portability.

Parts of the spec necessitate including a descriptor to a blob where some implementations of artifacts do not have associated content.
While an empty blob (`size` of 0) may be preferable, practice has shown that not to be ubiquitously supported.
The media type `application/vnd.oci.empty.v1+json` (`MediaTypeEmptyJSON`) has been specified for a descriptor that has no content for the implementation.
The blob payload is the most minimal content that is still a valid JSON object: `{}` (`size` of 2).
The blob digest of `{}` is `sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a`.
The data field is optional, and if included is the base64 encoding of `{}`: `e30=`.

The resulting descriptor shown here is also defined in reference code as `DescriptorEmptyJSON`:

```json,title=empty%20config&mediatype=application/vnd.oci.descriptor.v1%2Bjson
{
  "mediaType": "application/vnd.oci.empty.v1+json",
  "digest": "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
  "size": 2,
  "data": "e30="
}
```

## Guidelines for Artifact Usage

Content other than OCI container images MAY be packaged using the image manifest.
When this is done, the `config.mediaType` value MUST be set to a value specific to the artifact type or the [empty value](#guidance-for-an-empty-descriptor).
If the `config.mediaType` is set to the empty value, the `artifactType` MUST be defined.
If the artifact does not need layers, a single layer SHOULD be included with a non-zero size.
The suggested content for an unused `layers` array is the [empty descriptor](#guidance-for-an-empty-descriptor).

The design of the artifact depends on what content is being packaged with the artifact.
The decision tree below and the associated examples MAY be used to design new artifacts:

1. Does the artifact consist of at least one file or blob?
   If yes, continue to 2.
   If no, specify the `artifactType`, and set the `config` and a single `layers` element to the empty descriptor value.
   Here is an example of this with annotations included:

   ```json,title=Minimal%20artifact&mediatype=application/vnd.oci.image.manifest.v1%2Bjson
   {
     "schemaVersion": 2,
     "mediaType": "application/vnd.oci.image.manifest.v1+json",
     "artifactType": "application/vnd.example+type",
     "config": {
       "mediaType": "application/vnd.oci.empty.v1+json",
       "digest": "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
       "size": 2
     },
     "layers": [
       {
         "mediaType": "application/vnd.oci.empty.v1+json",
         "digest": "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
         "size": 2
       }
     ],
     "annotations": {
       "oci.opencontainers.image.created": "2023-01-02T03:04:05Z",
       "com.example.data": "payload"
     }
   }
   ```

2. Does the artifact have additional JSON formatted metadata as configuration?
   If yes, continue to 3.
   If no, specify the `artifactType`, include the artifact in the `layers`, and set `config` to the empty descriptor value.
   Here is an example of this with a single layer:

   ```json,title=Artifact%20without%20config&mediatype=application/vnd.oci.image.manifest.v1%2Bjson
   {
     "schemaVersion": 2,
     "mediaType": "application/vnd.oci.image.manifest.v1+json",
     "artifactType": "application/vnd.example+type",
     "config": {
       "mediaType": "application/vnd.oci.empty.v1+json",
       "digest": "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
       "size": 2
     },
     "layers": [
       {
         "mediaType": "application/vnd.example+type",
         "digest": "sha256:e258d248fda94c63753607f7c4494ee0fcbe92f1a76bfdac795c9d84101eb317",
         "size": 1234
       }
     ]
   }
   ```

3. For artifacts with a config blob, specify the `artifactType` to a common value for your artifact tooling, specify the `config` with the metadata for this artifact, and include the artifact in the `layers`.
   Here is an example of this:

   ```json,title=Artifact%20with%20config&mediatype=application/vnd.oci.image.manifest.v1%2Bjson
   {
     "schemaVersion": 2,
     "mediaType": "application/vnd.oci.image.manifest.v1+json",
     "artifactType": "application/vnd.example+type",
     "config": {
       "mediaType": "application/vnd.example.config.v1+json",
       "digest": "sha256:5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03",
       "size": 123
     },
     "layers": [
       {
         "mediaType": "application/vnd.example.data.v1.tar+gzip",
         "digest": "sha256:e258d248fda94c63753607f7c4494ee0fcbe92f1a76bfdac795c9d84101eb317",
         "size": 1234
       }
     ]
   }
   ```

_Implementers note:_ artifacts have historically been created without an `artifactType` field, and tooling to work with artifacts should fallback to the `config.mediaType` value.

[dag]:           https://en.wikipedia.org/wiki/Merkle_tree
[iana]:          https://www.iana.org/assignments/media-types/media-types.xhtml
[referrers-api]: https://github.com/opencontainers/distribution-spec/blob/main/spec.md#listing-referrers
[rfc6838]:       https://tools.ietf.org/html/rfc6838
[rfc6838-s4.2]:  https://tools.ietf.org/html/rfc6838#section-4.2
