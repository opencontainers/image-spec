# Guidance for Artifacts Authors

Content other than OCI container images MAY be packaged using the image manifest.
When this is done, the `config.mediaType` value MAY not be a known OCI image config [media type](media-types.md).
Historically, due to registry limitations, some tools have created non-OCI conformant artifacts using the `application/vnd.oci.image.config.v1+json` value for `config.mediaType` and values specific to the artifact in `layer[*].mediaType`.  In some cases, the `application/vnd.oci.image.config.v1+json` may still be appropriate if the artifact is to be run via a runtime.
Implementation details and examples are provided in the [image manifest specification](manifest.md#guidelines-for-artifact-usage).
