echo "Waiting for MySQL database..."
sudo ./wait-for-it.sh mysql:3306 --timeout=0

echo "Waiting for RabbitMQ..."
sudo ./wait-for-it.sh rabbitmq:5672 --timeout=0

echo "Waiting for Elastic Search..."
sudo ./wait-for-it.sh elastic:9200 --timeout=0

echo "Waiting for testnet..."
sudo ./wait-for-it.sh testnet:7950 --timeout=0

echo "Waiting for VM queries to work, this can take some minutes..."
sudo python3 ./wait-for-vm-query-ready.py

echo "Running api service..."
cd mx-api-service
cp ./config/config.devnet.yaml ./config/config.yaml
nest build
nest start
