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

var validTransactionStatuses = []domain.TransactionStatus{
	domain.TransactionStatusPending,
	domain.TransactionStatusPosted,
	domain.TransactionStatusVoided,
}

func Test_TransactionFaker_Make_Overrides(t *testing.T) {
	customID := xuid.MustNewSortable(domain.TransactionIDPrefix)
	customLedgerID := xuid.MustNewSortable(domain.LedgerIDPrefix)
	customReversesID := xuid.MustNewSortable(domain.TransactionIDPrefix)
	customTime := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	customDescription := "custom description"
	customMetadata := map[string]string{"k1": "v1", "k2": "v2"}
	customEntries := []domain.Entry{
		{
			ID:            xuid.MustNewSortable(domain.EntryIDPrefix),
			TransactionID: customID,
			LedgerID:      customLedgerID,
			AccountID:     xuid.MustNewSortable(domain.AccountIDPrefix),
			Line:          1,
			Direction:     domain.DirectionDebit,
			MonetaryValue: domain.MonetaryValue{
				Amount:   decimal.NewFromInt(100),
				Currency: "USD",
			},
		},
	}

	testCases := []struct {
		name      string
		givenOpts []domain.TransactionOption
		assertFn  func(t *testing.T, got *domain.Transaction)
	}{
		{
			name:      "with_id/success",
			givenOpts: []domain.TransactionOption{fakes.WithTransactionID(customID)},
			assertFn:  func(t *testing.T, got *domain.Transaction) { assert.Equal(t, customID, got.ID) },
		},
		{
			name:      "with_ledger_id/success",
			givenOpts: []domain.TransactionOption{fakes.WithTransactionLedgerID(customLedgerID)},
			assertFn:  func(t *testing.T, got *domain.Transaction) { assert.Equal(t, customLedgerID, got.LedgerID) },
		},
		{
			name:      "with_idempotency_key/success",
			givenOpts: []domain.TransactionOption{fakes.WithTransactionIdempotencyKey("idem_key_1")},
			assertFn:  func(t *testing.T, got *domain.Transaction) { assert.Equal(t, "idem_key_1", got.IdempotencyKey) },
		},
		{
			name:      "with_reference/success",
			givenOpts: []domain.TransactionOption{fakes.WithTransactionReference("order_123")},
			assertFn:  func(t *testing.T, got *domain.Transaction) { assert.Equal(t, "order_123", got.Reference) },
		},
		{
			name:      "with_status_overrides_random/success",
			givenOpts: []domain.TransactionOption{domain.WithTransactionStatus(domain.TransactionStatusPosted)},
			assertFn: func(t *testing.T, got *domain.Transaction) {
				assert.Equal(t, domain.TransactionStatusPosted, got.Status)
			},
		},
		{
			name:      "with_description_overrides_random/success",
			givenOpts: []domain.TransactionOption{domain.WithTransactionDescription(customDescription)},
			assertFn:  func(t *testing.T, got *domain.Transaction) { assert.Equal(t, customDescription, got.Description) },
		},
		{
			name:      "with_effective_at/success",
			givenOpts: []domain.TransactionOption{domain.WithTransactionEffectiveAt(customTime)},
			assertFn:  func(t *testing.T, got *domain.Transaction) { assert.True(t, customTime.Equal(got.EffectiveAt)) },
		},
		{
			name:      "with_posted_at/success",
			givenOpts: []domain.TransactionOption{domain.WithTransactionPostedAt(customTime)},
			assertFn: func(t *testing.T, got *domain.Transaction) {
				require.NotNil(t, got.PostedAt)
				assert.True(t, customTime.Equal(*got.PostedAt))
			},
		},
		{
			name:      "with_reverses_id/success",
			givenOpts: []domain.TransactionOption{domain.WithTransactionReversesID(customReversesID)},
			assertFn: func(t *testing.T, got *domain.Transaction) {
				require.NotNil(t, got.ReversesID)
				assert.Equal(t, customReversesID, *got.ReversesID)
			},
		},
		{
			name:      "with_metadata_overrides_default/success",
			givenOpts: []domain.TransactionOption{domain.WithTransactionMetadata(customMetadata)},
			assertFn:  func(t *testing.T, got *domain.Transaction) { assert.Equal(t, customMetadata, got.Metadata) },
		},
		{
			name:      "with_entries/success",
			givenOpts: []domain.TransactionOption{domain.WithTransactionEntries(customEntries)},
			assertFn:  func(t *testing.T, got *domain.Transaction) { assert.Equal(t, customEntries, got.Entries) },
		},
		{
			name:      "with_created_at/success",
			givenOpts: []domain.TransactionOption{domain.WithTransactionCreatedAt(customTime)},
			assertFn:  func(t *testing.T, got *domain.Transaction) { assert.True(t, customTime.Equal(got.CreatedAt)) },
		},
		{
			name:      "with_updated_at/success",
			givenOpts: []domain.TransactionOption{domain.WithTransactionUpdatedAt(customTime)},
			assertFn:  func(t *testing.T, got *domain.Transaction) { assert.True(t, customTime.Equal(got.UpdatedAt)) },
		},
		{
			name: "last_option_wins/success",
			givenOpts: []domain.TransactionOption{
				domain.WithTransactionStatus(domain.TransactionStatusPosted),
				domain.WithTransactionStatus(domain.TransactionStatusVoided),
			},
			assertFn: func(t *testing.T, got *domain.Transaction) {
				assert.Equal(t, domain.TransactionStatusVoided, got.Status)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			faker := fakes.NewTransactionFaker()

			// act
			result := faker.Make(tc.givenOpts...)

			// assert
			require.NotNil(t, result)
			tc.assertFn(t, result)
		})
	}
}

func Test_TransactionFaker_Make_Defaults(t *testing.T) {
	// arrange
	faker := fakes.NewTransactionFaker()

	// act
	result := faker.Make()

	// assert
	require.NotNil(t, result)
	assert.False(t, xuid.IsEmpty(result.ID))
	assert.False(t, xuid.IsEmpty(result.LedgerID))
	assert.NotEmpty(t, result.IdempotencyKey)
	assert.NotEmpty(t, result.Reference)
	assert.True(t, slices.Contains(validTransactionStatuses, result.Status), "unexpected status %q", result.Status)
	assert.NotEmpty(t, result.Description)
	assert.GreaterOrEqual(t, len(result.Metadata), 1)
	assert.LessOrEqual(t, len(result.Metadata), 10)
	assert.Nil(t, result.PostedAt)
	assert.Nil(t, result.ReversesID)
	assert.Nil(t, result.Entries)
	assert.WithinDuration(t, time.Now().UTC(), result.EffectiveAt, 5*time.Second)
	assert.WithinDuration(t, time.Now().UTC(), result.CreatedAt, 5*time.Second)
	assert.WithinDuration(t, time.Now().UTC(), result.UpdatedAt, 5*time.Second)
}

func Test_TransactionFaker_Make_GeneratesUniqueIDs(t *testing.T) {
	// arrange
	faker := fakes.NewTransactionFaker()

	// act
	first := faker.Make()
	second := faker.Make()

	// assert
	assert.NotEqual(t, first.ID, second.ID)
	assert.NotEqual(t, first.LedgerID, second.LedgerID)
	assert.NotEqual(t, first.IdempotencyKey, second.IdempotencyKey)
	assert.NotEqual(t, first.Reference, second.Reference)
}

func Test_TransactionFaker_Make_CoversAllRandomizedValues(t *testing.T) {
	// probabilistic: with 300 draws the chance of missing any value is negligible
	// arrange
	faker := fakes.NewTransactionFaker()
	statuses := map[domain.TransactionStatus]struct{}{}

	// act
	for range 300 {
		tx := faker.Make()
		statuses[tx.Status] = struct{}{}
	}

	// assert
	assert.Len(t, statuses, len(validTransactionStatuses))
}

func Test_TransactionFaker_MakeMany_Count(t *testing.T) {
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
			faker := fakes.NewTransactionFaker()

			// act
			result := faker.MakeMany(tc.givenCount)

			// assert
			require.Len(t, result, tc.wantedLen)
			ids := make(map[xuid.XUID]struct{}, len(result))
			for _, tx := range result {
				require.NotNil(t, tx)
				ids[tx.ID] = struct{}{}
			}
			assert.Len(t, ids, tc.wantedLen, "every transaction should have a unique id")
		})
	}
}

func Test_TransactionFaker_MakeMany_AppliesSharedOpts(t *testing.T) {
	// arrange
	faker := fakes.NewTransactionFaker()
	ledgerID := xuid.MustNewSortable(domain.LedgerIDPrefix)

	// act
	result := faker.MakeMany(3,
		fakes.WithTransactionLedgerID(ledgerID),
		domain.WithTransactionStatus(domain.TransactionStatusPosted),
	)

	// assert
	require.Len(t, result, 3)
	for _, tx := range result {
		assert.Equal(t, ledgerID, tx.LedgerID)
		assert.Equal(t, domain.TransactionStatusPosted, tx.Status)
	}
}

func Test_TransactionFaker_MakeMany_SharedIDOptCollides(t *testing.T) {
	// documents the known behavior: a pinned ID is shared by every transaction
	// arrange
	faker := fakes.NewTransactionFaker()
	id := xuid.MustNewSortable(domain.TransactionIDPrefix)

	// act
	result := faker.MakeMany(3, fakes.WithTransactionID(id))

	// assert
	require.Len(t, result, 3)
	for _, tx := range result {
		assert.Equal(t, id, tx.ID)
	}
}

func Test_TransactionOptions(t *testing.T) {
	id := xuid.MustNewSortable(domain.TransactionIDPrefix)

	testCases := []struct {
		name     string
		givenOpt domain.TransactionOption
		assertFn func(t *testing.T, got *domain.Transaction)
	}{
		{
			name:     "with_transaction_id",
			givenOpt: fakes.WithTransactionID(id),
			assertFn: func(t *testing.T, got *domain.Transaction) { assert.Equal(t, id, got.ID) },
		},
		{
			name:     "with_transaction_ledger_id",
			givenOpt: fakes.WithTransactionLedgerID(id),
			assertFn: func(t *testing.T, got *domain.Transaction) { assert.Equal(t, id, got.LedgerID) },
		},
		{
			name:     "with_transaction_idempotency_key",
			givenOpt: fakes.WithTransactionIdempotencyKey("idem_key_1"),
			assertFn: func(t *testing.T, got *domain.Transaction) { assert.Equal(t, "idem_key_1", got.IdempotencyKey) },
		},
		{
			name:     "with_transaction_reference",
			givenOpt: fakes.WithTransactionReference("order_123"),
			assertFn: func(t *testing.T, got *domain.Transaction) { assert.Equal(t, "order_123", got.Reference) },
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			transaction := &domain.Transaction{}

			// act
			tc.givenOpt(transaction)

			// assert
			tc.assertFn(t, transaction)
		})
	}
}
