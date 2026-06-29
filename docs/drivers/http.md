# http

Executes HTTP requests.

## Options

| Field | Type | Required | Description |
|---|---|---|---|
| `method` | string | no | HTTP method. Default: `GET` |
| `url` | string | yes | Full URL |
| `headers` | object | no | Request headers |
| `body` | any | no | Request body, serialized as JSON |

## Output

```json
{
  "status": 201,
  "body": { ... },
  "headers": {
    "Content-Type": "application/json"
  }
}
```

| Field | Type | Description |
|---|---|---|
| `status` | int | HTTP status code |
| `body` | any | Parsed JSON body, or raw string if not JSON |
| `headers` | object | Response headers |

## Example

```yaml
instances:
  api:
    driver: http

flows:
  - name: "example"
    actions:
      - id: create_user
        description: "POST new user"
        use: api
        options:
          method: POST
          url: "${vars.base_url}/users"
          headers:
            Authorization: "Bearer ${actions.login.body.token}"
            X-Request-Id: "${actions.gen.request_id}"
          body:
            name: "${actions.gen.name}"
            email: "${actions.gen.email}"
```

## Notes

- `Content-Type: application/json` is set automatically when `body` is present
- Non-2xx responses do **not** cause the action to fail — use the `validate` driver to assert status codes
