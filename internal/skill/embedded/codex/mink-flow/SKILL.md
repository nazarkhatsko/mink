Generate a mink flow YAML file based on the user's description.

## Check current docs first

This file is a static snapshot and may drift from the installed `mink` version. Before generating a flow, prefer running the CLI's embedded docs as the source of truth:

```bash
mink manual getting-started      # first flow walkthrough
mink manual configuration        # full YAML spec: vars, instances, flows, state, mutate_on
mink manual drivers               # list all available drivers
mink manual drivers <name>        # config/options/output for one driver
```

If anything below conflicts with `mink manual`, trust `mink manual`.

## Rules

- Always start with `version: "1.0"` and `info:` block
- Define all required drivers in `instances:` before using them in `flows:`
- Every flow must have `id:` (snake_case, same format as action `id`) and may include an optional `description:`
- Every action must have `id:`, `description:`, `instance:`, and `execute_with:`
- Add `timeout:` (milliseconds) on an action when it may hang (shell commands, slow endpoints)
- Use `${vars['key']}` for reusable values, `${env['KEY']}` for secrets
- Reference previous action outputs via `${actions['id']['field']}`
- Use the `state` dict (written via `mutate_on.done`/`mutate_on.fail`) to carry values across actions
- Use `gen` instance with `driver: generate` for generating fake data
- Use `check` instance with `driver: validate` for JSON Schema validation
- Use `sleep` instance with `driver: sleep` when an async delay is needed
- Use `sh` instance with `driver: shell` for setup/teardown or system-level steps
- Use `py` instance with `driver: python` for data transformation, hashing/signing, or other logic awkward as a single Starlark expression — print `json.dumps(...)` to get a structured `stdout`
- Use `llm` instance with `driver: claude` to call the Claude API from a flow
- HTTP actions must always include `method:` and `url:` in `execute_with:`
- All `${}` expressions are Starlark — use dict access `['key']`, not dot notation
- JSON numbers decode as `float64`; cast IDs with `int(...)` in `mutate_on` before interpolating them into a URL or string, or you'll get `"1.0"` instead of `"1"`
- Name the file `mink.yaml` for a single-suite project; for multiple suites use `<name>.mink.yaml` (e.g. `smoke.mink.yaml`)

## Available drivers

| driver | purpose |
|---|---|
| `http` | HTTP requests |
| `generate` | fake data generation |
| `validate` | JSON Schema validation |
| `sleep` | delay execution |
| `shell` | shell command execution |
| `python` | Python code/script execution |
| `claude` | Claude API messages |

## Generate field types

| type | value examples |
|---|---|
| `string` | `faker.name`, `faker.email`, `faker.uuid`, `faker.password`, `faker.phone` |
| `int` | `faker.age` |

## Output structure per driver

**http:**
```
req:
  method: string
  url: string
  headers: object
  body: any
resp:
  status: int
  headers: object
  body: any
```

**generate:** returns the generated object directly

**validate:** `{ valid: true }` on success; on failure the action errors and the flow stops — `actions['id']` is never populated, but `mutate_on.fail`'s `event['output']` still gets `{ valid: false, error: "..." }`

**sleep:** `{ slept_ms: int }`

**shell:** `{ exit_code: int, stdout: string, stderr: string, success: bool }` — a non-zero exit code does **not** fail the action, assert on `success` with `validate`

**python:** `{ exit_code: int, stdout: any, stderr: string, success: bool }` — `stdout` is parsed as JSON if valid, otherwise the raw string; a non-zero exit code does **not** fail the action, assert on `success` with `validate`

**claude:**
```
in:
  model: string
  messages: []object
  system: string
  max_tokens: int
out:
  id: string
  model: string
  role: string
  content: string       # concatenated text from response content blocks
  stop_reason: string
  usage: { input_tokens: int, output_tokens: int }
```

## Example

User: "Generate a user, post it to /users, validate the response has an id"

```yaml
version: "1.0"
info:
  name: "Create user flow"
  description: "Generate and post a user, validate response"

vars:
  base_url: "${env['BASE_URL']}"

instances:
  api:
    driver: http
  gen:
    driver: generate
  check:
    driver: validate

flows:
  - id: "create_user"
    description: "Create a user and validate the response"
    actions:
      - id: generate_user
        description: "Generate random user payload"
        instance: gen
        execute_with:
          schema:
            name:
              type: string
              value: faker.name
            email:
              type: string
              value: faker.email

      - id: create_user
        description: "POST user to API"
        instance: api
        execute_with:
          method: POST
          url: "${vars['base_url'] + '/users'}"
          body: "${actions['generate_user']}"
        mutate_on:
          done: |
            state["user_id"] = int(event["output"]["resp"]["body"]["id"])

      - id: validate_response
        description: "Validate response contains id"
        instance: check
        execute_with:
          value: "${actions['create_user']['resp']['body']}"
          schema:
            type: object
            required: [id]
```
