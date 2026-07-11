# firestore

Fetches a single Firestore document by path. Authorization is always through an explicit service-account key file — there is no Application Default Credentials fallback (no gcloud user login, no GCE/GKE attached metadata service account).

## Config

| Field | Type | Required | Description |
|---|---|---|---|
| `project_id` | string | yes | GCP project ID |
| `credentials_file` | string | yes | Path to a service-account JSON key file (falls back to `GOOGLE_APPLICATION_CREDENTIALS` env var) |
| `database_id` | string | no | Firestore database ID. Default: `(default)` |

## Options

| Field | Type | Required | Description |
|---|---|---|---|
| `path` | string | yes | Document path, e.g. `users/abc123` or `orgs/org1/users/user2` — must have an even number of non-empty segments |

## Output

```json
{
  "doc": { "name": "Ada Lovelace", "age": 36 },
  "exists": true
}
```

| Field | Type | Description |
|---|---|---|
| `doc` | object \| null | Document body, or `null` if the document doesn't exist |
| `exists` | bool | `true` if the document exists |

A missing document is **not** an error — it returns `exists: false, doc: null`. The action only errors on real failures: bad config, a malformed path, or an auth/permission/network problem talking to Firestore.

## Example

```yaml
vars:
  sa_key_path: "${env['FIRESTORE_SA_KEY_PATH']}"

instances:
  fs:
    driver: firestore
    config:
      project_id: "my-gcp-project"
      credentials_file: "${vars['sa_key_path']}"

  check:
    driver: validate

flows:
  - id: fetch_user
    actions:
      - id: fetch_user
        description: "Fetch a user document by its Firestore path"
        instance: fs
        execute_with:
          path: "users/abc123"

      - id: validate_user
        description: "Validate the user document exists and has the expected fields"
        instance: check
        execute_with:
          value: "${actions['fetch_user']}"
          schema:
            type: object
            properties:
              exists:
                type: boolean
                const: true
              doc:
                type: object
                required: [name]
```
