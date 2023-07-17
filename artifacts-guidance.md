# Guidance for Artifacts Authors

Content other than OCI container images MAY be packaged using the image manifest.
When this is done, the `config.mediaType` value MUST be set to a value specific to the artifact type or the [empty value](manifest.md#guidance-for-an-empty-descriptor).
Additional details and examples are provided in the [image manifest specification](manifest.md#guidelines-for-artifact-usage).
