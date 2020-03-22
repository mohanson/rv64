# RISC-V RV64IMAFDC Emulator

An outstanding RISC-V RV64IMAFDC(RV64GC) simulator.

# Install riscv-gnu-toolchain

First of all, riscv gnu toolchain must be installed. Source repo at [https://github.com/riscv/riscv-gnu-toolchain](https://github.com/riscv/riscv-gnu-toolchain), complete the build with the following commands:

```sh
$ ./configure --prefix=/opt/riscv --with-arch=rv64g
$ make
```

# Install rv64

```sh
$ mkdir bin
$ go build -o bin github.com/mohanson/rv64/cmd/make
$ ./bin/make
```

The binary file `rv64` will be located at the `./bin` directory. Could test it with the following command:

```sh
$ export RISCV=/opt/riscv
$ ./bin/make test
```

```
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoadd_d
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoadd_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoand_d
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoand_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amomax_d
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amomax_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amomaxu_d
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amomaxu_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amomin_d
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amomin_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amominu_d
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amominu_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoor_d
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoor_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoswap_d
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoswap_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoxor_d
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-amoxor_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ua-u-lrsc
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fadd
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fclass
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fcmp
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fcvt
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fcvt_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fdiv
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fmadd
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fmin
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-ldst
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-move
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-recoding
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-structural
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fadd
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fclass
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fcmp
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fcvt
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fcvt_w
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fdiv
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fmadd
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fmin
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-ldst
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-move
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-recoding
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-add
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-addi
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-addiw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-addw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-and
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-andi
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-auipc
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-beq
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-bge
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-bgeu
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-blt
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-bltu
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-bne
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-fence_i
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-jal
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-jalr
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-lb
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-lbu
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-ld
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-lh
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-lhu
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-lui
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-lw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-lwu
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-or
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-ori
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sb
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sd
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sh
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-simple
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sll
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-slli
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-slliw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sllw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-slt
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-slti
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sltiu
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sltu
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sra
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-srai
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sraiw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sraw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-srl
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-srli
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-srliw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-srlw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sub
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-subw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-sw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-xor
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64ui-u-xori
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-div
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-divu
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-divuw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-divw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-mul
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-mulh
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-mulhsu
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-mulhu
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-mulw
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-rem
[ok] $ ./bin/rv64 /tmp/riscv-tests/isa/rv64um-u-remu
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
$ ./bin/rv64 /tmp/fib
$ echo $?
# 55
```

# Licence

WTFPL.
