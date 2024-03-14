#!/bin/bash

function handle_sigterm {
  echo "Container is being stopped..."
  ./pause.sh
}

trap handle_sigterm SIGTERM
set -e

setup() {
  set -e

  rm -f results.json

  echo "Replacing localnet.toml..."
  sudo python3 replace_localnet_toml.py "$MX_LT_NUM_SHARDS"

  echo "Compiling nodes..."
  rm -rf localnet
  mxpy localnet config

  echo "Copying files..."

  # Workarounds because mxpy doesn't copy arm libs yet
  sudo cp /home/ubuntu/mx-chain-vm-go/wasmer2/libvmexeccapi.so /usr/lib/libvmexeccapi.so
  sudo cp /home/ubuntu/mx-chain-vm-go/wasmer2/libvmexeccapi_arm.so /usr/lib/libvmexeccapi_arm.so
  sudo cp /home/ubuntu/mx-chain-vm-go/wasmer/libwasmer_linux_amd64.so /usr/lib/libwasmer_linux_amd64.so
  sudo cp /home/ubuntu/mx-chain-vm-go/wasmer/libwasmer_linux_arm64_shim.so /usr/lib/libwasmer_linux_arm64_shim.so

  cp read_result.py localnet/read_result.py
  cp change_prefs.py localnet/change_prefs.py
  cp change_genesis.py localnet/change_genesis.py
  cp create_wallet.py localnet/create_wallet.py
  cp economics.toml localnet/economics.toml
  cp genesis.json localnet/genesis.json
  cp replace_economics.py localnet/replace_economics.py
  cp systemSmartContractsConfig.toml localnet/systemSmartContractsConfig.toml
  cp replace_system_contracts_config.py localnet/replace_system_contracts_config.py

  cp temp_replace_enable_epochs.py localnet/temp_replace_enable_epochs.py # Workaround while mxpy doesn't support v1.7.0 localnet
  cp temp_replace_genesis_smart_contracts.py localnet/temp_replace_genesis_smart_contracts.py # Workaround while mxpy doesn't support v1.7.0 localnet
  cp -R genesisContracts localnet/genesisContracts # Workaround while mxpy doesn't support v1.7.0 localnet

  echo "Changing nodes system contracts config..."
  cd localnet && sudo python3 replace_system_contracts_config.py && cd ..

  echo "Changing nodes preferences..."
  cd localnet && sudo python3 change_prefs.py "$MX_LT_ELASTIC_ENABLED" && cd ..

  echo "Changing nodes genesis..."
  cd localnet && sudo python3 change_genesis.py "$MX_LT_CUSTOM_EGLD_ADDRESS" "$MX_LT_NUM_SHARDS" && cd ..

  echo "Changing nodes economics..."
  cd localnet && sudo python3 replace_economics.py "$MX_RESULT_TOTAL_SUPPLY" && cd ..

  echo "Applying workaround for enableEpochs..."
  cd localnet && sudo python3 temp_replace_enable_epochs.py "$MX_LT_NUM_SHARDS" && cd ..

  echo "Applying workaround for genesisContracts..."
  cd localnet && sudo python3 temp_replace_genesis_smart_contracts.py && cd ..
}

if [ "$(python3 read_result.py "state")" != "paused" ]; then
  echo "Localnet is not paused, starting setup..."

  setup
fi

if [ "$MX_LT_ELASTIC_ENABLED" = "true" ]; then
  echo "Waiting for Elastic Search..."
  sudo bash ./wait-for-it.sh elastic:9200 --timeout=0
fi

sudo python3 add_result.py "state" "running"

echo "Running localnet..."
mxpy localnet start &
localnet_pid=$!

while kill -0 $localnet_pid > /dev/null 2>&1; do
  sleep 1
done
