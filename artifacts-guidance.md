# Guidance for Artifacts Authors

## Artifacts and Images

This specification is primarily concerned with packaging two kinds of content: Artifacts and Images.
Both are representing using a [manifest](manifest.md).
Images are defined in this specification as conformant content with a conformant [config](config.md), processed according to [conversion.md](conversion.md) to derive a [runtime-spec][] configuration blob.
Conversely, an Artifact is any other conformant content that **does not contain a config to be interpreted by a runtime-spec implementation using the conversion mechanism.

## Creating an Artifact

Content other than Images MAY be packaged using the [manifest]; this is otherwise known as an Artifact.
When this is done, the `artifactType` should be set to a custom media type, or the `config.mediaType` should not be a known Image config [media type](media-types.md).
Implementation details and examples are provided in the [image manifest specification](manifest.md#guidelines-for-artifact-usage).

Note: Historically, due to registry limitations, some tools have created non-conformant Artifacts using the `application/vnd.oci.image.config.v1+json` value for `config.mediaType`.

## Interacting with Artifacts

Software following the process described in [conversion.md](conversion.md) to create a [runtime-spec][] configuration blob SHOULD ignore unknown Artifacts (as determined by the presence of a descriptor `artifactType`) when selecting content from an [index](image-index.md).
It is possible that implementations may also be able to interpret known Artifact types; however that is outside the scope of this spec.

Artifacts can be detected at runtime using by checking two keys:
1. Is an `artifactType` present in the descriptor, or in the [manifest](manifest.md)?
2. Is the `config.mediaType` of the manifest something other than a [known media type](media-types.md) for [config](config.md)?

If either of these tests is true, then the content is an Artifact.

[runtime-spec]: https://github.com/opencontainers/runtime-spec/blob/main/spec.md

