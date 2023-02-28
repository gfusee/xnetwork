## Description

This tool is used to export the storage of an account at a given time. 

## Input

The input required by this tool is:
- the AccountsTrie database from a snapshot
- the State Root Hash used to recreate the trie
- the Address to export the storage for


## Output

If everything is successfully, the map containing the key-value pairs from the address storage will be printed in the console and also written in the `output.json`

Note that both the keys and the values are hex-encoded.

## How to use

1. compile the binary by issuing a `go build` command in elrond-tools-go/trieTools/accountStorageExporter directory
2. create a `db` directory and place inside directories `0`, `1` ... that contains the state data
3. start the app with the following parameters: 
   `./accountStorageExporter --log-level *:DEBUG --log-save --hex-roothash c93be73e9e1d8918ea240523372bc3094aa4bbc7221000300a493a6ae593b348 --address erd1qqqqqqqqqqqqqpgqhe8t5jewej70zupmh44jurgn29psua5l2jps3ntjj3` 
   
where `c93be73e9e1d8918ea240523372bc3094aa4bbc7221000300a493a6ae593b348` is the required trie hash to be checked and erd1qqqqqqqqqqqqqpgqhe8t5jewej70zupmh44jurgn29psua5l2jps3ntjj3 is the address to export the storage for
