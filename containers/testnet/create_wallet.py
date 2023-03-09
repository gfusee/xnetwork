import subprocess
import sys

mxpy_path = "/home/ubuntu/multiversx-sdk/mxpy"

wallet_name = sys.argv[1]

wallet_output = f"{wallet_name}/wallet.pem"


def create_wallet():
    raw_mnemonic = subprocess.check_output(f'{mxpy_path} wallet new', shell=True)
    # Remove "Mnemonic: " from the beginning of the string
    mnemonic = raw_mnemonic[9:].decode('utf-8').strip()
    subprocess.run(f"mkdir -p {wallet_name}", shell=True)
    subprocess.run(f"{mxpy_path} wallet derive {wallet_output} --mnemonic", shell=True, input=mnemonic, text=True)


create_wallet()
