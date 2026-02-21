# Project conventions for Claude Code

## Git / Commits

- **Never commit generated files.** Files produced by code generators (e.g. `src/api/api.gen.go` produced by `oapi-codegen`) are gitignored and must stay that way. Run the relevant `make` target to regenerate them locally when needed.
