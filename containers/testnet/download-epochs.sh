SHARD=$1
DL_OUTPUT_PATH=$2

cd mx-chain-db-archive-scripts
export URL="https://deep-history-archives.fra1.digitaloceanspaces.com/devnet/shard-${SHARD}"
PYTHONPATH=. python3 ./archive_py/download.py --folder="${DL_OUTPUT_PATH}/shard-${SHARD}" --url=${URL} --epochs=3904:3905 #--include-static
cd ..
