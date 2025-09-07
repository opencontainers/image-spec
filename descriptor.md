# OCI Content Descriptors

- An OCI image consists of several different components, arranged in a [Merkle Directed Acyclic Graph (DAG)](https://en.wikipedia.org/wiki/Merkle_tree).
- References between components in the graph are expressed through _Content Descriptors_.
- A Content Descriptor (or simply _Descriptor_) describes the disposition of the targeted content.
- A Content Descriptor includes the type of the content, a content identifier (_digest_), and the byte-size of the raw content.
  Optionally, it includes the type of artifact it is describing.
- Descriptors SHOULD be embedded in other formats to securely reference external content.
- Other formats SHOULD use descriptors to securely reference external content.
- When other formats contain multiple descriptors, unless otherwise specified, those descriptors are independent of each other, allowing fields like the `mediaType` and the algorithm for the `digest` to vary within that external content.

This section defines the `application/vnd.oci.descriptor.v1+json` [media type](media-types.md).

## Properties

A descriptor consists of a set of properties encapsulated in key-value fields.

The following fields contain the primary properties that constitute a Descriptor:

- **`mediaType`** *string*

  This REQUIRED property contains the media type of the referenced content.
  Values MUST comply with [RFC 6838][rfc6838], including the [naming requirements in its section 4.2][rfc6838-s4.2].

  The OCI image specification defines [several of its own MIME types](media-types.md) for resources defined in the specification.

- **`digest`** *string*

  This REQUIRED property is the _digest_ of the targeted content, conforming to the requirements outlined in [Digests](#digests).
  Retrieved content SHOULD be verified against this digest when consumed via untrusted sources.

- **`size`** *int64*

  This REQUIRED property specifies the size, in bytes, of the raw content.
  This property exists so that a client will have an expected size for the content before processing.
  If the length of the retrieved content does not match the specified length, the content SHOULD NOT be trusted.
  The size MUST NOT be negative.

- **`urls`** *array of strings*

  This OPTIONAL property specifies a list of URIs from which this object MAY be downloaded.
  Each entry MUST conform to [RFC 3986][rfc3986].
  Entries SHOULD use the `http` and `https` schemes, as defined in [RFC 7230][rfc7230-s2.7].

- **`annotations`** *string-string map*

  This OPTIONAL property contains arbitrary metadata for this descriptor.
  This OPTIONAL property MUST use the [annotation rules](annotations.md#rules).

- **`data`** *string*

  This OPTIONAL property contains an embedded representation of the referenced content.
  Values MUST conform to the Base 64 encoding, as defined in [RFC 4648][rfc4648-s4].
  The decoded data MUST be identical to the referenced content and SHOULD be verified against the [`digest`](#digests) and `size` fields by content consumers.
  See [Embedded Content](#embedded-content) for when this is appropriate.

- **`artifactType`** *string*

  This OPTIONAL property contains the type of an artifact when the descriptor points to an [Artifact](artifacts-guidance.md).
  This property MUST be the same as the `artifactType` of the referenced [manifest](manifest.md), or the `mediaType` of the `config` descriptor if an `artifactType` is not set.
  If defined, the value MUST comply with [RFC 6838][rfc6838], including the [naming requirements in its section 4.2][rfc6838-s4.2], and MAY be registered with [IANA][iana].

Descriptors referencing an Image (as defined in the [artifacts guidance](artifacts-guidance.md)) SHOULD include the extended field `platform`.
See [Image Index Property Descriptions](image-index.md#image-index-property-descriptions) for details.
Artifacts as defined in the artifacts guidance MAY include the extended field `platform` when it improves discoverability or interoperability of the artifact. 

### Reserved

Extended _Descriptor_ field additions proposed in other OCI specifications SHOULD first be considered for addition into this specification.

## Digests

The _digest_ property of a Descriptor acts as a content identifier, enabling [content addressability](https://en.wikipedia.org/wiki/Content-addressable_storage).
It uniquely identifies content by taking a [collision-resistant hash](https://en.wikipedia.org/wiki/Cryptographic_hash_function) of the bytes.
If the _digest_ can be communicated in a secure manner, one can verify content from an insecure source by recalculating the digest independently, ensuring the content has not been modified.

The value of the `digest` property is a string consisting of an _algorithm_ portion and an _encoded_ portion.
The _algorithm_ specifies the cryptographic hash function and encoding used for the digest; the _encoded_ portion contains the encoded result of the hash function.

A digest string MUST match the following [grammar](considerations.md#ebnf):

```ebnf
digest                ::= algorithm ":" encoded
algorithm             ::= algorithm-component (algorithm-separator algorithm-component)*
algorithm-component   ::= [a-z0-9]+
algorithm-separator   ::= [+._-]
encoded               ::= [a-zA-Z0-9=_-]+
```

Note that _algorithm_ MAY impose algorithm-specific restriction on the grammar of the _encoded_ portion.
See also [Registered Algorithms](#registered-algorithms).

Some example digest strings include the following:

| digest                                                                    | algorithm                   | Registered |
|---------------------------------------------------------------------------|-----------------------------|------------|
| `sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b` | [SHA-256](#sha-256)         | Yes        |
| `sha512:401b09eab3c013d4ca54922bb802bec8fd5318192b0a75f201d8b372742...`   | [SHA-512](#sha-512)         | Yes        |
| `multihash+base58:QmRZxt2b1FVZPNqd8hsiykDL3TdBDeTSPX9Kv46HmX4Gx8`         | Multihash                   | No         |
| `sha256+b64u:LCa0a2j_xo_5m0U8HTBBNBNCLXBkg7-g-YpeiGJm564`                 | SHA-256 with urlsafe base64 | No         |
| `blake3:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b` | [BLAKE3](#blake3)           | Yes        |

Please see [Registered Algorithms](#registered-algorithms) for a list of registered algorithms.

Implementations SHOULD allow digests with unrecognized algorithms to pass validation if they comply with the above grammar.
While `sha256` will only use hex encoded digests, separators in _algorithm_ and alphanumerics in _encoded_ are included to allow for extensions.
As an example, we can parameterize the encoding and algorithm as `multihash+base58:QmRZxt2b1FVZPNqd8hsiykDL3TdBDeTSPX9Kv46HmX4Gx8`, which would be considered valid but unregistered by this specification.

### Verification

Before consuming content targeted by a descriptor from untrusted sources, the byte content SHOULD be verified against the digest string.
Before calculating the digest, the size of the content SHOULD be verified to reduce hash collision space.
Heavy processing before calculating a hash SHOULD be avoided.
Implementations MAY employ [canonicalization](considerations.md#canonicalization) of the underlying content to ensure stable content identifiers.

### Digest calculations

A _digest_ is calculated by the following pseudo-code, where `H` is the selected hash algorithm, identified by string `<alg>`:

```text
let ID(C) = Descriptor.digest
let C = <bytes>
let D = '<alg>:' + Encode(H(C))
let verified = ID(C) == D
```

Above, we define the content identifier as `ID(C)`, extracted from the `Descriptor.digest` field.
Content `C` is a string of bytes.
Function `H` returns the hash of `C` in bytes and is passed to function `Encode` and prefixed with the algorithm to obtain the digest.
The result `verified` is true if `ID(C)` is equal to `D`, confirming that `C` is the content identified by `D`.
After verification, the following is true:

```text
D == ID(C) == '<alg>:' + Encode(H(C))
```

The _digest_ is confirmed as the content identifier by independently calculating the _digest_.

### Registered algorithms

While the _algorithm_ component of the digest string allows the use of a variety of cryptographic algorithms, compliant implementations SHOULD use [SHA-256](#sha-256).

The following algorithm identifiers are currently defined by this specification:

| algorithm identifier | algorithm           |
|----------------------|---------------------|
| `sha256`             | [SHA-256](#sha-256) |
| `sha512`             | [SHA-512](#sha-512) |
| `blake3`             | [BLAKE3](#blake3)   |

If a useful algorithm is not included in the above table, it SHOULD be submitted to this specification for registration.

#### SHA-256

[SHA-256][rfc4634-s4.1] is a collision-resistant hash function, chosen for ubiquity, reasonable size and secure characteristics.
Implementations MUST implement SHA-256 digest verification for use in descriptors.

When the _algorithm identifier_ is `sha256`, the _encoded_ portion MUST match `/[a-f0-9]{64}/`.
Note that `[A-F]` MUST NOT be used here.

#### SHA-512

[SHA-512][rfc4634-s4.2] is a collision-resistant hash function which [may be more performant][sha256-vs-sha512] than [SHA-256](#sha-256) on some CPUs.
Implementations MAY implement SHA-512 digest verification for use in descriptors.

When the _algorithm identifier_ is `sha512`, the _encoded_ portion MUST match `/[a-f0-9]{128}/`.
Note that `[A-F]` MUST NOT be used here.

#### BLAKE3

[BLAKE3][blake3] is a high performance, highly parallelizable, collision-resistant hash function which [is more performant][blake3-vs-sha2] than
[SHA-256][rfc4634-s4.1].
The hash output length MUST be 256 bits.
Implementations MAY implement BLAKE3 digest verification for use in descriptors.

When the _algorithm identifier_ is `blake3`, the _encoded_ portion MUST match `/[a-f0-9]{64}/`.
Note that `[A-F]` MUST NOT be used here.

## Embedded Content

In many contexts, such as when downloading content over a network, resolving a descriptor to its content has a measurable fixed "roundtrip" latency cost.
For large blobs, the fixed cost is usually inconsequential, as the majority of time will be spent actually fetching the content.
For very small blobs, the fixed cost can be quite significant.

Implementations MAY choose to embed small pieces of content directly within a descriptor to avoid roundtrips.

Implementations MUST NOT populate the `data` field in situations where doing so would modify existing content identifiers.
For example, a registry MUST NOT arbitrarily populate `data` fields within uploaded manifests, as that would modify the content identifier of those manifests.
In contrast, a client MAY populate the `data` field before uploading a manifest, because the manifest would not yet have a content identifier in the registry.

Implementations SHOULD consider portability when deciding whether to embed data, as some providers are known to refuse to accept or parse manifests that exceed a certain size.

## Examples

The following example describes a [_Manifest_](manifest.md#image-manifest) with a content identifier of "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270" and a size of 7682 bytes:

```json,title=Content%20Descriptor&mediatype=application/vnd.oci.descriptor.v1%2Bjson
{
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "size": 7682,
  "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270"
}
```

In the following example, the descriptor indicates that the referenced manifest is retrievable from a particular URL:

```json,title=Content%20Descriptor&mediatype=application/vnd.oci.descriptor.v1%2Bjson
{
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "size": 7682,
  "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270",
  "urls": [
    "https://example.com/example-manifest"
  ]
}
```

In the following example, the descriptor indicates the type of artifact it is referencing:

```json,title=Content%20Descriptor&mediatype=application/vnd.oci.descriptor.v1%2Bjson
{
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "size": 123,
  "digest": "sha256:87923725d74f4bfb94c9e86d64170f7521aad8221a5de834851470ca142da630",
  "artifactType": "application/vnd.example.sbom.v1"
}
```

[rfc3986]: https://tools.ietf.org/html/rfc3986
[rfc4634-s4.1]: https://tools.ietf.org/html/rfc4634#section-4.1
[rfc4634-s4.2]: https://tools.ietf.org/html/rfc4634#section-4.2
[rfc4648-s4]: https://tools.ietf.org/html/rfc4648#section-4
[rfc6838]: https://tools.ietf.org/html/rfc6838
[rfc6838-s4.2]: https://tools.ietf.org/html/rfc6838#section-4.2
[rfc7230-s2.7]: https://tools.ietf.org/html/rfc7230#section-2.7
[sha256-vs-sha512]: https://groups.google.com/a/opencontainers.org/forum/#!topic/dev/hsMw7cAwrZE
[iana]: https://www.iana.org/assignments/media-types/media-types.xhtml
[blake3]: https://github.com/C2SP/C2SP/blob/BLAKE3/v1.0.0/BLAKE3.md
[blake3-vs-sha2]: https://github.com/BLAKE3-team/BLAKE3-specs/blob/master/blake3.pdf
