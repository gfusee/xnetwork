# mx-chain-db-archive-scripts

Scripts useful for managing incremental archives of Node's database.

Features:
 - archive / extract
 - upload / download

## Prerequisites

```
pip3 install -r ./requirements.txt
```

## Archive `Epoch_*` and `Static`

```
PYTHONPATH=. python3 ./archive_py/archive.py --input-folder=~/db/T --output-folder=~/db-archived --epochs=0:10 --shard=0 --include-static
```

## Upload archives to DO Spaces or Amazon S3

```
export S3_ENDPOINT=https://fra1.digitaloceanspaces.com
export S3_REGION=fra1
export S3_ACCESS_KEY=DO000000000000000000
export S3_BUCKET=example-bucket

PYTHONPATH=. python3 ./archive_py/upload.py --folder=~/db-archived --endpoint=${S3_ENDPOINT} --region=${S3_REGION} --access-key=${S3_ACCESS_KEY} --bucket=${S3_BUCKET} --prefix=foo/bar
```

You will be prompted to enter the S3 secrey key.

Files will be uploaded with `ACL == public-read`.

## Download archives

```
export URL=https://example-bucket.fra1.digitaloceanspaces.com
PYTHONPATH=. python3 ./archive_py/download.py --folder=~/db-downloads --url=${URL}/foo/bar --epochs=0:10 --include-static
```

## Extract archives

```
PYTHONPATH=. python3 ./archive_py/extract.py --input-folder=~/db-downloads --output-folder=~/db-extracted
```

## Move & symlink epoch folders

```
PYTHONPATH=. python3 ./archive_py/move_link_epochs.py --input-folder=~/db/T --output-folder=~/place/db/T
```

## Remove old epochs

```
PYTHONPATH=. python3 ./archive_py/remove_epochs.py --folder=~/db/T --epochs=0:1000
```
