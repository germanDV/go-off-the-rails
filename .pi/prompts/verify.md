---
description: Verify that an implemented feature matches its spec file.
argument-hint: "{feature-name}"
---
Verify that the implementation of `specs/$1.md` is complete and correct.

Please:
1. Read `specs/$1.md`.
2. Compare every item in the spec against the actual code:
   - DB migrations exist and match
   - sqlc queries exist and match
   - Domain models and logic exist
   - Repositories exist with correct methods
   - Controller endpoints exist and are wired
   - Views/templates exist
   - Routes are registered with correct middleware
3. Report any gaps or mismatches.
4. Confirm if the feature is fully implemented or if steps remain.
