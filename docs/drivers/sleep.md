# sleep

Pauses flow execution for a given number of milliseconds. Useful when waiting for async operations to complete.

## Options

| Field | Type | Required | Description |
|---|---|---|---|
| `ms` | int | yes | Duration in milliseconds |

## Output

```json
{ "slept_ms": 500 }
```

## Example

```yaml
instances:
  sleep:
    driver: sleep

flows:
  - name: "example"
    actions:
      - id: trigger_job
        description: "Trigger async job"
        use: api
        run_with:
          method: POST
          url: "${vars['base_url'] + '/jobs'}"

      - id: wait
        description: "Wait for job to complete"
        use: sleep
        run_with:
          ms: 2000

      - id: check_job
        description: "Check job result"
        use: api
        run_with:
          method: GET
          url: "${vars['base_url'] + '/jobs/' + str(actions['trigger_job']['resp']['body']['id'])}"
```
