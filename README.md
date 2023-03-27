# xNetwork

This project offers an effortless way to run a local testnet of MultiversX and its corresponding API without any configuration needed.
While you can run a local testnet of MultiversX using the [MultiversX CLI](https://multiversx.com/), this project provides a more comprehensive environment, including a proxy and a complete API.

## Features
- Run your own local testnet of MultiversX
- Choose the number of shards you prefer
- Create an initial address with 1,000,000 $EGLD
- Run the corresponding API, identical to https://api.multiversx.com
- Run MxOps scenes at startup

## Requirements
- [Docker](https://docs.docker.com/get-started/) and [Docker Compose](https://docs.docker.com/compose/gettingstarted/).

## Usage

1. Install the package via NPM:

    ```bash
    npm install -g @gfusee/xnetwork
    ```

2. That's it! Run the following command to run the tool:

    ```bash
    xnetwork
    ```

# Features documentation

## Give 1,000,000 EGLD to a custom address

By default, xNetwork will generate a wallet for you and give it 1,000,000 EGLD, the secret key will be print in the console after the network is up.
But you can choose in the CLI to give the EGLD to a custom address.

## Run MxOps scenes at startup

When creating a new network, you can choose to run [MxOps](https://github.com/Catenscia/MxOps) scenes from a selected folder at startup. By doing so, each scene will run under a scenario called `xnetwork`, so make sure you allow it in your scenes AND that the network `LOCAL` is allowed too.

The `xnetwork` scenario includes some helpful variables:

- An account named `xnetwork_genesis`, which is the account that has the initial 1,000,000 EGLD (if you chose not to use a custom address to give them)

Note that you can put all the files and folders you want in the scenes folder (`mxops_config.ini`, `.wasm`, `.pem`, etc...), and they will all be available in scenes under the path 'mxops'.

Here is an example of a valid scene file, assuming you have the `.wasm` in `<selected folder in the CLI>/contract/ping-pong/output/ping-pong.wasm` : 

```yaml
allowed_networks:
    - LOCAL

allowed_scenario:
  - ".*"

steps:

  - type: ContractDeploy
    sender: xnetwork_genesis
    wasm_path: "./contract/ping-pong/output/ping-pong.wasm"
    contract_id: "egld-ping-pong"
    gas_limit: 60000000
    arguments:
        - 500000000000000000
        - 1
    upgradeable: True
    readable: False
    payable: False
    payable_by_sc: True

  - type: ContractCall
    sender: xnetwork_genesis
    contract: "egld-ping-pong"
    endpoint: ping
    gas_limit: 3000000
    value: 500000000000000000

  - type: ContractCall
    sender: xnetwork_genesis
    contract: "egld-ping-pong"
    endpoint: pong
    gas_limit: 3000000

```

# Contributing

If you encounter any issues or would like to contribute to this project, feel free to open a pull request or an issue on GitHub.
