# Contributing

## Getting started

```bash
git clone https://github.com/nazarkhatsko/mink
cd mink
go mod download
```

## Development

```bash
task fmt        # format code
task lint       # run linter
task build      # build binary
task run -- examples/simple-api/mink.yaml
```

## Adding a custom driver

Implement the `Driver` interface from `pkg/driver`:

```go
import "github.com/nazarkhatsko/mink/pkg/driver"

type MyDriver struct{}

func (d *MyDriver) Name() string { return "mydriver" }

func (d *MyDriver) Describe() driver.Doc {
    return driver.Doc{
        Description: "What this driver does",
        Options: []driver.FieldDoc{
            {Name: "key", Type: "string", Required: true, Description: "What this option does"},
        },
        Output: []driver.FieldDoc{
            {Name: "result", Type: "string", Description: "What this field contains"},
        },
    }
}

func (d *MyDriver) Execute(ctx context.Context, options map[string]any) (driver.Output, error) {
    // ...
    return driver.Output{"result": "ok"}, nil
}
```

Register it before running:

```go
func init() {
    driver.Register(&MyDriver{})
}
```

Use in config:

```yaml
instances:
  my:
    driver: mydriver
    config:
      key: "value"
```

## Branches

| Pattern | When to use |
|---|---|
| `feat/<name>` | New feature or driver |
| `fix/<name>` | Bug fix |
| `docs/<name>` | Documentation only |
| `refactor/<name>` | Refactoring without behavior change |

Always branch off `main`.

## Commits

Follow [Conventional Commits](https://www.conventionalcommits.org):

```
feat: add postgres driver
fix: resolve array index out of bounds in dotpath
docs: add sleep driver reference
refactor: extract envfile loader into separate package
```

- Use present tense imperative: `add`, `fix`, `update` — not `added`, `fixes`
- Keep the subject line under 72 characters
- Reference issues when relevant: `fix: handle empty body (#42)`

## Pull requests

- One feature or fix per PR — keep diffs small and focused
- Add an example in `examples/` if adding a new driver
- Run `task fmt` and `task lint` before submitting
- PR title follows the same Conventional Commits format as commits
- Fill in a short description of what changed and why
