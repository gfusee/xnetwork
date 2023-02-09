import os
import re

elasticNotEnabledString = """
[ElasticSearchConnector]
    ## We do not recommend to activate this indexer on a validator node since
    #the node might loose rating (even facing penalties) due to the fact that
    #the indexer is called synchronously and might block due to external causes.
    #Strongly suggested to activate this on a regular observer node.
    Enabled           = false
    IndexerCacheSize  = 0
    BulkRequestMaxSizeInBytes = 4194304 # 4MB
    URL               = "http://localhost:9200"
"""

elasticEnabledString = elasticNotEnabledString.replace("Enabled           = false", "Enabled           = true")
elasticFinalString = elasticEnabledString.replace('URL               = "http://localhost:9200"', 'URL               = "http://elastic:9200"')

def replace_in_external_toml():
    cwd = os.getcwd()
    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        external_toml_path = os.path.join(cwd, validator_dir, 'config', 'external.toml')
        if os.path.exists(external_toml_path):
            with open(external_toml_path, 'r') as file:
                contents = file.read()
            contents = contents.replace(elasticNotEnabledString, elasticFinalString)
            with open(external_toml_path, 'w') as file:
                file.write(contents)

replace_in_external_toml()
