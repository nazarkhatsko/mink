# validate

Validates a value against a JSON Schema. If validation fails, the flow stops and an error is reported.

## Options

| Field | Type | Required | Description |
|---|---|---|---|
| `value` | any | yes | Value to validate |
| `schema` | object | yes | JSON Schema definition |

## Output

```json
{ "valid": true }
```

On failure the action errors and stops the flow — no output is produced.

## Example

```yaml
instances:
  check:
    driver: validate

flows:
  - name: "example"
    actions:
      - id: validate_response
        description: "Validate API response"
        use: check
        options:
          value: "${actions.create_user.resp.body}"
          schema:
            type: object
            required: [id, name, email]
            properties:
              id:
                type: string
              name:
                type: string
              email:
                type: string
                format: email
```

## Validating status codes

```yaml
- id: validate_status
  description: "Assert 201 Created"
  use: check
  options:
    value: "${actions.create_user.resp}"
    schema:
      type: object
      properties:
        status:
          type: integer
          const: 201
```

## Notes

- Uses [JSON Schema Draft 2020-12](https://json-schema.org/draft/2020-12)
- The `schema` field accepts any valid JSON Schema object
