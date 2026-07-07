# 🧩 Skills

mink ships AI assistant skills (e.g. `mink-flow`, for AI-assisted flow generation) embedded in the binary. The `mink skill` command installs them into your project, or globally for your user.

⚙️ **Commands**
```bash
mink skill install --claude     # install skills for Claude Code
mink skill install --codex      # install skills for Codex
mink skill uninstall --claude   # remove installed skills
mink skill upgrade --claude     # upgrade installed skills to the embedded version
```

Exactly one of `--claude` or `--codex` is required — they target different skill directories and are mutually exclusive.

📂 **Install location**

By default skills are installed into the current project:

| Flag | Target |
|---|---|
| `--claude` | `.claude/skills/` |
| `--codex` | `.codex/skills/` |

Pass `--global` to install into your home directory instead (`~/.claude/skills/` or `~/.codex/skills/`).

```bash
mink skill install --claude --global
```

Without `--global`, the target subdirectory (`.claude/skills/` or `.codex/skills/`) must already exist in the current directory.

## 🔄 Versioning

Each installed skill is written as a `<name>/` directory containing `SKILL.md` and a `VERSION` file. `mink skill upgrade` compares `VERSION` against the version embedded in the binary and rewrites the skill if they differ. `mink skill install` is idempotent — it's a no-op when the installed skill is already up to date.

## 📦 Available skills

| Skill | Description |
|---|---|
| `mink-flow` | Generates a mink flow YAML file based on a natural-language description |

⬅️ [Back to docs](README.md)
