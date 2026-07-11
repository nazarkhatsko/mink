# shell

Executes a shell command via `sh -c`. Useful for database seeding, service health checks, file setup, or any system-level operation in a flow.

## Options

| Field | Type | Required | Description |
|---|---|---|---|
| `command` | string | yes | Shell command to run |
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

flows:
  - id: db_setup
    actions:
      - id: truncate
        description: "Truncate users table"
        instance: sh
        timeout: 10000
        execute_with:
          command: "psql -U admin -c 'TRUNCATE users;'"
          env:
            PGPASSWORD: "${vars['db_pass']}"

      - id: check_truncated
        description: "Assert truncate succeeded"
        instance: check
        execute_with:
          value: "${actions['truncate']}"
          schema:
            type: object
            properties:
              success:
                type: boolean
                const: true
```
