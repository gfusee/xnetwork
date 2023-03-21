import os


def replace_in_files():
    cwd = os.getcwd()
    external_toml_path = os.path.join(cwd, 'proxy', 'config', 'config.toml')
    if os.path.exists(external_toml_path):
        with open(external_toml_path, 'r') as file:
            contents = file.read()
        contents = contents.replace('EconomicsMetricsCacheValidityDurationSec = 600',
                                    'EconomicsMetricsCacheValidityDurationSec = 600')
        with open(external_toml_path, 'w') as file:
            file.write(contents)


replace_in_files()
