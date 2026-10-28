package domain_test

import (
	"ledger/internal/domain"
	"testing"
	"time"

	"github.com/47monad/xuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EntryMake(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		id1 := xuid.MustNewSortable("entry_")
		transactionID1 := xuid.MustNewSortable("transaction_")
		ledgerID1 := xuid.MustNewSortable("ledger_")
		accountID1 := xuid.MustNewSortable("account_")

		now := time.Now().UTC()
		customCreatedAt := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
		customEffectiveAt := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
		customLine := int32(3)
		description1 := "some_entry_description"
		metadata1 := map[string]string{"source": "checkout", "env": "test"}

		money1 := domain.MonetaryValue{Amount: decimal.NewFromInt(100), Currency: "USD"}
		fractionalMoney := domain.MonetaryValue{Amount: decimal.RequireFromString("12.34"), Currency: "EUR"}
		smallestMoney := domain.MonetaryValue{Amount: decimal.RequireFromString("0.01"), Currency: "USD"}

		testCases := []struct {
			name               string
			givenID            xuid.XUID
			givenTransactionID xuid.XUID
			givenLedgerID      xuid.XUID
			givenAccountID     xuid.XUID
			givenDirection     domain.Direction
			givenMoney         domain.MonetaryValue
			opts               []domain.EntryOption
			wanted             *domain.Entry
			wantedError        error
		}{
			// ---- success ----
			{
				name:               "valid_inputs/success",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionCredit,
					Line:          1,
					MonetaryValue: money1,
					EffectiveAt:   now,
					CreatedAt:     now,
				},
			},
			{
				name:               "debit_direction/success",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionDebit,
				givenMoney:         money1,
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionDebit,
					Line:          1,
					MonetaryValue: money1,
					EffectiveAt:   now,
					CreatedAt:     now,
				},
			},
			{
				name:               "fractional_amount/success",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionDebit,
				givenMoney:         fractionalMoney,
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionDebit,
					Line:          1,
					MonetaryValue: fractionalMoney,
					EffectiveAt:   now,
					CreatedAt:     now,
				},
			},
			{
				name:               "smallest_positive_amount/success",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         smallestMoney,
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionCredit,
					Line:          1,
					MonetaryValue: smallestMoney,
					EffectiveAt:   now,
					CreatedAt:     now,
				},
			},
			{
				name:               "with_line/success",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				opts:               []domain.EntryOption{domain.WithEntryLine(customLine)},
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionCredit,
					Line:          customLine,
					MonetaryValue: money1,
					EffectiveAt:   now,
					CreatedAt:     now,
				},
			},
			{
				name:               "with_effective_at/success",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				opts:               []domain.EntryOption{domain.WithEntryEffectiveAt(customEffectiveAt)},
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionCredit,
					Line:          1,
					MonetaryValue: money1,
					EffectiveAt:   customEffectiveAt,
					CreatedAt:     now,
				},
			},
			{
				name:               "with_description/success",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				opts:               []domain.EntryOption{domain.WithEntryDescription(description1)},
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionCredit,
					Line:          1,
					MonetaryValue: money1,
					Description:   description1,
					EffectiveAt:   now,
					CreatedAt:     now,
				},
			},
			{
				name:               "with_metadata/success",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				opts:               []domain.EntryOption{domain.WithEntryMetadata(metadata1)},
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionCredit,
					Line:          1,
					MonetaryValue: money1,
					Metadata:      metadata1,
					EffectiveAt:   now,
					CreatedAt:     now,
				},
			},
			{
				name:               "with_created_at/success",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				opts:               []domain.EntryOption{domain.WithEntryCreatedAt(customCreatedAt)},
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionCredit,
					Line:          1,
					MonetaryValue: money1,
					EffectiveAt:   now,
					CreatedAt:     customCreatedAt,
				},
			},
			{
				name:               "with_all_options/success",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionDebit,
				givenMoney:         fractionalMoney,
				opts: []domain.EntryOption{
					domain.WithEntryLine(customLine),
					domain.WithEntryEffectiveAt(customEffectiveAt),
					domain.WithEntryDescription(description1),
					domain.WithEntryMetadata(metadata1),
					domain.WithEntryCreatedAt(customCreatedAt),
				},
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionDebit,
					Line:          customLine,
					MonetaryValue: fractionalMoney,
					Description:   description1,
					Metadata:      metadata1,
					EffectiveAt:   customEffectiveAt,
					CreatedAt:     customCreatedAt,
				},
			},
			{
				name:               "with_repeated_line/last_option_wins",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				opts: []domain.EntryOption{
					domain.WithEntryLine(2),
					domain.WithEntryLine(customLine),
				},
				wanted: &domain.Entry{
					ID:            id1,
					TransactionID: transactionID1,
					LedgerID:      ledgerID1,
					AccountID:     accountID1,
					Direction:     domain.DirectionCredit,
					Line:          customLine,
					MonetaryValue: money1,
					EffectiveAt:   now,
					CreatedAt:     now,
				},
			},

			// ---- validation errors: each argument alone ----
			{
				name:               "empty_id/err_empty_entry_id",
				givenID:            xuid.XUID{},
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				wantedError:        domain.ErrEmptyEntryID,
			},
			{
				name:               "empty_transaction_id/err_empty_entry_transaction_id",
				givenID:            id1,
				givenTransactionID: xuid.XUID{},
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				wantedError:        domain.ErrEmptyEntryTransactionID,
			},
			{
				name:               "empty_ledger_id/err_empty_entry_ledger_id",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      xuid.XUID{},
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				wantedError:        domain.ErrEmptyEntryLedgerID,
			},
			{
				name:               "empty_account_id/err_empty_entry_account_id",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     xuid.XUID{},
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				wantedError:        domain.ErrEmptyEntryAccountID,
			},
			{
				name:               "empty_direction/err_empty_entry_direction",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     "",
				givenMoney:         money1,
				wantedError:        domain.ErrEmptyEntryDirection,
			},
			{
				name:               "unsupported_direction/Err_unsupported_entry_direction",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     "one_direction",
				givenMoney:         money1,
				wantedError:        domain.ErrUnsuppEntryDirection,
			},
			{
				name:               "zero_amount/err_invalid_entry_amount",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         domain.MonetaryValue{Amount: decimal.Zero, Currency: "USD"},
				wantedError:        domain.ErrInvalidEntryAmount,
			},
			{
				name:               "unset_amount/err_invalid_entry_amount",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         domain.MonetaryValue{Currency: "USD"},
				wantedError:        domain.ErrInvalidEntryAmount,
			},
			{
				name:               "negative_amount/err_invalid_entry_amount",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         domain.MonetaryValue{Amount: decimal.NewFromInt(-100), Currency: "USD"},
				wantedError:        domain.ErrInvalidEntryAmount,
			},
			{
				name:               "empty_currency/err_empty_entry_currency",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         domain.MonetaryValue{Amount: decimal.NewFromInt(100), Currency: ""},
				wantedError:        domain.ErrEmptyEntryCurrency,
			},

			// ---- validation errors: check order (id, transaction, ledger, account, direction, amount, currency) ----
			{
				name:               "empty_id_and_transaction_id/err_empty_entry_id",
				givenID:            xuid.XUID{},
				givenTransactionID: xuid.XUID{},
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				wantedError:        domain.ErrEmptyEntryID,
			},
			{
				name:               "empty_transaction_id_and_ledger_id/err_empty_entry_transaction_id",
				givenID:            id1,
				givenTransactionID: xuid.XUID{},
				givenLedgerID:      xuid.XUID{},
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				wantedError:        domain.ErrEmptyEntryTransactionID,
			},
			{
				name:               "empty_ledger_id_and_account_id/err_empty_entry_ledger_id",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      xuid.XUID{},
				givenAccountID:     xuid.XUID{},
				givenDirection:     domain.DirectionCredit,
				givenMoney:         money1,
				wantedError:        domain.ErrEmptyEntryLedgerID,
			},
			{
				name:               "empty_account_id_and_direction/err_empty_entry_account_id",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     xuid.XUID{},
				givenDirection:     "",
				givenMoney:         money1,
				wantedError:        domain.ErrEmptyEntryAccountID,
			},
			{
				name:               "empty_direction_and_zero_amount/err_empty_entry_direction",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     "",
				givenMoney:         domain.MonetaryValue{Amount: decimal.Zero, Currency: "USD"},
				wantedError:        domain.ErrEmptyEntryDirection,
			},
			{
				name:               "zero_amount_and_empty_currency/err_invalid_entry_amount",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         domain.MonetaryValue{Amount: decimal.Zero, Currency: ""},
				wantedError:        domain.ErrInvalidEntryAmount,
			},
			{
				name:               "all_empty/err_empty_entry_id",
				givenID:            xuid.XUID{},
				givenTransactionID: xuid.XUID{},
				givenLedgerID:      xuid.XUID{},
				givenAccountID:     xuid.XUID{},
				givenDirection:     "",
				givenMoney:         domain.MonetaryValue{},
				wantedError:        domain.ErrEmptyEntryID,
			},
			{
				name:               "with_valid_options/empty_currency/err_empty_entry_currency",
				givenID:            id1,
				givenTransactionID: transactionID1,
				givenLedgerID:      ledgerID1,
				givenAccountID:     accountID1,
				givenDirection:     domain.DirectionCredit,
				givenMoney:         domain.MonetaryValue{Amount: decimal.NewFromInt(100), Currency: ""},
				opts: []domain.EntryOption{
					domain.WithEntryLine(customLine),
					domain.WithEntryDescription(description1),
				},
				wantedError: domain.ErrEmptyEntryCurrency,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// act
				result, err := domain.NewEntryFactory().Make(
					tc.givenID, tc.givenTransactionID, tc.givenLedgerID, tc.givenAccountID,
					tc.givenDirection, tc.givenMoney, tc.opts...,
				)

				// assert
				if tc.wantedError != nil {
					require.ErrorIs(t, err, tc.wantedError)
					assert.Nil(t, result)
					return
				}

				require.NoError(t, err)
				require.NotNil(t, result)
				require.NotNil(t, tc.wanted)

				assert.Equal(t, tc.wanted.ID, result.ID)
				assert.Equal(t, tc.wanted.TransactionID, result.TransactionID)
				assert.Equal(t, tc.wanted.LedgerID, result.LedgerID)
				assert.Equal(t, tc.wanted.AccountID, result.AccountID)
				assert.Equal(t, tc.wanted.Direction, result.Direction)
				assert.Equal(t, tc.wanted.Line, result.Line)
				assert.Equal(t, tc.wanted.MonetaryValue.Currency, result.MonetaryValue.Currency)
				// decimal.Decimal must be compared with Equal, not assert.Equal:
				// equal values can have different internal representations (e.g. 100 vs 100.00).
				assert.True(t, tc.wanted.MonetaryValue.Amount.Equal(result.MonetaryValue.Amount),
					"amount: wanted %s, got %s", tc.wanted.MonetaryValue.Amount, result.MonetaryValue.Amount)
				assert.Equal(t, tc.wanted.Description, result.Description)
				assert.Equal(t, tc.wanted.Metadata, result.Metadata)
				assert.WithinDuration(t, tc.wanted.EffectiveAt, result.EffectiveAt, 1*time.Second)
				assert.WithinDuration(t, tc.wanted.CreatedAt, result.CreatedAt, 1*time.Second)
			})
		}
	})
}
