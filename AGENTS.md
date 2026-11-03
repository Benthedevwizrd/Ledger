# AGENTS.md

Double-entry ledger in Go (module `ledger`). Imports look like `ledger/internal/domain`.

## Commands
- Test everything: `go test ./...`
- One package: `go test ./internal/domain -v`
- One test/case: `go test ./internal/domain -run 'Test_AccountMake/valid_inputs' -v`
- Vet: `go vet ./...`
- Format: `gofmt -w .`
- Do NOT treat `go build ./...` as verification: `cmd/main.go` is `package main` with no `func main()` yet, so it fails. `go test ./...` and `go vet ./...` work.

## Workflow
- Branches follow `feat/implement-<feature>` (e.g. `feat/implement-account-faker`, `feat/implement-entry-faker`).

## Current state
- Real code lives in `internal/domain` (entities, enums, factories) and `testing/fakes` (random test data).
- `internal/appsvc`, `internal/repository`, `interface/shttp`, and `cmd` are empty placeholder packages. README's PostgreSQL/engine claims are not implemented; don't assume persistence exists.

## Domain conventions
- IDs are `xuid.XUID` from `xuid.NewSortable(domain.<Entity>IDPrefix)`; use the existing prefix constants (`account_`, `ledger_`, `transaction_`, `entry_`).
- Construction goes through `New<Entity>Factory().Make(required..., opts...)`. Required-arg validation returns package-level `Err...` sentinels (e.g. `ErrEmptyAccountID`), checked in argument order.
- Factories set defaults, then apply options in order, so later options win. Don't add constructors that bypass validation.
- `Entry` is append-only (no `UpdatedAt`); posted transactions are immutable and corrected via reversal.
- Money is `domain.MonetaryValue{Amount decimal.Decimal, Currency string}`. Compare amounts with `.Amount.Equal(...)`, not `==`.

## Testing conventions
- Tests are external packages (`domain_test`, `fakes_test`), table-driven, using `testify`: `require` for fatal setup/errors, `assert` for value checks.
- Timestamps are generated with `time.Now().UTC()`; assert with `assert.WithinDuration`, not equality.
- Each factory has a matching test file; add table cases for both success and sentinel errors when changing validation.

## Fakers (`testing/fakes`)
- Each faker has a private `make` that generates randomized defaults then appends caller options **last**, so caller options override randomness.
- Public API is `Make(opts...)` and `MakeMany(count, opts...)`; count < 1 is clamped to 1. Fakers panic on factory errors, so faker options may bypass validation intentionally.
- Faker-specific options (e.g. `WithAccountID`) live in `fakes`, not `domain`. `MakeMany` with a pinned ID deliberately reuses that ID.
- When adding an entity, mirror the existing factory + domain test + faker + faker test set.

## Notes
- No Makefile, CI, or lint config exists; use the Go commands above.
- Keep `gofmt` clean.
