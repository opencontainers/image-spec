<!--[metadata]>
+++
draft = true
+++
<![end-metadata]-->

# OpenContainers Image Manifest Specification

There are three main goals of the Image Manifest Specification.
The first goal is content-addressable images, by supporting an image model where the image's configuration can be hashed to generate a unique ID for the image and its components.
The second goal is to allow multi-architecture images, through a "fat manifest" which references image manifests for platform-specific versions of an image.
The third goal is to be translatable to the [OpenContainers/runtime-spec](https://github.com/opencontainers/runtime-spec)


# Manifest List

The manifest list is the "fat manifest" which points to specific image manifests for one or more platforms.
While the use of a manifest list is OPTIONAL for image providers, image consumers SHOULD be prepared to process them.
A client will distinguish a manifest list from an image manifest based on the Content-Type returned in the HTTP response.

## *Manifest List* Property Descriptions

- **`schemaVersion`** *int*

  This REQUIRED property specifies the image manifest schema version.
  For this version of the specification, this MUST be `2` to ensure backward compatibility with older versions of Docker. The value of this field will not change. This field MAY be removed in a future version of the specification.

- **`mediaType`** *string*

  This REQUIRED property contains the MIME type of the manifest list.
  For this version of the specification, this MUST be set to `application/vnd.oci.image.manifest.list.v1+json`.

- **`manifests`** *array*

  This REQUIRED property contains a list of manifests for specific platforms.
  While the property MUST be present, the size of the array MAY be zero.

  Each object in the manifest is a [descriptor](descriptor.md) with the following additional properties:

  - **`platform`** *object*

    This REQUIRED property describes the platform which the image in the manifest runs on.
    A full list of valid operating system and architecture values are listed in the [Go language documentation for `$GOOS` and `$GOARCH`](https://golang.org/doc/install/source#environment)

    - **`architecture`** *string*

        This REQUIRED property specified the CPU architecture, for example `amd64` or `ppc64le`.

    - **`os`** *string*

        This REQUIRED property specifies the operating system, for example `linux` or `windows`.

    - **`os.version`** *string*

        This OPTIONAL property specifies the operating system version, for example `10.0.10586`.

    - **`os.features`** *array*

        This OPTIONAL property specifies an array of strings, each specifying a mandatory OS feature (for example on Windows `win32k`).

    - **`variant`** *string*

        This OPTIONAL property specifies the variant of the CPU, for example `armv6l` to specify a particular CPU variant of the ARM CPU.

    - **`features`** *array*

        This OPTIONAL property specifies an array of strings, each specifying a mandatory CPU feature (for example `sse4` or `aes`).

- **`annotations`** *string-string hashmap*

    This OPTIONAL property contains arbitrary metadata for the manifest list.
    Annotations is a key-value, unordered hashmap.
    Keys are unique, and best practice is to namespace the keys.
    Common annotation keys include:
    * **created** date on which the image was built (string, timestamps type)
    * **authors** contact details of the people or organization responsible for the image (freeform string)
    * **homepage** URL to find more information on the image (string, must be a URL with scheme HTTP or HTTPS)
    * **documentation** URL to get documentation on the image (string, must be a URL with scheme HTTP or HTTPS)


## Example Manifest List

*Example showing a simple manifest list pointing to image manifests for two platforms:*
```json,title=Manifest%20List&mediatype=application/vnd.oci.image.manifest.list.v1%2Bjson
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.list.v1+json",
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
        "os": "linux",
        "features": [
          "sse4"
        ]
      }
    }
  ],
  "annotations": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

# Image Manifest

Unlike the [Manifest List](#manifest-list), which contains information about a set of images that can span a variety of architectures and operating systems, an image manifest provides a configuration and set of layers for a single container image for a specific architecture and operating system.

## *Image Manifest* Property Descriptions

- **`schemaVersion`** *int*

  This REQUIRED property specifies the image manifest schema version.
  For this version of the specification, this MUST be `2` to ensure backward compatibility with older versions of Docker. The value of this field will not change. This field MAY be removed in a future version of the specification.

- **`mediaType`** *string*

    This REQUIRED property contains the MIME type of the image manifest.
    For this version of the specification, this MUST be set to `application/vnd.oci.image.manifest.v1+json`.

- **`config`** *[descriptor](descriptor.md)*

    This REQUIRED property references a configuration object for a container, by digest.
    The referenced configuration object is a JSON blob that the runtime uses to set up the container, see [Image JSON Description](serialization.md#image-json-description).

- **`layers`** *array*

    Each item in the array MUST be a [descriptor](descriptor.md).
    The array MUST have the base image at index 0.
	Subsequent layers MUST then follow in the order in which they are to be layered on top of each other.
    The algorithm to create the final unpacked filesystem layout MUST be to first unpack the layer at index 0, then index 1, and so on.

- **`annotations`** *hashmap*

    This OPTIONAL property contains arbitrary metadata for the manifest list.
    Annotations is a key-value, unordered hashmap.
    Keys MUST be unique within the hashmap, and best practice is to namespace the keys.
    Common annotation keys include:
    * **created** date on which the image was built (string, timestamps type)
    * **authors** contact details of the people or organization responsible for the image (freeform string)
    * **homepage** URL to find more information on the image (string, must be a URL with scheme HTTP or HTTPS)
    * **documentation** URL to get documentation on the image (string, must be a URL with scheme HTTP or HTTPS)
    * **source** URL to get the source code for the binary files in the image (string, must be a URL with scheme HTTP or HTTPS)


## Example Image Manifest

*Example showing an image manifest:*
```json,title=Manifest&mediatype=application/vnd.oci.image.manifest.v1%2Bjson
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.serialization.config.v1+json",
    "size": 7023,
    "digest": "sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.serialization.rootfs.tar.gzip",
      "size": 32654,
      "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f"
    },
    {
      "mediaType": "application/vnd.oci.image.serialization.rootfs.tar.gzip",
      "size": 16724,
      "digest": "sha256:3c3a4604a545cdc127456d94e421cd355bca5b528f4a9c1905b15da2eb4a4c6b"
    },
    {
      "mediaType": "application/vnd.oci.image.serialization.rootfs.tar.gzip",
      "size": 73109,
      "digest": "sha256:ec4b8955958665577945c89419d1af06b5f7636b4ac3da7f12184802ad867736"
    }
  ],
  "annotations": {
    "key1": "value1",
    "key2": "value2"
  }
}
```
