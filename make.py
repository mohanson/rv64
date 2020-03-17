import glob
import subprocess

# for e in glob.glob("/tmp/riscv-tests/isa/rv64ua-u-amo*_d"):
#     print(e, subprocess.call(f"go run cmd/rv64/main.go {e}", shell=True))

subprocess.call("./bin/rv64 -d /tmp/riscv-tests/isa/rv64ud-u-fadd", shell=True)
