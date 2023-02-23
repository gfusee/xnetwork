def confirm_continuation(yes: bool = True):
    if (yes):
        return

    answer = input("Continue? (y/n)")
    if answer.lower() not in ["y", "yes"]:
        print("Confirmation not given. Will stop.")
        exit(1)


def ask_number(message: str):
    answer = input(f"{message}\n")
    return int(answer)
