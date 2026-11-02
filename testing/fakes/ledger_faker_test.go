package fakes_test

import (
	"testing"
	"time"

	"ledger/internal/domain"
	"ledger/testing/fakes"

	"github.com/47monad/xuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LedgerFaker_Make_Overrides(t *testing.T) {
	customID := xuid.MustNewSortable(domain.LedgerIDPrefix)
	customTime := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	customMetadata := map[string]string{"k1": "v1", "k2": "v2"}

	testCases := []struct {
		name      string
		givenOpts []domain.LedgerOption
		assertFn  func(t *testing.T, got *domain.Ledger)
	}{
		{
			name:      "with_id/success",
			givenOpts: []domain.LedgerOption{fakes.WithLedgerID(customID)},
			assertFn: func(t *testing.T, got *domain.Ledger) {
				assert.Equal(t, customID, got.ID)
			},
		},
		{
			name:      "with_name/success",
			givenOpts: []domain.LedgerOption{fakes.WithLedgerName("custom name")},
			assertFn: func(t *testing.T, got *domain.Ledger) {
				assert.Equal(t, "custom name", got.Name)
			},
		},
		{
			name:      "with_description/success",
			givenOpts: []domain.LedgerOption{fakes.WithLedgerDescription("custom description")},
			assertFn: func(t *testing.T, got *domain.Ledger) {
				assert.Equal(t, "custom description", got.Description)
			},
		},
		{
			name:      "with_metadata_overrides_default/success",
			givenOpts: []domain.LedgerOption{domain.WithLedgerMetadata(customMetadata)},
			assertFn: func(t *testing.T, got *domain.Ledger) {
				assert.Equal(t, customMetadata, got.Metadata)
			},
		},
		{
			name: "with_timestamps_override_defaults/success",
			givenOpts: []domain.LedgerOption{
				domain.WithLedgerCreatedAt(customTime),
				domain.WithLedgerUpdatedAt(customTime),
			},
			assertFn: func(t *testing.T, got *domain.Ledger) {
				assert.True(t, customTime.Equal(got.CreatedAt))
				assert.True(t, customTime.Equal(got.UpdatedAt))
			},
		},
		{
			name: "last_option_wins/success",
			givenOpts: []domain.LedgerOption{
				fakes.WithLedgerName("first"),
				fakes.WithLedgerName("second"),
			},
			assertFn: func(t *testing.T, got *domain.Ledger) {
				assert.Equal(t, "second", got.Name)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			faker := fakes.NewLedgerFaker()

			// act
			result := faker.Make(tc.givenOpts...)

			// assert
			require.NotNil(t, result)
			tc.assertFn(t, result)
		})
	}
}

func Test_LedgerFaker_Make_Defaults(t *testing.T) {
	// arrange
	faker := fakes.NewLedgerFaker()
	before := time.Now().UTC().Add(-time.Second)

	// act
	result := faker.Make()

	// assert
	require.NotNil(t, result)
	assert.False(t, xuid.IsEmpty(result.ID))
	assert.NotEmpty(t, result.Name)
	assert.NotEmpty(t, result.Description)
	assert.GreaterOrEqual(t, len(result.Metadata), 1)
	assert.LessOrEqual(t, len(result.Metadata), 10)
	assert.True(t, result.CreatedAt.Equal(result.UpdatedAt))
	assert.WithinDuration(t, time.Now().UTC(), result.CreatedAt, 5*time.Second)
	assert.False(t, result.CreatedAt.Before(before))
}

func Test_LedgerFaker_Make_GeneratesUniqueIDs(t *testing.T) {
	// arrange
	faker := fakes.NewLedgerFaker()

	// act
	first := faker.Make()
	second := faker.Make()

	// assert
	assert.NotEqual(t, first.ID, second.ID)
}

func Test_LedgerFaker_MakeMany_Count(t *testing.T) {
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
			faker := fakes.NewLedgerFaker()

			// act
			result := faker.MakeMany(tc.givenCount)

			// assert
			require.Len(t, result, tc.wantedLen)
			ids := make(map[xuid.XUID]struct{}, len(result))
			for _, l := range result {
				require.NotNil(t, l)
				ids[l.ID] = struct{}{}
			}
			assert.Len(t, ids, tc.wantedLen, "every ledger should have a unique id")
		})
	}
}

func Test_LedgerFaker_MakeMany_AppliesSharedOpts(t *testing.T) {
	// arrange
	faker := fakes.NewLedgerFaker()

	// act
	result := faker.MakeMany(3, fakes.WithLedgerName("shared name"))

	// assert
	require.Len(t, result, 3)
	for _, l := range result {
		assert.Equal(t, "shared name", l.Name)
	}
}

func Test_LedgerFaker_MakeMany_SharedIDOptCollides(t *testing.T) {
	// documents the known behavior: a pinned ID is shared by every ledger
	// arrange
	faker := fakes.NewLedgerFaker()
	id := xuid.MustNewSortable(domain.LedgerIDPrefix)

	// act
	result := faker.MakeMany(3, fakes.WithLedgerID(id))

	// assert
	require.Len(t, result, 3)
	for _, l := range result {
		assert.Equal(t, id, l.ID)
	}
}

func Test_LedgerOptions(t *testing.T) {
	id := xuid.MustNewSortable(domain.LedgerIDPrefix)

	testCases := []struct {
		name     string
		givenOpt domain.LedgerOption
		assertFn func(t *testing.T, got *domain.Ledger)
	}{
		{
			name:     "with_ledger_id",
			givenOpt: fakes.WithLedgerID(id),
			assertFn: func(t *testing.T, got *domain.Ledger) { assert.Equal(t, id, got.ID) },
		},
		{
			name:     "with_ledger_name",
			givenOpt: fakes.WithLedgerName("n"),
			assertFn: func(t *testing.T, got *domain.Ledger) { assert.Equal(t, "n", got.Name) },
		},
		{
			name:     "with_ledger_description",
			givenOpt: fakes.WithLedgerDescription("d"),
			assertFn: func(t *testing.T, got *domain.Ledger) { assert.Equal(t, "d", got.Description) },
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange
			ledger := &domain.Ledger{}

			// act
			tc.givenOpt(ledger)

			// assert
			tc.assertFn(t, ledger)
		})
	}
}
