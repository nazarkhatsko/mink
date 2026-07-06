# 🐍 python-step

Demonstrates the `python` driver — no local server needed.

▶️ **Run**
```bash
# Inline python code, structured output via json.dumps
mink run examples/python-step/mink.yaml --flow inline-code

# Running an existing .py script file with args
mink run examples/python-step/mink.yaml --flow script-file
```

📋 **Flows**
| Flow | What it does |
|---|---|
| `inline-code` | Runs inline Python code that prints JSON, validates the structured `stdout` |
| `script-file` | Runs `transform.py` with `args`, validates its JSON output |

⬅️ [Back to examples](../README.md)
