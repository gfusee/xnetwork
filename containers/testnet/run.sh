if [ "$MX_LT_ELASTIC_ENABLED" = "true" ]; then
  echo "Waiting for Elastic Search..."
  sudo bash ./wait-for-it.sh elastic:9200 --timeout=0
fi

echo "Replacing testnet.toml..."
sudo python3 replace-testnet-toml.py "$MX_LT_NUM_SHARDS"

echo "Compiling nodes..."
rm -rf testnet
mxpy testnet config

echo "Copying files..."
cp change-prefs.py testnet/change-prefs.py
cp change-genesis.py testnet/change-genesis.py
cp create-wallet.py testnet/create-wallet.py
cp economics.toml testnet/economics.toml
cp genesis.json testnet/genesis.json
cp replace-economics.py testnet/replace-economics.py
cp systemSmartContractsConfig.toml testnet/systemSmartContractsConfig.toml
cp replace-system-contracts-config.py testnet/replace-system-contracts-config.py

echo "Changing nodes economics..."
cd testnet && sudo python3 replace-economics.py && cd ..

echo "Changing nodes system contracts config..."
cd testnet && sudo python3 replace-system-contracts-config.py && cd ..

echo "Changing nodes preferences..."
cd testnet && sudo python3 change-prefs.py "$MX_LT_ELASTIC_ENABLED" && cd ..

echo "Changing nodes genesis..."
cd testnet && sudo python3 change-genesis.py "$MX_LT_CUSTOM_EGLD_ADDRESS" && cd ..

echo "Running testnet..."
mxpy testnet start
