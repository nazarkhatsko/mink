# Configuration

A mink config file is a YAML file with the following top-level structure:

```yaml
version: "1.0"
info: ...
vars: ...
instances: ...
flows: ...
```

A single file is conventionally named `mink.yaml`. If a project has several (e.g. one per suite), name them `<name>.mink.yaml` — `smoke.mink.yaml`, `regression.mink.yaml` — so the descriptive part sorts naturally in a file listing while `.mink.yaml` stays a greppable, tooling-recognizable suffix.

## `version`

Schema version. Currently `"1.0"`.

## `info`

Metadata about the file.

| Field | Type | Description |
|---|---|---|
| `name` | string | Name of the test suite |
| `description` | string | Optional description |

## Environment variables

Use `--env` to load a `.env` file before execution:

```bash
mink run mink.yaml --env .env
mink run mink.yaml --env .env.staging
```

`.env` format — `KEY=VALUE` per line, `#` comments and blank lines are ignored:

```
# .env
BASE_URL=https://api.example.com
API_KEY=secret-key
```

Variables loaded this way are available as `${env['KEY']}` in the config.

## `vars`

Named string values reusable across the config. Resolved before execution.

```yaml
vars:
  base_url: "https://api.example.com"
  admin_email: "admin@example.com"
  api_url: "${env['API_URL']}"
```

Reference with `${vars['name']}`.

## `instances`

Named driver instances used by actions. Each key is the instance name.

```yaml
instances:
  api:
    driver: http

  db:
    driver: postgres
    config:
      dsn: "${env['DATABASE_URL']}"
```

| Field | Type | Description |
|---|---|---|
| `driver` | string | Driver name (built-in or custom) |
| `config` | object | Driver-specific configuration |

`config` is optional for stateless drivers (`sleep`, `generate`, `validate`).

## `flows`

List of flows to execute sequentially.

```yaml
flows:
  - id: "create_user"
    description: "Create a user and validate the payload"
    actions:
      - id: generate_user
        description: "Generate random payload"
        instance: gen
        execute_with:
          schema:
            name:
              type: string
              value: faker.name
```

### Flow fields

| Field | Type | Description |
|---|---|---|
| `id` | string | Unique flow identifier, used with `--flow-id` flag; must be snake_case (`^[a-z][a-z0-9_]*$`) |
| `description` | string | Optional human-readable description |
| `actions` | list | Ordered list of actions |

### Action fields

| Field | Type | Description |
|---|---|---|
| `id` | string | Unique action identifier within the flow; must be snake_case (`^[a-z][a-z0-9_]*$`) |
| `description` | string | Human-readable description |
| `instance` | string | Instance name from `instances` |
| `timeout` | int | Timeout in milliseconds; cancels the action if exceeded (0 = no timeout) |
| `execute_with` | object | Driver-specific options, merged over instance `config` |
| `mutate_on.done` | string | Starlark script executed after successful action; has access to `event['output']` |
| `mutate_on.fail` | string | Starlark script executed on failure; has access to `event['output']` (partial, if any) and `event['error']`; flow still stops |

## Variable resolution

All `${}` expressions are evaluated as [Starlark](https://github.com/google/starlark-go) expressions before execution.

The following objects are available inside any `${}`:

| Object | Type | Description |
|---|---|---|
| `vars` | dict | Values from `vars` section |
| `env` | dict | Environment variables |
| `actions` | dict | Outputs of all previous actions, keyed by action id |
| `state` | dict | Mutable flow-level state, written via `mutate_on` |

```yaml
url: "${vars['base_url'] + '/users'}"
token: "${env['API_KEY']}"
user_id: "${actions['create_user']['resp']['body']['data']['id']}"
endpoint: "${vars['base_url'] + '/users/' + str(state['user_id'])}"
```

Since expressions are full Starlark, you can use built-in functions and list comprehensions:

```yaml
label: "${vars['env'].upper() + '-' + str(len(actions['list']['resp']['body']['data']))}"
```

## `mutate_on`

Starlark scripts that run after an action to update `state`. The `state` dict is shared across all actions in the flow and readable in any `${}` expression.

```yaml
- id: login
  instance: api
  execute_with:
    method: POST
    url: "${vars['base_url'] + '/auth'}"
  mutate_on:
    done: |
      body = event["output"]["resp"]["body"]
      state["token"]   = body["token"]
      state["user_id"] = body["data"]["id"]
    fail: |
      state["failed_action"] = event["action_id"]
      state["failed_code"]   = event["error"]["code"]
      state["failed_error"]  = event["error"]["message"]
```

### `event` object

`done` and `fail` receive the same shape — only `error` differs. Check `event['error'] == None` to tell them apart.

| Field | Available in | Description |
|---|---|---|
| `event['action_id']` | `done`, `fail` | ID of the current action |
| `event['output']` | `done`, `fail` | Driver output (same as `actions['id']`); on `fail` this is whatever partial output the driver returned, or `None` if it returned none |
| `event['error']` | `fail` | `{"code": "...", "message": "..."}`; `None` on `done` |

#### `event['error']['code']`

| Code | Set when |
|---|---|
| `config` | The driver rejected a missing or invalid option before doing anything |
| `transport` | An external call or process failed (network error, non-`ExitError` exec failure, RPC error) |
| `timeout` | The action's `timeout` was exceeded — takes priority over any other code a driver set |
| `internal` | Anything else: marshaling/parsing failures, and driver-specific assertion failures (e.g. `validate` schema mismatch) |

Branch on `code`, not `message` — `message` is free-form and driver-specific, `code` is not.

`mutate_on.fail` always runs before the flow stops — it cannot prevent termination.

### Numeric IDs from JSON

JSON numbers are always decoded as `float64` in Go. When you store an ID from a response and later use it in a URL or comparison, cast it explicitly with `int()` to avoid `"1.0"` in string interpolation:

```yaml
mutate_on:
  done: |
    state["order_id"] = int(event["output"]["resp"]["body"]["data"]["id"])
```

```yaml
url: "${vars['base_url'] + '/orders/' + str(state['order_id'])}"
```

Without the `int()` cast, `str(1.0)` produces `"1.0"` and the URL will not match.

## Merge rules

When an action uses an instance with `config`, the action's `execute_with` is **shallow merged** on top of `config`. Fields in `execute_with` override the same fields in `config`.

```yaml
instances:
  api:
    driver: http
    config:
      headers:
        X-Api-Key: "secret"    # always sent

flows:
  - id: "example"
    actions:
      - id: create
        instance: api
        execute_with:
          method: POST           # merged on top
          url: "https://..."
```
