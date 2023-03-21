setup() {
  rm results.json

  echo "Replacing testnet.toml..."
  sudo python3 replace_testnet_toml.py "$MX_LT_NUM_SHARDS"

  echo "Compiling nodes..."
  rm -rf testnet
  mxpy testnet config

  echo "Copying files..."
  cp change_proxy_config.py testnet/change_proxy_config.py
  cp read_result.py testnet/read_result.py
  cp change_prefs.py testnet/change_prefs.py
  cp change_genesis.py testnet/change_genesis.py
  cp change_enable_epochs.py testnet/change_enable_epochs.py
  cp create_wallet.py testnet/create_wallet.py
  cp economics.toml testnet/economics.toml
  cp genesis.json testnet/genesis.json
  cp replace_economics.py testnet/replace_economics.py
  cp systemSmartContractsConfig.toml testnet/systemSmartContractsConfig.toml
  cp replace_system_contracts_config.py testnet/replace_system_contracts_config.py

  echo "Changing proxy config..."
  cd testnet && sudo python3 change_proxy_config.py && cd ..

  echo "Changing nodes system contracts config..."
  cd testnet && sudo python3 replace_system_contracts_config.py && cd ..

  echo "Changing nodes preferences..."
  cd testnet && sudo python3 change_prefs.py "$MX_LT_ELASTIC_ENABLED" && cd ..

  echo "Changing nodes genesis..."
  cd testnet && sudo python3 change_genesis.py "$MX_LT_CUSTOM_EGLD_ADDRESS" "$MX_LT_NUM_SHARDS" && cd ..

  echo "Changing nodes enable epochs..."
  cd testnet && sudo python3 change_enable_epochs.py "$MX_LT_ENABLE_EPOCHS" && cd ..

  echo "Changing nodes economics..."
  cd testnet && sudo python3 replace_economics.py "$MX_RESULT_TOTAL_SUPPLY" && cd ..
}

if [ "$(python3 read_result.py "state")" != "paused" ]; then
  echo "Testnet is not paused, starting setup..."

  setup
fi

if [ "$MX_LT_ELASTIC_ENABLED" = "true" ]; then
  echo "Waiting for Elastic Search..."
  sudo bash ./wait-for-it.sh elastic:9200 --timeout=0
fi

sudo python3 add_result.py "state" "running"

echo "Running testnet..."
mxpy testnet start
