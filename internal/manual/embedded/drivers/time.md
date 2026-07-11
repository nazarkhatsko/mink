# time

Pauses flow execution for a given number of milliseconds. Useful when waiting for async operations to complete.

## Methods

### `wait`

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
  api:
    driver: http
    methods: [get, post]

  time:
    driver: time
    methods: [wait]

flows:
  - id: example
    actions:
      - id: trigger_job
        description: "Trigger async job"
        instance: api
        method: post
        execute_with:
          url: "${vars['base_url'] + '/jobs'}"

      - id: wait
        description: "Wait for job to complete"
        instance: time
        method: wait
        execute_with:
          ms: 2000

      - id: check_job
        description: "Check job result"
        instance: api
        method: get
        execute_with:
          url: "${vars['base_url'] + '/jobs/' + str(actions['trigger_job']['resp']['body']['id'])}"
```
