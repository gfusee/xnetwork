import os
import re
import sys

should_use_elastic = sys.argv[1] == 'true'

prefs_to_add = f"""
OverridableConfigTomlValues = [
   {{ File = "external.toml", Path = "ElasticSearchConnector.Enabled", Value = "{should_use_elastic}" }},
   {{ File = "external.toml", Path = "ElasticSearchConnector.URL", Value = "http://elastic:9200" }}
]
"""


def replace_in_files():
    cwd = os.getcwd()
    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        external_toml_path = os.path.join(cwd, validator_dir, 'config', 'prefs.toml')
        if os.path.exists(external_toml_path):
            with open(external_toml_path, 'r') as file:
                contents = file.read()
            contents = contents + prefs_to_add
            with open(external_toml_path, 'w') as file:
                file.write(contents)


replace_in_files()
