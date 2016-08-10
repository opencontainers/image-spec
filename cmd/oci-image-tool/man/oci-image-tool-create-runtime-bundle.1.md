% OCI(1) OCI-IMAGE-TOOL User Manuals
% OCI Community
% JULY 2016
# NAME
oci-image-tool-create-runtime-bundle \- Create an OCI image runtime bundle

# SYNOPSIS
**oci-image-tool create-runtime-bundle** [src] [dest] [flags]

# DESCRIPTION
`oci-image-tool create-runtime-bundle` generates an [OCI bundle](https://github.com/opencontainers/runtime-spec/blob/master/bundle.md) from an [OCI image layout](https://github.com/opencontainers/image-spec/blob/master/image-layout.md).


# OPTIONS
**--help**
  Print usage statement

**--ref**
  The ref pointing to the manifest of the OCI image. This must be present in the "refs" subdirectory of the image. (default "v1.0")

**--rootfs**
  A directory representing the root filesystem of the container in the OCI runtime bundle. It is strongly recommended to keep the default value. (default "rootfs")

**--type**
  Type of the file to unpack. If unset, oci-image-tool will try to auto-detect the type. One of "imageLayout,image"

# EXAMPLES
```
$ skopeo copy docker://busybox oci:busybox-oci
$ mkdir busybox-bundle
$ oci-image-tool create-runtime-bundle --ref latest busybox-oci busybox-bundle
$ cd busybox-bundle && sudo runc start busybox
[...]
```

# SEE ALSO
**oci-image-tool(1)**, **runc**(1), **skopeo**(1)

# HISTORY
July 2016, Originally compiled by Antonio Murdaca (runcom at redhat dot com)
