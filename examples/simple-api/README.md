# 🔑 simple-api

A local Go HTTP server with CRUD `/users` endpoints and `X-Api-Key` auth. Demonstrates a full end-to-end user lifecycle.

🚀 **Start the server**
```bash
go run examples/simple-api/main.go
```

▶️ **Run the flows**
```bash
# Check that unauthenticated requests return 401
mink run examples/simple-api/mink.yaml --flow-id unauthorized_request

# Full create → get → list → delete lifecycle
mink run examples/simple-api/mink.yaml --flow-id user_lifecycle
```

📋 **Flows**
| Flow | What it does |
|---|---|
| `unauthorized_request` | Sends a request without an API key, validates 401 response |
| `user_lifecycle` | Generates a random user, creates it, fetches, lists, deletes, and confirms 404 |

⬅️ [Back to examples](../README.md)
