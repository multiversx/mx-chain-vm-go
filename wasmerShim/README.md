# Shim for wasmer 1 (legacy)

## Build

On Linux AMD64: not applicable (not needed).

On MacOS AMD64: not applicable (not needed).

On MacOS ARM64:

```
go build -buildmode=c-shared -ldflags="-w" -o ../wasmer/libwasmer_darwin_arm64_shim.dylib .

install_name_tool -id @rpath/libwasmer_darwin_arm64_shim.dylib ../wasmer/libwasmer_darwin_arm64_shim.dylib
```

On Linux ARM64:

```
go build -buildmode=c-shared -o ../wasmer/libwasmer_linux_arm64_shim.so .

patchelf --set-soname libwasmer_linux_arm64_shim.so ../wasmer/libwasmer_linux_arm64_shim.so
```
