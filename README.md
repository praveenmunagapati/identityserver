# Public website


## Contribute

When you want to contribute to the development, follow the [contribution guidelines](contributing.md).

This website is intended to be used by the identityserver. For the identityserver to pick up the changes all html files and assets are packed in bytecode in go source files in the `packaged` folder.

In order to make the html files and assets available for the identityserver make sure you have go-bindata installed:
```
go get -u github.com/jteeuwen/go-bindata/...
```

After this execute `go generate` in the root of this repository. Check in the overwritten go files in the packaged folder.
