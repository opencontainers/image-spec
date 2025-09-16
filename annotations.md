# Annotations

Several components of the specification, like [Image Manifests](manifest.md) and [Descriptors](descriptor.md), feature an optional annotations property, whose format is common and defined in this section.

This property contains arbitrary metadata.

## Rules

- Annotations MUST be a key-value map where both the key and value MUST be strings.
- While the value MUST be present, it MAY be an empty string.
- Keys MUST be unique within this map, and best practice is to namespace the keys.
- Keys SHOULD be named using a reverse domain notation - e.g. `com.example.myKey`.
- The prefix `org.opencontainers` is reserved for keys defined in Open Container Initiative (OCI) specifications and MUST NOT be used by other specifications and extensions.
- Keys using the `org.opencontainers.image` namespace are reserved for use in the OCI Image Specification and MUST NOT be used by other specifications and extensions, including other OCI specifications.
- If there are no annotations then this property MUST either be absent or be an empty map.
- Consumers MUST NOT generate an error if they encounter an unknown annotation key.

## Pre-Defined Annotation Keys

This specification defines the following annotation keys, intended for but not limited to [image index](image-index.md), image [manifest](manifest.md), and [descriptor](descriptor.md) authors.

- **org.opencontainers.image.created** date and time on which the image was built, conforming to [RFC 3339][rfc3339].
- **org.opencontainers.image.authors** contact details of the people or organization responsible for the image (freeform string)
- **org.opencontainers.image.url** URL to find more information on the image (string)
- **org.opencontainers.image.documentation** URL to get documentation on the image (string)
- **org.opencontainers.image.source** URL to get source code for building the image (string)
- **org.opencontainers.image.version** version of the packaged software
  - The version MAY match a label or tag in the source code repository
  - version MAY be [Semantic versioning-compatible](https://semver.org/)
- **org.opencontainers.image.revision** Source control revision identifier for the packaged software.
- **org.opencontainers.image.vendor** Name of the distributing entity, organization or individual.
- **org.opencontainers.image.licenses** License(s) under which contained software is distributed as an [SPDX License Expression][spdx-license-expression].
- **org.opencontainers.image.ref.name** Name of the reference for a target (string).
  - SHOULD only be considered valid when on descriptors on `index.json` within [image layout](image-layout.md).
  - Character set of the value SHOULD conform to alphanum of `A-Za-z0-9` and separator set of `-._:@/+`
  - A valid reference matches the following [grammar](considerations.md#ebnf):

    ```ebnf
    ref       ::= component ("/" component)*
    component ::= alphanum (separator alphanum)*
    alphanum  ::= [A-Za-z0-9]+
    separator ::= [-._:@+] | "--"
    ```

- **org.opencontainers.image.title** Human-readable title of the image (string)
- **org.opencontainers.image.description** Human-readable description of the software packaged in the image (string)
- **org.opencontainers.image.base.digest** [Digest](descriptor.md#digests) of the image this image is based on (string)
  - This SHOULD be the immediate image sharing zero-indexed layers with the image, such as from a Dockerfile `FROM` statement.
  - This SHOULD NOT reference any other images used to generate the contents of the image (e.g., multi-stage Dockerfile builds).
- **org.opencontainers.image.base.name** Image reference of the image this image is based on (string)
  - This SHOULD be image references in the format defined by [distribution/distribution][distribution-reference].
  - This SHOULD be a fully qualified reference name, without any assumed default registry. (e.g., `registry.example.com/my-org/my-image:tag` instead of `my-org/my-image:tag`).
  - This SHOULD be the immediate image sharing zero-indexed layers with the image, such as from a Dockerfile `FROM` statement.
  - This SHOULD NOT reference any other images used to generate the contents of the image (e.g., multi-stage Dockerfile builds).
  - If the `image.base.name` annotation is specified, the `image.base.digest` annotation SHOULD be the digest of the manifest referenced by the `image.ref.name` annotation.
- **org.opencontainers.image.referrer.subject** Digest of the subject referenced by the referrers response (string)
  - This SHOULD only be considered valid when on descriptors on `index.json` within [image layout](image-layout.md).
  - The descriptor SHOULD be the referrers response for the subject digest.
- **org.opencontainers.image.referrer.convert** Defined and set to `true` when tooling has converted any referrers from the fallback tag to using the `org.opencontainers.image.referrer.subject` annotation.
  - This SHOULD only be considered valid when on the manifest of the `index.json` within [image layout](image-layout.md).
  - Tooling that reads an [image layout](image-layout.md) MAY skip the referrers conversion process when the annotation is detected.

## Back-compatibility with Label Schema

[Label Schema][label-schema] defined a number of conventional labels for container images, and these are now superseded by annotations with keys starting **org.opencontainers.image**.

While users are encouraged to use the **org.opencontainers.image** keys, tools MAY choose to support compatible annotations using the **org.label-schema** prefix as follows.

| `org.opencontainers.image` prefix | `org.label-schema` prefix | Compatibility notes |
|---------------------------|-------------------------|---------------------|
| `created` | `build-date` | Compatible |
| `url` | `url` | Compatible |
| `source` | `vcs-url` | Compatible |
| `version` | `version` | Compatible |
| `revision` | `vcs-ref` | Compatible |
| `vendor` | `vendor` | Compatible |
| `title` | `name` | Compatible |
| `description` | `description` | Compatible |
| `documentation` | `usage` | Value is compatible if the documentation is located by a URL |
| `authors` |  | No equivalent in Label Schema |
| `licenses` | | No equivalent in Label Schema |
| `ref.name` | | No equivalent in Label Schema |
| | `schema-version`| No equivalent in the OCI Image Spec |
| | `docker.*`, `rkt.*` | No equivalent in the OCI Image Spec |

[distribution-reference]: https://github.com/distribution/distribution/blob/d0deff9cd6c2b8c82c6f3d1c713af51df099d07b/reference/reference.go
[label-schema]: https://github.com/label-schema/label-schema.org/blob/gh-pages/rc1.md
[rfc3339]:     https://tools.ietf.org/html/rfc3339#section-5.6
[spdx-license-expression]: https://spdx.github.io/spdx-spec/v2.3/SPDX-license-expressions/
