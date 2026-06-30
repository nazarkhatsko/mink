# Getting Started

## Installation

```bash
go install github.com/nazarkhatsko/mink@latest
```

Or build from source:

```bash
git clone https://github.com/nazarkhatsko/mink
cd mink
task build
```

## Your first flow

Create a file `mink.yaml`:

```yaml
version: "1.0"
info:
  name: "My first flow"
  description: "Ping httpbin and validate response"

instances:
  api:
    driver: http

  check:
    driver: validate

flows:
  - name: "ping"
    actions:
      - id: ping
        description: "GET httpbin"
        use: api
        options:
          method: GET
          url: "https://httpbin.org/get"

      - id: validate
        description: "Validate response status"
        use: check
        options:
          value: "${actions.ping.resp}"
          schema:
            type: object
            properties:
              status:
                type: integer
                const: 200
```

## Run it

```bash
mink run mink.yaml
```

## Commands

```bash
mink run <file>                      # run all flows
mink run <file> --flow <name>        # run a specific flow
mink run <file> --reporter compact   # change output format
mink run <file> --env .env           # load environment variables from file
mink validate <file>                 # validate config without running
mink doc drivers                     # list available drivers
mink doc drivers <name>              # show documentation for a driver
mink skill install --claude          # install bundled skills for Claude Code
mink skill install --codex           # install bundled skills for Codex
mink skill uninstall --claude        # remove installed skills
mink skill upgrade --claude          # upgrade installed skills to embedded version
mink version                         # print mink version
```

Skill commands accept `--global` to target `~/.<platform>/skills/` instead of the project directory.
