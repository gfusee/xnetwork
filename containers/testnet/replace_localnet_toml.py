import os
import sys

num_shards = sys.argv[1]

toml_file_content = f"""
[networking]
port_proxy = 7950

[shards]
num_shards = 3

[software.mx_chain_go]
resolution = "local"
local_path = "/home/ubuntu/mx-chain-go"

[software.mx_chain_proxy_go]
resolution = "local"
local_path = "/home/ubuntu/mx-chain-proxy-go"
"""


def replace_toml_file():
    # Delete testnet.toml if it exists
    testnet_toml_path = os.path.join(os.getcwd(), 'localnet.toml')

    # Create testnet.toml
    with open(testnet_toml_path, 'w') as file:
        file.write(toml_file_content)


replace_toml_file()
