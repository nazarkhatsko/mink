# 🚨 error-handling

Four flows, each deliberately failing one action, to show how to read `event['error']['code']` in `mutate_on.fail` — no local server needed.

▶️ **Run** (each flow fails by design, so run them one at a time)
```bash
mink run examples/error-handling/mink.yaml --flow config-error     --reporter json
mink run examples/error-handling/mink.yaml --flow transport-error  --reporter json
mink run examples/error-handling/mink.yaml --flow timeout-error    --reporter json
mink run examples/error-handling/mink.yaml --flow assertion-error  --reporter json
```

📋 **Flows**
| Flow | What fails | `event['error']['code']` |
|---|---|---|
| `config-error` | `http` action missing the required `url` option | `config` |
| `transport-error` | `http` request to `127.0.0.1:1` — nothing listens there, connection refused | `transport` |
| `timeout-error` | `http` request to the unroutable `10.255.255.1`, with a 200ms action `timeout` | `timeout` |
| `assertion-error` | `validate` action whose value doesn't match its schema | `internal` |

🧠 **Why `assertion-error` is `internal`, not a fifth code**

`validate` failing to match a schema isn't a broken flow — it's the assertion doing its job (see [flows-and-state.md](../../docs/flows-and-state.md)). It falls back to the generic `internal` code, same as `config`/`transport`/`timeout` would if nothing more specific applied. The useful bit isn't the code here, it's `event['output']`: `mutate_on.fail` gets the driver's actual output (`{valid: false, error: "..."}`) even though the action errored, so you can log or assert on *what* didn't match — not just that something didn't.

Every flow's `mutate_on.fail` script writes `event['error']['code']` and `event['error']['message']` into `state`, so `--reporter json`'s final `flow_done` event shows exactly what a script sees:

```json
{"type":"action_fail","flow":"transport-error","action_id":"connection_refused","duration_ms":3,"output":null,"error":{"code":"transport","message":"http: do request: ... connection refused"},"state":{"code":"transport","message":"..."}}
```

See [manual/configuration](../../internal/manual/embedded/configuration.md#mutate_on) for the full `event` reference and what each error code means.

⬅️ [Back to examples](../README.md)
