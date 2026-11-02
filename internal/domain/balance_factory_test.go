package domain_test

import (
	"testing"
	"time"

	"ledger/internal/domain"

	"github.com/47monad/xuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_MakeBalance(t *testing.T) {
	accountID1 := xuid.MustNewSortable("balance_")
	ledgerID1 := xuid.MustNewSortable(domain.LedgerIDPrefix)
	entryID1 := xuid.MustNewSortable("entry_") // adjust prefix to your entry ID prefix constant
	updatedAt := time.Date(2026, time.November, 1, 12, 0, 0, 0, time.UTC)

	// Adapt to however MonetaryValue is constructed in your domain package.
	debitsPosted := newMonetaryValue(t, 100, "USD")
	debitsPending := newMonetaryValue(t, 50, "USD")
	creditsPosted := newMonetaryValue(t, 200, "USD")
	creditsPending := newMonetaryValue(t, 25, "USD")

	var emptyID xuid.XUID

	testCases := []struct {
		name           string
		givenAccountID xuid.XUID
		givenLedgerID  xuid.XUID
		givenOpts      []domain.BalanceOption
		wanted         *domain.Balance
		wantedError    error
	}{
		{
			name:           "valid_input/success",
			givenAccountID: accountID1,
			givenLedgerID:  ledgerID1,
			wanted: &domain.Balance{
				AccountID: accountID1,
				LedgerID:  ledgerID1,
			},
		},
		{
			name:           "empty_account_id/error",
			givenAccountID: emptyID,
			givenLedgerID:  ledgerID1,
			wantedError:    domain.ErrEmptyBalanceAccountID,
		},
		{
			name:           "empty_ledger_id/error",
			givenAccountID: accountID1,
			givenLedgerID:  emptyID,
			wantedError:    domain.ErrEmptyBalanceLedgerID,
		},
		{
			name:           "both_ids_empty/account_id_error_first",
			givenAccountID: emptyID,
			givenLedgerID:  emptyID,
			wantedError:    domain.ErrEmptyBalanceAccountID,
		},
		{
			name:           "with_debits_posted/success",
			givenAccountID: accountID1,
			givenLedgerID:  ledgerID1,
			givenOpts:      []domain.BalanceOption{domain.WithBalanceDebitsPosted(debitsPosted)},
			wanted: &domain.Balance{
				AccountID:    accountID1,
				LedgerID:     ledgerID1,
				DebitsPosted: debitsPosted,
			},
		},
		{
			name:           "with_debits_pending/success",
			givenAccountID: accountID1,
			givenLedgerID:  ledgerID1,
			givenOpts:      []domain.BalanceOption{domain.WithBalanceDebitsPending(debitsPending)},
			wanted: &domain.Balance{
				AccountID:     accountID1,
				LedgerID:      ledgerID1,
				DebitsPending: debitsPending,
			},
		},
		{
			name:           "with_credits_posted/success",
			givenAccountID: accountID1,
			givenLedgerID:  ledgerID1,
			givenOpts:      []domain.BalanceOption{domain.WithBalanceCreditsPosted(creditsPosted)},
			wanted: &domain.Balance{
				AccountID:     accountID1,
				LedgerID:      ledgerID1,
				CreditsPosted: creditsPosted,
			},
		},
		{
			name:           "with_credits_pending/success",
			givenAccountID: accountID1,
			givenLedgerID:  ledgerID1,
			givenOpts:      []domain.BalanceOption{domain.WithBalanceCreditsPending(creditsPending)},
			wanted: &domain.Balance{
				AccountID:      accountID1,
				LedgerID:       ledgerID1,
				CreditsPending: creditsPending,
			},
		},
		{
			name:           "with_last_entry_id/success",
			givenAccountID: accountID1,
			givenLedgerID:  ledgerID1,
			givenOpts:      []domain.BalanceOption{domain.WithBalanceLastEntryID(entryID1)},
			wanted: &domain.Balance{
				AccountID:   accountID1,
				LedgerID:    ledgerID1,
				LastEntryID: &entryID1,
			},
		},
		{
			name:           "with_updated_at/success",
			givenAccountID: accountID1,
			givenLedgerID:  ledgerID1,
			givenOpts:      []domain.BalanceOption{domain.WithBalanceUpdatedAt(updatedAt)},
			wanted: &domain.Balance{
				AccountID: accountID1,
				LedgerID:  ledgerID1,
				UpdatedAt: updatedAt,
			},
		},
		{
			name:           "with_version/success",
			givenAccountID: accountID1,
			givenLedgerID:  ledgerID1,
			givenOpts:      []domain.BalanceOption{domain.WithBalanceVersion(7)},
			wanted: &domain.Balance{
				AccountID: accountID1,
				LedgerID:  ledgerID1,
				Version:   7,
			},
		},
		{
			name:           "with_all_options/success",
			givenAccountID: accountID1,
			givenLedgerID:  ledgerID1,
			givenOpts: []domain.BalanceOption{
				domain.WithBalanceDebitsPosted(debitsPosted),
				domain.WithBalanceDebitsPending(debitsPending),
				domain.WithBalanceCreditsPosted(creditsPosted),
				domain.WithBalanceCreditsPending(creditsPending),
				domain.WithBalanceLastEntryID(entryID1),
				domain.WithBalanceUpdatedAt(updatedAt),
				domain.WithBalanceVersion(3),
			},
			wanted: &domain.Balance{
				AccountID:      accountID1,
				LedgerID:       ledgerID1,
				DebitsPosted:   debitsPosted,
				DebitsPending:  debitsPending,
				CreditsPosted:  creditsPosted,
				CreditsPending: creditsPending,
				LastEntryID:    &entryID1,
				UpdatedAt:      updatedAt,
				Version:        3,
			},
		},
		{
			name:           "last_option_wins/success",
			givenAccountID: accountID1,
			givenLedgerID:  ledgerID1,
			givenOpts: []domain.BalanceOption{
				domain.WithBalanceVersion(1),
				domain.WithBalanceVersion(2),
			},
			wanted: &domain.Balance{
				AccountID: accountID1,
				LedgerID:  ledgerID1,
				Version:   2,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			factory := domain.NewBalanceFactory()

			// act
			result, err := factory.Make(tc.givenAccountID, tc.givenLedgerID, tc.givenOpts...)

			// assert
			if tc.wantedError != nil {
				assert.ErrorIs(t, err, tc.wantedError)
				assert.Nil(t, result)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, result)

			assert.Equal(t, tc.wanted.AccountID, result.AccountID)
			assert.Equal(t, tc.wanted.LedgerID, result.LedgerID)

			assert.True(t, tc.wanted.DebitsPosted.Amount.Equal(result.DebitsPosted.Amount),
				"DebitsPosted: want %s, got %s", tc.wanted.DebitsPosted.Amount, result.DebitsPosted.Amount)
			assert.Equal(t, tc.wanted.DebitsPosted.Currency, result.DebitsPosted.Currency)

			assert.True(t, tc.wanted.DebitsPending.Amount.Equal(result.DebitsPending.Amount),
				"DebitsPending: want %s, got %s", tc.wanted.DebitsPending.Amount, result.DebitsPending.Amount)
			assert.Equal(t, tc.wanted.DebitsPending.Currency, result.DebitsPending.Currency)

			assert.True(t, tc.wanted.CreditsPosted.Amount.Equal(result.CreditsPosted.Amount),
				"CreditsPosted: want %s, got %s", tc.wanted.CreditsPosted.Amount, result.CreditsPosted.Amount)
			assert.Equal(t, tc.wanted.CreditsPosted.Currency, result.CreditsPosted.Currency)

			assert.True(t, tc.wanted.CreditsPending.Amount.Equal(result.CreditsPending.Amount),
				"CreditsPending: want %s, got %s", tc.wanted.CreditsPending.Amount, result.CreditsPending.Amount)
			assert.Equal(t, tc.wanted.CreditsPending.Currency, result.CreditsPending.Currency)

			assert.Equal(t, tc.wanted.LastEntryID, result.LastEntryID)
			assert.True(t, tc.wanted.UpdatedAt.Equal(result.UpdatedAt),
				"UpdatedAt: want %s, got %s", tc.wanted.UpdatedAt, result.UpdatedAt)
			assert.Equal(t, tc.wanted.Version, result.Version)
		})
	}
}

// newMonetaryValue is a placeholder: replace the body with your real constructor.
func newMonetaryValue(t *testing.T, amount int64, currency string) domain.MonetaryValue {
	t.Helper()
	// e.g. return domain.MustNewMonetaryValue(amount)
	return domain.MonetaryValue{
		Amount:   decimal.NewFromInt(amount),
		Currency: currency,
	}
}
