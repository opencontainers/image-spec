% OCI(1) OCI-IMAGE-TOOL User Manuals
% OCI Community
% AUGUST 2016
# NAME
oci-image-tool-cas \- Content-addressable storage manipulation

# SYNOPSIS
**oci-image-tool cas** [command]

# DESCRIPTION
`oci-image-tool cas` manipulates content-addressable storage.

# OPTIONS
**--help**
  Print usage statement

# COMMANDS
**get**
  Retrieve a blob from the store.
  See **oci-image-tool-cas-get**(1) for full documentation on the **get** command.

**put**
  Write a blob to the store.
  See **oci-image-tool-cas-put**(1) for full documentation on the **put** command.

# EXAMPLES
```
$ oci-image-tool init image-layout image.tar
$ echo hello | oci-image-tool cas put image.tar
sha256:5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03
$ oci-image-tool cas get image.tar sha256:5891b5b522d5df086d0ff0b110fbd9d21bb4fc7163af34d08286a2e846f6be03
hello
```

# SEE ALSO
**oci-image-tool**(1), **oci-image-tool-cas-get**(1), **oci-image-tool-cas-put**(1), **oci-image-tool-init**(1)

# HISTORY
August 2016, Originally compiled by W. Trevor King (wking at tremily dot us)
