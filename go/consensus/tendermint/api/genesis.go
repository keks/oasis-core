package api

import (
	"encoding/json"
	"fmt"
	"time"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/oasisprotocol/oasis-core/go/common/node"
	"github.com/oasisprotocol/oasis-core/go/common/quantity"
	"github.com/oasisprotocol/oasis-core/go/consensus/tendermint/crypto"
	genesis "github.com/oasisprotocol/oasis-core/go/genesis/api"
	cmdFlags "github.com/oasisprotocol/oasis-core/go/oasis-node/cmd/common/flags"
	registry "github.com/oasisprotocol/oasis-core/go/registry/api"
	scheduler "github.com/oasisprotocol/oasis-core/go/scheduler/api"
	staking "github.com/oasisprotocol/oasis-core/go/staking/api"
)

// GenesisProvider is a tendermint specific genesis document provider.
type GenesisProvider interface {
	GetTendermintGenesisDocument() (*tmtypes.GenesisDoc, error)
}

// GetTendermintGenesisDocument returns the Tendermint genesis document corresponding to the Oasis
// genesis document specified by the given genesis provider.
func GetTendermintGenesisDocument(provider genesis.Provider) (*tmtypes.GenesisDoc, error) {
	doc, err := provider.GetGenesisDocument()
	if err != nil {
		return nil, fmt.Errorf("tendermint: failed to obtain genesis document: %w", err)
	}

	var tmGenDoc *tmtypes.GenesisDoc
	if tmProvider, ok := provider.(GenesisProvider); ok {
		// This is a single node config, because the genesis document was
		// missing, probably in unit tests.
		tmGenDoc, err = tmProvider.GetTendermintGenesisDocument()
	} else {
		tmGenDoc, err = genesisToTendermint(doc)
	}
	if err != nil {
		return nil, fmt.Errorf("tendermint: failed to create genesis document: %w", err)
	}

	// HACK: Certain test cases use TimeoutCommit < 1 sec, and care about the
	// BFT view of time pulling ahead.
	timeoutCommit := doc.Consensus.Parameters.TimeoutCommit
	tmGenDoc.ConsensusParams.Block.TimeIotaMs = int64(timeoutCommit / time.Millisecond)

	return tmGenDoc, nil
}

// genesisToTendermint converts the Oasis genesis block to Tendermint's format.
func genesisToTendermint(d *genesis.Document) (*tmtypes.GenesisDoc, error) {
	// WARNING: The AppState MUST be encoded as JSON since its type is
	// json.RawMessage which requires it to be valid JSON. It may appear
	// to work until you try to restore from an existing data directory.
	//
	// The runtime library sorts map keys, so the output of json.Marshal
	// should be deterministic.
	b, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("tendermint: failed to serialize genesis doc: %w", err)
	}

	// Translate special "disable block gas limit" value as Tendermint uses
	// -1 for some reason (as if a zero limit makes sense) and we use 0.
	maxBlockGas := int64(d.Consensus.Parameters.MaxBlockGas)
	if maxBlockGas == 0 {
		maxBlockGas = -1
	}

	// Automatically compute evidence parameters based on debonding period.
	debondingInterval := int64(d.Staking.Parameters.DebondingInterval)
	if debondingInterval == 0 && cmdFlags.DebugDontBlameOasis() {
		// Use a default of 1 epoch in case debonding is disabled and we are using debug mode. If
		// not in debug mode, this will just cause startup to fail which is good.
		debondingInterval = 1
	}
	epochInterval := d.EpochTime.Parameters.Interval
	if epochInterval == 0 && cmdFlags.DebugDontBlameOasis() && d.EpochTime.Parameters.DebugMockBackend {
		// Use a default of 100 blocks in case epoch interval is unset and we are using debug mode.
		epochInterval = 100
	}

	var evCfg tmproto.EvidenceParams
	evCfg.MaxNum = d.Consensus.Parameters.MaxEvidenceNum
	evCfg.MaxAgeNumBlocks = debondingInterval * epochInterval
	evCfg.MaxAgeDuration = time.Duration(evCfg.MaxAgeNumBlocks) * (d.Consensus.Parameters.TimeoutCommit + 1*time.Second)
	evCfg.ProofTrialPeriod = evCfg.MaxAgeNumBlocks / 2

	doc := tmtypes.GenesisDoc{
		ChainID:     d.ChainContext()[:tmtypes.MaxChainIDLen],
		GenesisTime: d.Time,
		ConsensusParams: &tmproto.ConsensusParams{
			Block: tmproto.BlockParams{
				MaxBytes:   int64(d.Consensus.Parameters.MaxBlockSize),
				MaxGas:     maxBlockGas,
				TimeIotaMs: 1000,
			},
			Evidence: evCfg,
			Validator: tmproto.ValidatorParams{
				PubKeyTypes: []string{tmtypes.ABCIPubKeyTypeEd25519},
			},
		},
		AppState: b,
	}

	var tmValidators []tmtypes.GenesisValidator
	for _, v := range d.Registry.Nodes {
		var openedNode node.Node
		if err = v.Open(registry.RegisterGenesisNodeSignatureContext, &openedNode); err != nil {
			return nil, fmt.Errorf("tendermint: failed to verify validator: %w", err)
		}
		// TODO: This should cross check that the entity is valid.
		if !openedNode.HasRoles(node.RoleValidator) {
			continue
		}

		var power int64
		if d.Scheduler.Parameters.DebugBypassStake {
			power = 1
		} else {
			var stake *quantity.Quantity
			acctAddr := staking.NewAddress(openedNode.EntityID)
			if account, ok := d.Staking.Ledger[acctAddr]; ok {
				stake = account.Escrow.Active.Balance.Clone()
			} else {
				// If all balances and stuff are zero, it's permitted not to have an account in the ledger at all.
				stake = &quantity.Quantity{}
			}
			power, err = scheduler.VotingPowerFromStake(stake)
			if err != nil {
				return nil, fmt.Errorf("tendermint: computing voting power for entity %s with account %s and stake %v: %w",
					openedNode.EntityID,
					acctAddr,
					stake,
					err,
				)
			}
		}

		pk := crypto.PublicKeyToTendermint(&openedNode.Consensus.ID)
		validator := tmtypes.GenesisValidator{
			Address: pk.Address(),
			PubKey:  pk,
			Power:   power,
			Name:    "oasis-validator-" + openedNode.ID.String(),
		}
		tmValidators = append(tmValidators, validator)
	}

	doc.Validators = tmValidators

	return &doc, nil
}
