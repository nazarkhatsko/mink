# http

Executes HTTP requests.

## Config

| Field | Type | Required | Description |
|---|---|---|---|
| `timeout` | int | no | Request timeout in milliseconds (default: no timeout) |

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
  "req": {
    "method": "POST",
    "url": "https://api.example.com/users",
    "headers": { "Content-Type": "application/json" },
    "body": { ... }
  },
  "resp": {
    "status": 201,
    "headers": { "Content-Type": "application/json" },
    "body": { ... }
  }
}
```

| Field | Type | Description |
|---|---|---|
| `req` | object | Outgoing request: `method`, `url`, `headers`, `body` |
| `resp` | object | Incoming response: `status`, `headers`, `body` |
| `resp.status` | int | HTTP status code |
| `resp.body` | any | Parsed JSON body, or raw string if not JSON |
| `resp.headers` | object | Response headers |

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
        run_with:
          method: POST
          url: "${vars.base_url}/users"
          headers:
            Authorization: "Bearer ${actions.login.resp.body.token}"
            X-Request-Id: "${actions.gen.request_id}"
          body:
            name: "${actions.gen.name}"
            email: "${actions.gen.email}"
```

## Notes

- `Content-Type: application/json` is set automatically when `body` is present
- Non-2xx responses do **not** cause the action to fail — use the `validate` driver to assert status codes
