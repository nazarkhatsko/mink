# 🐍 python-step

Demonstrates the `python` driver — no local server needed.

▶️ **Run**
```bash
# Inline python code, structured output via json.dumps
mink run examples/python-step/mink.yaml --flow-id inline_code

# Running an existing .py script file with args
mink run examples/python-step/mink.yaml --flow-id script_file
```

📋 **Flows**
| Flow | What it does |
|---|---|
| `inline_code` | Runs inline Python code that prints JSON, validates the structured `stdout` |
| `script_file` | Runs `transform.py` with `args`, validates its JSON output |

⬅️ [Back to examples](../README.md)
