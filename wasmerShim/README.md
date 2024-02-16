# Shim for wasmer 1

## Build

On MacOS:

```
go build -buildmode=c-shared -ldflags="-w" -o ../wasmer/libwasmer_darwin_arm64_shim.dylib .

install_name_tool -id @rpath/libwasmer_darwin_arm64_shim.dylib ../wasmer/libwasmer_darwin_arm64_shim.dylib
```
