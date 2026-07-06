import argparse
import json

parser = argparse.ArgumentParser()
parser.add_argument("--name", required=True)
args = parser.parse_args()

print(json.dumps({"greeting": f"Hello, {args.name}!"}))
