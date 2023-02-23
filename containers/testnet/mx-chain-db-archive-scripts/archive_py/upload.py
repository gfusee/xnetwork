import time
from argparse import ArgumentParser
from getpass import getpass
from pathlib import Path
from typing import Any, List

from boto3.s3.transfer import TransferConfig
from boto3.session import Session

from archive_py.args import parse_epochs_arg
from archive_py.constants import HELP_STRING_EPOCHS, ONE_GB, ONE_MB
from archive_py.io import list_archives
from archive_py.ux import confirm_continuation

MULTI_PART_THRESHOLD = 4 * ONE_GB


def main():
    parser = ArgumentParser()
    parser.add_argument("--folder", required=True, help="Folder with archives")
    parser.add_argument("--endpoint", required=True, default="https://fra1.digitaloceanspaces.com")
    parser.add_argument("--region", required=True, default="fra1")
    parser.add_argument("--access-key", required=True, help="The access key. This is NOT the secret key.")
    parser.add_argument("--bucket", required=True, help="For Digital Ocean, it's the space name.")
    parser.add_argument("--prefix", required=True, help="E.g. foo/bar")
    parser.add_argument("--epochs", required=False, help=HELP_STRING_EPOCHS)
    parser.add_argument("--include-static", action="store_true", default=False)
    args = parser.parse_args()

    folder = Path(args.folder).expanduser()
    endpoint = args.endpoint
    region = args.region
    access_key = args.access_key
    bucket = args.bucket
    prefix = args.prefix
    with_specified_archives = args.epochs or args.include_static
    specified_epochs = parse_epochs_arg(args.epochs)
    specified_static = args.include_static
    print("Specified epochs to upload:", specified_epochs)
    print("Specified to include 'Static':", specified_static)

    if with_specified_archives:
        archives: List[Path] = [folder / f"Epoch_{epoch}.tar" for epoch in specified_epochs]
        archives.extend([folder / "Static.tar"] if specified_static else [])
    else:
        archives: List[Path] = list_archives(folder)

    print("Folder:", folder)
    print("Endpoint:", endpoint)
    print("Region:", region)
    print("Access key:", access_key)
    print("Bucket:", bucket)
    print("Prefix:", prefix)
    print(f"Will upload {len(archives)} files.")

    confirm_continuation()

    secret_key = getpass("Enter S3 secrey key:\n")
    print(f"Entered secret key of length {len(secret_key)}")

    confirm_continuation()

    # See: https://docs.digitalocean.com/reference/api/spaces-api/
    # The new session validates your request and directs it to your Space's specified endpoint using the AWS SDK.
    session: Any = Session()
    client = session.client("s3",
                            endpoint_url=endpoint,
                            region_name=region,
                            aws_access_key_id=access_key,
                            aws_secret_access_key=secret_key)

    # See https://boto3.amazonaws.com/v1/documentation/api/latest/reference/customizations/s3.html#boto3.s3.transfer.TransferConfig
    transfer_config = TransferConfig(multipart_threshold=MULTI_PART_THRESHOLD)

    index = 0
    for file in archives:
        index += 1
        print(f"({index} / {len(archives)}) Uploading {file}")
        start = time.time()

        file_path = folder / file
        file_size = file_path.stat().st_size
        print("Size:", int(file_size / ONE_MB), "MB")

        # See:
        # https://boto3.amazonaws.com/v1/documentation/api/latest/reference/services/s3.html#S3.Client.upload_file
        # https://github.com/boto/s3transfer/blob/develop/s3transfer/manager.py
        client.upload_file(
            Filename=str(file_path),
            Bucket=bucket,
            ExtraArgs={"ACL": "public-read"},
            Key=f"{prefix}/{file.name}",
            Config=transfer_config
        )

        end = time.time()
        print("Took", end - start, "seconds")


if __name__ == "__main__":
    main()
