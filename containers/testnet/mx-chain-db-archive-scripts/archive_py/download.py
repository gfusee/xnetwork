
from argparse import ArgumentParser
from pathlib import Path

from archive_py.args import parse_epochs_arg
from archive_py.constants import HELP_STRING_EPOCHS
from archive_py.downloader import download_file
from archive_py.io import ensure_folder
from archive_py.ux import confirm_continuation


def main():
    parser = ArgumentParser()
    parser.add_argument("--folder", required=True, help="Download folder")
    parser.add_argument("--url", required=True, help="E.g. https://example-bucket.fra1.digitaloceanspaces.com/foo/bar")
    parser.add_argument("--epochs", required=False, help=HELP_STRING_EPOCHS)
    parser.add_argument("--include-static", action="store_true", default=False)
    args = parser.parse_args()

    folder = Path(args.folder).expanduser()
    base_url = args.url
    epochs = parse_epochs_arg(args.epochs)
    include_static = args.include_static

    ensure_folder(folder)

    print("Folder:", folder)
    print("Url:", base_url)
    print("Epochs to download:", epochs)
    print("Include 'Static':", include_static)

    confirm_continuation()

    for epoch in epochs:
        filename = f"Epoch_{epoch}.tar"
        download_file(f"{base_url}/{filename}", str(folder / filename))

    if include_static:
        filename = "Static.tar"
        download_file(f"{base_url}/{filename}", str(folder / filename))


if __name__ == "__main__":
    main()
