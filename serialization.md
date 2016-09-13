# OCI Image Serialization

An *Image* is an ordered collection of root filesystem changes and the corresponding execution parameters for use within a container runtime.
This specification outlines the JSON format describing images for use with a container runtime and execution tool and its relationship to filesystem changesets, described in [Layers](layer.md).

## Terminology

This specification uses the following terms:

<dl>
    <dt>
        <a href="layer.md">Layer</a>
    </dt>
    <dd>
        Image filesystems are composed of <i>layers</i>.
        Each layer represents a set of filesystem changes in a tar-based <a href="layer.md">layer format</a>, recording files to be added, changed, or deleted relative to its parent layer.
	Layers do not have configuration metadata such as environment variables or default arguments - these are properties of the image as a whole rather than any particular layer.
        Using a layer-based or union filesystem such as AUFS, or by computing the diff from filesystem snapshots, the filesystem changeset can be used to present a series of image layers as if they were one cohesive filesystem.
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
        Layer DiffID
    </dt>
    <dd>
	A layer DiffID is a SHA256 digest over the layer's uncompressed tar archive and serialized in the descriptor digest format, e.g., <code>sha256:a9561eb1b190625c9adb5a9513e72c4dedafc1cb2d4c5236c9a6957ec7dfd5a9</code>.
	Layers must be packed and unpacked reproducibly to avoid changing the layer ID, for example by using tar-split to save the tar headers.
	NOTE: the DiffID is different than the digest in the manifest list because the manifest digest is taken over the gzipped layer for <code>application/vnd.oci.image.layer.tar+gzip</code> types.
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
</dl>

## Image JSON Description

Here is an example image JSON file:

```json,title=Image%20JSON&mediatype=application/vnd.oci.image.config.v1%2Bjson
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
    "/etc/some-config.d/": {}
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
            <code>type</code> which MUST be set to <code>layers</code>.
            Implementations MUST generate an error if they encounter a unknown value while verifying or unpacking an image.
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
