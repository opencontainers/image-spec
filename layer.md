# Creating an Image Filesystem Changeset

An example of creating an Image Filesystem Changeset follows.

An image root filesystem is first created as an empty directory.
Here is the initial empty directory structure for a changeset using the randomly-generated directory name `c3167915dc9d` ([actual layer DiffIDs are generated based on the content](#id_desc)).

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

This example change is going to add a configuration directory at `/etc/my-app.d` which contains a default config file.
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

When processing the second layer, `a/.wh..wh..opq` is applied first, before creating the new version of `a/b`, regardless of the ordering in which the whiteout file was encountered.
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
