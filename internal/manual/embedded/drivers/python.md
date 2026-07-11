# python

Executes Python code as a flow step — either inline via `run_code` or an existing `.py` file via `run_script`. Useful for data transformation, signing/hashing, or any logic that's awkward to express as a single Starlark expression.

## Config

| Field | Type | Required | Description |
|---|---|---|---|
| `interpreter` | string | no | Python executable to run. Default: `python3` (point at a venv, e.g. `.venv/bin/python`) |
| `env` | object | no | Base environment variables merged with the current process environment |
| `dir` | string | no | Default working directory |

## Methods

### `run_code`

| Field | Type | Required | Description |
|---|---|---|---|
| `code` | string | yes | Inline Python source (use YAML `\|` for multiline scripts) |
| `args` | []any | no | Arguments passed to the script (`sys.argv[1:]`) |
| `env` | object | no | Extra environment variables, merged over the instance-level `env` |
| `dir` | string | no | Overrides the working directory for this action |

### `run_script`

| Field | Type | Required | Description |
|---|---|---|---|
| `script` | string | yes | Path to an existing `.py` file |
| `args` | []any | no | Arguments passed to the script (`sys.argv[1:]`) |
| `env` | object | no | Extra environment variables, merged over the instance-level `env` |
| `dir` | string | no | Overrides the working directory for this action |

> `timeout` is a top-level action field, not a driver option — see "Action fields" in `mink manual configuration`.

## Output

```json
{
  "exit_code": 0,
  "stdout": { "sum": 108, "count": 6 },
  "stderr": "",
  "success": true
}
```

| Field | Type | Description |
|---|---|---|
| `exit_code` | int | Exit code of the process |
| `stdout` | any | Standard output — parsed as JSON if valid, otherwise the raw string |
| `stderr` | string | Standard error |
| `success` | bool | `true` if `exit_code == 0` |

Print `json.dumps(...)` from the script to get a structured `stdout` usable in later actions (e.g. `actions['step']['stdout']['field']`). A non-zero exit code does **not** fail the action — it returns `success: false`. Use the `validate` driver to assert on the result.

## Example

```yaml
instances:
  py:
    driver: python
    methods: [run_code, run_script]

  check:
    driver: validate
    methods: [schema]

flows:
  - id: signature_check
    actions:
      - id: compute_signature
        description: "Compute an HMAC signature in Python and return it as structured data"
        instance: py
        method: run_code
        execute_with:
          code: |
            import hashlib, hmac, json, os

            secret = os.environ["SIGNING_SECRET"].encode()
            payload = "order-42"
            signature = hmac.new(secret, payload.encode(), hashlib.sha256).hexdigest()

            print(json.dumps({"payload": payload, "signature": signature}))
          env:
            SIGNING_SECRET: "${vars['signing_secret']}"

      - id: validate_signature
        description: "Assert the script produced a signature"
        instance: check
        method: schema
        execute_with:
          value: "${actions['compute_signature']}"
          schema:
            type: object
            properties:
              success:
                type: boolean
                const: true
              stdout:
                type: object
                required: [signature]
```
