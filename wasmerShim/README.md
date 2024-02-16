# Shim for wasmer 1

## Build

On Linux:

```
go build -buildmode=c-shared -o ../wasmer/libwasmer_linux_amd64.so .
go build -buildmode=c-shared -o ../wasmer/libwasmer_linux_arm64.so .
```

On MacOS:

```
go build -buildmode=c-shared -o ../wasmer/libwasmer_darwin_amd64.dylib .
go build -buildmode=c-shared -ldflags="-w" -o ../wasmer/libwasmer_darwin_arm64.dylib .
```
