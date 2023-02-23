import os
import subprocess
import sys

current_dir = os.path.dirname(os.path.realpath(__file__))

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 extract-epochs.py <last epoch downloaded>", flush=True)
        return

    last_epoch_dl = int(sys.argv[1])
    print("Last epoch downloaded: ", last_epoch_dl, flush=True)

    epochs_for_import_db = range(last_epoch_dl - 2, last_epoch_dl + 1)

    db_downloads_path = "db-downloads"
    if not os.path.exists(db_downloads_path):
        print("db-downloads folder does not exist", flush=True)
        return

    os.chdir(db_downloads_path)

    # Create the folders
    for i in range(4):
        validator_folder = current_dir + "/extracted" + "/validator" + str(i).zfill(2)
        os.makedirs(validator_folder + "/db/D", exist_ok=True)
        os.makedirs(validator_folder + "/import-db/db/D", exist_ok=True)

    # Extract the shards
    validator_counter = 0
    for folder in os.listdir():
        validator_name = "validator" + str(validator_counter).zfill(2)
        validator_folder = current_dir + "/extracted" + "/" + validator_name
        if folder.startswith("shard-"):
            print("Extracting " + folder + " inside " + validator_folder + " ...", flush=True)
            os.chdir(folder)
            for f in os.listdir():
                if f.endswith(".tar"):
                    # Check whether the file is named Epoch_XXX and if XXX is in the range of epochs for import-db
                    final_output_folder = ""
                    import_db_folder = validator_folder + "/import-db/db/D"
                    db_folder = validator_folder + "/db/D"
                    if f.startswith("Epoch_") and int(f[6:-4]) in epochs_for_import_db:
                        print("--- Extracting inside import-db : " + f + "...", flush=True)
                        final_output_folder = import_db_folder
                    elif f.startswith("Epoch_"):
                        print("--- Extracting inside db : " + f + "...", flush=True)
                        final_output_folder = db_folder
                    # Extract both in db and import-db if the file is named Static
                    elif f.startswith("Static"):
                        print("--- Extracting inside import-db and db : " + f + "...", flush=True)
                        final_output_folder = db_folder

                    subprocess.run(["tar", "-xf", f])
                    f_name = f[:-4]
                    subprocess.run(["cp", "-R", f_name, final_output_folder + "/" + f_name])
                    subprocess.run(["rm", "-rf", f_name])

            os.chdir("..")
            validator_counter += 1

    print("Extract finished", flush=True)


if __name__ == "__main__":
    main()
