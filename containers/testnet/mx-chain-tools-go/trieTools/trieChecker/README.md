## Description

This tool is used to check the correct loading of a state trie. Works with user accounts main tries 
(it will automatically check the referenced data tries), data tries and peer account tries.

# How to use

1. compile the binary by issuing a `go build` command in elrond-tools-go/trieTools/trieChecker directory
2. create a `db` directory and place inside directories `0`, `1` ... that contains the state data, alternatively, you can place a randomly named directory and use that solely to load the data
3. start the app with the following parameters: `./trieChecker -log-level *:DEBUG -log-save -hex-roothash c93be73e9e1d8918ea240523372bc3094aa4bbc7221000300a493a6ae593b348` where `c93be73e9e1d8918ea240523372bc3094aa4bbc7221000300a493a6ae593b348` is the required trie hash to be checked
