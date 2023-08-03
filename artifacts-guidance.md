# Guidance for Artifacts Authors

Content other than OCI container images MAY be packaged using the image manifest.
When this is done, 

 - `image.artifactType` value MUST be set to a value specific to the artifact type.
 - `config.mediaType` value must be set to [empty value](manifest.md#guidance-for-an-empty-descriptor).

Additional details and examples are provided in the [image manifest specification](manifest.md#guidelines-for-artifact-usage).
