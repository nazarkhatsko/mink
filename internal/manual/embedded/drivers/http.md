# http

Executes HTTP requests. Each method pins a fixed HTTP verb, except `request`
which takes it as an option — for verbs without a dedicated method (`HEAD`,
`OPTIONS`, custom).

## Config

| Field | Type | Required | Description |
|---|---|---|---|
| `headers` | object | no | Headers merged into every request from this instance |

## Methods

### `get`, `post`, `put`, `patch`, `delete`

| Field | Type | Required | Description |
|---|---|---|---|
| `url` | string | yes | Full URL |
| `headers` | object | no | Request headers |
| `body` | any | no | Request body, serialized as JSON |

### `request`

Escape hatch for HTTP verbs without a dedicated method.

| Field | Type | Required | Description |
|---|---|---|---|
| `method` | string | yes | HTTP method, e.g. `HEAD`, `OPTIONS` |
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
    methods: [get, post]

flows:
  - id: example
    actions:
      - id: create_user
        description: "POST new user"
        instance: api
        method: post
        execute_with:
          url: "${vars['base_url'] + '/users'}"
          headers:
            Authorization: "${'Bearer ' + state['token']}"
            X-Request-Id: "${actions['gen']['request_id']}"
          body:
            name: "${actions['gen']['name']}"
            email: "${actions['gen']['email']}"

      - id: fetch_user
        description: "GET the created user"
        instance: api
        method: get
        execute_with:
          url: "${vars['base_url'] + '/users/' + str(state['user_id'])}"
```

## Notes

- `Content-Type: application/json` is set automatically when `body` is present
- Non-2xx responses do **not** cause the action to fail — use the `validate` driver to assert status codes
- `execute_with` is strictly validated per method — e.g. `method:` is only a valid key under `request`, not under `get`/`post`/`put`/`patch`/`delete` (their verb is fixed by which method you call)
