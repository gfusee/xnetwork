import os
import re
import subprocess
import sys
import json

genesis_egld_wallet = 'genesis-egld-wallet'

validator01_object = {
    "nickname": "validator01",
    "address": "erd1lgvcwjt3udzpuzajedz9malydaazwh0mg7apdd33j4a6qthyfz0q6nl2qv",
    "supply": "2500000000000000000000",
    "balance": "0",
    "stakingvalue": "2500000000000000000000",
    "delegation": {
        "address": "",
        "value": "0"
    }
}

validator02_object = {
    "nickname": "validator02",
    "address": "erd1fdxl8gkjcawmc7rw3s7fsdddxtqayjl85uf4vfa3ftqr0nmg3avqscfd7z",
    "supply": "2500000000000000000000",
    "balance": "0",
    "stakingvalue": "2500000000000000000000",
    "delegation": {
        "address": "",
        "value": "0"
    }
}

validator03_object = {
    "nickname": "validator03",
    "address": "erd12lqeap62f8qjfwktqxtv5avusu679xfnmcsl02u0wlx5em6nptqqfe8w49",
    "supply": "2500000000000000000000",
    "balance": "0",
    "stakingvalue": "2500000000000000000000",
    "delegation": {
        "address": "",
        "value": "0"
    }
}

validator04_object = {
    "nickname": "validator04",
    "address": "erd1qar9zjh2nlnsva4rpln7gnx6jgd4csce679au2h0xp64plav3mtslq9lrf",
    "supply": "2500000000000000000000",
    "balance": "0",
    "stakingvalue": "2500000000000000000000",
    "delegation": {
        "address": "",
        "value": "0"
    }
}

validator05_object = {
    "nickname": "validator05",
    "address": "erd1gwny3yg5zz6hrx0qg68thd88xvg2frqs3r6f3h6zdzq8gdewhf7qkwxtad",
    "supply": "2500000000000000000000",
    "balance": "0",
    "stakingvalue": "2500000000000000000000",
    "delegation": {
        "address": "",
        "value": "0"
    }
}

validator06_object = {
    "nickname": "validator06",
    "address": "erd1t5fhu2gwujzucs8jt35veax6tncmt436csr6wwl5n2kvg9faahmsztzhpd",
    "supply": "2500000000000000000000",
    "balance": "0",
    "stakingvalue": "2500000000000000000000",
    "delegation": {
        "address": "",
        "value": "0"
    }
}

validator07_object = {
    "nickname": "validator07",
    "address": "erd1t8k0awt4g64q5nme6r3dy8cr7yarm72z7pqx6zrayrj3g0t5eweshjpp9j",
    "supply": "2500000000000000000000",
    "balance": "0",
    "stakingvalue": "2500000000000000000000",
    "delegation": {
        "address": "",
        "value": "0"
    }
}

validator08_object = {
    "nickname": "validator08",
    "address": "erd1d86n3sn42ruyn0r324ayr83fweh69mh27tz5tt4r5cdyjma8wv4q0ncfel",
    "supply": "2500000000000000000000",
    "balance": "0",
    "stakingvalue": "2500000000000000000000",
    "delegation": {
        "address": "",
        "value": "0"
    }
}

validator09_object = {
    "nickname": "validator09",
    "address": "erd12vtnenfm384jdmejets9kqdsmfagwvvdhppwqcc07pynl3rkrqwqx6qtk9",
    "supply": "2500000000000000000000",
    "balance": "0",
    "stakingvalue": "2500000000000000000000",
    "delegation": {
        "address": "",
        "value": "0"
    }
}

validator10_object = {
    "nickname": "validator10",
    "address": "erd19g006qg4x72098cfykcqtqtczjzzlllc6d3wht6r67ky7n6uzuaszknrms",
    "supply": "2500000000000000000000",
    "balance": "0",
    "stakingvalue": "2500000000000000000000",
    "delegation": {
        "address": "",
        "value": "0"
    }
}

validators_objects = [validator01_object, validator02_object, validator03_object, validator04_object, validator05_object, validator06_object, validator07_object, validator08_object, validator09_object, validator10_object]

genesis_object = {
    "nickname": "genesis",
    "address": "${MX_LT_GENESIS_ADDRESS}",
    "supply": "10000000000000000000000000",
    "balance": "10000000000000000000000000",
    "stakingvalue": "0",
    "delegation": {
        "address": "",
        "value": "0"
    }
}


def replace_in_files():
    cwd = os.getcwd()
    genesis_address = sys.argv[1]
    shards_count = int(sys.argv[2])

    add_result_command = 'python3 /home/ubuntu/add_result.py'

    if len(genesis_address) == 0:
        print("No genesis address provided, generating a new one...")
        genesis_pem_path = os.path.join(cwd, genesis_egld_wallet, 'wallet.pem')
        subprocess.run(f'python3 create_wallet.py "{genesis_egld_wallet}"', shell=True)
        genesis_address = subprocess.check_output(f"mxpy wallet pem-address {genesis_pem_path}", shell=True).decode('utf-8').strip()
        print("Generated genesis address: " + genesis_address)
        subprocess.run(f'{add_result_command} "genesisEgldPemPath" "{genesis_pem_path}"', shell=True)

    subprocess.run(f'{add_result_command} "genesisEgldAddress" "{genesis_address}"', shell=True)

    genesis_object['address'] = genesis_address

    total_supply = 0

    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        external_json_path = os.path.join(cwd, validator_dir, 'config', 'genesis.json')
        subprocess.run(f"rm {external_json_path}", shell=True)
        subprocess.run(f"cp genesis.json {external_json_path}", shell=True)

        # Open the existing JSON file for reading
        with open(external_json_path, 'r') as f:
            existing_data = json.load(f)

        existing_data = validators_objects[0:shards_count] + existing_data

        existing_data.append(genesis_object)

        # Open the JSON file for writing and write the updated data
        with open(external_json_path, 'w') as f:
            json.dump(existing_data, f)

        if total_supply == 0:
            new_total_supply = 0
            for existing_object in existing_data:
                new_total_supply += int(existing_object['supply'])

            total_supply = new_total_supply

            subprocess.run(f'{add_result_command} "totalSupply" "{str(new_total_supply)}"', shell=True)


replace_in_files()
