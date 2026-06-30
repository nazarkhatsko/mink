Generate a mink flow YAML file based on the user's description.

## Rules

- Always start with `version: "1.0"` and `info:` block
- Define all required drivers in `instances:` before using them in `flows:`
- Every action must have `id:`, `description:`, `use:`, and `options:`
- Use `${vars.x}` for reusable values, `${env.X}` for secrets
- Reference previous action outputs via `${actions.<id>.<field>}`
- Use `gen` instance with `driver: generate` for generating fake data
- Use `check` instance with `driver: validate` for JSON Schema validation
- Use `sleep` instance with `driver: sleep` when async delay is needed
- HTTP actions must always include `method:` and `url:` in `options:`

## Available drivers

| driver | purpose |
|---|---|
| `http` | HTTP requests |
| `generate` | fake data generation |
| `validate` | JSON Schema validation |
| `sleep` | delay execution |

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

## Example

User: "Generate a user, post it to /users, validate the response has an id"

```yaml
version: "1.0"
info:
  name: "Create user flow"
  description: "Generate and post a user, validate response"

vars:
  base_url: "${env.BASE_URL}"

instances:
  api:
    driver: http
  gen:
    driver: generate
  check:
    driver: validate

flows:
  - name: "create-user"
    actions:
      - id: generate_user
        description: "Generate random user payload"
        use: gen
        options:
          schema:
            name:
              type: string
              value: faker.name
            email:
              type: string
              value: faker.email

      - id: create_user
        description: "POST user to API"
        use: api
        options:
          method: POST
          url: "${vars.base_url}/users"
          body: "${actions.generate_user}"

      - id: validate_response
        description: "Validate response contains id"
        use: check
        options:
          value: "${actions.create_user.resp.body}"
          schema:
            type: object
            required: [id]
```
