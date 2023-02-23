import tarfile
from argparse import ArgumentParser
from pathlib import Path
from typing import Optional

from archive_py.args import parse_epochs_arg
from archive_py.constants import HELP_STRING_DB, HELP_STRING_EPOCHS
from archive_py.io import ensure_folder
from archive_py.ux import confirm_continuation


def main():
    parser = ArgumentParser()
    parser.add_argument("--input-folder", required=True, help=HELP_STRING_DB)
    parser.add_argument("--output-folder", required=True)
    parser.add_argument("--epochs", required=False, help=HELP_STRING_EPOCHS)
    parser.add_argument("--shard", required=True)
    parser.add_argument("--include-static", action="store_true", default=False)

    args = parser.parse_args()

    input_folder = Path(args.input_folder).expanduser()
    output_folder = Path(args.output_folder).expanduser()
    epochs = parse_epochs_arg(args.epochs)
    include_static = args.include_static
    shard = args.shard

    ensure_folder(output_folder)

    print("Input folder (Node DB):", input_folder)
    print("Output folder (destination of archives):", output_folder)
    print("Shard:", shard)
    print("Epochs to archive:", epochs)
    print("Include 'Static':", include_static)

    confirm_continuation()

    for epoch in epochs:
        archive_file = output_folder / f"Epoch_{epoch}.tar"
        relative_path = Path(f"Epoch_{epoch}") / f"Shard_{shard}"
        data_folder = input_folder / relative_path

        print("Archive:", archive_file)
        print("Data folder:", data_folder)

        tar = tarfile.open(archive_file, "w|")
        tar.add(data_folder, arcname=relative_path)
        tar.close()

    if include_static:
        print("Archiving folder 'Static' (with and without dblookup extensions)")

        archive_file = output_folder / "Static.tar"
        archive_file_min = output_folder / "Static.min.tar"
        data_folder = input_folder / "Static"

        print("Archive:", archive_file)
        print("Archive (min):", archive_file_min)
        print("Data folder:", data_folder)

        tar = tarfile.open(archive_file, "w|")
        tar.add(data_folder, arcname="Static")
        tar.close()

        def min_filter(info: tarfile.TarInfo) -> Optional[tarfile.TarInfo]:
            return None if "DbLookupExtensions" in info.name else info

        tar = tarfile.open(archive_file_min, "w|")
        tar.add(data_folder, arcname="Static", filter=min_filter)
        tar.close()


if __name__ == "__main__":
    main()
