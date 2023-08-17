# Guidance for Artifacts Authors

Content other than OCI container images MAY be packaged using the image manifest.
When this is done, the `config.mediaType` value should not be a known OCI image config [media type](media-types.md).
Historically, due to registry limitations, some tools have created non-OCI conformant artifacts using the `application/vnd.oci.image.config.v1+json` value for `config.mediaType` and values specific to the artifact in `layer[*].mediaType`.
Implementation details and examples are provided in the [image manifest specification](manifest.md#guidelines-for-artifact-usage).
