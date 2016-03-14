# Contributing

## Conventions

Update the documentation when creating or modifying features. Test your
documentation changes for clarity, concision, and correctness, as well as a
clean documentation build.

Write clean code. Universally formatted code promotes ease of writing, reading,
and maintenance.

Pull request descriptions should be as clear as possible and include a reference
to all the issues that they address.

Commit messages must start with a capitalized and short summary (max. 50 chars)
written in the imperative, followed by an optional, more detailed explanatory
text which is separated from the summary by an empty line.

Code review comments may be added to your pull request. Discuss, then make the
suggested modifications and push additional commits to your feature branch. Post
a comment after pushing. New commits show up in the pull request automatically,
but the reviewers are notified only when you comment.

Pull requests must be cleanly rebased on top of master without multiple branches
mixed into the PR.

**Git tip**: If your PR no longer merges cleanly, use `rebase master` in your
feature branch to update your pull request rather than `merge master`.

Before you make a pull request, squash your commits into logical units of work
using `git rebase -i` and `git push -f`. A logical unit of work is a consistent
set of patches that should be reviewed together: for example, upgrading the
version of a vendored dependency and taking advantage of its now available new
feature constitute two separate units of work. Implementing a new function and
calling it in another file constitute a single logical unit of work. The very
high majority of submissions should have a single commit, so if in doubt: squash
down to one.

After every commit, make sure the test suite passes. Include documentation
changes in the same pull request so that a revert would remove all traces of
the feature or fix.

Include an issue reference like `Closes #XXXX` or `Fixes #XXXX` in commits that
close an issue. Including references automatically closes the issue on a merge.

## Packaging

This website is intended to be used by the identityserver. For the identityserver to pick up the changes,
all html files and assets are packed in go source files in the `packaged` folder.

In order to make the html files and assets available for the identityserver make sure you have go-bindata installed:
```
go get -u github.com/jteeuwen/go-bindata/...
```

After this execute `go generate` in the root of this repository. Check in the overwritten go files in the packaged folder.

During development it can be easier if the files are served directly, execute go-bindata with the -debug flag:
```
go-bindata -debug -pkg assets -prefix assets -o ./packaged/assets/assets.go assets/...
go-bindata -debug -pkg thirdpartyassets -prefix thirdpartyassets -o ./packaged/thirdpartyassets/thirdpartyassets.go thirdpartyassets/...
go-bindata -debug -pkg components -prefix components -o ./packaged/components/components.go components/...
go-bindata -debug -pkg html -o ./packaged/html/html.go index.html registration.html login.html home.html error.html apidocumentation.html

```

## Bower dependencies

Although 3rd party dependencies are installed through bower,
only the relevant files should be checked in and be in the thirdpartyassets folder when packaging using `go generate`.
