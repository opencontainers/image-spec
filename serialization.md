# OCI Image Serialization

An *Image* is an ordered collection of root filesystem changes and the corresponding execution parameters for use within a container runtime.
This specification outlines the format of these filesystem changes and corresponding parameters and describes how to create and use them for use with a container runtime and execution tool.

## Terminology

This specification uses the following terms:

<dl>
    <dt>
        Layer
    </dt>
    <dd>
        Images are composed of <i>layers</i>.
	Each layer is a set of filesystem changes.
	Layers do not have configuration metadata such as environment variables or default arguments - these are properties of the image as a whole rather than any particular layer.
    </dd>
    <dt>
        Image JSON
    </dt>
    <dd>
        Each image has an associated JSON structure which describes some basic information about the image such as date created, author, and the ID of its parent image as well as execution/runtime configuration like its entrypoint, default arguments, CPU/memory shares, networking, and volumes.
	The JSON structure also references a cryptographic hash of each layer used by the image, and provides history information for those layers.
	This JSON is considered to be immutable, because changing it would change the computed ImageID.
	Changing it means creating a new derived image, instead of changing the existing image.
    </dd>
    <dt>
        Image Filesystem Changeset
    </dt>
    <dd>
        Each layer has an archive of the files which have been added, changed, or deleted relative to its parent layer.
	Using a layer-based or union filesystem such as AUFS, or by computing the diff from filesystem snapshots, the filesystem changeset can be used to present a series of image layers as if they were one cohesive filesystem.
    </dd>
    <dt>
        Layer DiffID
    </dt>
    <dd>
	A layer DiffID is a SHA256 digest over the layer's uncompressed tar archive and serialized in the descriptor digest format, e.g., <code>sha256:a9561eb1b190625c9adb5a9513e72c4dedafc1cb2d4c5236c9a6957ec7dfd5a9</code>.
	Layers must be packed and unpacked reproducibly to avoid changing the layer ID, for example by using tar-split to save the tar headers.
	NOTE: the DiffID is different than the digest in the manifest list because the manifest digest is taken over the gzipped layer for `application/vnd.oci.image.serialization.rootfs.tar.gzip` types.
    </dd>
    <dt>
        Layer ChainID
    </dt>
    <dd>
        For convenience, it is sometimes useful to refer to a stack of layers with a single identifier.
	This is called a <code>ChainID</code>.
	For a
        single layer (or the layer at the bottom of a stack), the
        <code>ChainID</code> is equal to the layer's <code>DiffID</code>.
        Otherwise the <code>ChainID</code> is given by the formula:
        <code>ChainID(layerN) = SHA256hex(ChainID(layerN-1) + " " + DiffID(layerN))</code>.
    </dd>
    <dt>
        ImageID <a name="id_desc"></a>
    </dt>
    <dd>
        Each image's ID is given by the SHA256 hash of its configuration JSON.
	It is represented as a hexadecimal encoding of 256 bits, e.g., <code>sha256:a9561eb1b190625c9adb5a9513e72c4dedafc1cb2d4c5236c9a6957ec7dfd5a9</code>.
	Since the configuration JSON that gets hashed references hashes of each layer in the image, this formulation of the ImageID makes images content-addresable.
    </dd>
    <dt>
        Tag
    </dt>
    <dd>
        A tag serves to map a descriptive, user-given name to any single image ID.
	Tag values are limited to the set of characters <code>[a-zA-Z_0-9]</code>.
    </dd>
    <dt>
        Repository
    </dt>
    <dd>
        A collection of tags grouped under a common prefix (the name component before <code>:</code>).
	For example, in an image tagged with the name <code>my-app:3.1.4</code>, <code>my-app</code> is the <i>Repository</i> component of the name.
	A repository name is made up of slash-separated name components, optionally prefixed by a DNS hostname.
	The hostname must follow comply with standard DNS rules, but may not contain <code>_</code> characters.
	If a hostname is present, it may optionally be followed by a port number in the format <code>:8080</code>.
	Name components may contain lowercase characters, digits, and separators.
	A separator is defined as a period, one or two underscores, or one or more dashes.
	A name component may not start or end with a separator.
    </dd>
</dl>

## Image JSON Description

Here is an example image JSON file:

```json,title=Image%20JSON&mediatype=application/vnd.oci.image.serialization.config.v1%2Bjson
{  
    "created": "2015-10-31T22:22:56.015925234Z",
    "author": "Alyssa P. Hacker <alyspdev@example.com>",
    "architecture": "amd64",
    "os": "linux",
    "config": {
        "User": "alice",
        "Memory": 2048,
        "MemorySwap": 4096,
        "CpuShares": 8,
        "ExposedPorts": {  
            "8080/tcp": {}
        },
        "Env": [  
            "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
            "FOO=oci_is_a",
            "BAR=well_written_spec"
        ],
        "Entrypoint": [
            "/bin/my-app-binary"
        ],
        "Cmd": [
            "--foreground",
            "--config",
            "/etc/my-app.d/default.cfg"
        ],
        "Volumes": {
            "/var/job-result-data": {},
            "/var/log/my-app-logs": {}
        },
        "WorkingDir": "/home/alice"
    },
    "rootfs": {
      "diff_ids": [
        "sha256:c6f988f4874bb0add23a778f753c65efe992244e148a1d2ec2a8b664fb66bbd1",
        "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"
      ],
      "type": "layers"
    },
    "history": [
      {
        "created": "2015-10-31T22:22:54.690851953Z",
        "created_by": "/bin/sh -c #(nop) ADD file:a3bc1e842b69636f9df5256c49c5374fb4eef1e281fe3f282c65fb853ee171c5 in /"
      },
      {
        "created": "2015-10-31T22:22:55.613815829Z",
        "created_by": "/bin/sh -c #(nop) CMD [\"sh\"]",
        "empty_layer": true
      }
    ]
}
```

Note: whitespace has been added to this example for clarity. Whitespace is OPTIONAL and implementations MAY have compact JSON with no whitespace.

### Image JSON Field Descriptions

<dl>
    <dt>
        created <code>string</code>
    </dt>
    <dd>
        ISO-8601 formatted combined date and time at which the image was created.
    </dd>
    <dt>
        author <code>string</code>
    </dt>
    <dd>
        Gives the name and/or email address of the person or entity which created and is responsible for maintaining the image.
    </dd>
    <dt>
        architecture <code>string</code>
    </dt>
    <dd>
        The CPU architecture which the binaries in this image are built to run on.
	Possible values include:
        <ul>
            <li>386</li>
            <li>amd64</li>
            <li>arm</li>
        </ul>
        More values may be supported in the future and any of these may or may not be supported by a given container runtime implementation.
        New entries SHOULD be submitted to this specification for standardization and be inspired by the <a href=https://golang.org/doc/install/source#environment>Go language documentation for $GOOS and $GOARCH</a>.
    </dd>
    <dt>
        os <code>string</code>
    </dt>
    <dd>
        The name of the operating system which the image is built to run on.
        Possible values include:
        <ul>
            <li>darwin</li>
            <li>freebsd</li>
            <li>linux</li>
        </ul>
        More values may be supported in the future and any of these may or may not be supported by a given container runtime implementation.
        New entries SHOULD be submitted to this specification for standardization and be inspired by the <a href=https://golang.org/doc/install/source#environment>Go language documentation for $GOOS and $GOARCH</a>.
    </dd>
    <dt>
        config <code>struct</code>
    </dt>
    <dd>
        The execution parameters which should be used as a base when running a container using the image.
	This field can be <code>null</code>, in which case any execution parameters should be specified at creation of the container.

        <h4>Container RunConfig Field Descriptions</h4>

        <dl>
            <dt>
                User <code>string</code>
            </dt>
            <dd>
                <p>
		The username or UID which the process in the container should run as.
		This acts as a default value to use when the value is not specified when creating a container.
		</p>

                <p>All of the following are valid:</p>

                <ul>
                    <li><code>user</code></li>
                    <li><code>uid</code></li>
                    <li><code>user:group</code></li>
                    <li><code>uid:gid</code></li>
                    <li><code>uid:group</code></li>
                    <li><code>user:gid</code></li>
                </ul>

                <p>
		If <code>group</code>/<code>gid</code> is not specified, the default group and supplementary groups of the given <code>user</code>/<code>uid</code> in <code>/etc/passwd</code> from the container are applied.
		</p>
            </dd>
            <dt>
                Memory <code>integer</code>
            </dt>
            <dd>
                Memory limit (in bytes).
		This acts as a default value to use when the value is not specified when creating a container.
            </dd>
            <dt>
                MemorySwap <code>integer</code>
            </dt>
            <dd>
                Total memory usage (memory + swap); set to <code>-1</code> to disable swap.
		This acts as a default value to use when the value is not specified when creating a container.
            </dd>
            <dt>
                CpuShares <code>integer</code>
            </dt>
            <dd>
                CPU shares (relative weight vs. other containers).
		This acts as a default value to use when the value is not specified when creating a container.
            </dd>
            <dt>
                ExposedPorts <code>struct</code>
            </dt>
            <dd>
                A set of ports to expose from a container running this image.
                This JSON structure value is unusual because it is a direct JSON serialization of the Go type <code>map[string]struct{}</code> and is represented in JSON as an object mapping its keys to an empty object.
		Here is an example:

<pre>{
    "8080": {},
    "53/udp": {},
    "2356/tcp": {}
}</pre>

                Its keys can be in the format of:
                <ul>
                    <li>
                        <code>"port/tcp"</code>
                    </li>
                    <li>
                        <code>"port/udp"</code>
                    </li>
                    <li>
                        <code>"port"</code>
                    </li>
                </ul>
                with the default protocol being <code>"tcp"</code> if not specified.

                These values act as defaults and are merged with any specified when creating a container.
            </dd>
            <dt>
                Env <code>array of strings</code>
            </dt>
            <dd>
                Entries are in the format of <code>VARNAME="var value"</code>.
                These values act as defaults and are merged with any specified when creating a container.
            </dd>
            <dt>
                Entrypoint <code>array of strings</code>
            </dt>
            <dd>
                A list of arguments to use as the command to execute when the container starts.
		This value acts as a  default and is replaced by an entrypoint specified when creating a container. This field MAY be "null".
            </dd>
            <dt>
                Cmd <code>array of strings</code>
            </dt>
            <dd>
                Default arguments to the entrypoint of the container.
		These values act as defaults and are replaced with any specified when creating a container.
		If an <code>Entrypoint</code> value is not specified, then the first entry of the <code>Cmd</code> array should be interpreted as the executable to run. This field MAY be "null".
            </dd>
            <dt>
                Volumes <code>struct</code>
            </dt>
            <dd>
                A set of directories which should be created as data volumes in a container running this image. This field MAY be "null".
                <p>
                If a file or folder exists within the image with the same path as a data volume, that file or folder is replaced with the data volume and is never merged.
                </p>
		This JSON structure value is unusual because it is a direct JSON serialization of the Go type <code>map[string]struct{}</code> and is represented in JSON as an object mapping its keys to an empty object.
		Here is an example:
<pre>{
    "/var/my-app-data/": {},
    "/etc/some-config.d/": {},
}</pre>
            </dd>
            <dt>
                WorkingDir <code>string</code>
            </dt>
            <dd>
                Sets the current working directory of the entrypoint process in the container.
		This value acts as a default and is replaced by a working directory specified when creating a container.
            </dd>
        </dl>
    </dd>
    <dt>
        rootfs <code>struct</code>
    </dt>
    <dd>
        The rootfs key references the layer content addresses used by the image.
	This makes the image config hash depend on the filesystem hash.
        rootfs has two subkeys:

        <ul>
          <li>
            <code>type</code> is usually set to <code>layers</code>.
	    There is also a Windows-specific value <code>layers+base</code> that allows a base layer to be specified in a field of <code>rootfs</code> called <code>base_layer</code>.
          </li>
          <li>
            <code>diff_ids</code> is an array of layer content hashes (<code>DiffIDs</code>), in order from bottom-most to top-most.
          </li>
        </ul>


        Here is an example rootfs section:

<pre>"rootfs": {
  "diff_ids": [
    "sha256:c6f988f4874bb0add23a778f753c65efe992244e148a1d2ec2a8b664fb66bbd1",
    "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef",
    "sha256:13f53e08df5a220ab6d13c58b2bf83a59cbdc2e04d0a3f041ddf4b0ba4112d49"
  ],
  "type": "layers"
}</pre>
    </dd>
    <dt>
        history <code>struct</code>
    </dt>
    <dd>
        <code>history</code> is an array of objects describing the history of each layer.
	The array is ordered from bottom-most layer to top-most layer.
	The object has the following fields.

        <ul>
          <li>
            <code>created</code>: Creation time, expressed as a ISO-8601 formatted
            combined date and time
          </li>
          <li>
            <code>author</code>: The author of the build point
          </li>
          <li>
            <code>created_by</code>: The command which created the layer
          </li>
          <li>
            <code>comment</code>: A custom message set when creating the layer
          </li>
          <li>
            <code>empty_layer</code>: This field is used to mark if the history item created a filesystem diff.
	    It is set to true if this history item doesn't correspond to an actual layer in the rootfs section (for example, a command like ENV which results in no change to the filesystem).
          </li>
        </ul>

Here is an example history section:

<pre>"history": [
  {
    "created": "2015-10-31T22:22:54.690851953Z",
    "created_by": "/bin/sh -c #(nop) ADD file:a3bc1e842b69636f9df5256c49c5374fb4eef1e281fe3f282c65fb853ee171c5 in /"
  },
  {
    "created": "2015-10-31T22:22:55.613815829Z",
    "created_by": "/bin/sh -c #(nop) CMD [\"sh\"]",
    "empty_layer": true
  }
]</pre>
    </dd>
</dl>

Any extra fields in the Image JSON struct are considered implementation specific and should be ignored by any implementations which are unable to interpret them.

## Creating an Image Filesystem Changeset

An example of creating an Image Filesystem Changeset follows.

An image root filesystem is first created as an empty directory.
Here is the initial empty directory structure for the a changeset using the randomly-generated directory name `c3167915dc9d` ([actual layer DiffIDs are generated based on the content](#id_desc)).

```
c3167915dc9d/
```

Files and directories are then created:

```
c3167915dc9d/
    etc/
        my-app-config
    bin/
        my-app-binary
        my-app-tools
```

The `c3167915dc9d` directory is then committed as a plain Tar archive with entries for the following files:

```
etc/my-app-config
bin/my-app-binary
bin/my-app-tools
```

To make changes to the filesystem of this container image, create a new directory, such as `f60c56784b83`, and initialize it with a snapshot of the parent image's root filesystem, so that the directory is identical to that of `c3167915dc9d`.
NOTE: a copy-on-write or union filesystem can make this very efficient:

```
f60c56784b83/
    etc/
        my-app-config
    bin/
        my-app-binary
        my-app-tools
```

This example change is going add a configuration directory at `/etc/my-app.d` which contains a default config file.
There's also a change to the `my-app-tools` binary to handle the config layout change.
The `f60c56784b83` directory then looks like this:

```
f60c56784b83/
    etc/
        .wh.my-app-config
        my-app.d/
            default.cfg
    bin/
        my-app-binary
        my-app-tools
```

This reflects the removal of `/etc/my-app-config` and creation of a file and directory at `/etc/my-app.d/default.cfg`.
`/bin/my-app-tools` has also been replaced with an updated version.
Before committing this directory to a changeset, because it has a parent image, it is first compared with the directory tree of the parent snapshot, `f60c56784b83`, looking for files and directories that have been added, modified, or removed.
The following changeset is found:

```
Added:      /etc/my-app.d/default.cfg
Modified:   /bin/my-app-tools
Deleted:    /etc/my-app-config
```

A Tar Archive is then created which contains *only* this changeset:

- Added and modified files and directories in their entirety
- Deleted files or directory marked with a whiteout file

A whiteout file is an empty file that prefixes the deleted paths basename `.wh.`.
When a whiteout is found in the upper changeset of a filesystem, any matching name in the lower changeset is ignored, and the whiteout itself is also hidden.
As files prefixed with `.wh.` are special whiteout tombstones it is not possible to create a filesystem which has a file or directory with a name beginning with `.wh.`.

The resulting Tar archive for `f60c56784b83` has the following entries:

```
/etc/my-app.d/default.cfg
/bin/my-app-tools
/etc/.wh.my-app-config
```

Whiteout files MUST only apply to resources in lower layers.
Files that are present in the same layer as a whiteout file can only be hidden by whiteout files in subsequent layers.
The following is a base layer with several resources:

```
a/
a/b/
a/b/c/
a/b/c/bar
```

When the next layer is created, the original `a/b` directory is deleted and recreated with `a/b/c/foo`:

```
a/
a/.wh..wh..opq
a/b/
a/b/c/
a/b/c/foo
```

When processing the second layer, `a/.wh..wh..opq` is applied first, before creating the new version of `a/b, regardless of the ordering in which the whiteout file was encountered.
For example, the following layer is equivalent to the layer above:

```
a/
a/b/
a/b/c/
a/b/c/foo
a/.wh..wh..opq
```

Implementations SHOULD generate layers such that the whiteout files appear before sibling directory entries.

In addition to expressing that a single entry should be removed from a lower layer, layers may remove all of the children using an opaque whiteout entry.
An opaque whiteout entry is a file with the name `.wh..wh..opq` indicating that all siblings are hidden in the lower layer.
Let's take the following base layer as an example:

```
etc/
	my-app-config
bin/
	my-app-binary
	my-app-tools
	tools/
		my-app-tool-one
```

If all children of `bin/` are removed, the next layer would have the following:

```
bin/
	.wh..wh..opq
```

This is called _opaque whiteout_ format.
An _opaque whiteout_ file hides _all_ children of the `bin/` including sub-directories and all descendents.
Using _explicit whiteout_ files, this would be equivalent to the following:

```
bin/
	.wh.my-app-binary
	.wh.my-app-tools
	.wh.tools
```

In this case, a unique whiteout file is generated for each entry.
If there were more children of `bin/` in the base layer, there would be an entry for each.
Note that this opaque file will apply to _all_ children, including sub-directories, other resources and all descendents.

Implementations SHOULD generate layers using _explicit whiteout_ files, but MUST accept both.

Any given image is likely to be composed of several of these Image Filesystem Changeset tar archives.
