# RISC-V RV64IMAFDC Emulator

An outstanding RISC-V RV64IMAFDC(RV64GC) simulator.

# Install riscv-gnu-toolchain

First of all, riscv gnu toolchain must be installed. Source repo is at [https://github.com/riscv/riscv-gnu-toolchain](https://github.com/riscv/riscv-gnu-toolchain), complete the build with the following commands:

```sh
$ apt install autoconf automake autotools-dev curl python3 libmpc-dev libmpfr-dev libgmp-dev gawk build-essential bison flex texinfo gperf libtool patchutils bc zlib1g-dev libexpat-dev
```

```sh
$ git clone --recursive https://github.com/riscv/riscv-gnu-toolchain
$ cd riscv-gnu-toolchain
$ git submodule update --init --recursive
$ ./configure --prefix=/opt/riscv --with-arch=rv64gc
$ make
```

# Install rv64

```sh
$ git clone https://github.com/mohanson/rv64
$ cd rv64
$ go run cmd/make/main.go
```

The binary file `rv64` will be located at the `./bin` directory. Could test it with the following command:

```sh
$ export RISCV=/opt/riscv
$ ./bin/make test

[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoadd_d
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoadd_w
# Many lines ...
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-remuw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-remw
```

# Lick it

Let's compile a simple C file, which implements the fibonacci function.

```c
int fib(int n) {
    if (n <= 1) {
        return n;
    }
    return fib(n - 1) + fib(n - 2);
}

int main() {
    return fib(10);
}
```

```sh
$ /opt/riscv/bin/riscv64-unknown-elf-gcc -o /tmp/fib ./res/program/fib.c
$ ./bin/rv64 -- /tmp/fib
$ echo $?
# 55
```

# Licence

WTFPL.
