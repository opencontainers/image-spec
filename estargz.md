# eStargz: Standard-Compatible Extension to Container Image Layers for Lazy Pulling

This doc describes the extension to gzip layers of container images `application/vnd.oci.image.layer.v1.tar+gzip` for *lazy pulling*.
The extension is called *eStargz*.

eStargz is a *backward-compatible extension* which means that images can be pushed to the extension-agnostic registry and can run on extension-agnostic runtimes.

This extension is based on stargz (stands for *seekable tar.gz*) proposed by [Google CRFS](https://github.com/google/crfs) project (initially [discussed in Go community](https://github.com/golang/go/issues/30829)).
eStargz extends stargz for chunk-level verification, runtime performance optimization and RFC1952 compliance ([Extra field](https://tools.ietf.org/html/rfc1952#section-2.3.1.1)).

## Overview

Lazy pulling is a technique of pulling container images aiming at the faster cold start.
This allows a container to startup without waiting for the entire image layer contents to be locally available.
Instead, necessary files (or chunks for large files) in the layer are fetched *on-demand* during running the container.

For achieving this, runtimes need to fetch and extract each file in a layer independently.
However, layer without eStargz extension doesn't allow this because of the following reasons,

1. The entire layer blob needs to be extracted even for getting a single file entry.
2. Digests aren't provided for each file so it cannot be verified independently.

eStargz solves these issues and enables lazy pulling.

Additionally, it supports prefetching of files.
This can be used to mitigate runtime performance drawbacks caused by the on-demand fetching of each file.

This extension is a backward-compatible so the eStargz-formatted image can be pushed to the registry and can run even on eStargz-agnostic runtimes.

## The structure

![The structure of eStargz](/img/estargz-structure.png)

eStargz is a gzip-compressed tar archive of files and a metadata component called *TOC* (described in the later section).
In an eStargz-formatted blob, each non-empty regular file and each metadata component MUST be separately compressed as gzip.
This structure is inherited from [stargz](https://github.com/google/crfs).

Therefore, the gzip headers MUST locate at the following locations.

- The top of the blob
- The top of the payload of each non-empty regular file tar entry except *TOC*
- The top of *TOC* tar header
- The top of *footer* (described in the later section)

Large regular files in an eStargz blob MAY be chunked into several smaller gzip members.
Each chunked member is called *chunk* in this doc.

Therefore, gzip headers MAY locate at the following locations.

- Arbitrary location within the payload of non-empty regular file entry

An eStargz-formatted blob is the concatenation of these gzip members, which is a still valid gzip blob.

## TOC, TOCEntries and Footer

### TOC and TOCEntries

eStargz contains a regular file called *TOC* which records metadata (e.g. name, file type, owners, offset etc) of all file entries in eStargz, except TOC itself.
Container runtimes MAY use TOC to mount the container's filesystem without downloading the entire layer contents.

TOC MUST be a JSON file contained as the last tar entry and MUST be named `stargz.index.json`.

The following fields contain the primary properties that constitute a TOC.

- **`version`** *int*

   This REQUIRED property contains the version of the TOC. This value MUST be `1`.

- **`entries`** *array of objects*

   This property MUST contain an array of *TOCEntry* of all tar entries and chunks in the blob, except `stargz.index.json`.

*TOCEntry* consists of metadata of a file or chunk in eStargz.
If metadata in a TOCEntry of a file differs from the corresponding tar entry, TOCEntry SHOULD be respected.

The following fields contain the primary properties that constitute a TOCEntry.
Properties other than `chunkDigest` are inherited from [stargz](https://github.com/google/crfs).

- **`name`** *string*

  This REQUIRED property contains the name of the tar entry.
  This MUST be the complete path stored in the tar file.

- **`type`** *string*

  This REQUIRED property contains the type of tar entry.
  This MUST be either of the following.
  - `dir`: directory
  - `reg`: regular file
  - `symlink`: symbolic link
  - `hardlink`: hard link
  - `char`: character device
  - `block`: block device
  - `fifo`: fifo
  - `chunk`: a chunk of regular file data
  As described in the above section, a regular file can be divided into several chunks.
  TOCEntry MUST be created for each chunk.
  TOCEntry of the first chunk of that file MUST be typed as `reg`.
  TOCEntry of each chunk after 2nd MUST be typed as `chunk`.
  `chunk` TOCEntry MUST set *offset*, *chunkOffset* and *chunkSize* properties.

- **`size`** *uint64*

  This OPTIONAL property contains the uncompressed size of the regular file.
  Non-empty `reg` file MUST set this property.

- **`modtime`** *string*

  This OPTIONAL property contains the modification time of the tar entry.
  Empty means zero or unknown.
  Otherwise, the value is in UTC RFC3339 format.

- **`linkName`** *string*

  This OPTIONAL property contains the link target.
  `symlink` and `hardlink` MUST set this property.

- **`mode`** *int64*

  This REQUIRED property contains the permission and mode bits.

- **`uid`** *uint*

  This REQUIRED property contains the user ID of the owner of this file.

- **`gid`** *uint*

  This REQUIRED property contains the group ID of the owner of this file.

- **`userName`** *string*

  This OPTIONAL property contains the username of the owner.

- **`groupName`** *string*

  This OPTIONAL property contains the groupname of the owner.

- **`devMajor`** *int*

  This OPTIONAL property contains the major device number of device files.
  `char` and `block` files MUST set this property.

- **`devMinor`** *int*

  This OPTIONAL property contains the minor device number of device files.
  `char` and `block` files MUST set this property.

- **`xattrs`** *string-bytes map*

  This OPTIONAL property contains the extended attribute for the tar entry.

- **`digest`** *string*

  This OPTIONAL property contains the digest of the regular file contents.

- **`offset`** *int64*

  This OPTIONAL property contains the offset of the gzip header of the regular file or chunk in the blob.
  TOCEntries of non-empty `reg` and `chunk` MUST set this property.

- **`chunkOffset`** *int64*

  This OPTIONAL property contains the offset of this chunk in the decompressed regular file payload.
  TOCEntries of `chunk` type MUST set this property.

- **`chunkSize`** *int64*

  This OPTIONAL property contains the decompressed size of this chunk.
  The last `chunk` in a `reg` file or `reg` file that isn't chunked MUST set this property to zero.
  Other `reg` and `chunk` MUST set this property.

- **`chunkDigest`** *string*

  This OPTIONAL property contains a digest of this chunk.
  TOCEntries of non-empty `reg` and `chunk` MUST set this property.
  This MAY be used for verifying the data of the chunk.

### Footer

At the end of the blob, a *footer* MUST be appended.
This MUST be an empty gzip member whose [Extra field](https://tools.ietf.org/html/rfc1952#section-2.3.1.1) contains the offset of TOC in the blob.
The footer MUST be the following 51 bytes (1 byte = 8 bits in gzip).

```
- 10 bytes  gzip header
- 2  bytes  XLEN (length of Extra field) = 26 (4 bytes header + 16 hex digits + len("STARGZ"))
- 2  bytes  Extra: SI1 = 'S', SI2 = 'G'
- 2  bytes  Extra: LEN = 22 (16 hex digits + len("STARGZ"))
- 22 bytes  Extra: subfield = fmt.Sprintf("%016xSTARGZ", offsetOfTOC)
- 5  bytes  flate header: BFINAL = 1(last block), BTYPE = 0(non-compressed block), LEN = 0
- 8  bytes  gzip footer
(End of eStargz)
```

Runtimes MAY first read and parse the footer to get the offset of TOC.

Each file's metadata is recorded in the TOC so runtimes don't need to extract other parts of the archive as long as it only uses file metadata.
If runtime needs to get a regular file's content, it can get the size and offset of that content from the TOC and extract that range without scanning the entire blob.
By combining this with HTTP Range Request supported by [OCI Distribution Spec](https://github.com/opencontainers/distribution-spec/blob/ef28f81727c3b5e98ab941ae050098ea664c0960/detail.md), runtimes can selectively download file entries from the registry.

### Notes on compatibility with stargz

eStargz is designed aiming to compatibility with gzip layers.
For achieving this, eStargz's footer structure is incompatible with [stargz's one](https://github.com/google/crfs/blob/71d77da419c90be7b05d12e59945ac7a8c94a543/stargz/stargz.go#L36-L49).
eStargz adds SI1, SI2 and LEN fields to the footer to make it compliant to [Extra field definition in RFC1952](https://tools.ietf.org/html/rfc1952#section-2.3.1.1).
TOC, TOCEntry and the position of gzip headers are still compatible with stargz.

## Prioritized Files and Landmark Files

![Prioritized files and landmark files](/img/estargz-landmark.png)

Lazy pulling can cause runtime performance overhead by on-demand fetching of each file.
eStargz mitigates this by supporting prefetching of important files called *prioritized files*.

eStargz encodes the information about prioritized files to the *order* of file entries with some *landmark* file entries.

File entries in eStargz are grouped into the following groups,

- A. *prioritized files*
- B. non *prioritized files*

If no files are belonging to A, a landmark file *no-prefetch landmark* MUST be contained in the archive.

If one or more files are belonging to A, eStargz MUST consist of two separated areas corresponding to these groups and a landmark file *prefetch landmark* MUST be contained at the boundary between these two areas.

The Landmark file MUST be a regular file entry with 4 bits contents 0xf in eStargz.
It MUST be recorded to TOC as a TOCEntry. Prefetch landmark MUST be named `.prefetch.landmark`. No-prefetch landmark MUST be named `.no.prefetch.landmark`.

### Example use-case of prioritized files: workload-based image optimization in Stargz Snapshotter

Stargz Snapshotter makes use of eStargz's prioritized files for *workload-based* optimization to mitigate the overhead of reading files.
The *workload* of the image is the runtime configuration defined in the Dockerfile, including entrypoint command, environment variables and user.

Stargz snapshotter provides an image converter command `ctr-remote images optimize` to create optimized eStargz images.
When converting the image, this command runs the specified workload in a sandboxed environment and profiles all file accesses.
This command treats all accessed files as prioritized files.
Then it constructs eStargz by

- putting prioritized files from the top of the archive, sorting them by the accessed order,
- putting *prefetch landmark* file entry at the end of this range, and
- putting all other files (non-prioritized files) after the prefetch landmark.

Before running the container, stargz snapshotter prefetches and pre-caches the range where prioritized files are contained, by a single HTTP Range Request supported by the registry.
This can increase the cache hit rate for the specified workload and can mitigate runtime overheads.

## Content Verification in eStargz

The goal of the content verification in eStargz is to ensure the downloaded metadata and contents of all files are the expected ones, based on the calculated digests.
The verification of other components in the image including image manifests is out-of-scope of eStargz.
On the verification step of an eStargz layer, we assume that the manifest that references this eStargz layer is already verified (using digest tag, etc).

![the overview of the verification](/img/estargz-verification.png)

A non-eStargz layer can be verified by recalculating the digest and comparing it with the one written in the layer descriptor referencing that layer in the verified manifest.
However, an eStargz layer is *lazily* pulled from the registry in file (or chunk if that file is large) granularity so each one needs to be independently verified every time fetched.

The following describes how the verification of eStargz is done using the verified manifest.

eStargz consists of the following components to be verified:

- TOC (a set of metadata of all files contained in the layer)
- chunks of contents of each regular file

TOC contains metadata (name, type, mode, etc.) of all files and chunks in the blob.
On mounting eStargz, filesystem fetches the TOC from the registry.
For making the TOC verifiable using the verified manifest, we define an annotation `org.opencontainers.image.estargz.toc.digest`.
The value of this annotation is the digest of the TOC and this MUST be contained in the descriptor that references this eStargz layer.
Using this annotation, filesystem can verify the TOC by recalculating the digest and comparing it to the annotation value.

Each file's metadata is encoded to a TOCEntry in the TOC.
TOCEntry is created also for each chunk of regular files.
For making the contents of each file and chunk verifiable using the verified manifest, TOCEntry has a property *chunkDigest*.
*chunkDigest* contains the digest of the content of the `reg` or `chunk` entry.
As mentioned above, the TOC is verifiable using the special annotation.
Using *chunkDigest* fields written in the verified TOC, each file and chunk can be independently verified by recalculating the digest and comparing it to the property.

As the conclusion, eStargz MUST contain the following metadata:

- `org.opencontainers.image.estargz.toc.digest` annotation in the descriptor that references eStargz layer: The value is the digest of the TOC.
- *chunkDigest* properties of non-empty `reg` or `chunk` TOCEntry: The value is the digest of the contents of the file or chunk.

### Example usecase: Content verification in Stargz Snapshotter

Stargz Snapshotter verifies eStargz layers leveraging the above metadata.
As mentioned above, the verification of other image components including the manifests is out-of-scope of the snapshotter.
When this snapshotter mounts an eStargz layer, the manifest that references this layer must be verified in advance and the TOC digest annotation written in the verified manifest must be passed down to this snapshotter.

On mounting a layer, stargz snapshotter fetches the TOC from the registry.
Then it verifies the TOC by recalculating the digest and comparing it with the one written in the manifest.
After the TOC is verified, the snapshotter mounts this layer using the metadata recorded in the TOC.

During runtime of the container, this snapshotter fetches chunks of regular file contents lazily.
Before providing a chunk to the filesystem user, snapshotter recalculates the digest and checks it matches the one recorded in the corresponding TOCEntry.

## Example of TOC

Here is an example TOC JSON:

```json
{
  "version": 1,
  "entries": [
    {
      "name": "bin/",
      "type": "dir",
      "modtime": "2019-08-20T10:30:43Z",
      "mode": 16877,
      "NumLink": 0
    },
    {
      "name": "bin/busybox",
      "type": "reg",
      "size": 833104,
      "modtime": "2019-06-12T17:52:45Z",
      "mode": 33261,
      "offset": 126,
      "NumLink": 0,
      "digest": "sha256:8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521f",
      "chunkDigest": "sha256:8b7c559b8cccca0d30d01bc4b5dc944766208a53d18a03aa8afe97252207521f"
    },
    {
      "name": "lib/",
      "type": "dir",
      "modtime": "2019-08-20T10:30:43Z",
      "mode": 16877,
      "NumLink": 0
    },
    {
      "name": "lib/ld-musl-x86_64.so.1",
      "type": "reg",
      "size": 580144,
      "modtime": "2019-08-07T07:15:30Z",
      "mode": 33261,
      "offset": 512427,
      "NumLink": 0,
      "digest": "sha256:45c6ee3bd1862697eab8058ec0e462f5a760927331c709d7d233da8ffee40e9e",
      "chunkDigest": "sha256:45c6ee3bd1862697eab8058ec0e462f5a760927331c709d7d233da8ffee40e9e"
    },
    {
      "name": ".prefetch.landmark",
      "type": "reg",
      "size": 1,
      "offset": 886633,
      "NumLink": 0,
      "digest": "sha256:dc0e9c3658a1a3ed1ec94274d8b19925c93e1abb7ddba294923ad9bde30f8cb8",
      "chunkDigest": "sha256:dc0e9c3658a1a3ed1ec94274d8b19925c93e1abb7ddba294923ad9bde30f8cb8"
    },
... (omit) ...
```
