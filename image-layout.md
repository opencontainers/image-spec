# OCI Image Layout Specification

- The OCI Image Layout is the directory structure for OCI content-addressable blobs and [location-addressable](https://en.wikipedia.org/wiki/Content-addressable_storage#Content-addressed_vs._location-addressed) references (refs).
- This layout MAY be used in a variety of different transport mechanisms: archive formats (e.g. tar, zip), shared filesystem environments (e.g. nfs), or networked file fetching (e.g. http, ftp, rsync).

Given an image layout and a ref, a tool can create an [OCI Runtime Specification bundle](https://github.com/opencontainers/runtime-spec/blob/main/bundle.md) by:

- Following the ref to find a [manifest](manifest.md#image-manifest), possibly via an [image index](image-index.md)
- [Applying the filesystem layers](layer.md#applying-changesets) in the specified order
- Converting the [image configuration](config.md) into an [OCI Runtime Specification `config.json`](https://github.com/opencontainers/runtime-spec/blob/main/config.md)

## Content

The image layout is as follows:

- `blobs` directory
  - Contains content-addressable blobs
  - A blob has no schema and SHOULD be considered opaque
  - Directory MUST exist and MAY be empty
  - See [blobs](#blobs) section
- `oci-layout` file
  - It MUST exist
  - It MUST be a JSON object
  - It MUST contain an `imageLayoutVersion` field
  - See [oci-layout file](#oci-layout-file) section
  - It MAY include additional fields
- `index.json` file
  - It MUST exist
  - It MUST be an [image index](image-index.md) JSON object.
  - See [index.json](#indexjson-file) section

**Implementor's Note:**
For extensibility and future expansion, additional files may be included in the directory.
Implementations should not error when encountering unknown files.
A common usage includes the `manifest.json` file associated with a backwards compatible `docker save` format.

## Example Layout

This is an example image layout:

```shell
$ cd example.com/app/
$ find . -type f
./index.json
./oci-layout
./blobs/sha256/3588d02542238316759cbf24502f4344ffcc8a60c803870022f335d1390c13b4
./blobs/sha256/4b0bc1c4050b03c95ef2a8e36e25feac42fd31283e8c30b3ee5df6b043155d3c
./blobs/sha256/7968321274dc6b6171697c33df7815310468e694ac5be0ec03ff053bb135e768
```

Blobs are named by their contents:

```shell
$ shasum -a 256 ./blobs/sha256/afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51
afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51 ./blobs/sha256/afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51
```

## Blobs

- Object names in the `blobs` subdirectories are composed of a directory for each hash algorithm, the children of which will contain the actual content.
- The content of `blobs/<alg>/<encoded>` MUST match the digest `<alg>:<encoded>` (referenced per [descriptor](descriptor.md#digests)). For example, the content of `blobs/sha256/da39a3ee5e6b4b0d3255bfef95601890afd80709` MUST match the digest `sha256:da39a3ee5e6b4b0d3255bfef95601890afd80709`.
- The character set of the entry name for `<alg>` and `<encoded>` MUST match the respective grammar elements described in [descriptor](descriptor.md#digests).
- The blobs directory MAY contain blobs which are not referenced by any of the [refs](#indexjson-file).
- The blobs directory MAY be missing referenced blobs, in which case the missing blobs SHOULD be fulfilled by an external blob store.

### Example Blobs

```shell
$ cat ./blobs/sha256/9b97579de92b1c195b85bb42a11011378ee549b02d7fe9c17bf2a6b35d5cb079 | jq
{
  "schemaVersion": 2,
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 7143,
      "digest": "sha256:afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51",
      "platform": {
        "architecture": "ppc64le",
        "os": "linux"
      }
    },
...
```

```shell
$ cat ./blobs/sha256/afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51 | jq
{
  "schemaVersion": 2,
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 7023,
    "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 32654,
      "digest": "sha256:9834876dcfb05cb167a5c24953eba58c4ac89b1adf57f28f2f9d09af107ee8f0"
    },
...
```

```shell
$ cat ./blobs/sha256/5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270 | jq
{
  "architecture": "amd64",
  "author": "Alyssa P. Hacker <alyspdev@example.com>",
  "config": {
    "Hostname": "8dfe43d80430",
    "Domainname": "",
    "User": "",
    "AttachStdin": false,
    "AttachStdout": false,
    "AttachStderr": false,
    "Tty": false,
    "OpenStdin": false,
    "StdinOnce": false,
    "Env": [
      "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    ],
    "Cmd": null,
    "Image": "sha256:6986ae504bbf843512d680cc959484452034965db15f75ee8bdd1b107f61500b",
...
```

```shell
$ cat ./blobs/sha256/9834876dcfb05cb167a5c24953eba58c4ac89b1adf57f28f2f9d09af107ee8f0
[gzipped tar stream]
```

## oci-layout file

This JSON object serves as a marker for the base of an Open Container Image Layout and to provide the version of the image-layout in use.
The `imageLayoutVersion` value will align with the OCI Image Specification version at the time changes to the layout are made, and will pin a given version until changes to the image layout are required.
This section defines the `application/vnd.oci.layout.header.v1+json` [media type](media-types.md).

### oci-layout Example

```json,title=OCI%20Layout&mediatype=application/vnd.oci.layout.header.v1%2Bjson
{
    "imageLayoutVersion": "1.0.0"
}
```

## index.json file

This REQUIRED file is the entry point for references and descriptors of the image-layout.
The [image index](image-index.md) is a multi-descriptor entry point.

This index provides an established path (`/index.json`) to have an entry point for an image-layout and to discover auxiliary descriptors.

- No semantic restriction is given for the "org.opencontainers.image.ref.name" annotation of descriptors.
- In general the `mediaType` of each [descriptor][descriptors] object in the `manifests` field will be either `application/vnd.oci.image.index.v1+json` or `application/vnd.oci.image.manifest.v1+json`.
- Future versions of the spec MAY use a different mediatype (i.e. a new versioned format).
- An encountered `mediaType` that is unknown MUST NOT generate an error.

**Implementor's Note:**
A common use case of descriptors with a "org.opencontainers.image.ref.name" annotation is representing a "tag" for a container image.
For example, an image may have a tag for different versions or builds of the software.
In the wild you often see "tags" like "v1.0.0-vendor.0", "2.0.0-debug", etc.
Those tags will often be represented in an image-layout repository with matching "org.opencontainers.image.ref.name" annotations like "v1.0.0-vendor.0", "2.0.0-debug", etc.

**Referrers Support:**
Referrers MAY be referenced using the fallback tag if the "org.opencontainers.image.referrer.convert" annotation is not set to "true".
Before writing descriptors with the "org.opencontainers.image.referrer.subject" annotation, implementations MUST ensure the "org.opencontainers.image.referrer.convert" annotation is set to "true" and convert any existing content referenced with the fallback tag if the annotation was not set.
If the "org.opencontainers.image.referrer.convert" annotation is set to "true", implementations MAY skip the conversion of referrers stored with the fallback tag and depend on the "org.opencontainers.image.referrer.subject" annotation to detect any referrers.

### Index Example

```json,title=Image%20Index&mediatype=application/vnd.oci.image.index.v1%2Bjson
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.index.v1+json",
      "size": 7143,
      "digest": "sha256:0228f90e926ba6b96e4f39cf294b2586d38fbb5a1e385c05cd1ee40ea54fe7fd",
      "annotations": {
        "org.opencontainers.image.ref.name": "stable-release"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 7143,
      "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
      "platform": {
        "architecture": "ppc64le",
        "os": "linux"
      },
      "annotations": {
        "org.opencontainers.image.ref.name": "v1.0"
      }
    },
    {
      "mediaType": "application/xml",
      "size": 7143,
      "digest": "sha256:b3d63d132d21c3ff4c35a061adf23cf43da8ae054247e32faa95494d904a007e",
      "annotations": {
        "org.freedesktop.specifications.metainfo.version": "1.0",
        "org.freedesktop.specifications.metainfo.type": "AppStream"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.index.v1+json",
      "size": 7143,
      "digest": "sha256:1efe7ab979c486a5af7a29d2c4603d84a3b934a7253d61b37e8573afecf47c03",
      "annotations": {
        "org.opencontainers.image.referrer.subject": "sha256:0228f90e926ba6b96e4f39cf294b2586d38fbb5a1e385c05cd1ee40ea54fe7fd"
      }
    }
  ],
  "annotations": {
    "org.opencontainers.image.referrer.convert": "true",
    "com.example.index.revision": "r124356"
  }
}
```

This illustrates an index that provides two named references and an auxiliary mediatype for this image layout.

The first named reference (`stable-release`) points to another index that might contain multiple references with distinct platforms and annotations.
Note that the [`org.opencontainers.image.ref.name` and `org.opencontainers.image.referrer.subject` annotations](annotations.md) SHOULD only be considered valid when on descriptors on `index.json`.
The [`org.opencontainers.image.referrer.convert` annotation](annotations.md) SHOULD only be considered valid when on manifest of the `index.json`.

The second named reference (`v1.0`) points to a manifest that is specific to the linux/ppc64le platform.

[descriptors]: ./descriptor.md
