import os
import re
import subprocess

prefs_to_add = """
OverridableConfigTomlValues = [
   { File = "external.toml", Path = "ElasticSearchConnector.Enabled", Value = "false" },
   { File = "external.toml", Path = "ElasticSearchConnector.URL", Value = "http://elastic:9200" },
   { File = "config.toml", Path = "GeneralSettings.StartInEpochEnabled", Value = "false" },
   { File = "config.toml", Path = "GeneralSettings.ChainID", Value = "D" }
]
"""


def replace_in_files():
    cwd = os.getcwd()
    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        # Replace in config/nodesSetup.json the chainID property to "D"
        nodes_setup_path = os.path.join(cwd, validator_dir, 'config', 'nodesSetup.json')
        json_content = ""
        with open(nodes_setup_path, 'r') as file:
            json_content = file.read()
            json_content = json_content.replace('"chainID": "local-testnet"', '"chainID": "D"')
            json_content = re.sub(r'"startTime": \d+', '"startTime": 1648551600', json_content)
            # Do the same as commented lines below, but with regex
            # json_content = json_content.replace('"consensusGroupSize": 1', '"consensusGroupSize": 21')
            # json_content = json_content.replace('"minNodesPerShard": 1', '"minNodesPerShard": 58')
            # json_content = json_content.replace('"metaChainConsensusGroupSize": 1', '"metaChainConsensusGroupSize": 58')
            # json_content = json_content.replace('"metaChainMinNodes": 1', '"metaChainMinNodes": 58')

            json_content = re.sub(r'"consensusGroupSize": \d+', '"consensusGroupSize": 1', json_content)
            json_content = re.sub(r'"minNodesPerShard": \d+', '"minNodesPerShard": 1', json_content)
            json_content = re.sub(r'"metaChainConsensusGroupSize": \d+', '"metaChainConsensusGroupSize": 1', json_content)
            json_content = re.sub(r'"metaChainMinNodes": \d+', '"metaChainMinNodes": 1', json_content)
        with open(nodes_setup_path, 'w') as file:
            file.write(json_content)

        external_toml_path = os.path.join(cwd, validator_dir, 'config', 'prefs.toml')
        if os.path.exists(external_toml_path):
            with open(external_toml_path, 'r') as file:
                contents = file.read()
            contents = contents + prefs_to_add
            with open(external_toml_path, 'w') as file:
                file.write(contents)

        validator_number = re.match(r'validator(\d\d)', validator_dir).group(1)
        backup_dir = f'/home/ubuntu/testnet/backup/validator{validator_number}/config'
        subprocess.run(["cp", os.path.join(cwd, backup_dir, 'validatorKey.pem'), os.path.join(cwd, validator_dir, 'config')])
        subprocess.run(["cp", os.path.join(cwd, backup_dir, 'p2p.toml'), os.path.join(cwd, validator_dir, 'config')])


replace_in_files()
