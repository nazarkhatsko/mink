# 📡 Reporters

Control output format with the `--reporter` flag:

```bash
mink run mink.yaml --reporter <name>
```

📋 **Reporters**
| Reporter | Description |
|---|---|
| `pretty` (default) | Human-readable output with icons, descriptions and timing |
| `compact` | One line per action, no descriptions |
| `json` | NDJSON — one JSON object per line, suited to CI pipelines and log aggregators |
| `silent` | No output; exit code reflects success or failure |

## 🎨 pretty

<details>
<summary>Example output</summary>

```
▶ flow: user_lifecycle — Generate a user, create it, fetch, list, delete, and confirm 404
  · generate_user — Generate random user payload
  ✓ generate_user (0ms)
  · create_user — Create new user via API
  ✓ create_user (43ms)
  · validate_response — Validate response schema
  ✗ validate_response (1ms): validation failed: ...
✓ flow done: user_lifecycle
```

</details>

## 📎 compact

<details>
<summary>Example output</summary>

```
▶ user_lifecycle
  ✓ generate_user (0ms)
  ✓ create_user (43ms)
  ✗ validate_response (1ms): validation failed
✓ user_lifecycle
```

</details>

## 🧾 json

<details>
<summary>Example output</summary>

```json
{"type":"flow_start","flow_id":"user_lifecycle","description":"Generate a user, create it, fetch, list, delete, and confirm 404"}
{"type":"action_done","flow_id":"user_lifecycle","action_id":"generate_user","duration_ms":0,"output":{...},"error":null,"state":{}}
{"type":"action_done","flow_id":"user_lifecycle","action_id":"create_user","duration_ms":43,"output":{...},"error":null,"state":{"user_id":"0001"}}
{"type":"action_fail","flow_id":"user_lifecycle","action_id":"validate_response","duration_ms":1,"output":{"valid":false,"error":"..."},"error":{"code":"internal","message":"validation failed: ..."},"state":{"user_id":"0001"}}
{"type":"flow_done","flow_id":"user_lifecycle","state":{"user_id":"0001"}}
```

</details>

**Event types**
| `type` | Fields |
|---|---|
| `flow_start` | `flow_id`, `description` |
| `action_done` | `flow_id`, `action_id`, `duration_ms`, `output`, `error`, `state` |
| `action_fail` | `flow_id`, `action_id`, `duration_ms`, `output`, `error`, `state` |
| `flow_done` | `flow_id`, `state` |

`action_done` and `action_fail` share the same field set — `error` is always `null` on `action_done`, and always an object `{"code": "...", "message": "..."}` on `action_fail` (`code` is one of `config`, `transport`, `timeout`, or `internal` — see [manual/configuration](../internal/manual/embedded/configuration.md#mutate_on) for what each means). `output` on `action_fail` is whatever (possibly partial) output the driver returned before failing — `null` if it returned none.

`state` reflects the flow-level state **after** the action's `mutate_on` script has run, so each event shows the accumulated state at that point in time. On `action_fail` the state includes any mutations made by `mutate_on.fail` before the flow stopped.

## 🔇 silent

No output. Exit code reflects success or failure.

```bash
mink run mink.yaml --reporter silent && echo "ok"
```

⬅️ [Back to docs](README.md)
