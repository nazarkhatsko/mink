# Configuration

A mink config file is a YAML file with the following top-level structure:

```yaml
version: "1.0"
info: ...
vars: ...
instances: ...
flows: ...
```

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

Variables loaded this way are available as `${env.KEY}` in the config.

## `vars`

Named string values reusable across the config. Resolved before execution.

```yaml
vars:
  base_url: "https://api.example.com"
  admin_email: "admin@example.com"
  api_url: "${env.API_URL}"
```

Reference with `${vars.name}`.

## `instances`

Named driver instances used by actions. Each key is the instance name.

```yaml
instances:
  api:
    driver: http

  db:
    driver: postgres
    config:
      dsn: "${env.DATABASE_URL}"
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
  - name: "create-user"
    actions:
      - id: generate_user
        description: "Generate random payload"
        use: gen
        run_with:
          schema:
            name:
              type: string
              value: faker.name
```

### Flow fields

| Field | Type | Description |
|---|---|---|
| `name` | string | Unique flow name, used with `--flow` flag |
| `actions` | list | Ordered list of actions |

### Action fields

| Field | Type | Description |
|---|---|---|
| `id` | string | Unique action identifier within the flow |
| `description` | string | Human-readable description |
| `use` | string | Instance name from `instances` |
| `timeout` | int | Timeout in milliseconds; cancels the action if exceeded (0 = no timeout) |
| `run_with` | object | Driver-specific options, merged over instance `config` |

## Variable resolution

Expressions inside `${}` are resolved before execution.

| Syntax | Resolves to |
|---|---|
| `${vars.name}` | Value from `vars` section |
| `${env.NAME}` | `os.Getenv("NAME")` |
| `${actions.id.field}` | Field from a previous action's output |
| `${actions.id}` | Full output of a previous action |
| `${actions.id.arr[0]}` | Array index access |

Expressions can appear in any string value inside `run_with`.

## Merge rules

When an action uses an instance with `config`, the action's `run_with` is **shallow merged** on top of `config`. Fields in `run_with` override the same fields in `config`.

```yaml
instances:
  api:
    driver: http
    config:
      headers:
        X-Api-Key: "secret"    # always sent

flows:
  - name: "example"
    actions:
      - id: create
        use: api
        run_with:
          method: POST           # merged on top
          url: "https://..."
```
