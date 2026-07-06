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

func (d *MyDriver) Execute(ctx context.Context, options map[string]any) (driver.Output, error) {
    key, ok := options["key"].(string)
    if !ok {
        return nil, driver.NewError(driver.ErrConfig, "mydriver: key is required")
    }
    // ... call out to whatever mydriver wraps ...
    return driver.Output{"result": "ok"}, nil
}
```

Return errors via `driver.NewError`/`driver.Wrap` instead of a bare `fmt.Errorf`, so `mutate_on.fail`'s `event['error']['code']` is meaningful. Pick the code by what actually failed: `driver.ErrConfig` for a missing/invalid option, `driver.ErrTransport` for a failed external call/process (use `driver.Wrap` to keep the original error reachable via `errors.Is`/`errors.As`), `driver.ErrInternal` for anything else (marshaling, parsing). `driver.ErrTimeout` is applied automatically by the engine when `ctx`'s deadline is exceeded — don't set it yourself.

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

Add `internal/manual/embedded/drivers/mydriver.md` (config/options/output fields, an example) and list it in `internal/manual/embedded/drivers.md` — this is the only source for `mink manual drivers mydriver`.

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
