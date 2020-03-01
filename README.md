# RISC-V Emulator

\[English\] \[[中文](./README_CN.md)\]

I wrote this simulator to understand riscv more accurately. RISC-V is awesome, but also very young, I hope my work can provide reference value for latecomers.

# Install riscv-gnu-toolchain

Repo: [https://github.com/riscv/riscv-gnu-toolchain](https://github.com/riscv/riscv-gnu-toolchain)

```sh
$ ./configure --prefix=/opt/riscv --with-arch=rv64g
$ make
```

# Install riscv-tests

Repo: [https://github.com/riscv/riscv-tests](https://github.com/riscv/riscv-tests)

```sh
$ export RISCV=/opt/riscv
```

```sh
$ git clone https://github.com/riscv/riscv-tests
$ cd riscv-tests
$ git submodule update --init --recursive
$ autoconf
$ ./configure --prefix=$RISCV/target
$ make
$ make install
```
