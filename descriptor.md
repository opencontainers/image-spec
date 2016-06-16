# OpenContainers Content Descriptors

OCI have several components that come to together to describe an image.
References between components form a [Merkle Directed Acyclic Graph (DAG)](https://en.wikipedia.org/wiki/Merkle_tree).
The references in the _Merkle DAG_ are expressed through _Content Descriptors_.
A _Content Descriptor_ or _Descriptor_, describes the disposition of targeted content.
A _Descriptor_ includes the type of content, an independently-verifiable content identifier, known as a "digest" and the byte-size of the raw content.

Descriptors SHOULD be embedded in other formats to securely reference external content.

Other formats SHOULD use descriptors to securely reference external content.

## Properties

The following describe the primary set of properties that make up a _Descriptor_.

- **`mediaType`** *string*

  This REQUIRED property contains the MIME type of the referenced object.

- **`digest`** *string*

  This REQUIRED property is the _digest_ of the targeted content, meeting the requirements outlined in [Digests and Verification](#digests-and—verification).
  Retrieved content SHOULD be verified against this digest when consumed via untrusted sources.

- **`size`** *int64*
  This REQUIRED property specifies the size in bytes of the blob.
  This property exists so that a client will have an expected size for the content before validating.
  If the length of the retrieved content does not match the specified length, the content SHOULD NOT be trusted.

### Reserved

The following are field keys that MUST NOT be used in descriptors specified in other OCI specifications:

- **`urls`** *array*

  This key is RESERVED for future versions of the specification.

- **`data`** *string*

  This key is RESERVED for futures versions of the specification.

All other fields may be included in other OCI specifications.
Extended _Descriptor_ field additions proposed in other OCI specifications SHOULD first be considered for addition into this specification.

## Digests and Verification

The _digest_ component of a _Descriptor_ acts as a content identifier, employing [content addressability](http://en.wikipedia.org/wiki/Content-addressable_storage) for the OCI image format.
It uniquely identifies content by taking a collision-resistant hash of the bytes.
Such an identifier can be independently calculated and verified by selection of a common _algorithm_.
If such an identifier can be communicated in a secure manner, one can retrieve the content from an insecure source, calculate it independently and be certain that the correct content was obtained.
Put simply, the identifier is a property of the content.

To disambiguate from other concepts, we call this identifier a _digest_.
A _digest_ is a serialized hash result, consisting of a _algorithm_ and _hex_ portion.
The _algorithm_ identifies the methodology used to calculate the digest, which are shared by implementations.
The _hex_ portion is the hex-encoded result of the hash.

We define a _digest_ string to match the following grammar:

```
digest      := algorithm ":" hex
algorithm   := /[A-Fa-f0-9_+.-]+/
hex         := /[A-Fa-f0-9]+/
```

Some examples of _digests_ include the following:

digest                                                                            | description                                   |
----------------------------------------------------------------------------------|------------------------------------------------
sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b           | Common sha256 based digest                    |

Before consuming content targeted by a descriptor from untrusted sources, the byte content SHOULD be verified against the _digest_.
The size of the content SHOULD be verified to reduce hash collision space.
Heavy processing of before calculating a hash SHOULD be avoided.
Implementations MAY employ some canonicalization to ensure stable content identifiers.

### Algorithms

While the _algorithm_ does allow one to implement a wide variety of algorithms, compliant implementations SHOULD use [SHA-256](#SHA-256).

Let's use a simple example in pseudo-code to demonstrate a digest calculation:
A _digest_ is calculated by the following pseudo-code, where `H` is the selected hash algorithm, identified by string `<alg>`:
```
let ID(C) = Descriptor.digest
let C = <bytes>
let D = '<alg>:' + EncodeHex(H(C))
let verified = ID(C) == D
```
Above, we define the content identifier as `ID(C)`, extracted from the `Descriptor.digest` field.
Content `C` is a string of bytes.
Function `H` returns a the hashs of `C` in bytes and is passed to function `EncodeHex` to obtain the _digest_.
The result `verified` is true if `ID(C)` is equal to `D`, confirming that `C` is the content identified by `D`.
After verification, the following is true:

```
D == ID(C) == '<alg>:' + EncodeHex(H(C))
```

The _digest_ is confirmed as the content identifier by independently calculating the _digest_.

#### SHA-256

[SHA-256](https://tools.ietf.org/html/rfc4634#page-7) is a collision-resistant hash function, chosen for ubiquity, reasonable size and secure characteristics.
Implementations MUST implement SHA-256 digest verification for use in descriptors.

## Examples

The following example describes a [_Manifest_](manifest.md#image-manifest) with a content identifier of "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270", of size 7682 bytes:

```json,title=Content%20Descriptor&mediatype=application/vnd.oci.descriptor.v1%2Bjson
{
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "size": 7682,
  "digest": "sha256:5b0bcabd1ed22e9fb1310cf6c2dec7cdef19f0ad69efa1f392e94a4333501270"
}
```
