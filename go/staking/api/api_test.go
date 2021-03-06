package api

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oasisprotocol/oasis-core/go/common/quantity"
)

func TestConsensusParameters(t *testing.T) {
	require := require.New(t)

	// Default consensus parameters.
	var emptyParams ConsensusParameters
	require.Error(emptyParams.SanityCheck(), "default consensus parameters should be invalid")

	// Valid thresholds.
	validThresholds := map[ThresholdKind]quantity.Quantity{
		KindEntity:            *quantity.NewQuantity(),
		KindNodeValidator:     *quantity.NewQuantity(),
		KindNodeCompute:       *quantity.NewQuantity(),
		KindNodeStorage:       *quantity.NewQuantity(),
		KindNodeKeyManager:    *quantity.NewQuantity(),
		KindRuntimeCompute:    *quantity.NewQuantity(),
		KindRuntimeKeyManager: *quantity.NewQuantity(),
	}
	validThresholdsParams := ConsensusParameters{
		Thresholds:         validThresholds,
		FeeSplitWeightVote: mustInitQuantity(t, 1),
	}
	require.NoError(validThresholdsParams.SanityCheck(), "consensus parameters with valid thresholds should be valid")

	// NOTE: There is currently no way to construct invalid thresholds.

	// Degenerate fee split.
	degenerateFeeSplit := ConsensusParameters{
		Thresholds:                validThresholds,
		FeeSplitWeightPropose:     mustInitQuantity(t, 0),
		FeeSplitWeightVote:        mustInitQuantity(t, 0),
		FeeSplitWeightNextPropose: mustInitQuantity(t, 0),
	}
	require.Error(degenerateFeeSplit.SanityCheck(), "consensus parameters with degenerate fee split should be invalid")
}

func TestThresholdKind(t *testing.T) {
	require := require.New(t)

	for k := ThresholdKind(0); k <= KindMax; k++ {
		enc, err := k.MarshalText()
		require.NoError(err, "MarshalText")

		var d ThresholdKind
		err = d.UnmarshalText(enc)
		require.NoError(err, "UnmarshalText")
		require.Equal(k, d, "threshold kind should round-trip")
	}
}

func TestStakeThreshold(t *testing.T) {
	require := require.New(t)

	// Empty stake threshold is invalid.
	st := StakeThreshold{}
	_, err := st.Value(nil)
	require.Error(err, "empty stake threshold is invalid")

	// Global threshold reference is resolved correctly.
	tm := map[ThresholdKind]quantity.Quantity{
		KindEntity: *quantity.NewFromUint64(1_000),
	}
	kind := KindEntity
	st = StakeThreshold{Global: &kind}
	v, err := st.Value(tm)
	require.NoError(err, "global threshold reference should be resolved correctly")
	q := tm[kind]
	require.True(q.Cmp(v) == 0, "global threshold reference should be resolved correctly")

	// Constant threshold is resolved correctly.
	c := *quantity.NewFromUint64(5_000)
	st = StakeThreshold{Constant: &c}
	v, err = st.Value(tm)
	require.NoError(err, "constant threshold should be resolved correctly")
	require.True(c.Cmp(v) == 0, "constant threshold should be resolved correctly")

	// Equality checks.
	kind2 := KindEntity
	kind3 := KindNodeCompute
	c2 := *quantity.NewFromUint64(1_000)

	for _, t := range []struct {
		a     StakeThreshold
		b     StakeThreshold
		equal bool
	}{
		{StakeThreshold{Global: &kind}, StakeThreshold{Global: &kind}, true},
		{StakeThreshold{Global: &kind}, StakeThreshold{Global: &kind2}, true},
		{StakeThreshold{Global: &kind}, StakeThreshold{Global: &kind3}, false},
		{StakeThreshold{Global: &kind}, StakeThreshold{Constant: &c2}, false},
		{StakeThreshold{Constant: &c2}, StakeThreshold{Constant: &c2}, true},
		{StakeThreshold{Constant: &c}, StakeThreshold{Constant: &c2}, false},
		{StakeThreshold{}, StakeThreshold{Constant: &c2}, false},
		{StakeThreshold{}, StakeThreshold{}, false},
	} {
		require.True(t.a.Equal(&t.b) == t.equal, "stake threshold equality should work (a == b)")
		require.True(t.b.Equal(&t.a) == t.equal, "stake threshold equality should work (b == a)")
	}
}

func TestStakeAccumulator(t *testing.T) {
	require := require.New(t)

	thresholds := map[ThresholdKind]quantity.Quantity{
		KindEntity:            *quantity.NewFromUint64(1_000),
		KindNodeValidator:     *quantity.NewFromUint64(10_000),
		KindNodeCompute:       *quantity.NewFromUint64(5_000),
		KindNodeStorage:       *quantity.NewFromUint64(2_000),
		KindNodeKeyManager:    *quantity.NewFromUint64(50_000),
		KindRuntimeCompute:    *quantity.NewFromUint64(100_000),
		KindRuntimeKeyManager: *quantity.NewFromUint64(1_000_000),
	}

	// Empty escrow account tests.
	var acct EscrowAccount
	err := acct.CheckStakeClaims(thresholds)
	require.NoError(err, "empty escrow account should check out")
	err = acct.RemoveStakeClaim(StakeClaim("dummy claim"))
	require.Error(err, "removing a non-existing claim should return an error")
	err = acct.AddStakeClaim(thresholds, StakeClaim("claim1"), GlobalStakeThresholds(KindEntity, KindNodeValidator))
	require.Error(err, "adding a stake claim with insufficient stake should fail")
	require.Equal(err, ErrInsufficientStake)
	require.EqualValues(EscrowAccount{}, acct, "account should be unchanged after failure")

	// Add some stake into the account.
	acct.Active.Balance = *quantity.NewFromUint64(3_000)
	err = acct.CheckStakeClaims(thresholds)
	require.NoError(err, "escrow account with no claims should check out")

	err = acct.AddStakeClaim(thresholds, StakeClaim("claim1"), GlobalStakeThresholds(KindEntity, KindNodeCompute))
	require.Error(err, "adding a stake claim with insufficient stake should fail")
	require.Equal(err, ErrInsufficientStake)

	err = acct.AddStakeClaim(thresholds, StakeClaim("claim1"), GlobalStakeThresholds(KindEntity))
	require.NoError(err, "adding a stake claim with sufficient stake should work")
	err = acct.CheckStakeClaims(thresholds)
	require.NoError(err, "escrow account should check out")

	// Update an existing claim.
	err = acct.AddStakeClaim(thresholds, StakeClaim("claim1"), GlobalStakeThresholds(KindEntity, KindNodeCompute))
	require.Error(err, "updating a stake claim with insufficient stake should fail")
	require.Equal(err, ErrInsufficientStake)

	err = acct.AddStakeClaim(thresholds, StakeClaim("claim1"), GlobalStakeThresholds(KindEntity, KindNodeStorage))
	require.NoError(err, "updating a stake claim with sufficient stake should work")

	err = acct.AddStakeClaim(thresholds, StakeClaim("claim1"), GlobalStakeThresholds(KindEntity, KindNodeStorage))
	require.NoError(err, "updating a stake claim with sufficient stake should work")
	err = acct.CheckStakeClaims(thresholds)
	require.NoError(err, "escrow account should check out")

	// Add another claim.
	err = acct.AddStakeClaim(thresholds, StakeClaim("claim2"), GlobalStakeThresholds(KindNodeStorage))
	require.Error(err, "updating a stake claim with insufficient stake should fail")
	require.Equal(err, ErrInsufficientStake)

	acct.Active.Balance = *quantity.NewFromUint64(13_000)

	err = acct.AddStakeClaim(thresholds, StakeClaim("claim2"), GlobalStakeThresholds(KindNodeStorage))
	require.NoError(err, "adding a stake claim with sufficient stake should work")
	err = acct.CheckStakeClaims(thresholds)
	require.NoError(err, "escrow account should check out")

	require.Len(acct.StakeAccumulator.Claims, 2, "stake accumulator should contain two claims")

	err = acct.AddStakeClaim(thresholds, StakeClaim("claim3"), GlobalStakeThresholds(KindNodeValidator))
	require.Error(err, "adding a stake claim with insufficient stake should fail")
	require.Equal(err, ErrInsufficientStake)

	// Add constant claim.
	q1 := *quantity.NewFromUint64(10)
	err = acct.AddStakeClaim(thresholds, StakeClaim("claimC1"), []StakeThreshold{{Constant: &q1}})
	require.NoError(err, "adding a constant stake claim with sufficient stake should work")
	err = acct.CheckStakeClaims(thresholds)
	require.NoError(err, "escrow account should check out")

	q2 := *quantity.NewFromUint64(10_000)
	err = acct.AddStakeClaim(thresholds, StakeClaim("claimC2"), []StakeThreshold{{Constant: &q2}})
	require.Error(err, "adding a constant stake claim with insufficient stake should fail")
	require.Equal(err, ErrInsufficientStake)

	// Remove an existing claim.
	err = acct.RemoveStakeClaim(StakeClaim("claim2"))
	require.NoError(err, "removing an existing claim should work")
	require.Len(acct.StakeAccumulator.Claims, 2, "stake accumulator should contain two claims")

	err = acct.RemoveStakeClaim(StakeClaim("claimC1"))
	require.NoError(err, "removing an existing claim should work")
	require.Len(acct.StakeAccumulator.Claims, 1, "stake accumulator should contain one claim")

	err = acct.AddStakeClaim(thresholds, StakeClaim("claim3"), GlobalStakeThresholds(KindNodeValidator))
	require.NoError(err, "adding a stake claim with sufficient stake should work")
	require.Len(acct.StakeAccumulator.Claims, 2, "stake accumulator should contain two claims")
	err = acct.CheckStakeClaims(thresholds)
	require.NoError(err, "escrow account should check out")

	// Reduce stake.
	acct.Active.Balance = *quantity.NewFromUint64(5_000)
	err = acct.CheckStakeClaims(thresholds)
	require.Error(err, "escrow account should no longer check out")
	require.Equal(err, ErrInsufficientStake)
}
