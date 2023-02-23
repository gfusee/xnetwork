import tarfile
from argparse import ArgumentParser
from pathlib import Path

from archive_py.constants import HELP_STRING_DB
from archive_py.io import ensure_folder, list_archives
from archive_py.ux import confirm_continuation


def main():
    parser = ArgumentParser()
    parser.add_argument("--input-folder", required=True)
    parser.add_argument("--output-folder", required=True, help=HELP_STRING_DB)
    args = parser.parse_args()

    input_folder = Path(args.input_folder).expanduser()
    output_folder = Path(args.output_folder).expanduser()

    ensure_folder(output_folder)

    print("Input folder (with archives):", input_folder)
    print("Output folder (extraction destination):", output_folder)

    confirm_continuation()

    archives = list_archives(input_folder)

    for archive_path in archives:
        print("Archive:", archive_path)

        tar = tarfile.open(archive_path, "r")
        tar.extractall(output_folder)
        tar.close()


if __name__ == "__main__":
    main()
