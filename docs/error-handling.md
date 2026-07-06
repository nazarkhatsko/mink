# 🚨 Error Handling

How to react to a failed action — for the exact `event` field reference see [manual/configuration](../internal/manual/embedded/configuration.md#mutate_on).

## 🧭 Two outcomes, one shape

Every action either succeeds or fails. Either way, `mutate_on` sees the same `event` — `event['output']` is the driver's output (partial on failure, if the driver returned one), and `event['error']` is `None` on success, or `{"code": "...", "message": "..."}` on failure:

```yaml
- id: create_order
  use: api
  run_with:
    method: POST
    url: "${vars['base_url'] + '/orders'}"
  mutate_on:
    done: |
      state["order_id"] = event["output"]["resp"]["body"]["data"]["id"]
    fail: |
      state["last_failure"] = {
          "action":  event["action_id"],
          "code":    event["error"]["code"],
          "message": event["error"]["message"],
      }
```

Check `event['error'] == None` if a single script needs to handle both cases; usually you just write separate `done`/`fail` scripts instead.

## 🏷️ Branch on `code`, not `message`

`message` is for humans reading logs — it can change wording between mink versions. `event['error']['code']` won't:

| Code | What it means |
|---|---|
| `config` | The action's options were missing or invalid — a flow bug, not a runtime surprise |
| `transport` | The external call or process itself failed (connection refused, DNS failure, non-zero exec that isn't a normal exit) |
| `timeout` | The action's `timeout` was exceeded |
| `internal` | Anything else, including a failed `validate` assertion (see below) |

```yaml
mutate_on:
  fail: |
    state["retryable"] = (event["error"]["code"] == "timeout")
```

`mutate_on` scripts run as a Starlark module, not a function body, so `if`/`for` aren't allowed at the top level (only inside a `def`) — express branching as an expression like the ternary/comparison above instead.

`mutate_on.fail` runs before the flow stops — it can log or record state for whatever runs next (a CI step reading the final `state`, for instance), but it can't stop the failure from ending the flow.

## 🎯 Why assertion failures aren't a fifth code

`validate` failing to match a schema is exactly what it's supposed to do when a response is wrong — it's the test assertion, not an unexpected error. It reports through the same `config`/`transport`/`timeout`/`internal` set (falling into `internal`) rather than getting its own code, because the useful signal there isn't the code, it's `event['output']`: the driver still hands you `{valid: false, error: "..."}` so you can log or branch on *what* didn't match, not just that something didn't.

```yaml
mutate_on:
  fail: |
    state["schema_error"] = event["output"]["error"]
```

## 🧪 Try it

[`examples/error-handling`](../examples/error-handling/README.md) runs four flows, each failing one action on purpose — one per code — and prints the `event` fields a `mutate_on.fail` script would see.

⬅️ [Back to docs](README.md)
