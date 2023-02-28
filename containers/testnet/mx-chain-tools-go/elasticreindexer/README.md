# Elastic reindexer tool

- This tool can be used to copy an Elasticsearch index/indices from a cluster to another. It can also copy the indices 
mappings (we recommend to use it without copying the mappings because most often they are not exposed). The indices mappings part can be done with the `indices-creator` tool.
Also, by using the `indices-creator` tool, some properties of the indices can be changed (depending on the needs).

- The main scope of these tools is to copy all the information from an Elasticsearch cluster to another (all the indices that were
populated by the Elrond nodes from the genesis until the current time).

## How to use
### STEP 1
- In order to reindex all the information from an Elasticsearch cluster to another, one has to create all the indices mappings. 
To do that one has to use the `indices-creator` tool:

    ```
    cd elasticreindexer
    cd cmd/indices-creator/
    go build 
    ```
  
- After the build is done one has to update the `config/cluster.toml` file with the information about the Elasticsearch cluster. In the `cluster.toml` one has to set the URL 
of the Elasticsearch cluster, and what indices to be used to create the mappings.

- Optionally, the mappings can be also customized (open file with mappings for every index and set different settings, based on the needs).

- Run `./indices-creator` in order to create all the indices and mappings.

_**note:** STEP 1 can be skipped for the clusters that already have the information indexed._ 

***

### STEP 2
- After the mappings and indices were created we can start to copy all the information. In order to do this one has to use the `elasticreindexer` tool that 
will copy every index one by one.

    ```
    cd elasticreindexer
    cd cmd/elasticreindexer/
    go build 
    ```
- Update the Elasticsearch instance configuration for both source and destination inside `config.toml`. In the `config.toml` file one has to set the 
URL of the Elasticsearch `input` instance ( the one from where we have to copy all the information ) and the `output` instance (the one where all the information 
from the `input` instance will be copied).

- The `config.toml` file contains by default all the Elasticsearch indices that are populated by an Elrond observing-squad.

- Also, if you want to copy indices with timestamp you have to set the `blockchain-start-time` in the `config.toml` file (by default is the one from the mainnet).

- Instances that include a timestamp are already defined in the configuration file. Their cloning will be much faster due to the parallel execution on batches split depending on timestamp.

- Run `./elasticreindexer --skip-mappings` (will start to reindex all the information from the input cluster in the output cluster based on the `config.toml` file).


_**WARN**: Start the observing-squad only after the indices `accounts`, `accountsesdt` and `tokens` are copied._

***

#### SPEED UP STEP 2
- The `STEP 2` will take lot of time because `hundreds of gigabytes` of data have to be fetched. In order to speed up the process we can do the 
next things:

1. Run an instance of `elasticreindexer` only with the indices from the list `indices-no-timestamp`.
2. Run multiple instances of `elasticreindexer` each with the next configurations for `indices-with-timestamp`:

    `a.` `indices-with-timestamp = ["accountsesdt", "tokens"]`

    `b.` `indices-with-timestamp = ["blocks", "receipts", "rounds"]`

    `c.` `indices-with-timestamp = ["transactions", "miniblocks", "scdeploys"]`

    `d.` `indices-with-timestamp = ["accountshistory", "scresults"]`

    `e.` `indices-with-timestamp = ["accountsesdthistory"]`

    `f.` `indices-with-timestamp = ["logs"]`

    `g.` `indices-with-timestamp = ["operations"]`


_**note:** For all the configs from `a-g` the field `indices-no-timestamp` has to be empty or commented._

_**WARN:** Start the observing-squad only after the instance with `indices-not-timetamp` finished to copy `accounts` index 
and the instance from the point `a.` finished to copy `accountsesdt` and `tokens` indices._

***

## Audience

This tool should be as generic as possible, and it shouldn't have any custom code related to Elrond instances
of Elasticsearch.
