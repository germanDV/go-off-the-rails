---
description: Implement a feature from its spec file, full tool access enabled.
argument-hint: "{feature-name}"
---
Implement the feature described in `specs/$1.md`.

Please:
1. Read `specs/$1.md` carefully.
2. Exit plan mode (restore full tool access if restricted).
3. Implement the spec step by step in order.
4. Mark completed steps in your responses using `[DONE:n]` where `n` is the step number.
5. If you discover during implementation that the spec needs to change, suggest an edit to the spec before proceeding.
6. Run `sqlc generate` if new queries were added.
7. Run `go build ./...` after changes to ensure everything compiles.
8. Run `go test ./...` to ensure all tests pass.
