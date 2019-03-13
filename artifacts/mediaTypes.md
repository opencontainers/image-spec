# Artifact Media Types

Artifact tools require the ability to understand and differentiate their specific artifacts from other artifacts stored in a registry. 

Registry products and projects will need to understand the different artifact types they present to their users, providing additional information about the artifact.

To uniquely identify artifact types within a registry, the [`manifest.config` property descriptor](https://github.com/opencontainers/image-spec/blob/master/manifest.md#image-manifest-property-descriptions) supports additional `mediaTypes`. 

## Config Media Types

By supporting new `manifest.config.mediaTypes`, registry tools MAY understand the types stored. 

A simplified example of a registry listing may look like:

| artifact reference | icon | type |
|-|-|-|
| `samples/image/hello-world:1.0` |![](../img/oci-container.png)| container image | 
| `samples/image/hello-world-doc:1.0`|![](../img/doc.png)| doc |
| `samples/helm/hello-world:1.0` |![](../img/helm.png)| helm chart | 

## Authoring New Artifact Types

Artifact authors have the freedom to define their artifacts as meets their needs. To conform to the OCI Image spec, and be stored within an OCI distribution based registry, Artifact authors MUST follow the following requirements:

- Artifact authors MUST provide unique `config.mediaTypes` to represent their artifact. 
- Artifact authors MAY submit a PR for their `config.mediaType` to [./mediaTypeMappings.json]() to make their new artifact type well known.
- Artifact authors MAY provide new schemas for the config.json object. 
- Artifacts and tooling that use existing `config.mediaTypes` MUST conform to the artifacts spec.
- Tools MAY choose to support new `config.mediaTypes`.
- Tools MUST be able to ignore `config.mediaTypes` they don't support. 
- Distribution implementations MAY choose to support new `config.mediaTypes`. See the [distribution-spec] for more details on processing of the `config.mediaTypes`

### mediaTypes formatting

New mediaTypes MUST use the following convention:

`application/vnd.`[governingOrg/vendor]`.`[artifactName]`.config.`[version]`+json` 

**Example:** 
The Helm 3 artifact may be represented as: `application/vnd.cncf.helm/config.v3+json`

## Config Parsing

New artifacts MUST choose from 3 options for how the `config` object is processed and parsed:

- Support the [`application/vnd.oci.image.config.v1+json`](../config.md) schema
- Define a new schema, which artifact specific tools would need to support
- Define a `null` schema, where the artifact only uses the `config.mediaType` to identify it's unique artifact type. 

Artifact authors MAY publicize their schema in [./mediaTypeMappings.json]() through a pull request.

Artifact authors understand supporting additional artifact types is optional, and distribution parsing of new config schemas is optional. However, new config schemas may be useful to that artifact tooling experience. 

## Layer mediaTypes

While OCI container images have ordinal layers, supporting overlaying of a unified file system, new artifact authors MAY define how they persist their layers, which may be individual files, or collections of files. The tooling specific to that artifact type owns the processing of the files and how they are extracted on the client. 

The [layer mediaType descriptor in the OCI Image spec](https://github.com/opencontainers/image-spec/blob/master/manifest.md#image-manifest-property-descriptions) identifies the `mediaType` MUST support one of the existing layer formats. 

- **`layers`** *array of objects*

    Each item in the array MUST be a [descriptor](../descriptor.md).
    
    Beyond the [descriptor requirements](../descriptor.md#properties), the value has the following additional restrictions:

    - **`mediaType`** *string*

        This [descriptor property](descriptor.md#properties) has additional restrictions for `layers[]`.
        Implementations MUST support at least the following media types:

        - [`application/vnd.oci.image.layer.v1.tar`](layer.md)
        - [`application/vnd.oci.image.layer.v1.tar+gzip`](layer.md#gzip-media-types)
    
Artifact authors can extend these mediaTypes to identify their unique usage, using the following format: 

`application/vnd.`[governingOrg/vendor]`.`[artifactName]`.layer.`[version][`.tar`|`+tar+gzip`]

**Helm Example**

`application/vnd.cncf.helm.chart.layer.v3.tar`
`application/vnd.cncf.helm.chart.values.v3.tar`

## Example Manifests

For illustrative purposes, OCI Image, and mockups for Helm and Doc are expressed below. 

**Note:** While MUST will be included in the Helm and DOC examples, these are examples of what these specs *MAY* look like. As Helm and Doc formalize, and this PR gets reviewed the specific references will be udpated to reflect the Helm and Doc specs. 

Artifact authors have the freedom to define their artifact representation using `manifest.config.mediaTypes` and `layer.mediaTypes`

### OCI image manifest example

The following is an example of the image manifest, used for comparison to other artifact types.
```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 16078,
    "digest": "sha256:bd698aa18aa02a2f083292b9448130a3afaa9a3e5fba24fc0aef7845c336b8ad"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 23133155,
      "digest": "sha256:9a1a13172ed974323f7c35153e8b23b8fa1c85355b6b26cc3127e640e45ef0aa"
    }
  ]
```

### Doc manifest example
The following is an example for a new _hypothetical_ doc artifact type, used to provide local/offline readme contents. 

The manifest uses `config.mediaType`=`application/vnd.tbd.doc.config.v1+json` to define the doc artifact type. 

`layer.mediaType`=`application/vnd.tbd.doc.markdown.layer.v1.tar` defines a markdown content layer.

```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.tbd.doc.config.v1+json",
    "size": 4142,
    "digest": "sha256:7d6da8aa18aa02a2f083292b9448130a3afaa9a3e5fba24fc0aef7845c336c6bd"
  },
  "layers": [
    {
      "mediaType": "application/vnd.tbd.doc.markdown.layer.v1.tar",
      "size": 29123,
      "digest": "sha256:8a7d13172ed974323f7c35153e8b23b8fa1c85355b6b26cc3127e640712927dd"
    }
  ]
```


## OCI image index

An OCI index is a higher-level manifest which contains a collection of image manifests, typically used for one or more platforms. In this example, the image index is used to group an image an it's associated docs. 

> **Note:**  The artifact spec does not account for a new typed index. The proposal focuses on manifest being the unique artifact type.

image:tag = `hello-world:1.0` 
```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.index.v1+json",
  "manifests": [
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 4288,
      "digest": "sha256:bd698aa18aa02a2f083292b9448130a3afaa9a3e5fba24fc0aef7845c336b8ad",
      "platform": {
        "architecture": "amd64",
        "os": "linux"
      }
    },
    {
      "mediaType": "application/vnd.oci.image.manifest.v1+json",
      "size": 4288,
      "digest": "sha256:7d6da8aa18aa02a2f083292b9448130a3afaa9a3e5fba24fc0aef7845c336c6bd"
    }
  ]
```


## Helm Chart
> **Note:** *this is just an illustrative example of mediaTypes.*
> The persistance format has yet to be decided by the Helm community, the following represents what a Helm chart *could* be expressed as.

Unlike OCI Image manifests which represents a unified file system, Helm charts use layers to represent different elements of a chart. By separating the chart files from values, multiple deployments of the same chart, only differentiated by the values, may be tracked and evaluated by new Helm tools.

| mediaType | usage |
|-|-|
|`application/vnd.cncf.helm.config.v3`|Helm Config - defining the Helm artifact type|
|`application/vnd.cncf.helm.chart.layer.v3`|Helm chart layer |
|`application/vnd.cncf.helm.values.layer.v3`|Helm values layer |

**Helm manifest example**

`helm chart pull example.com/samples/helm/wordpress:0.1.0`

```json
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.cncf.helm.config.v3+json",
    "size": 7023,
    "digest": "sha256:b5b2b2c507a0944348e0303114d8d93aaaa081732b86451d9bce1f432a537bc7",
    "annotations": {
      "io.cncf.helm.appVersion": "0.1.0"
    }
  },
  "layers": [
    {
      "mediaType": "application/vnd.cncf.helm.chart.layer.v3.tar+gzip",
      "size": 3155,
      "digest": "sha256:9a1a13172ed974323f7c35153e8b23b8fa1c85355b6b26cc3127e640e45ef0aa"
    },    
    {
      "mediaType": "application/vnd.cncf.helm.values.layer.v3+tar",
      "size": 1292,
      "digest": "sha256:62af13172ed974323f7c35153e8b23b8fa1c85355b6b26cc3127e640e45eeaf7"
    }
  ]
}
```

[distribution-spec]: https://github.com/opencontainers/distribution-spec/