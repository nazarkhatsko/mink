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
▶ flow: user-lifecycle
  · generate_user — Generate random user payload
  ✓ generate_user (0ms)
  · create_user — Create new user via API
  ✓ create_user (43ms)
  · validate_response — Validate response schema
  ✗ validate_response (1ms): validation failed: ...
✓ flow done: user-lifecycle
```

</details>

## 📎 compact

<details>
<summary>Example output</summary>

```
▶ user-lifecycle
  ✓ generate_user (0ms)
  ✓ create_user (43ms)
  ✗ validate_response (1ms): validation failed
✓ user-lifecycle
```

</details>

## 🧾 json

<details>
<summary>Example output</summary>

```json
{"type":"flow_start","flow":"user-lifecycle"}
{"type":"action_done","flow":"user-lifecycle","action_id":"generate_user","duration_ms":0,"output":{...},"error":null,"state":{}}
{"type":"action_done","flow":"user-lifecycle","action_id":"create_user","duration_ms":43,"output":{...},"error":null,"state":{"user_id":"0001"}}
{"type":"action_fail","flow":"user-lifecycle","action_id":"validate_response","duration_ms":1,"output":{"valid":false,"error":"..."},"error":{"code":"internal","message":"validation failed: ..."},"state":{"user_id":"0001"}}
{"type":"flow_done","flow":"user-lifecycle","state":{"user_id":"0001"}}
```

</details>

**Event types**
| `type` | Fields |
|---|---|
| `flow_start` | `flow` |
| `action_done` | `flow`, `action_id`, `duration_ms`, `output`, `error`, `state` |
| `action_fail` | `flow`, `action_id`, `duration_ms`, `output`, `error`, `state` |
| `flow_done` | `flow`, `state` |

`action_done` and `action_fail` share the same field set — `error` is always `null` on `action_done`, and always an object `{"code": "...", "message": "..."}` on `action_fail` (`code` is one of `config`, `transport`, `timeout`, or `internal` — see [manual/configuration](../internal/manual/embedded/configuration.md#mutate_on) for what each means). `output` on `action_fail` is whatever (possibly partial) output the driver returned before failing — `null` if it returned none.

`state` reflects the flow-level state **after** the action's `mutate_on` script has run, so each event shows the accumulated state at that point in time. On `action_fail` the state includes any mutations made by `mutate_on.fail` before the flow stopped.

## 🔇 silent

No output. Exit code reflects success or failure.

```bash
mink run mink.yaml --reporter silent && echo "ok"
```

⬅️ [Back to docs](README.md)
