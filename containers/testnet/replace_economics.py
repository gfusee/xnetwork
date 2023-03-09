import os
import re
import subprocess
import read_result

total_supply = read_result.get_result('totalSupply')


def replace_file():
    cwd = os.getcwd()
    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        external_toml_path = os.path.join(cwd, validator_dir, 'config', 'economics.toml')
        subprocess.run(f"rm {external_toml_path}", shell=True)
        subprocess.run(f"cp economics.toml {external_toml_path}", shell=True)

        with open(external_toml_path, 'r') as f:
            filedata = f.read()

        print("Total supply: " + str(total_supply))

        filedata = filedata.replace('GenesisTotalSupply = "${MX_RESULT_TOTAL_SUPPLY}"',
                                    f'GenesisTotalSupply = "{total_supply}"')

        with open(external_toml_path, 'w') as f:
            f.write(filedata)


replace_file()
