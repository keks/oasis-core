package e2e

import (
	"context"
	"errors"
	"fmt"
	"time"

	consensus "github.com/oasisprotocol/oasis-core/go/consensus/api"
	"github.com/oasisprotocol/oasis-core/go/oasis-test-runner/env"
	"github.com/oasisprotocol/oasis-core/go/oasis-test-runner/oasis"
	"github.com/oasisprotocol/oasis-core/go/oasis-test-runner/scenario"
)

// EarlyQuery is the early query scenario where we query a validator node before the network
// has started and there are no committed blocks.
var EarlyQuery scenario.Scenario = &earlyQueryImpl{
	E2E: *NewE2E("early-query"),
}

type earlyQueryImpl struct {
	E2E
}

func (sc *earlyQueryImpl) Clone() scenario.Scenario {
	return &earlyQueryImpl{
		E2E: sc.E2E.Clone(),
	}
}

func (sc *earlyQueryImpl) Fixture() (*oasis.NetworkFixture, error) {
	f, err := sc.E2E.Fixture()
	if err != nil {
		return nil, err
	}

	// Only one validator should actually start to prevent the network from committing any blocks.
	f.Validators[1].NoAutoStart = true
	f.Validators[2].NoAutoStart = true

	return f, nil
}

func (sc *earlyQueryImpl) Run(childEnv *env.Env) error {
	// Start the network.
	var err error
	if err = sc.Net.Start(); err != nil {
		return err
	}

	// Perform some queries.
	cs := sc.Net.Controller().Consensus
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// StateToGenesis.
	_, err = cs.StateToGenesis(ctx, consensus.HeightLatest)
	if !errors.Is(err, consensus.ErrNoCommittedBlocks) {
		return fmt.Errorf("StateToGenesis query should fail with ErrNoCommittedBlocks (got: %s)", err)
	}
	// GetBlock.
	_, err = cs.GetBlock(ctx, consensus.HeightLatest)
	if !errors.Is(err, consensus.ErrNoCommittedBlocks) {
		return fmt.Errorf("GetBlock query should fail with ErrNoCommittedBlocks (got: %s)", err)
	}
	// GetTransactions.
	_, err = cs.GetTransactions(ctx, consensus.HeightLatest)
	if !errors.Is(err, consensus.ErrNoCommittedBlocks) {
		return fmt.Errorf("GetTransactions query should fail with ErrNoCommittedBlocks (got: %s)", err)
	}
	// GetTransactionsWithResults.
	_, err = cs.GetTransactionsWithResults(ctx, consensus.HeightLatest)
	if !errors.Is(err, consensus.ErrNoCommittedBlocks) {
		return fmt.Errorf("GetTransactionsWithResults query should fail with ErrNoCommittedBlocks (got: %s)", err)
	}

	// GetStatus.
	status, err := sc.Net.Controller().GetStatus(ctx)
	if err != nil {
		return fmt.Errorf("failed to get status for node: %w", err)
	}
	if status.Consensus.LatestHeight != 0 {
		return fmt.Errorf("node reports non-zero latest height before chain is initialized")
	}
	if !status.Consensus.IsValidator {
		return fmt.Errorf("node does not report itself to be a validator at genesis")
	}

	return nil
}
