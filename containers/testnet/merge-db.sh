rm -rf /home/ubuntu/testnet/validator00/db/D/Static
mkdir -p /home/ubuntu/testnet/validator00/db/D/Static
mx-chain-tools-go/dbMerger/cmd/generalDBMerger/generalDBMerger -dest=/home/ubuntu/testnet/validator00/db/D/Static -sources=/home/ubuntu/testnet/validator00/db/D/Epoch_3903,/home/ubuntu/testnet/validator00/db/D/Epoch_3904,/home/ubuntu/testnet/validator00/db/D/Epoch_3905
