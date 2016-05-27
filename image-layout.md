## Open Container Initiative Image Layout Specification

The OCI Image Layout is a slash separated layout of OCI content-addressable blobs and [location-addressable](https://en.wikipedia.org/wiki/Content-addressable_storage#Content-addressed_vs._location-addressed) references (refs).
This layout MAY be used in a variety of different transport mechanisms: archive formats (e.g. tar, zip), shared filesystem environments (e.g. nfs), or networked file fetching (e.g. http, ftp, rsync).
Given an image layout a tool can convert a given ref into a runnable OCI Image Format by finding an appopriate manifest from the manifest list, unpacking the filesystem serializations in the correct order, and then converting the image configuration into an OCI Runtime config.json.

The image layout has two top level directories:

- "blobs" contains content-addressable blobs. A blob has no schema and should be considered opaque.
- "refs" contains descriptors pointing to an image manifest list

It also contains a file that is used to identify the layout version:

- "oci-layout" MUST contain a JSON object with a version field `{"imageLayoutVersion": "1.0.0"}` and MAY include additional fields.

This is an example image layout:

```
$ cd example.com/app/
$ find .
.
./blobs
./blobs/sha256-afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51
./blobs/sha256-5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270
./blobs/sha256-e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f
./oci-layout
./refs
./refs/v1.0
./refs/v1.1
```

Blobs are named by their contents:

```
$ shasum -a 256 ./blobs/sha256-afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51
afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51 ./blobs/sha256-afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51
```

Object names in the refs and blobs MUST NOT include characters outside of the set of "A" to "Z", "a" to "z", the hyphen `-`, the dot `.`, and the underscore `_`.
Hash algorithm identifiers containing the colon `:` will be converted to the hyphen `-`.
For example `sha256:5b` will map to the layout `blobs/sha256-5b`.
The blobs directory MAY contain blobs which are not referenced by any of the refs.
The blobs directory MAY be missing referenced blobs, in which case the missing blobs SHOULD be fulfilled by an external blob store.

Each object in the refs subdirectory MUST be of type `application/vnd.oci.descriptor.v1+json`.
In general the `mediatype` of this descriptor object will be either `application/vnd.oci.image.manifest.list.v1+json` or `application/vnd.oci.image.manifest.v1+json` although future versions of the spec MAY use a different mediatype.

This illustrates the expected contents of a given ref and the manifest list it points to.

```
$ cat ./refs/v1.0
{"size": 4096, "digest": "sha256:afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51", "mediatype": "application/vnd.oci.image.manifest.list.v1+json"}
```
```
$ cat ./blobs/sha256-afff3924849e458c5ef237db5f89539274d5e609db5db935ed3959c90f1f2d51
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.list.v1+json",
  "manifests": [
    {
...
```
