# Getting Started

For an illustrated walkthrough, see `docs/getting-started.md`. This page is the terminal-only technical reference.

## Installation

```bash
go install github.com/nazarkhatsko/mink@latest      # or: git clone + task build
```

## Minimal flow

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

```bash
mink run mink.yaml
```

## Commands

```bash
mink run <file>                         # run all flows
mink run <file> --flow <name>           # run a specific flow
mink run <file> --reporter compact      # change output format
mink run <file> --env .env              # load environment variables from file
mink validate <file>                    # validate config without running
mink manual configuration               # show configuration documentation
mink manual getting-started             # show this guide
mink manual drivers                     # list available drivers
mink manual drivers <name>              # show documentation for a driver
mink skill install --claude             # install bundled skills for Claude Code
mink skill install --codex              # install bundled skills for Codex
mink skill uninstall --claude           # remove installed skills
mink skill upgrade --claude             # upgrade installed skills to embedded version
mink version                            # print mink version
```

Skill commands accept `--global` to target `~/.<platform>/skills/` instead of the project directory.
