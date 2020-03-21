import glob
import subprocess
import sys

# for e in glob.glob("/tmp/riscv-tests/isa/rv64ua-u-amo*_d"):
#     print(e, subprocess.call(f"go run cmd/rv64/main.go {e}"))


def call(text):
    print(text)
    r = subprocess.call(text, shell=True)
    if r != 0:
        print('Failed')
        sys.exit(r)


call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fadd")
# call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fclass")
# call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fcmp")
# call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fcvt")
# call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fcvt_w")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fdiv")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fmadd")
# call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-fmin")
# call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-ldst")
# call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-move")
# call("./bin/rv64 /tmp/riscv-tests/isa/rv64uf-u-recoding")

call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fadd")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fclass")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fcmp")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fcvt")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fcvt_w")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fdiv")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fmadd")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-fmin")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-ldst")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-move")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-recoding")
call("./bin/rv64 /tmp/riscv-tests/isa/rv64ud-u-structural")
