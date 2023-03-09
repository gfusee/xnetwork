import json
import argparse
import os

# Parse command line arguments
parser = argparse.ArgumentParser()
parser.add_argument("key", help="the key to add to the JSON object")
parser.add_argument("value", help="the value to add to the JSON object")
args = parser.parse_args()

file_path = '/home/ubuntu/results.json'

# Create the JSON object
data = {}
if os.path.exists(file_path):
    with open(file_path, 'r') as f:
        data = json.load(f)

# Add the new key-value pair
data[args.key] = args.value

# Write the updated JSON object to file
with open(file_path, 'w') as f:
    json.dump(data, f)
