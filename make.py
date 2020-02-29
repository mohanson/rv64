import subprocess

subprocess.call("go run cmd/riscv_emu/main.go ./build/res/program/minimal", shell=True)
