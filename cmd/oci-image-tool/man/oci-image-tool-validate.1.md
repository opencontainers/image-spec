% OCI(1) OCI-IMAGE-TOOL User Manuals
% OCI Community
% JULY 2016
# NAME
oci-image-tool-validate \- Validate one or more image files

# SYNOPSIS
**oci-image-tool validate** FILE... [flags]

# DESCRIPTION
`oci-image-tool validate` validates the given file(s) against the OCI image specification.


# OPTIONS
**--help**
  Print usage statement

**--ref**
  The ref pointing to the manifest to be validated. This must be present in the "refs" subdirectory of the image. Only applicable if type is image or imageLayout. (default "v1.0")

**--type**
  Type of the file to validate. If unset, oci-image-tool will try to auto-detect the type. One of "imageLayout,image,manifest,manifestList,config"

# EXAMPLES
```
$ skopeo copy docker://busybox oci:busybox-oci
$ oci-image-tool validate --type imageLayout --ref latest busybox-oci
busybox-oci: OK
```

# SEE ALSO
**oci-image-tool(1)**, **skopeo**(1)

# HISTORY
July 2016, Originally compiled by Antonio Murdaca (runcom at redhat dot com)
