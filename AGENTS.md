# GOTR Project Guidelines

## Workflow

Always use the three-phase workflow:

1. **Discuss** — `/think [topic]` — read-only exploration, no file modifications
2. **Spec** — `/spec feature-name` — write detailed plan to `specs/{feature-name}.md`
3. **Build** — `/build feature-name` — implement from spec, mark steps with `[DONE:n]`
4. **Verify** — `/verify feature-name` — audit completeness against spec

Use `specs/_template.md` as the scaffold for every new spec.

## Verification Checklist

After implementing any feature:
- [ ] Migrations exist and match the spec
- [ ] sqlc queries generated (`sqlc generate`)
- [ ] Domain model and validation in place
- [ ] Repository wraps generated queries, maps to domain
- [ ] Controller endpoints registered with correct middleware
- [ ] Views/templates render correctly
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
