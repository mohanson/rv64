# RISC-V RV64IMAFDC Emulator

I wrote this simulator to understand riscv more accurately. RISC-V is awesome, but also very young, I hope my work can provide reference value for latecomers.

# Install riscv-gnu-toolchain

First of all, riscv gnu toolchain must be installed. Source repo at [https://github.com/riscv/riscv-gnu-toolchain](https://github.com/riscv/riscv-gnu-toolchain), complete the build with the following commands:

```sh
$ ./configure --prefix=/opt/riscv --with-arch=rv64g
$ make
```

# Installation

```sh
$ mkdir bin
$ go build -o bin github.com/mohanson/rv64/cmd/make
$ ./bin/make make
```

The binary file `rv64` will be located at the `./bin` directory. Could test it with the following command:

```sh
$ export RISCV=/opt/riscv
$ ./bin/make test
```

# Licence

WTFPL.
