from typing import List


def parse_epochs_arg(epochs_arg: str) -> List[int]:
    if not epochs_arg:
        return []

    # Handle specific epochs
    try:
        parts = epochs_arg.split(",")
        return [int(part) for part in parts]
    except Exception:
        pass

    # Handle ranges. E.g. 7:9.
    try:
        parts = epochs_arg.split(":")
        return list(range(int(parts[0]), int(parts[1]) + 1))
    except Exception:
        pass

    raise Exception(f"Cannot parse epochs: {epochs_arg}")
