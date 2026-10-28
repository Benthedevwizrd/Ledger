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

func Test_TransactionMake(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		id1 := xuid.MustNewSortable("transaction_")
		ledgerID1 := xuid.MustNewSortable("ledger_")
		reversedID1 := xuid.MustNewSortable("transaction_")
		idempotencyKey1 := "idem_key_1"
		reference1 := "order_123"

		now := time.Now().UTC()
		customCreatedAt := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
		customUpdatedAt := time.Date(2024, 6, 7, 8, 9, 10, 0, time.UTC)
		customEffectiveAt := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
		customPostedAt := time.Date(2024, 1, 3, 4, 5, 6, 0, time.UTC)
		description1 := "some_transaction_description"
		metadata1 := map[string]string{"source": "checkout", "env": "test"}
		entries1 := []domain.Entry{
			{
				ID:            xuid.MustNewSortable("entry_"),
				TransactionID: id1,
				LedgerID:      ledgerID1,
				AccountID:     xuid.MustNewSortable("account_"),
				Line:          1,
				Direction:     domain.DirectionDebit,
				MonetaryValue: domain.MonetaryValue{
					Amount:   decimal.NewFromInt(100),
					Currency: "USD",
				},
			},
			{
				ID:            xuid.MustNewSortable("entry_"),
				TransactionID: id1,
				LedgerID:      ledgerID1,
				AccountID:     xuid.MustNewSortable("account_"),
				Line:          2,
				Direction:     domain.DirectionCredit,
				MonetaryValue: domain.MonetaryValue{
					Amount:   decimal.NewFromInt(100),
					Currency: "USD",
				},
			},
		}

		testCases := []struct {
			name                string
			givenID             xuid.XUID
			givenLedgerID       xuid.XUID
			givenIdempotencyKey string
			givenReference      string
			opts                []domain.TransactionOption
			wanted              *domain.Transaction
			wantedError         error
		}{
			// ---- success ----
			{
				name:                "valid_inputs/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Status:         domain.TransactionStatusPending,
					EffectiveAt:    now,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			{
				name:                "with_description/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts:                []domain.TransactionOption{domain.WithTransactionDescription(description1)},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Description:    description1,
					Status:         domain.TransactionStatusPending,
					EffectiveAt:    now,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			{
				name:                "with_status/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts:                []domain.TransactionOption{domain.WithTransactionStatus(domain.TransactionStatusPosted)},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Status:         domain.TransactionStatusPosted,
					EffectiveAt:    now,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			{
				name:                "with_effective_at/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts:                []domain.TransactionOption{domain.WithTransactionEffectiveAt(customEffectiveAt)},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Status:         domain.TransactionStatusPending,
					EffectiveAt:    customEffectiveAt,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			{
				name:                "with_posted_at/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts:                []domain.TransactionOption{domain.WithTransactionPostedAt(customPostedAt)},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Status:         domain.TransactionStatusPending,
					EffectiveAt:    now,
					PostedAt:       &customPostedAt,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			{
				name:                "with_reverses_id/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts:                []domain.TransactionOption{domain.WithTransactionReversesID(reversedID1)},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Status:         domain.TransactionStatusPending,
					EffectiveAt:    now,
					ReversesID:     &reversedID1,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			{
				name:                "with_metadata/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts:                []domain.TransactionOption{domain.WithTransactionMetadata(metadata1)},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Status:         domain.TransactionStatusPending,
					EffectiveAt:    now,
					Metadata:       metadata1,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			{
				name:                "with_entries/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts:                []domain.TransactionOption{domain.WithTransactionEntries(entries1)},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Status:         domain.TransactionStatusPending,
					EffectiveAt:    now,
					Entries:        entries1,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},
			{
				name:                "with_created_at/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts:                []domain.TransactionOption{domain.WithTransactionCreatedAt(customCreatedAt)},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Status:         domain.TransactionStatusPending,
					EffectiveAt:    now,
					CreatedAt:      customCreatedAt,
					UpdatedAt:      now,
				},
			},
			{
				name:                "with_updated_at/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts:                []domain.TransactionOption{domain.WithTransactionUpdatedAt(customUpdatedAt)},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Status:         domain.TransactionStatusPending,
					EffectiveAt:    now,
					CreatedAt:      now,
					UpdatedAt:      customUpdatedAt,
				},
			},
			{
				name:                "with_all_options/success",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts: []domain.TransactionOption{
					domain.WithTransactionDescription(description1),
					domain.WithTransactionStatus(domain.TransactionStatusPosted),
					domain.WithTransactionEffectiveAt(customEffectiveAt),
					domain.WithTransactionPostedAt(customPostedAt),
					domain.WithTransactionReversesID(reversedID1),
					domain.WithTransactionMetadata(metadata1),
					domain.WithTransactionEntries(entries1),
					domain.WithTransactionCreatedAt(customCreatedAt),
					domain.WithTransactionUpdatedAt(customUpdatedAt),
				},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Description:    description1,
					Status:         domain.TransactionStatusPosted,
					EffectiveAt:    customEffectiveAt,
					PostedAt:       &customPostedAt,
					ReversesID:     &reversedID1,
					Metadata:       metadata1,
					Entries:        entries1,
					CreatedAt:      customCreatedAt,
					UpdatedAt:      customUpdatedAt,
				},
			},
			{
				name:                "with_repeated_status/last_option_wins",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				opts: []domain.TransactionOption{
					domain.WithTransactionStatus(domain.TransactionStatusPosted),
					domain.WithTransactionStatus(domain.TransactionStatusVoided),
				},
				wanted: &domain.Transaction{
					ID:             id1,
					LedgerID:       ledgerID1,
					IdempotencyKey: idempotencyKey1,
					Reference:      reference1,
					Status:         domain.TransactionStatusVoided,
					EffectiveAt:    now,
					CreatedAt:      now,
					UpdatedAt:      now,
				},
			},

			// ---- validation errors ----
			{
				name:                "empty_id/err_empty_transaction_id",
				givenID:             xuid.XUID{},
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				wantedError:         domain.ErrEmptyTransactonID,
			},
			{
				name:                "empty_ledger_id/err_empty_transaction_ledger_id",
				givenID:             id1,
				givenLedgerID:       xuid.XUID{},
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				wantedError:         domain.ErrEmptyTransactionLedgerID,
			},
			{
				name:                "empty_idempotency_key/err_empty_transaction_idempotency_key",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: "",
				givenReference:      reference1,
				wantedError:         domain.ErrEmptyTransactionIdempotencyKey,
			},
			{
				name:                "empty_reference/err_empty_transaction_reference",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      "",
				wantedError:         domain.ErrEmptyTransactionReferece,
			},
			{
				name:                "empty_id_and_ledger_id/err_empty_transaction_id",
				givenID:             xuid.XUID{},
				givenLedgerID:       xuid.XUID{},
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      reference1,
				wantedError:         domain.ErrEmptyTransactonID,
			},
			{
				name:                "empty_ledger_id_and_idempotency_key/err_empty_transaction_ledger_id",
				givenID:             id1,
				givenLedgerID:       xuid.XUID{},
				givenIdempotencyKey: "",
				givenReference:      reference1,
				wantedError:         domain.ErrEmptyTransactionLedgerID,
			},
			{
				name:                "empty_idempotency_key_and_reference/err_empty_transaction_idempotency_key",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: "",
				givenReference:      "",
				wantedError:         domain.ErrEmptyTransactionIdempotencyKey,
			},
			{
				name:                "all_empty/err_empty_transaction_id",
				givenID:             xuid.XUID{},
				givenLedgerID:       xuid.XUID{},
				givenIdempotencyKey: "",
				givenReference:      "",
				wantedError:         domain.ErrEmptyTransactonID,
			},
			{
				name:                "with_valid_options/empty_reference/err_empty_transaction_reference",
				givenID:             id1,
				givenLedgerID:       ledgerID1,
				givenIdempotencyKey: idempotencyKey1,
				givenReference:      "",
				opts: []domain.TransactionOption{
					domain.WithTransactionDescription(description1),
					domain.WithTransactionMetadata(metadata1),
				},
				wantedError: domain.ErrEmptyTransactionReferece,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// act
				result, err := domain.NewTransactionFactory().Make(tc.givenID, tc.givenLedgerID, tc.givenIdempotencyKey, tc.givenReference, tc.opts...)

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
				assert.Equal(t, tc.wanted.LedgerID, result.LedgerID)
				assert.Equal(t, tc.wanted.IdempotencyKey, result.IdempotencyKey)
				assert.Equal(t, tc.wanted.Reference, result.Reference)
				assert.Equal(t, tc.wanted.Description, result.Description)
				assert.Equal(t, tc.wanted.Status, result.Status)
				assert.Equal(t, tc.wanted.PostedAt, result.PostedAt)
				assert.Equal(t, tc.wanted.ReversesID, result.ReversesID)
				assert.Equal(t, tc.wanted.Metadata, result.Metadata)
				assert.Equal(t, tc.wanted.Entries, result.Entries)
				assert.WithinDuration(t, tc.wanted.EffectiveAt, result.EffectiveAt, 1*time.Second)
				assert.WithinDuration(t, tc.wanted.CreatedAt, result.CreatedAt, 1*time.Second)
				assert.WithinDuration(t, tc.wanted.UpdatedAt, result.UpdatedAt, 1*time.Second)
			})
		}
	})
}
