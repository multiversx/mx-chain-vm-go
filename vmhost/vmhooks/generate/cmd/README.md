# VM Hooks Code generator

The code generator generates boilerplate code for both this Go VM, and for [the executor repository](https://github.com/multiversx/mx-vm-executor-rs)

For it to automatically copy files there, create a file called `wasm-vm-executor-rs-path.txt` here, in the `cmd` folder, contianing your local path to that repository, on your disk.

Finally, simply run `go generate` in `vmhost/vmhooks`.
