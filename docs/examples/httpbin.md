# httpbin

Runs against [httpbin.org](https://httpbin.org) — no local server needed.

```bash
mink run examples/httpbin/mink.yaml
```

Flows:
- `create-and-verify-user` — generates a fake user payload, POSTs it to httpbin, waits 200ms, validates the echoed response shape
