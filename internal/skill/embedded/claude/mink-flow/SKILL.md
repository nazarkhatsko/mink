Generate a mink flow YAML file based on the user's description.

## Check current docs first

This file is a static snapshot and may drift from the installed `mink` version. Before generating a flow, prefer running the CLI's embedded docs as the source of truth:

```bash
mink manual getting-started      # first flow walkthrough
mink manual configuration        # full YAML spec: vars, instances, flows, state, mutate_on
mink manual drivers               # list all available drivers
mink manual drivers <name>        # config/methods/output for one driver
```

If anything below conflicts with `mink manual`, trust `mink manual`.

## Rules

- Always start with `version: "1.0"` and `info:` block
- Define all required drivers in `instances:` before using them in `flows:`
- Every instance needs `driver:` and a required, non-empty `methods:` list — the subset of that driver's methods this instance may call (see "Available drivers" below)
- Every flow must have `id:` (snake_case, same format as action `id`) and may include an optional `description:`
- Every action must have `id:`, `description:`, `instance:`, `method:`, and `execute_with:`
- `execute_with` is strictly validated per method — only the keys listed for that method in `mink manual drivers <name>` are allowed, and required ones must be present (in `execute_with` or the instance's `config`)
- Add `timeout:` (milliseconds) on an action when it may hang (shell commands, slow endpoints)
- Use `${vars['key']}` for reusable values, `${env['KEY']}` for secrets
- Reference previous action outputs via `${actions['id']['field']}`
- Use the `state` dict (written via `mutate_on.done`/`mutate_on.fail`) to carry values across actions
- Use a `gen` instance with `driver: generate`, `methods: [object]` for generating fake data
- Use a `check` instance with `driver: validate`, `methods: [schema]` for JSON Schema validation
- Use a `time` instance with `driver: time`, `methods: [wait]` when an async delay is needed
- Use an `sh` instance with `driver: shell`, `methods: [run_code, run_script]` for setup/teardown or system-level steps — print status via exit code, not stdout parsing (shell's `stdout` is always a raw string, unlike `python`'s)
- Use a `py` instance with `driver: python`, `methods: [run_code, run_script]` for data transformation, hashing/signing, or other logic awkward as a single Starlark expression — print `json.dumps(...)` to get a structured `stdout`
- Use an `llm` instance with `driver: claude`, `methods: [message]` to call the Claude API from a flow
- `http` actions pick their verb via `method:` (`get`/`post`/`put`/`patch`/`delete`) — `execute_with` for these never contains a `method` key, only `url:` (required), `headers:`, `body:`. Only the `request` method takes `method:` as an `execute_with` option, for verbs without a dedicated one (`HEAD`, `OPTIONS`)
- All `${}` expressions are Starlark — use dict access `['key']`, not dot notation
- JSON numbers decode as `float64`; cast IDs with `int(...)` in `mutate_on` before interpolating them into a URL or string, or you'll get `"1.0"` instead of `"1"`
- Name the file `mink.yaml` for a single-suite project; for multiple suites use `<name>.mink.yaml` (e.g. `smoke.mink.yaml`)

## Available drivers

| driver | methods | purpose |
|---|---|---|
| `http` | `get`, `post`, `put`, `patch`, `delete`, `request` | HTTP requests |
| `generate` | `object` | fake data generation |
| `validate` | `schema` | JSON Schema validation |
| `time` | `wait` | delay execution |
| `shell` | `run_code`, `run_script` | shell command execution |
| `python` | `run_code`, `run_script` | Python code/script execution |
| `claude` | `message` | Claude API messages |

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

**time:** `{ slept_ms: int }`

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
    methods: [post]
  gen:
    driver: generate
    methods: [object]
  check:
    driver: validate
    methods: [schema]

flows:
  - id: create_user
    description: "Create a user and validate the response"
    actions:
      - id: generate_user
        description: "Generate random user payload"
        instance: gen
        method: object
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
        method: post
        execute_with:
          url: "${vars['base_url'] + '/users'}"
          body: "${actions['generate_user']}"
        mutate_on:
          done: |
            state["user_id"] = int(event["output"]["resp"]["body"]["id"])

      - id: validate_response
        description: "Validate response contains id"
        instance: check
        method: schema
        execute_with:
          value: "${actions['create_user']['resp']['body']}"
          schema:
            type: object
            required: [id]
```
