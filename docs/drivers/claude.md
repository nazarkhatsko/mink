# claude

Sends messages to the Claude API.

## Config

| Field | Type | Required | Description |
|---|---|---|---|
| `api_key` | string | yes | Anthropic API key (falls back to `ANTHROPIC_API_KEY` env var) |
| `model` | string | no | Model name. Default: `claude-sonnet-4-6` |

## Options

| Field | Type | Required | Description |
|---|---|---|---|
| `messages` | []object | yes | Conversation messages: `[{role, content}]` |
| `system` | string | no | System prompt |
| `max_tokens` | int | no | Maximum tokens to generate. Default: `1024` |
| `temperature` | float | no | Sampling temperature |

## Output

```json
{
  "in": {
    "model": "claude-sonnet-4-6",
    "messages": [ ... ],
    "system": "...",
    "max_tokens": 16
  },
  "out": {
    "id": "msg_...",
    "model": "claude-sonnet-4-6",
    "role": "assistant",
    "content": "bug",
    "stop_reason": "end_turn",
    "usage": { "input_tokens": 42, "output_tokens": 3 }
  }
}
```

| Field | Type | Description |
|---|---|---|
| `in` | object | Outgoing request: `model`, `messages`, `system`, `max_tokens` |
| `out` | object | Parsed response: `id`, `model`, `role`, `content`, `stop_reason`, `usage` |
| `out.content` | string | Concatenated text from the response's text content blocks |
| `out.usage` | object | `input_tokens`, `output_tokens` |

## Example

```yaml
vars:
  api_key: "${env['ANTHROPIC_API_KEY']}"

instances:
  llm:
    driver: claude
    config:
      api_key: "${vars['api_key']}"
      model: "claude-sonnet-4-6"

flows:
  - name: "llm-sanity-check"
    actions:
      - id: ask_claude
        description: "Ask Claude to classify a support ticket"
        use: llm
        run_with:
          system: "You are a support ticket triage assistant. Reply with one word: bug, question, or feature."
          messages:
            - role: user
              content: "The app crashes every time I upload a PDF over 10MB."
          max_tokens: 16

      - id: validate_classification
        description: "Validate Claude classified it as a bug"
        use: check
        run_with:
          value: "${actions['ask_claude']['out']}"
          schema:
            type: object
            properties:
              content:
                type: string
                const: "bug"
```
