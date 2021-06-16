<!-- Copyright 2021 Adarga Limited -->
<!-- SPDX-License-Identifier: Apache-2.0 -->

# Contributing to pachyderm-pipeline-controller

We are glad you want to contribute to an Adarga open source project!
This document will help answer common questions you may have.

## Submitting Issues

## Contribution Process

We have a 3-step process for contributions:

1. Commit changes to a git branch, making sure to sign-off those changes for the Developer Certificate of Origin
2. Create a GitHub Pull Request for your change, following the instructions in the pull request template.
3. Perform a Code Review with the project maintainers on the pull request.

### Pull Request Requirements

In order to maintain a high quality code base, we require all pull requests meet these specifications:

1. Code is clean; do `make check` to run the code base through various lint and formatting checks.
2. Change includes unit tests.
3. PRs should be conflict free unless there has been some prior discussion/agreement with Adarga maintainers that conflicts will be resolved in merge.

### Code Review Process

Code review takes place in GitHub pull requests. Once you open a pull request, project maintainers will review
your code and respond to your pull request with any feedback they might have. The process is as follows

1. One or more of the Adarga maintainers must approve your PR.
2. Your change will be merged into the project's `main` branch.
3. _[when we add this] the version will be incremented automatically and the binaries built._

### Developer Certification of Origin (DCO)

Licensing is very important to open source projects. It helps ensure
the software continues to be available under the terms that the
author desired.

This project uses the Apache-2.0 license to strike a balance between open
contribution and allowing you to use the software however you would
like to.

The license tells you what rights you have that are provided by the
copyright holder. It is important that the contributor fully
understands what rights they are licensing and agrees to them.
Sometimes the copyright holder isn't the contributor, such as when
the contributor is doing work on behalf of a company.

To make a good faith effort to ensure these criteria are met, Adarga
requires the Developer Certificate of Origin (DCO) process to be
followed.

The DCO is an attestation attached to every contribution made by
every developer. In the commit message of the contribution, the
developer simply adds a `Signed-off-by` statement and thereby agrees
to the DCO, which you can find below or at
http://developercertificate.org/.

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
1 Letterman Drive
Suite D4700
San Francisco, CA, 94129

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

### DCO Sign-Off Methods

The DCO requires a sign-off message in the following format appear on each commit in the pull request:

```
Signed-off-by: Your Name <your-email@address.com>
```

The DCO text can either be manually added to your commit body, or
you can add either `-s` or `--signoff` to your usual git commit
commands. If you are using the GitHub UI to make a change you can
add the sign-off message directly to the commit message when creating
the pull request. If you forget to add the sign-off you can also
amend a previous commit with the sign-off by running `git commit
--amend -s`. If you've pushed your changes to GitHub already you'll
need to force push your branch after this with `git push -f`.
