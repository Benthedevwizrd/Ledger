package fakes_test

import (
	"slices"
	"testing"
	"time"

	"ledger/internal/domain"
	"ledger/testing/fakes"

	"github.com/47monad/xuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EntryFaker_Make_Overrides(t *testing.T) {
	customID := xuid.MustNewSortable(domain.EntryIDPrefix)
	customTransactionID := xuid.MustNewSortable(domain.TransactionIDPrefix)
	customLedgerID := xuid.MustNewSortable(domain.LedgerIDPrefix)
	customAccountID := xuid.MustNewSortable(domain.AccountIDPrefix)
	customMoney := domain.MonetaryValue{
		Amount:   decimal.RequireFromString("12.34"),
		Currency: "EUR",
	}
	customTime := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	customMetadata := map[string]string{"k1": "v1", "k2": "v2"}
	customLine := int32(3)

	testCases := []struct {
		name      string
		givenOpts []domain.EntryOption
		assertFn  func(t *testing.T, got *domain.Entry)
	}{
		{
			name:      "with_id/success",
			givenOpts: []domain.EntryOption{fakes.WithEntryID(customID)},
			assertFn:  func(t *testing.T, got *domain.Entry) { assert.Equal(t, customID, got.ID) },
		},
		{
			name:      "with_transaction_id/success",
			givenOpts: []domain.EntryOption{fakes.WithEntryTransactionID(customTransactionID)},
			assertFn: func(t *testing.T, got *domain.Entry) {
				assert.Equal(t, customTransactionID, got.TransactionID)
			},
		},
		{
			name:      "with_ledger_id/success",
			givenOpts: []domain.EntryOption{fakes.WithEntryLedgerID(customLedgerID)},
			assertFn:  func(t *testing.T, got *domain.Entry) { assert.Equal(t, customLedgerID, got.LedgerID) },
		},
		{
			name:      "with_account_id/success",
			givenOpts: []domain.EntryOption{fakes.WithEntryAccountID(customAccountID)},
			assertFn:  func(t *testing.T, got *domain.Entry) { assert.Equal(t, customAccountID, got.AccountID) },
		},
		{
			name:      "with_direction_overrides_random/success",
			givenOpts: []domain.EntryOption{fakes.WithEntryDirection(domain.DirectionDebit)},
			assertFn: func(t *testing.T, got *domain.Entry) {
				assert.Equal(t, domain.DirectionDebit, got.Direction)
			},
		},
		{
			name:      "with_monetary_value_overrides_random/success",
			givenOpts: []domain.EntryOption{fakes.WithEntryMonetaryValue(customMoney)},
			assertFn: func(t *testing.T, got *domain.Entry) {
				assert.True(t, customMoney.Amount.Equal(got.MonetaryValue.Amount))
				assert.Equal(t, customMoney.Currency, got.MonetaryValue.Currency)
			},
		},
		{
			name:      "with_line_overrides_random/success",
			givenOpts: []domain.EntryOption{domain.WithEntryLine(customLine)},
			assertFn:  func(t *testing.T, got *domain.Entry) { assert.Equal(t, customLine, got.Line) },
		},
		{
			name:      "with_effective_at/success",
			givenOpts: []domain.EntryOption{domain.WithEntryEffectiveAt(customTime)},
			assertFn:  func(t *testing.T, got *domain.Entry) { assert.True(t, customTime.Equal(got.EffectiveAt)) },
		},
		{
			name:      "with_description_overrides_random/success",
			givenOpts: []domain.EntryOption{domain.WithEntryDescription("custom description")},
			assertFn:  func(t *testing.T, got *domain.Entry) { assert.Equal(t, "custom description", got.Description) },
		},
		{
			name:      "with_metadata_overrides_default/success",
			givenOpts: []domain.EntryOption{domain.WithEntryMetadata(customMetadata)},
			assertFn:  func(t *testing.T, got *domain.Entry) { assert.Equal(t, customMetadata, got.Metadata) },
		},
		{
			name:      "with_created_at/success",
			givenOpts: []domain.EntryOption{domain.WithEntryCreatedAt(customTime)},
			assertFn:  func(t *testing.T, got *domain.Entry) { assert.True(t, customTime.Equal(got.CreatedAt)) },
		},
		{
			name: "last_option_wins/success",
			givenOpts: []domain.EntryOption{
				domain.WithEntryDescription("first"),
				domain.WithEntryDescription("second"),
			},
			assertFn: func(t *testing.T, got *domain.Entry) { assert.Equal(t, "second", got.Description) },
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			faker := fakes.NewEntryFaker()

			// act
			result := faker.Make(tc.givenOpts...)

			// assert
			require.NotNil(t, result)
			tc.assertFn(t, result)
		})
	}
}

func Test_EntryFaker_Make_Defaults(t *testing.T) {
	// arrange
	faker := fakes.NewEntryFaker()

	// act
	result := faker.Make()

	// assert
	require.NotNil(t, result)
	assert.False(t, xuid.IsEmpty(result.ID))
	assert.False(t, xuid.IsEmpty(result.TransactionID))
	assert.False(t, xuid.IsEmpty(result.LedgerID))
	assert.False(t, xuid.IsEmpty(result.AccountID))
	assert.True(t, slices.Contains(validDirections, result.Direction), "unexpected direction %q", result.Direction)
	assert.True(t, result.MonetaryValue.Amount.IsPositive(), "amount should be positive, got %s", result.MonetaryValue.Amount)
	assert.NotEmpty(t, result.MonetaryValue.Currency)
	assert.GreaterOrEqual(t, result.Line, int32(1))
	assert.NotEmpty(t, result.Description)
	assert.GreaterOrEqual(t, len(result.Metadata), 1)
	assert.LessOrEqual(t, len(result.Metadata), 10)
	assert.WithinDuration(t, time.Now().UTC(), result.EffectiveAt, 5*time.Second)
	assert.WithinDuration(t, time.Now().UTC(), result.CreatedAt, 5*time.Second)
}

func Test_EntryFaker_Make_GeneratesUniqueIDs(t *testing.T) {
	// arrange
	faker := fakes.NewEntryFaker()

	// act
	first := faker.Make()
	second := faker.Make()

	// assert
	assert.NotEqual(t, first.ID, second.ID)
	assert.NotEqual(t, first.TransactionID, second.TransactionID)
	assert.NotEqual(t, first.LedgerID, second.LedgerID)
	assert.NotEqual(t, first.AccountID, second.AccountID)
}

func Test_EntryFaker_Make_CoversAllRandomizedValues(t *testing.T) {
	// probabilistic: with 300 draws the chance of missing any value is negligible
	// arrange
	faker := fakes.NewEntryFaker()
	directions := map[domain.Direction]struct{}{}
	lines := map[int32]struct{}{}

	// act
	for range 300 {
		e := faker.Make()
		directions[e.Direction] = struct{}{}
		lines[e.Line] = struct{}{}
	}

	// assert
	assert.Len(t, directions, len(validDirections))
	assert.Len(t, lines, 10)
}

func Test_EntryFaker_MakeMany_Count(t *testing.T) {
	testCases := []struct {
		name       string
		givenCount int
		wantedLen  int
	}{
		{name: "positive_count/success", givenCount: 5, wantedLen: 5},
		{name: "one/success", givenCount: 1, wantedLen: 1},
		{name: "zero_clamped_to_one/success", givenCount: 0, wantedLen: 1},
		{name: "negative_clamped_to_one/success", givenCount: -5, wantedLen: 1},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			faker := fakes.NewEntryFaker()

			// act
			result := faker.MakeMany(tc.givenCount)

			// assert
			require.Len(t, result, tc.wantedLen)
			ids := make(map[xuid.XUID]struct{}, len(result))
			for _, e := range result {
				require.NotNil(t, e)
				ids[e.ID] = struct{}{}
			}
			assert.Len(t, ids, tc.wantedLen, "every entry should have a unique id")
		})
	}
}

func Test_EntryFaker_MakeMany_AppliesSharedOpts(t *testing.T) {
	// arrange
	faker := fakes.NewEntryFaker()
	ledgerID := xuid.MustNewSortable(domain.LedgerIDPrefix)
	accountID := xuid.MustNewSortable(domain.AccountIDPrefix)

	// act
	result := faker.MakeMany(3, fakes.WithEntryLedgerID(ledgerID), fakes.WithEntryAccountID(accountID))

	// assert
	require.Len(t, result, 3)
	for _, e := range result {
		assert.Equal(t, ledgerID, e.LedgerID)
		assert.Equal(t, accountID, e.AccountID)
	}
}

func Test_EntryFaker_MakeMany_SharedIDOptCollides(t *testing.T) {
	// documents the known behavior: a pinned ID is shared by every entry
	// arrange
	faker := fakes.NewEntryFaker()
	id := xuid.MustNewSortable(domain.EntryIDPrefix)

	// act
	result := faker.MakeMany(3, fakes.WithEntryID(id))

	// assert
	require.Len(t, result, 3)
	for _, e := range result {
		assert.Equal(t, id, e.ID)
	}
}

func Test_EntryOptions(t *testing.T) {
	id := xuid.MustNewSortable(domain.EntryIDPrefix)
	money := domain.MonetaryValue{
		Amount:   decimal.NewFromInt(100),
		Currency: "USD",
	}

	testCases := []struct {
		name     string
		givenOpt domain.EntryOption
		assertFn func(t *testing.T, got *domain.Entry)
	}{
		{
			name:     "with_entry_id",
			givenOpt: fakes.WithEntryID(id),
			assertFn: func(t *testing.T, got *domain.Entry) { assert.Equal(t, id, got.ID) },
		},
		{
			name:     "with_entry_transaction_id",
			givenOpt: fakes.WithEntryTransactionID(id),
			assertFn: func(t *testing.T, got *domain.Entry) { assert.Equal(t, id, got.TransactionID) },
		},
		{
			name:     "with_entry_ledger_id",
			givenOpt: fakes.WithEntryLedgerID(id),
			assertFn: func(t *testing.T, got *domain.Entry) { assert.Equal(t, id, got.LedgerID) },
		},
		{
			name:     "with_entry_account_id",
			givenOpt: fakes.WithEntryAccountID(id),
			assertFn: func(t *testing.T, got *domain.Entry) { assert.Equal(t, id, got.AccountID) },
		},
		{
			name:     "with_entry_direction",
			givenOpt: fakes.WithEntryDirection(domain.DirectionDebit),
			assertFn: func(t *testing.T, got *domain.Entry) { assert.Equal(t, domain.DirectionDebit, got.Direction) },
		},
		{
			name:     "with_entry_monetary_value",
			givenOpt: fakes.WithEntryMonetaryValue(money),
			assertFn: func(t *testing.T, got *domain.Entry) {
				assert.True(t, money.Amount.Equal(got.MonetaryValue.Amount))
				assert.Equal(t, money.Currency, got.MonetaryValue.Currency)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			entry := &domain.Entry{}

			// act
			tc.givenOpt(entry)

			// assert
			tc.assertFn(t, entry)
		})
	}
}
