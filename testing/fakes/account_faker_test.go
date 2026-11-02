package fakes_test

import (
	"slices"
	"testing"

	"ledger/internal/domain"
	"ledger/testing/fakes"

	"github.com/47monad/xuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	validAccountTypes = []domain.AccountType{
		domain.AccountTypeAsset,
		domain.AccountTypeEquity,
		domain.AccountTypeExpense,
		domain.AccountTypeLiability,
		domain.AccountTypeRevenue,
	}
	validDirections = []domain.Direction{
		domain.DirectionCredit,
		domain.DirectionDebit,
	}
	validAccountStatuses = []domain.AccountStatus{
		domain.AccountStatusActive,
		domain.AccountStatusClosed,
		domain.AccountStatusFrozen,
	}
)

func Test_AccountFaker_Make_Overrides(t *testing.T) {
	customID := xuid.MustNewSortable(domain.AccountIDPrefix)
	customLedgerID := xuid.MustNewSortable(domain.LedgerIDPrefix)
	customMetadata := map[string]string{"k1": "v1", "k2": "v2"}

	testCases := []struct {
		name      string
		givenOpts []domain.AccountOption
		assertFn  func(t *testing.T, got *domain.Account)
	}{
		{
			name:      "with_id/success",
			givenOpts: []domain.AccountOption{fakes.WithAccountID(customID)},
			assertFn:  func(t *testing.T, got *domain.Account) { assert.Equal(t, customID, got.ID) },
		},
		{
			name:      "with_ledger_id/success",
			givenOpts: []domain.AccountOption{fakes.WithAccountLedgerID(customLedgerID)},
			assertFn:  func(t *testing.T, got *domain.Account) { assert.Equal(t, customLedgerID, got.LedgerID) },
		},
		{
			name:      "with_code/success",
			givenOpts: []domain.AccountOption{fakes.WithAccountCode("1001")},
			assertFn:  func(t *testing.T, got *domain.Account) { assert.Equal(t, "1001", got.Code) },
		},
		{
			name:      "with_name/success",
			givenOpts: []domain.AccountOption{fakes.WithAccountName("cash")},
			assertFn:  func(t *testing.T, got *domain.Account) { assert.Equal(t, "cash", got.Name) },
		},
		{
			name:      "with_type_overrides_random/success",
			givenOpts: []domain.AccountOption{domain.WithAccountType(domain.AccountTypeAsset)},
			assertFn: func(t *testing.T, got *domain.Account) {
				assert.Equal(t, domain.AccountTypeAsset, got.Type)
			},
		},
		{
			name:      "with_normal_balance_overrides_random/success",
			givenOpts: []domain.AccountOption{domain.WithAccountNormalBalance(domain.DirectionDebit)},
			assertFn: func(t *testing.T, got *domain.Account) {
				assert.Equal(t, domain.DirectionDebit, got.NormalBalance)
			},
		},
		{
			name:      "with_status_overrides_random/success",
			givenOpts: []domain.AccountOption{domain.WithAccountStatus(domain.AccountStatusActive)},
			assertFn: func(t *testing.T, got *domain.Account) {
				assert.Equal(t, domain.AccountStatusActive, got.Status)
			},
		},
		{
			name:      "with_metadata_overrides_default/success",
			givenOpts: []domain.AccountOption{domain.WithAccountMetaData(customMetadata)},
			assertFn: func(t *testing.T, got *domain.Account) {
				assert.Equal(t, customMetadata, got.Metadata)
			},
		},
		{
			name: "last_option_wins/success",
			givenOpts: []domain.AccountOption{
				fakes.WithAccountName("first"),
				fakes.WithAccountName("second"),
			},
			assertFn: func(t *testing.T, got *domain.Account) { assert.Equal(t, "second", got.Name) },
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			faker := fakes.NewAccountFaker()

			// act
			result := faker.Make(tc.givenOpts...)

			// assert
			require.NotNil(t, result)
			tc.assertFn(t, result)
		})
	}
}

func Test_AccountFaker_Make_Defaults(t *testing.T) {
	// arrange
	faker := fakes.NewAccountFaker()

	// act
	result := faker.Make()

	// assert
	require.NotNil(t, result)
	assert.False(t, xuid.IsEmpty(result.ID))
	assert.False(t, xuid.IsEmpty(result.LedgerID))
	assert.NotEmpty(t, result.Code)
	assert.NotEmpty(t, result.Name)
	assert.True(t, slices.Contains(validAccountTypes, result.Type), "unexpected type %q", result.Type)
	assert.True(t, slices.Contains(validDirections, result.NormalBalance), "unexpected direction %q", result.NormalBalance)
	assert.True(t, slices.Contains(validAccountStatuses, result.Status), "unexpected status %q", result.Status)
	assert.GreaterOrEqual(t, len(result.Metadata), 1)
	assert.LessOrEqual(t, len(result.Metadata), 10)
}

func Test_AccountFaker_Make_GeneratesUniqueIDs(t *testing.T) {
	// arrange
	faker := fakes.NewAccountFaker()

	// act
	first := faker.Make()
	second := faker.Make()

	// assert
	assert.NotEqual(t, first.ID, second.ID)
	assert.NotEqual(t, first.LedgerID, second.LedgerID)
}

func Test_AccountFaker_Make_CoversAllRandomizedValues(t *testing.T) {
	// probabilistic: with 300 draws the chance of missing any value is negligible
	// arrange
	faker := fakes.NewAccountFaker()
	types := map[domain.AccountType]struct{}{}
	directions := map[domain.Direction]struct{}{}
	statuses := map[domain.AccountStatus]struct{}{}

	// act
	for range 300 {
		a := faker.Make()
		types[a.Type] = struct{}{}
		directions[a.NormalBalance] = struct{}{}
		statuses[a.Status] = struct{}{}
	}

	// assert
	assert.Len(t, types, len(validAccountTypes))
	assert.Len(t, directions, len(validDirections))
	assert.Len(t, statuses, len(validAccountStatuses))
}

func Test_AccountFaker_MakeMany_Count(t *testing.T) {
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
			faker := fakes.NewAccountFaker()

			// act
			result := faker.MakeMany(tc.givenCount)

			// assert
			require.Len(t, result, tc.wantedLen)
			ids := make(map[xuid.XUID]struct{}, len(result))
			for _, a := range result {
				require.NotNil(t, a)
				ids[a.ID] = struct{}{}
			}
			assert.Len(t, ids, tc.wantedLen, "every account should have a unique id")
		})
	}
}

func Test_AccountFaker_MakeMany_AppliesSharedOpts(t *testing.T) {
	// arrange
	faker := fakes.NewAccountFaker()
	ledgerID := xuid.MustNewSortable(domain.LedgerIDPrefix)

	// act
	result := faker.MakeMany(3, fakes.WithAccountLedgerID(ledgerID))

	// assert
	require.Len(t, result, 3)
	for _, a := range result {
		assert.Equal(t, ledgerID, a.LedgerID)
	}
}

func Test_AccountFaker_MakeMany_SharedIDOptCollides(t *testing.T) {
	// documents the known behavior: a pinned ID is shared by every account
	// arrange
	faker := fakes.NewAccountFaker()
	id := xuid.MustNewSortable(domain.AccountIDPrefix)

	// act
	result := faker.MakeMany(3, fakes.WithAccountID(id))

	// assert
	require.Len(t, result, 3)
	for _, a := range result {
		assert.Equal(t, id, a.ID)
	}
}

func Test_AccountOptions(t *testing.T) {
	id := xuid.MustNewSortable(domain.AccountIDPrefix)
	ledgerID := xuid.MustNewSortable(domain.LedgerIDPrefix)

	testCases := []struct {
		name     string
		givenOpt domain.AccountOption
		assertFn func(t *testing.T, got *domain.Account)
	}{
		{
			name:     "with_account_id",
			givenOpt: fakes.WithAccountID(id),
			assertFn: func(t *testing.T, got *domain.Account) { assert.Equal(t, id, got.ID) },
		},
		{
			name:     "with_account_ledger_id",
			givenOpt: fakes.WithAccountLedgerID(ledgerID),
			assertFn: func(t *testing.T, got *domain.Account) { assert.Equal(t, ledgerID, got.LedgerID) },
		},
		{
			name:     "with_account_code",
			givenOpt: fakes.WithAccountCode("c"),
			assertFn: func(t *testing.T, got *domain.Account) { assert.Equal(t, "c", got.Code) },
		},
		{
			name:     "with_account_name",
			givenOpt: fakes.WithAccountName("n"),
			assertFn: func(t *testing.T, got *domain.Account) { assert.Equal(t, "n", got.Name) },
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			account := &domain.Account{}

			// act
			tc.givenOpt(account)

			// assert
			tc.assertFn(t, account)
		})
	}
}
