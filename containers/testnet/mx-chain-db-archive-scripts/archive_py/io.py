import os
from pathlib import Path
from typing import List


def ensure_folder(folder: Path):
    folder.mkdir(parents=True, exist_ok=True)


def list_archives(folder: Path) -> List[Path]:
    archives: List[Path] = [folder / file for file in os.listdir(folder)]
    archives = [file for file in archives if file.is_file() and file.suffix == ".tar"]
    archives.sort()
    return archives
