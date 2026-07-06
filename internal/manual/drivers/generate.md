# generate

Generates fake data based on a schema. Useful for creating randomized test payloads without hardcoding values.

## Options

| Field | Type | Required | Description |
|---|---|---|---|
| `schema` | object | yes | Object schema defining fields to generate |

### Schema field definition

Each field in `schema` is an object with:

| Field | Type | Required | Description |
|---|---|---|---|
| `type` | string | yes | Go type for the generated value |
| `value` | string | yes | Generator or static value |

### Types

| `type` | Go type |
|---|---|
| `string` | `string` |
| `int` | `int64` |
| `float` | `float64` |
| `bool` | `bool` |

### Values

| `value` | Description |
|---|---|
| `faker.name` | Full name |
| `faker.email` | Email address |
| `faker.password` | Random password |
| `faker.uuid` | UUID v4 |
| `faker.phone` | Phone number |
| `faker.age` | Integer age (18–80) |
| `"<literal>"` | Static value |
| `${env.VAR}` | Value from environment |

## Output

Returns the generated object directly.

## Example

```yaml
instances:
  gen:
    driver: generate

flows:
  - name: "example"
    actions:
      - id: payload
        description: "Generate user payload"
        use: gen
        run_with:
          schema:
            name:
              type: string
              value: faker.name
            email:
              type: string
              value: faker.email
            age:
              type: int
              value: faker.age
            role:
              type: string
              value: "user"
```

Output:
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 34,
  "role": "user"
}
```

## Nested objects

Fields without `type`/`value` are treated as nested objects:

```yaml
schema:
  address:
    city:
      type: string
      value: faker.name
    zip:
      type: string
      value: "00000"
```
