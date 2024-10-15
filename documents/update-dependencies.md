#  Update Dependencies

This document is about how to update the immortal repository dependencies to latest version.

### Packages

First of all you need to update golang dependencies to latest version using this commands:

```sh
go get -u ./...
go mod tidy
```
Once all packages got updated, make sure you run `make build` and `make test` commands to make sure
none of previous behaviors are broken. If any packages had breaking changes or some of them are deprecated, you need to
update the code and use new methods or use another package.

### Dev tools

After packages, you need to update dev tools such as `golangci-lint`, etc.

You can go to root [make file](../makefile) and find all dev tools on devtools part.
You have to find latest version of dev tools and replace them here.

> Note: consider breaking changes and deprecated packages for devtools too.

### Go version

You have to update the go version to latest release in [go.mod](../go.mod).
Make sure you are updating version of Golang on [Dockerfile](../dockerfile).

> Note: you must run `make build` after this change to make sure everything works smoothly.

### CI/CD and GitHub workflows

You need to go to [workflows](../.github/workflows) directory and update old GitHub actions to latest version.
You can find the latest version by searching the action name on GitHub.

