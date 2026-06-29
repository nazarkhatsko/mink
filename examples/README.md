# Examples


## httpbin

Runs against [httpbin.org](https://httpbin.org) — no local server needed.

```bash
mink run examples/httpbin/mink.yaml
```

Flows:
- `create-and-verify-user` — generates a fake user payload, POSTs it to httpbin, waits 200ms, validates the echoed response shape


---


## simple-api

A local Go HTTP server with CRUD `/users` endpoints and `X-Api-Key` auth. Demonstrates a full end-to-end user lifecycle.

**Start the server:**
```bash
go run examples/simple-api/main.go
```

**Run the flows:**
```bash
# Check that unauthenticated requests return 401
mink run examples/simple-api/mink.yaml --flow unauthorized-request

# Full create → get → list → delete lifecycle
mink run examples/simple-api/mink.yaml --flow user-lifecycle
```

Flows:
- `unauthorized-request` — sends a request without an API key, validates 401 response
- `user-lifecycle` — generates a random user, creates it, fetches, lists, deletes, and confirms 404
