import os
import re
import subprocess
import shutil

genesis_smart_contracts_content = """
[
  {
    "owner": "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
    "filename": "./config/genesisContracts/delegation.wasm",
    "vm-type": "0500",
    "init-parameters": "%validator_sc_address%@03E8@00@030D40@030D40",
    "type": "delegation",
    "version": "0.4.*"
  },
  {
    "owner": "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
    "filename": "./config/genesisContracts/dns.wasm",
    "vm-type": "0500",
    "init-parameters": "00",
    "type": "dns",
    "version": "0.2.*"
  }
]
"""

def temp_replace_genesis_smart_contracts():
    cwd = os.getcwd()

    genesis_contracts_dir_path = os.path.join(cwd, 'genesisContracts')

    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        genesis_smart_contracts_json_path = os.path.join(cwd, validator_dir, 'config', 'genesisSmartContracts.json')
        with open(genesis_smart_contracts_json_path, 'w') as file:
            file.write(genesis_smart_contracts_content)

        validator_genesis_contracts_dir_path = os.path.join(cwd, validator_dir, 'config', 'genesisContracts')

        if os.path.exists(validator_genesis_contracts_dir_path):
            shutil.rmtree(validator_genesis_contracts_dir_path)

        shutil.copytree(genesis_contracts_dir_path, validator_genesis_contracts_dir_path)

temp_replace_genesis_smart_contracts()
