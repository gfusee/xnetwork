import os
import shutil
from argparse import ArgumentParser
from pathlib import Path

from archive_py.constants import HELP_STRING_DB
from archive_py.io import ensure_folder
from archive_py.ux import ask_number, confirm_continuation


def main():
    parser = ArgumentParser()
    parser.add_argument("--input-folder", required=True, help=HELP_STRING_DB)
    parser.add_argument("--output-folder", required=True, help=HELP_STRING_DB)
    args = parser.parse_args()

    input_folder = Path(args.input_folder).expanduser()
    output_folder = Path(args.output_folder).expanduser()

    ensure_folder(output_folder)

    print("Input folder:", input_folder)
    print("Output folder:", output_folder)

    confirm_continuation()

    first_epoch_to_move = find_first_movable_epoch(input_folder)
    print("First epoch to move:", first_epoch_to_move)
    confirm_continuation()

    num_epochs_to_move = ask_number(f"""How many epochs to move
\tfrom {input_folder}
\tto {output_folder}?""")
    last_epoch_to_move = first_epoch_to_move + num_epochs_to_move
    epochs_to_move = list(range(first_epoch_to_move, last_epoch_to_move))

    print("Epochs to move:", epochs_to_move)

    confirm_continuation()

    for epoch in epochs_to_move:
        print("Moving epoch", epoch)

        source = input_folder / f"Epoch_{epoch}"
        if not is_movable(source):
            print("Not movable. Will stop here.")
            break

        destination = output_folder / f"Epoch_{epoch}"
        shutil.move(str(source), destination)
        os.symlink(destination, source)


def find_first_movable_epoch(folder: Path) -> int:
    epoch = 0

    while True:
        epoch_folder = folder / f"Epoch_{epoch}"
        if is_movable(epoch_folder):
            return epoch

        epoch += 1


def is_movable(folder: Path) -> bool:
    return folder.is_dir() and not folder.is_symlink()


if __name__ == "__main__":
    main()
