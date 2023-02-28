echo "Waiting for Elastic Search..."
sudo bash ./wait-for-it.sh elastic:9200 --timeout=0

echo "Compiling nodes..."
sudo chmod -R 777 testnet
cd testnet && rm -rf * && cd ..
mxpy testnet config

echo "Using devnet config..."
cp replace-config.py testnet/replace-config.py
cd testnet && python3 replace-config.py && cd ..
echo "Importing database..."
./copy-from-extracted.sh
./merge-db.sh
#./import-db.sh

echo "Changing preferences in order to act as a local testnet..."
cp change-prefs.py testnet/change-prefs.py
cd testnet && python3 change-prefs.py && cd ..

echo "Running testnet..."
cp start-testnet.py testnet/start-testnet.py
cd testnet && python3 start-testnet.py
