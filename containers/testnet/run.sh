echo "Waiting for Elastic Search..."
sudo bash ./wait-for-it.sh elastic:9200 --timeout=0

echo "Compiling nodes..."
rm -rf testnet
mxpy testnet config

echo "Copying python script..."
cp enable-elastic.py testnet/enable-elastic.py

echo "Enabling elastic search..."
cd testnet && python3 enable-elastic.py && cd ..

echo "Running testnet..."
mxpy testnet start
