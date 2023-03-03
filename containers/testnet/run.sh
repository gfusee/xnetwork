if [ "$MX_LT_ELASTIC_ENABLED" = "true" ]; then
  echo "Waiting for Elastic Search..."
  sudo bash ./wait-for-it.sh elastic:9200 --timeout=0
fi

echo "Replacing testnet.toml..."
sudo python3 replace-testnet-toml.py "$MX_LT_NUM_SHARDS"

echo "Compiling nodes..."
rm -rf testnet
mxpy testnet config

echo "Copying python script..."
cp change-prefs.py testnet/change-prefs.py

echo "Changing nodes preferences..."
cd testnet && sudo python3 change-prefs.py "$MX_LT_ELASTIC_ENABLED" && cd ..

echo "Running testnet..."
mxpy testnet start
