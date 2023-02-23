SHARD=$1
DL_OUTPUT_PATH=$2
START_EPOCH=$3
END_EPOCH=$4

cd mx-chain-db-archive-scripts
export URL="https://deep-history-archives.fra1.digitaloceanspaces.com/devnet/shard-${SHARD}"
PYTHONPATH=. python3 ./archive_py/download.py --folder="${DL_OUTPUT_PATH}/shard-${SHARD}" --url=${URL} --epochs=${START_EPOCH}:${END_EPOCH} --include-static
cd ..
