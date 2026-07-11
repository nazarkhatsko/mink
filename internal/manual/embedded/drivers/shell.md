# shell

Executes a shell command as a flow step — either inline `code` or an existing script file via `script`, both run through `sh`. Useful for database seeding, service health checks, file setup, or any system-level operation in a flow.

## Methods

### `run_code`

| Field | Type | Required | Description |
|---|---|---|---|
| `code` | string | yes | Inline shell source (use YAML `\|` for multiline scripts) |
| `args` | []any | no | Arguments passed to the script (`$1`, `$2`, ...) |
| `env` | object | no | Extra environment variables merged with the current process environment |
| `dir` | string | no | Working directory for the command (default: current process directory) |

### `run_script`

| Field | Type | Required | Description |
|---|---|---|---|
| `script` | string | yes | Path to an existing shell script file |
| `args` | []any | no | Arguments passed to the script (`$1`, `$2`, ...) |
| `env` | object | no | Extra environment variables merged with the current process environment |
| `dir` | string | no | Working directory for the command (default: current process directory) |

> `timeout` is a top-level action field, not a driver option — see "Action fields" in `mink manual configuration`.

## Output

```json
{
  "exit_code": 0,
  "stdout": "...",
  "stderr": "",
  "success": true
}
```

| Field | Type | Description |
|---|---|---|
| `exit_code` | int | Exit code of the process |
| `stdout` | string | Standard output |
| `stderr` | string | Standard error |
| `success` | bool | `true` if `exit_code == 0` |

A non-zero exit code does **not** fail the action — it returns `success: false`. Use the `validate` driver to assert on the result.

## Example

```yaml
instances:
  sh:
    driver: shell
    methods: [run_code, run_script]

  check:
    driver: validate
    methods: [schema]

flows:
  - id: db_setup
    actions:
      - id: truncate
        description: "Truncate users table"
        instance: sh
        method: run_code
        timeout: 10000
        execute_with:
          code: "psql -U admin -c 'TRUNCATE users;'"
          env:
            PGPASSWORD: "${vars['db_pass']}"

      - id: check_truncated
        description: "Assert truncate succeeded"
        instance: check
        method: schema
        execute_with:
          value: "${actions['truncate']}"
          schema:
            type: object
            properties:
              success:
                type: boolean
                const: true

      - id: seed
        description: "Run an existing seed.sh script"
        instance: sh
        method: run_script
        execute_with:
          script: "scripts/seed.sh"
          args: ["--env", "test"]
```
