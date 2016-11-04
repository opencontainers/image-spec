# Name assertions

The bulk of OCI image content is transmitted via a [Merkle DAG][merkle] with [descriptor](descriptor.md) references.
The location-addressable [refs](image-layout.md) provide convenient access into that DAG, but they introduce a degree of uncertainty.
Publishers MAY add signed name-assertions to the DAG to provide consumers with a way of mitigating ref uncertainty.
The meaning of the assertion is that the name [`name`](#payload) applies to the blob referenced by [`blob`](#payload).

This section defines the `application/vnd.oci.name.assertion.v1` [media type](media-types.md).

## Media type header and payload

All assertions MUST begin with the [US-ASCII][rfc5234-b.2] representation of their [media type][rfc6838] followed a [CRLF][rfc5234-b.1].
The content after the header is the *payload*.

## Payload

The [payload](#media-type-header-and-payload) for `application/vnd.oci.name.assertion.v1` MUST consist of [UTF-8][] [JSON][] object with the following properties:

* **`name`** (string, REQUIRED) the name being asserted.
* **`blob`** ([descriptor](descriptor.md), REQUIRED) the blob being named.

There is a JSON Schema for this payload at [`schema/name-assertion-schema.json`](schema/name-assertion-schema.json) and a Go type at [`specs-go/v1/assertions.go`](specs-go/v1/assertions.go).

## Example

The name assertion:

```title=Name%20assertion&mediatype=application/vnd.oci.name.assertion.v1
application/vnd.oci.name.assertion.v1
{
  "name": "foo-bar v1.0",
  "blob": {
    "size": 4096,
    "digest": "sha256:e692418e4cbaf90ca69d05a66403747baa33ee08806650b51fab815ad7fc331f",
    "mediatype": "application/vnd.oci.image.manifest.list.v1+json"
  }
}
```

asserts that [`sha256:e692…`](descriptor.md#digests-and-verification) is a 4096-byte [manifest list](manifest-list.md) named `foo-bar v1.0`.

## Use

You can use name assertions to name arbitrary blobs, but there are pros and cons to the various OCI [media types](media-types.md).

### Manifest lists

Naming a [manifest lists](manifest-list.md) locks in the full OCI-defined DAG.
The manifest list's platform information ([`architecture`, `os`, `os.version`, `os.features`, etc.][manifest-list-properties]) route consumers to their appropriate manifest and are more detailed than platform information further down the DAG (the only [configuration](config.md) properties are [`architecture` and `os`](config.md#properties)).
A malicious intermediate who alters a manifest list can disrupt the manifest-lookup process.

Even with signed name assertions on the referenced [manifest](manifest.md) or [configuration](config.md), the unsigned manifest list allows the intermediate to deny service (e.g. by removing all the manifest entries) or by misdirecting lookups (e.g. by altering platform fields on existing entries).
When the intermediate is only altering the manifest-list-specific platform information (e.g. [`features`][manifest-list-properties]), the change may not be detected during unpacking.

Directing consumers at specially-crafted blobs may also exploit vulnerabilities in the consumer's blob-handling logic, much like [decompression bombs][decompression-bomb] and [tarbombs][]).
Consumers who do not expect gzipped or tar-based media types in `manifests` can limit their exposure by restricting the set of supported media types for referenced blobs.
Consumers who restrict their support for referenced-blob media types to [UTF-8][] [JSON][] are still vulnerable to attacks on their [UTF-8][]- and [JSON][]-parsing implementation, although they were already exposed to those attacks when they parsed the manifest list itself.

[`annotations`][manifest-list-properties] is weakly structured, so the impact of malicious alterations is hard to assess.

Signed name assertions on the manifest list allow consumers to avoid these disruptions.

However, signing name assertions on manifest lists may require coordination among several parties (e.g. the developer producing the amd64 Linux version and the build farm producing the arm Linux version), and the original signature will not apply to new manifests resulting from any change to the downstream DAG.
As an assertion on the highest-level OCI type, manifest list assertions are the most affected by these issues.

### Manifests

Naming a [manifest](manifest.md) locks in a single platform's DAG.

A malicious intermediate who can perform the same lookup misdirection and reference bombing with [`config` and `layers`][manifest-properties] that the [manifest list attacker](#manifest-lists) could perform on `manifests`.
`layers` alterations has a higher exposure to [decompression bombs][decompression-bomb] and [tarbombs][] because consumers will expect layers to be [gzipped tar archives](layer.md#distribution-format).
With a signed name assertion on the referenced [configuration](configuration.md), the consumer can use the configuration's [`rootfs.diff_ids`](config.md#properties) to detect tar alteration (although that will not protect you from decompression bombs).

[`annotations`][manifest-properties] is weakly structured, so the impact of malicious alterations is hard to assess.

Signed name assertions on the manifest allow consumers to avoid these disruptions.

However, signing manifests may require coordination among several parties (e.g. the original publisher and the CAS admin who is replacing equivalent layers with a [consistently-compressed version](canonicalization.md)).
As an assertion on a mid-level OCI type, manifest assertions are moderately affected by these issues.

### Configurations

Naming a [configuration](configuration.md) locks in a single platform's unpacked filesystem.

With a signed name assertion on the [configuration](configuration.md), the consumer can be sure that a successfully unpacked filesystem contains an image of the asserted name.
And except for the [manifest-list-specific platform fields](#manifest-lists), the consumer can be sure that the unpacked filesystem matches the intended platform.
Without a signed name assertion on the configuration or the DAG leading to it, a malicious intermediate may have complete control over the unpacked filesystem.

Signing name assertions on the configuration rarely requires coordination among several parties.
As an assertion on a low-level OCI type, configuration assertions are minimally affected by these issues.

### Layers

Naming a [layer](layer.md) locks in the blob, but there are no references to other blobs (not even the configuration's quasi-reference [`rootfs.diff_ids`](config.md#properties)).
Signing assertions on layers directly may be useful for image assemblers in frameworks that use layers as a package-management system, with package managers publishing layers with signed name assertions and integrators assembling those layers into images.
However, the asserted name will rarely match the image-consumer's requested name (e.g. requesting an `nginx-1.10.1` image might pull in layers for `glibc-2.22` and `zlib-1.2.8`).

Signing name assertions on layers may require coordination among several parties (e.g. the original publisher and the CAS admin who is replacing equivalent layers with a [consistently-compressed version](canonicalization.md)).
As an assertion on a low-level OCI type whose stability may be relaxed by `rootfs.diff_ids`, layer assertions are moderately affected by these issues.

[JSON]: https://tools.ietf.org/html/rfc7159
[decompression-bombs]: https://en.wikipedia.org/wiki/Zip_bomb
[tarbombs]: https://en.wikipedia.org/wiki/Tar_%28computing%29#Tarbomb
[manifest-list-properties]: manifest-list.md#manifest-list-property-descriptions
[manifest-properties]: manifest.md#image-manifest-property-descriptions
[merkle]: https://en.wikipedia.org/wiki/Merkle_tree
[rfc5234-b.1]: https://tools.ietf.org/html/rfc5234#appendix-B.1
[rfc5234-b.2]: https://tools.ietf.org/html/rfc5234#appendix-B.2
[rfc6838]: https://tools.ietf.org/html/rfc6838
[UTF-8]: http://www.unicode.org/versions/Unicode8.0.0/ch03.pdf
