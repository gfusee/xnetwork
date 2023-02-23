import os
import re
import subprocess

prefs_to_add = """

OverridableConfigTomlValues = [
   { File = "external.toml", Path = "ElasticSearchConnector.Enabled", Value = "true" },
   { File = "external.toml", Path = "ElasticSearchConnector.URL", Value = "http://elastic:9200" },
   { File = "config.toml", Path = "GeneralSettings.StartInEpochEnabled", Value = "false" },
   { File = "config.toml", Path = "GeneralSettings.ChainID", Value = "D" }
]
"""

FORK_CONFIG_PATH = "/home/ubuntu/mx-chain-devnet-config"


# Replace in all validators directory the folder named config by the one from the fork config
def replace_config_folder():
    cwd = os.getcwd()
    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        config_dir = os.path.join(cwd, validator_dir, 'config')
        if os.path.exists(config_dir):
            subprocess.run(["rm", "-rf", config_dir])
        subprocess.run(["cp", "-R", FORK_CONFIG_PATH, config_dir])


def replace_in_files():
    cwd = os.getcwd()
    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        # # Replace in config/nodesSetup.json the chainID property to "D"
        # nodes_setup_path = os.path.join(cwd, validator_dir, 'config', 'nodesSetup.json')
        # json_content = ""
        # with open(nodes_setup_path, 'r') as file:
        #     json_content = file.read()
        #     json_content = json_content.replace('"chainID": "local-testnet"', '"chainID": "D"')
        #     json_content = re.sub(r'"startTime": \d+', '"startTime": 1648551600', json_content)
        #     # json_content = json_content.replace('"consensusGroupSize": 1', '"consensusGroupSize": 21')
        #     # json_content = json_content.replace('"minNodesPerShard": 1', '"minNodesPerShard": 58')
        #     # json_content = json_content.replace('"metaChainConsensusGroupSize": 1', '"metaChainConsensusGroupSize": 58')
        #     # json_content = json_content.replace('"metaChainMinNodes": 1', '"metaChainMinNodes": 58')
        # with open(nodes_setup_path, 'w') as file:
        #     file.write(json_content)

        external_toml_path = os.path.join(cwd, validator_dir, 'config', 'prefs.toml')
        if os.path.exists(external_toml_path):
            with open(external_toml_path, 'r') as file:
                contents = file.read()
            contents = contents + prefs_to_add
            with open(external_toml_path, 'w') as file:
                file.write(contents)


replace_config_folder()
replace_in_files()
