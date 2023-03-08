import os
import re
import subprocess
import sys

mxpy_path = "/home/ubuntu/multiversx-sdk/mxpy"

genesis_egld_wallet = 'genesis-egld-wallet'


def replace_in_files():
    cwd = os.getcwd()
    genesis_address = sys.argv[1]

    if len(genesis_address) == 0:
        print("No genesis address provided, generating a new one...")
        genesis_pem_path = os.path.join(cwd, genesis_egld_wallet, 'wallet.pem')
        subprocess.run(f'python3 create-wallet.py "{genesis_egld_wallet}"', shell=True)
        genesis_address = subprocess.check_output(f"{mxpy_path} wallet pem-address {genesis_pem_path}", shell=True).decode('utf-8').strip()
        print("Generated genesis address: " + genesis_address)
        subprocess.run(f'python3 /home/ubuntu/add-result.py "genesisEgldPemPath" "{genesis_pem_path}"', shell=True)

    subprocess.run(f'python3 /home/ubuntu/add-result.py "genesisEgldAddress" "{genesis_address}"', shell=True)

    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        external_json_path = os.path.join(cwd, validator_dir, 'config', 'genesis.json')
        subprocess.run(f"rm {external_json_path}", shell=True)
        subprocess.run(f"cp genesis.json {external_json_path}", shell=True)
        if os.path.exists(external_json_path):
            with open(external_json_path, 'r') as file:
                contents = file.read()
            contents = contents.replace('"address": "${MX_LT_GENESIS_ADDRESS}",', f'"address": "{genesis_address}",')
            with open(external_json_path, 'w') as file:
                file.write(contents)


replace_in_files()
