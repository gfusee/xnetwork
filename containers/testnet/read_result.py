import json
import os
import sys

file_path = '/home/ubuntu/results.json'


def get_result(key):
    # Create the JSON object
    data = {}
    if os.path.exists(file_path):
        with open(file_path, 'r') as f:
            data = json.load(f)

    return data[key]


if __name__ == '__main__':
    value = sys.argv[1]
    print(get_result(value))
