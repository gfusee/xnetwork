import subprocess
import sys

wallet_name = sys.argv[1]

wallet_output = f"{wallet_name}/wallet.pem"

def create_wallet():
    subprocess.run(f"mkdir -p {wallet_name}", shell=True)
    raw_mnemonic = subprocess.check_output(f'mxpy wallet new --format pem --outfile {wallet_output}', shell=True)

create_wallet()
