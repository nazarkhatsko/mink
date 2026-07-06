# 🌐 httpbin

Runs against [httpbin.org](https://httpbin.org) — no local server needed.

▶️ **Run**
```bash
mink run examples/httpbin/mink.yaml
```

📋 **Flows**
| Flow | What it does |
|---|---|
| `create-and-verify-user` | Generates a fake user payload, POSTs it to httpbin, waits 200ms, validates the echoed response shape |

⬅️ [Back to examples](../README.md)
