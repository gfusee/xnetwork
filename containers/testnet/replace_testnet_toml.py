import os
import sys

num_shards = sys.argv[1]

toml_file_content = f"""
[networking]
port_proxy = 7950

[shards]
num_shards = {num_shards}

[software]
resolution = "local_prebuilt_cmd_folders"

[software.local_prebuilt_cmd_folders]
# Has to contain "node" binary, "config" folder, wasmer libraries
mx_chain_go_node = "/home/ubuntu/mx-chain-go/cmd/node"
# Has to contain "seednode" binary, "config" folder, wasmer libraries
mx_chain_go_seednode = "/home/ubuntu/mx-chain-go/cmd/seednode"
# Has to contain "proxy" binary, "config" folder, wasmer libraries
mx_chain_proxy_go = "/home/ubuntu/mx-chain-proxy-go/cmd/proxy"
"""


def replace_toml_file():
    # Delete testnet.toml if it exists
    testnet_toml_path = os.path.join(os.getcwd(), 'testnet.toml')

    # Create testnet.toml
    with open(testnet_toml_path, 'w') as file:
        file.write(toml_file_content)


replace_toml_file()
