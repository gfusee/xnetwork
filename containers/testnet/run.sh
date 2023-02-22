echo "Waiting for Elastic Search..."
sudo bash ./wait-for-it.sh elastic:9200 --timeout=0

echo "Compiling nodes..."
sudo chmod -R 777 testnet
cd testnet && rm -rf * && cd ..
mxpy testnet config

echo "Changing nodes prefs..."
cp change-prefs.py testnet/change-prefs.py
cd testnet && python3 change-prefs.py && cd ..

echo "Importing database..."
python3 extract-epochs.py 3908
./import-db.sh

echo "Running testnet..."
mxpy testnet start
