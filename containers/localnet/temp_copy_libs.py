import os
import re
import subprocess

def temp_copy_libs():
    cwd = os.getcwd()

    wasmer2_lib_files = [
        "libvmexeccapi.so",
        "libvmexeccapi_arm.so",
        "libvmexeccapi_arm.dylib",
    ]

    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        for wasmer2_lib_file in wasmer2_lib_files:
            source_path = os.path.join("/home/ubuntu/mx-chain-vm-go/wasmer2", wasmer2_lib_file)
            dest_path = os.path.join(cwd, validator_dir, wasmer2_lib_file)
            subprocess.run(f"cp {source_path} {dest_path}", shell=True)



temp_copy_libs()
