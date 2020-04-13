import glob
import subprocess
import sys

def call(text):
    print(text)
    r = subprocess.call(text, shell=True)
    if r != 0:
        print('Failed')
        sys.exit(r)

call("./bin/rv64 -d -- res/riscv-tests/isa/rv64uc-u-rvc")
