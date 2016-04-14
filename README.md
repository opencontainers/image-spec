# Open Container Distributable Image Specification

The [Open Container Initiative](http://www.opencontainers.org/) develops specifications for standards on Operating System process and application containers.

This OCI project is tasked with creating a software shipping container image format spec (OCI Image Format) with security and naming as components.

## Initial Proposal

This new OCI project intends to start with the Docker v2.2 specification, improve any remaining technical concerns, and standardize and improve the understood properties of a container image format. This new project will have the objectives of:

* A serialized image format (base layer)
* A process of hashing the image format for integrity and content-addressing (base layer)
* Signatures that are based on signing image content address (optional layer)
* Naming that is federated based on DNS and can be delegated (optional layer)

## Cooperation with OCI Runtime Project

The [OCI Runtime Spec project](https://github.com/opencontainers/runtime-spec) is developing a specification for the lifecycle of a running container. The OCI Image Format Spec project should work with the OCI Runtime Spec project so that the image can support the UX that users have come to expect from container engines like Docker and rkt: primarily, the ability to run an image with no additional arguments:

* docker run example.com/org/app:v1.0.0
* rkt run example.com/org/app,version=v1.0.0

This implies that the OCI Image Format must contain sufficient information to launch the application on the target platform (e.g. command, arguments, environment variables, etc).

## FAQ

**Q: Why doesn't this project mention distribution?**

A: Distribution, for example using HTTP as both Docker v2.2 and AppC do today, is currently out of scope on the [OCI Scope Table](https://www.opencontainers.org/governance/oci-scope-table). There has been [some discussion on the TOB mailing list]( https://groups.google.com/a/opencontainers.org/d/msg/tob/A3JnmI-D-6Y/tLuptPDHAgAJ) to make distribution an optional layer but this topic is a work in progress.

**Q: Why a new project?**

A: The first OCI spec centered around defining the run side of a container. This is generally seen to be an orthogonal concern to the shipping container component. As practical examples of this separation you see many organizations separating these concerns into different teams and organizations: the Docker Distribution project and the Docker containerd project; Amazon ECS and Amazon EC2 Container Registry, etc.

**Q: Why start this work now?**

A: We are seeing many independent implementations of container image handling including build systems, registries, and image analysis tools. As an organization we would like to encourage this growth and bring people together to ensure a technically correct and open specification continues to evolve reflecting the OCI values.

**Q: What happens to AppC or Docker Image Formats?**

A: Existing formats can continue to be a proving ground for technologies, as needed. The OCI Image Format project should strive to provide a dependable open specification that can be shared between different tools and be evolved for years or decades of compatibility; as the deb and rpm format have.

## Roadmap

The current roadmap can be found in the [GitHub milestones](https://github.com/opencontainers/image-spec/milestones)

* April v0.0.0
  * Import Docker v2.2 format
* April v0.1.0
  * Spec factored for top to bottom reading with three audiences in-mind:
    * Build system creators
    * Image registry creators
    * Container engine creators
* May v0.2.0
  * Release version of spec with improvements from two independent experimental implementations from OCI members e.g. Amazon Container Registry and rkt
* June v1.0.0
  * Release initial version of spec with two independent non-experimental implementations from OCI members

# Contributing

Development happens on GitHub for the spec.
Issues are used for bugs and actionable items and longer discussions can happen on the [mailing list](#mailing-list).

The specification and code is licensed under the Apache 2.0 license found in the `LICENSE` file of this repository.

## Code of Conduct

Participation in the OpenContainers community is governed by [OpenContainer's Code of Conduct](https://github.com/opencontainers/tob/blob/master/code-of-conduct.md).

## Discuss your design

The project welcomes submissions, but please let everyone know what you are working on.

Before undertaking a nontrivial change to this specification, send mail to the [mailing list](#mailing-list) to discuss what you plan to do.
This gives everyone a chance to validate the design, helps prevent duplication of effort, and ensures that the idea fits.
It also guarantees that the design is sound before code is written; a GitHub pull-request is not the place for high-level discussions.

Typos and grammatical errors can go straight to a pull-request.
When in doubt, start on the [mailing-list](#mailing-list).

## Mailing List

You can subscribe and join the mailing list on [Google Groups](https://groups.google.com/a/opencontainers.org/forum/#!forum/dev).

## IRC

OCI discussion happens on #opencontainers on Freenode ([logs][irc-logs]).

## Markdown style

To keep consistency throughout the Markdown files in the Open Container spec all files should be formatted one sentence per line.
This fixes two things: it makes diffing easier with git and it resolves fights about line wrapping length.
For example, this paragraph will span three lines in the Markdown source.

## Git commit

### Sign your work

The sign-off is a simple line at the end of the explanation for the patch, which certifies that you wrote it or otherwise have the right to pass it on as an open-source patch.
The rules are pretty simple: if you can certify the below (from [developercertificate.org](http://developercertificate.org/)):

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
660 York Street, Suite 102,
San Francisco, CA 94110 USA

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.


Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

then you just add a line to every git commit message:

    Signed-off-by: Joe Smith <joe@gmail.com>

using your real name (sorry, no pseudonyms or anonymous contributions.)

You can add the sign off when creating the git commit via `git commit -s`.

### Commit Style

Simple house-keeping for clean git history.
Read more on [How to Write a Git Commit Message](http://chris.beams.io/posts/git-commit/) or the Discussion section of [`git-commit(1)`](http://git-scm.com/docs/git-commit).

1. Separate the subject from body with a blank line
2. Limit the subject line to 50 characters
3. Capitalize the subject line
4. Do not end the subject line with a period
5. Use the imperative mood in the subject line
6. Wrap the body at 72 characters
7. Use the body to explain what and why vs. how
  * If there was important/useful/essential conversation or information, copy or include a reference
8. When possible, one keyword to scope the change in the subject (i.e. "README: ...", "runtime: ...")

[irc-logs]: http://ircbot.wl.linuxfoundation.org/eavesdrop/%23opencontainers/
