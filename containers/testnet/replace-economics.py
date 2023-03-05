import os
import re
import subprocess


def replace_file():
    cwd = os.getcwd()
    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        external_toml_path = os.path.join(cwd, validator_dir, 'config', 'economics.toml')
        subprocess.run(f"rm {external_toml_path}", shell=True)
        subprocess.run(f"cp economics.toml {external_toml_path}", shell=True)


replace_file()
