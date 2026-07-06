# 🚀 Getting Started

Write API tests in YAML, run them with one command.

## 📦 Install

```bash
go install github.com/nazarkhatsko/mink@latest
```

## ✍️ Write a flow

Create `mink.yaml`:

```yaml
version: "1.0"
info:
  name: "My first flow"

instances:
  api:
    driver: http
  check:
    driver: validate

flows:
  - name: "ping"
    actions:
      - id: ping
        use: api
        run_with:
          method: GET
          url: "https://httpbin.org/get"

      - id: validate
        use: check
        run_with:
          value: "${actions['ping']['resp']}"
          schema:
            type: object
            properties:
              status: { type: integer, const: 200 }
```

## ▶️ Run it

```bash
mink run mink.yaml
```

<details>
<summary>Expected output</summary>

```
▶ flow: ping
  ✓ ping (120ms)
  ✓ validate (0ms)
✓ flow done: ping
```

</details>

## 🧭 Where next

| | Guide |
|---|---|
| 🧠 | [Flows & State](flows-and-state.md) — how `actions`, `state`, and `${}` fit together |
| 🧪 | [Examples](../examples/README.md) — runnable flows to copy from |
| 🔌 | Run `mink manual drivers` for the full driver reference |

⬅️ [Back to docs](README.md)
