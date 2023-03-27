import subprocess

import read_result

scenario_name = "xnetwork"


def run():
    genesis_wallet = None
    try:
        genesis_wallet = read_result.get_result('genesisEgldPemPath')
    except Exception as e:
        pass

    init_scenes = [
        "../mxops-init/genesis_account.yaml"
    ]

    if genesis_wallet is not None:
        subprocess.run(f"python3 -m mxops execute -n LOCAL -s {scenario_name} mxops-init/genesis_account.yaml",
                       shell=True)

    init_scenes_string = " ".join(init_scenes)

    subprocess.run(f"python3 -m mxops execute -n LOCAL -s {scenario_name} {init_scenes_string} .", shell=True, cwd="mxops")


if __name__ == '__main__':
    run()
