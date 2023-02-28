import os
import re
import subprocess

FORK_CONFIG_PATH = "/home/ubuntu/mx-chain-devnet-config"


# Replace in all validators directory the folder named config by the one from the fork config
def replace_config_folder():
    cwd = os.getcwd()
    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        config_dir = os.path.join(cwd, validator_dir, 'config')
        if os.path.exists(config_dir):
            # save validatorKey.pem and p2p.toml inside /home/ubuntu/testnet/backup/validatorXX/config
            validator_number = re.match(r'validator(\d\d)', validator_dir).group(1)
            backup_dir = f'/home/ubuntu/testnet/backup/validator{validator_number}/config'
            if not os.path.exists(backup_dir):
                os.makedirs(backup_dir)
                subprocess.run(["cp", os.path.join(config_dir, 'validatorKey.pem'), backup_dir])
                subprocess.run(["cp", os.path.join(config_dir, 'p2p.toml'), backup_dir])
            subprocess.run(["rm", "-rf", config_dir])
        subprocess.run(["cp", "-R", FORK_CONFIG_PATH, config_dir])


replace_config_folder()
