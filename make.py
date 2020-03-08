import glob
import subprocess

for e in glob.glob("/tmp/riscv-tests/isa/rv64ua-u-amo*_d"):
    print(e, subprocess.call(f"go run cmd/rv64/main.go {e}", shell=True))

# subprocess.call("go run cmd/make/main.go make test", shell=True)
