import subprocess

subprocess.call("go run cmd/riscv_emu/main.go ./bin/fuzz", shell=True)
