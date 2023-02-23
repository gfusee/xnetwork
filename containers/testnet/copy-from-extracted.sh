echo "Copying validator00's db folder..."
cp -R /home/ubuntu/extracted/validator00/db /home/ubuntu/testnet/validator00/db
echo "Copying validator00's import-db folder..."
cp -R /home/ubuntu/extracted/validator00/import-db /home/ubuntu/testnet/validator00/import-db
echo "Duplicating db into import-db for validator00..."
cp -R /home/ubuntu/testnet/validator00/db/. /home/ubuntu/testnet/validator00/import-db/db/

#echo "Copying validator01's db folder..."
#cp -R /home/ubuntu/extracted/validator01/db /home/ubuntu/testnet/validator01/db
#echo "Copying validator01's import-db folder..."
#cp -R /home/ubuntu/extracted/validator01/import-db /home/ubuntu/testnet/validator01/import-db
#
#echo "Copying validator02's db folder..."
#cp -R /home/ubuntu/extracted/validator02/db /home/ubuntu/testnet/validator02/db
#echo "Copying validator02's import-db folder..."
#cp -R /home/ubuntu/extracted/validator02/import-db /home/ubuntu/testnet/validator02/import-db
#
#echo "Copying validator03's db folder..."
#cp -R /home/ubuntu/extracted/validator03/db /home/ubuntu/testnet/validator03/db
#echo "Copying validator03's import-db folder..."
#cp -R /home/ubuntu/extracted/validator03/import-db /home/ubuntu/testnet/validator03/import-db

#echo "Duplicating Static for validator01..."
#cp -R /home/ubuntu/testnet/validator01/import-db/db/D/Static /home/ubuntu/testnet/validator01/db/D/Static
#
#echo "Duplicating Static for validator02..."
#cp -R /home/ubuntu/testnet/validator02/import-db/db/D/Static /home/ubuntu/testnet/validator02/db/D/Static
#
#echo "Duplicating Static for validator03..."
#cp -R /home/ubuntu/testnet/validator03/import-db/db/D/Static /home/ubuntu/testnet/validator03/db/D/Static
