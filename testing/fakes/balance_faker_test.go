package fakes_test

import (
	"testing"
	"time"

	"ledger/internal/domain"
	"ledger/testing/fakes"

	"github.com/47monad/xuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BalanceFaker_Make_Overrides(t *testing.T) {
	customAccountID := xuid.MustNewSortable(domain.AccountIDPrefix)
	customLedgerID := xuid.MustNewSortable(domain.LedgerIDPrefix)
	customEntryID := xuid.MustNewSortable(domain.EntryIDPrefix)
	customTime := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)

	debitsPosted := domain.MonetaryValue{Amount: decimal.NewFromInt(100), Currency: "USD"}
	debitsPending := domain.MonetaryValue{Amount: decimal.NewFromInt(50), Currency: "USD"}
	creditsPosted := domain.MonetaryValue{Amount: decimal.NewFromInt(200), Currency: "USD"}
	creditsPending := domain.MonetaryValue{Amount: decimal.NewFromInt(25), Currency: "USD"}

	testCases := []struct {
		name      string
		givenOpts []domain.BalanceOption
		assertFn  func(t *testing.T, got *domain.Balance)
	}{
		{
			name:      "with_account_id/success",
			givenOpts: []domain.BalanceOption{fakes.WithBalanceAccountID(customAccountID)},
			assertFn:  func(t *testing.T, got *domain.Balance) { assert.Equal(t, customAccountID, got.AccountID) },
		},
		{
			name:      "with_ledger_id/success",
			givenOpts: []domain.BalanceOption{fakes.WithBalanceLedgerID(customLedgerID)},
			assertFn:  func(t *testing.T, got *domain.Balance) { assert.Equal(t, customLedgerID, got.LedgerID) },
		},
		{
			name:      "with_debits_posted/success",
			givenOpts: []domain.BalanceOption{domain.WithBalanceDebitsPosted(debitsPosted)},
			assertFn:  func(t *testing.T, got *domain.Balance) { assertMonetaryValueEqual(t, debitsPosted, got.DebitsPosted) },
		},
		{
			name:      "with_debits_pending/success",
			givenOpts: []domain.BalanceOption{domain.WithBalanceDebitsPending(debitsPending)},
			assertFn:  func(t *testing.T, got *domain.Balance) { assertMonetaryValueEqual(t, debitsPending, got.DebitsPending) },
		},
		{
			name:      "with_credits_posted/success",
			givenOpts: []domain.BalanceOption{domain.WithBalanceCreditsPosted(creditsPosted)},
			assertFn:  func(t *testing.T, got *domain.Balance) { assertMonetaryValueEqual(t, creditsPosted, got.CreditsPosted) },
		},
		{
			name:      "with_credits_pending/success",
			givenOpts: []domain.BalanceOption{domain.WithBalanceCreditsPending(creditsPending)},
			assertFn: func(t *testing.T, got *domain.Balance) {
				assertMonetaryValueEqual(t, creditsPending, got.CreditsPending)
			},
		},
		{
			name:      "with_last_entry_id/success",
			givenOpts: []domain.BalanceOption{domain.WithBalanceLastEntryID(customEntryID)},
			assertFn: func(t *testing.T, got *domain.Balance) {
				require.NotNil(t, got.LastEntryID)
				assert.Equal(t, customEntryID, *got.LastEntryID)
			},
		},
		{
			name:      "with_updated_at/success",
			givenOpts: []domain.BalanceOption{domain.WithBalanceUpdatedAt(customTime)},
			assertFn:  func(t *testing.T, got *domain.Balance) { assert.True(t, customTime.Equal(got.UpdatedAt)) },
		},
		{
			name:      "with_version/success",
			givenOpts: []domain.BalanceOption{domain.WithBalanceVersion(7)},
			assertFn:  func(t *testing.T, got *domain.Balance) { assert.Equal(t, int64(7), got.Version) },
		},
		{
			name: "last_option_wins/success",
			givenOpts: []domain.BalanceOption{
				domain.WithBalanceVersion(1),
				domain.WithBalanceVersion(2),
			},
			assertFn: func(t *testing.T, got *domain.Balance) { assert.Equal(t, int64(2), got.Version) },
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			faker := fakes.NewBalanceFaker()

			// act
			result := faker.Make(tc.givenOpts...)

			// assert
			require.NotNil(t, result)
			tc.assertFn(t, result)
		})
	}
}

func Test_BalanceFaker_Make_Defaults(t *testing.T) {
	// arrange
	faker := fakes.NewBalanceFaker()

	// act
	result := faker.Make()

	// assert
	require.NotNil(t, result)
	assert.False(t, xuid.IsEmpty(result.AccountID))
	assert.False(t, xuid.IsEmpty(result.LedgerID))
	assert.NotEmpty(t, result.DebitsPosted.Currency)
	assert.True(t, result.DebitsPosted.Amount.IsPositive())
	assert.True(t, result.DebitsPending.Amount.IsPositive())
	assert.True(t, result.CreditsPosted.Amount.IsPositive())
	assert.True(t, result.CreditsPending.Amount.IsPositive())
	assert.Equal(t, result.DebitsPosted.Currency, result.DebitsPending.Currency)
	assert.Equal(t, result.DebitsPosted.Currency, result.CreditsPosted.Currency)
	assert.Equal(t, result.DebitsPosted.Currency, result.CreditsPending.Currency)
	require.NotNil(t, result.LastEntryID)
	assert.False(t, xuid.IsEmpty(*result.LastEntryID))
	assert.GreaterOrEqual(t, result.Version, int64(0))
	assert.WithinDuration(t, time.Now().UTC(), result.UpdatedAt, 5*time.Second)
}

func Test_BalanceFaker_Make_GeneratesUniqueIDs(t *testing.T) {
	// arrange
	faker := fakes.NewBalanceFaker()

	// act
	first := faker.Make()
	second := faker.Make()

	// assert
	assert.NotEqual(t, first.AccountID, second.AccountID)
	assert.NotEqual(t, first.LedgerID, second.LedgerID)
	require.NotNil(t, first.LastEntryID)
	require.NotNil(t, second.LastEntryID)
	assert.NotEqual(t, *first.LastEntryID, *second.LastEntryID)
}

func Test_BalanceFaker_MakeMany_Count(t *testing.T) {
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
			faker := fakes.NewBalanceFaker()

			// act
			result := faker.MakeMany(tc.givenCount)

			// assert
			require.Len(t, result, tc.wantedLen)
			ids := make(map[xuid.XUID]struct{}, len(result))
			for _, b := range result {
				require.NotNil(t, b)
				ids[b.AccountID] = struct{}{}
			}
			assert.Len(t, ids, tc.wantedLen, "every balance should have a unique account id")
		})
	}
}

func Test_BalanceFaker_MakeMany_AppliesSharedOpts(t *testing.T) {
	// arrange
	faker := fakes.NewBalanceFaker()
	accountID := xuid.MustNewSortable(domain.AccountIDPrefix)
	ledgerID := xuid.MustNewSortable(domain.LedgerIDPrefix)

	// act
	result := faker.MakeMany(3,
		fakes.WithBalanceAccountID(accountID),
		fakes.WithBalanceLedgerID(ledgerID),
	)

	// assert
	require.Len(t, result, 3)
	for _, b := range result {
		assert.Equal(t, accountID, b.AccountID)
		assert.Equal(t, ledgerID, b.LedgerID)
	}
}

func Test_BalanceFaker_MakeMany_SharedIDOptCollides(t *testing.T) {
	// documents the known behavior: a pinned id is shared by every balance
	// arrange
	faker := fakes.NewBalanceFaker()
	accountID := xuid.MustNewSortable(domain.AccountIDPrefix)

	// act
	result := faker.MakeMany(3, fakes.WithBalanceAccountID(accountID))

	// assert
	require.Len(t, result, 3)
	for _, b := range result {
		assert.Equal(t, accountID, b.AccountID)
	}
}

func Test_BalanceOptions(t *testing.T) {
	id := xuid.MustNewSortable(domain.AccountIDPrefix)

	testCases := []struct {
		name     string
		givenOpt domain.BalanceOption
		assertFn func(t *testing.T, got *domain.Balance)
	}{
		{
			name:     "with_balance_account_id",
			givenOpt: fakes.WithBalanceAccountID(id),
			assertFn: func(t *testing.T, got *domain.Balance) { assert.Equal(t, id, got.AccountID) },
		},
		{
			name:     "with_balance_ledger_id",
			givenOpt: fakes.WithBalanceLedgerID(id),
			assertFn: func(t *testing.T, got *domain.Balance) { assert.Equal(t, id, got.LedgerID) },
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			balance := &domain.Balance{}

			// act
			tc.givenOpt(balance)

			// assert
			tc.assertFn(t, balance)
		})
	}
}

func assertMonetaryValueEqual(t *testing.T, want, got domain.MonetaryValue) {
	t.Helper()
	assert.True(t, want.Amount.Equal(got.Amount), "amount: want %s, got %s", want.Amount, got.Amount)
	assert.Equal(t, want.Currency, got.Currency)
}
