# 🧠 Flows & State

The mental model behind actions, `state`, and `${}` expressions — for full field-by-field syntax see [manual/configuration](../internal/manual/embedded/configuration.md).

## 🔗 How a flow runs

A flow is an **ordered list of actions**. Each action runs through a driver instance, and its result becomes available to every action after it:

```
action 1 ──▶ action 2 ──▶ action 3 ──▶ ...
   │             │
   ▼             ▼
 result       result + state
 (actions)    (actions + state)
```

- 📥 **`actions['id']`** — the raw output of any earlier action, keyed by its `id`
- 🧳 **`state`** — a shared dict that travels across the whole flow, written explicitly via `mutate_on`

The difference: `actions` is automatic and read-only; `state` is opt-in and mutable — you decide what's worth carrying forward.

## ✍️ Writing to state

`mutate_on` runs a [Starlark](https://github.com/google/starlark-go) script right after an action finishes:

```yaml
- id: login
  use: api
  run_with:
    method: POST
    url: "${vars['base_url'] + '/auth'}"
  mutate_on:
    done: |
      state["token"] = event["output"]["resp"]["body"]["token"]
    fail: |
      state["failed_action"] = event["action_id"]
      state["failed_code"]   = event["error"]["code"]
```

`done` runs on success, `fail` runs on failure (and can't stop the flow from halting). Both scripts see the same `event` shape — `event['output']` is the driver's output either way (partial on failure, if the driver returned one), and `event['error']` is only non-null on failure (check `event['error'] == None` to tell them apart). See [Error Handling](error-handling.md) for how to branch on `event['error']['code']`, or [manual/configuration](../internal/manual/embedded/configuration.md#mutate_on) for the full `event` field reference. Anything you write to `state` here is readable by every later action.

## 📖 Reading it back

Any `${}` expression can reach into four objects:

| | Object | Scope |
|---|---|---|
| ⚙️ | `vars` | Config-wide constants from `vars:` |
| 🌱 | `env` | Environment variables (`--env` file or shell) |
| 📥 | `actions` | Outputs of every earlier action, keyed by `id` |
| 🧳 | `state` | Whatever `mutate_on` has written so far |

```yaml
url: "${vars['base_url'] + '/orders/' + str(state['order_id'])}"
```

## 🧩 Why this shape

Keeping `actions` automatic and `state` explicit means a flow's data flow is traceable just by reading it top to bottom: if a later step depends on something, it's either the direct output of a named action, or a value someone deliberately promoted into `state`. No hidden globals.

⬅️ [Back to docs](README.md)
