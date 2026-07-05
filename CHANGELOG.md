# Changelog

## Unreleased

### Drivers
- Added `shell` driver — executes shell commands via `sh -c`; options: `command` (required), `env`, `dir`; output: `exit_code`, `stdout`, `stderr`, `success`

### Core
- Added `timeout` as a top-level action field (milliseconds); applies to any driver via `context.WithTimeout`
- **Breaking:** `${}` expressions now use [Starlark](https://github.com/google/starlark-go) instead of custom dot-path resolver — use dict access `vars['key']`, `actions['id']['field']`, `env['KEY']` instead of dot notation
- Added `state` — mutable dict shared across all actions in a flow, readable in any `${}` expression
- Added `mutate_on` action field with `done` and `fail` Starlark scripts; `done` has access to `event['result']`, `fail` has access to `event['error']`

### Documentation
- `docs/drivers/shell.md` — new driver reference
- `docs/drivers/README.md` — added `shell` to built-in driver table
- `docs/configuration.md` — documented `timeout`, `mutate_on`, `state`, and new Starlark expression syntax
- Updated all examples and driver docs to new `${}` syntax

### CLI
- Renamed embedded skill package `skills/` → `skill/` (`package skills` → `package skill`) to match the `mink skill` command and `internal/skill` package naming

---


### CLI
- Restructured CLI as `cmd/cmd.go` + per-command packages (`cmd/run`, `cmd/validate`, `cmd/doc`, `cmd/skill`, `cmd/version`), built via cobra; entrypoint moved from `cmd/mink/main.go` to root `main.go`
- `mink list-drivers` replaced by `mink doc drivers [name]` — lists drivers, or shows config/options/output for a single driver
- Added `mink skill install|uninstall|upgrade --claude|--codex [--global]` to manage bundled AI assistant skills
- Added `mink version`

### Drivers
- Drivers now implement `Describe() driver.Doc` (`pkg/driver`), powering `mink doc drivers`
- `http` driver output reshaped: `{status, body, headers}` → `{req: {method, url, headers, body}, resp: {status, headers, body}}`
- `http` driver gained `timeout` config option (request timeout in milliseconds)

### Skills
- `skills/skills.go` embeds skill content per platform (`claude/`, `codex/`) with `Get`/`GetVersion`/`List`
- `internal/skill` package backs `mink skill` install/uninstall/upgrade, version-tracked via `<name>.version` sidecar files
- Manual `skills/claude/generate-flow/install.sh` and `skills/README.md` removed in favor of `mink skill install`
- Added `skills/codex/generate-flow` — Codex equivalent of the Claude `generate-flow` skill

### Documentation
- `docs/getting-started.md` — updated install path and command list
- `docs/drivers/http.md` — updated for `req`/`resp` output shape and `timeout` config
- `docs/skills.md` — new page documenting `mink skill` usage
- `CONTRIBUTING.md` — custom driver example now includes required `Describe()` method

## v1.0.0 — 2026-06-29

Initial release.

### Project
- Go module `github.com/nazarkhatsko/mink`
- Project structure: `cmd/`, `internal/`, `pkg/`, `examples/`, `docs/`, `skills/`
- `Taskfile.yaml` with `fmt`, `lint`, `lint:fix`, `build`, `clean`, `run` tasks
- MIT License
- `README.md` with banner, links to docs and changelog
- `CONTRIBUTING.md` and `CHANGELOG.md`

### Core
- Engine with sequential flow and action execution
- Variable resolution: `${vars['x']}`, `${env['X']}`, `${actions['id']['field']}` via Starlark expressions
- `instances` block for named driver configuration
- Shallow merge of action `run_with` over instance `config`
- Driver registry with public SDK (`pkg/driver`) for custom drivers

### Drivers
- `http` — HTTP requests with `method`, `url`, `headers`, `body`; output: `status`, `body`, `headers`
- `generate` — Fake data generation via schema (`type` + `value`), supports `faker.*` generators
- `validate` — JSON Schema validation (Draft 2020-12), fails action on invalid
- `sleep` — Execution delay via `ms`; output: `slept_ms`

### CLI
- `mink run <file>` — run all flows from a config file
- `mink run <file> --flow <name>` — run a specific flow by name
- `mink run <file> --reporter <name>` — output format: `pretty`, `compact`, `json`, `silent`
- `mink run <file> --env <file>` — load environment variables from a `.env` file
- `mink validate <file>` — validate config without running
- `mink list-drivers` — list registered drivers

### Reporters
- `pretty` — human-readable output with icons, descriptions, and timing (default)
- `compact` — one line per action with timing
- `json` — NDJSON output for CI pipelines
- `silent` — no output, exit code only

### Examples
- `examples/simple-api` — local Go HTTP server with CRUD `/users` and API key auth; flows: `unauthorized-request`, `user-lifecycle`
- `examples/httpbin` — flow against httpbin.org; no local server needed

### Documentation
- `docs/getting-started.md` — installation and first flow
- `docs/configuration.md` — full YAML spec and variable resolution rules
- `docs/reporters.md` — reporter formats with example output
- `docs/drivers/` — per-driver reference: `http`, `generate`, `validate`, `sleep`

### Skills
- `skills/claude/generate-flow` — Claude Code skill for AI-assisted flow generation
