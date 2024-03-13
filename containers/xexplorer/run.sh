echo "Waiting for MySQL database..."
sudo ./wait-for-it.sh mysql:3306 --timeout=0

echo "Waiting for RabbitMQ..."
sudo ./wait-for-it.sh rabbitmq:5672 --timeout=0

echo "Waiting for Elastic Search..."
sudo ./wait-for-it.sh elastic:9200 --timeout=0

echo "Waiting for localnet..."
sudo ./wait-for-it.sh localnet:7950 --timeout=0

echo "Waiting for api..."
sudo ./wait-for-it.sh api:3001 --timeout=0

echo "Running xExplorer dApp..."
cd mx-explorer-dapp
VITE_APP_USE_HTTPS=false VITE_APP_SHARE_PREFIX=localnet- yarn run start
