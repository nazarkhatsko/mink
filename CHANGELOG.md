# Changelog

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
- Variable resolution: `${vars.x}`, `${env.X}`, `${actions.id.field}`, `${actions.id.arr[0]}`
- Dot notation + array index traversal for action outputs
- `instances` block for named driver configuration
- Shallow merge of action `options` over instance `config`
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
