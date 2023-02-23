import shutil
from argparse import ArgumentParser
from pathlib import Path

from archive_py.args import parse_epochs_arg
from archive_py.constants import HELP_STRING_DB, HELP_STRING_EPOCHS
from archive_py.ux import confirm_continuation


def main():
    parser = ArgumentParser()
    parser.add_argument("--folder", required=True, help=HELP_STRING_DB)
    parser.add_argument("--epochs", required=False, help=HELP_STRING_EPOCHS)
    args = parser.parse_args()

    folder = Path(args.folder).expanduser()
    epochs_to_remove = parse_epochs_arg(args.epochs)

    print("Folder:", folder)
    print("Epochs to remove:", epochs_to_remove)

    confirm_continuation()

    for epoch in epochs_to_remove:
        print(f"Removing epoch {epoch}...")
        epoch_folder = folder / f"Epoch_{epoch}"

        try:
            shutil.rmtree(str(epoch_folder))
        except FileNotFoundError:
            print(f"Epoch {epoch} does not exist.")


if __name__ == "__main__":
    main()
