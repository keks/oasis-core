package fixtures

import (
	"fmt"
	"math"
	"time"

	"github.com/spf13/viper"

	"github.com/oasisprotocol/oasis-core/go/common"
	"github.com/oasisprotocol/oasis-core/go/common/node"
	"github.com/oasisprotocol/oasis-core/go/common/quantity"
	"github.com/oasisprotocol/oasis-core/go/common/sgx"
	"github.com/oasisprotocol/oasis-core/go/oasis-test-runner/oasis"
	registry "github.com/oasisprotocol/oasis-core/go/registry/api"
	"github.com/oasisprotocol/oasis-core/go/staking/api"
	staking "github.com/oasisprotocol/oasis-core/go/staking/api"
)

const (
	cfgDeterministicIdentities = "fixture.default.deterministic_entities"
	cfgEpochtimeMock           = "fixture.default.epochtime_mock"
	cfgHaltEpoch               = "fixture.default.halt_epoch"
	cfgKeymanagerBinary        = "fixture.default.keymanager.binary"
	cfgNodeBinary              = "fixture.default.node.binary"
	cfgNumEntities             = "fixture.default.num_entities"
	cfgRuntimeBinary           = "fixture.default.runtime.binary"
	cfgRuntimeGenesisState     = "fixture.default.runtime.genesis_state"
	cfgRuntimeLoader           = "fixture.default.runtime.loader"
	cfgSetupRuntimes           = "fixture.default.setup_runtimes"
	cfgTEEHardware             = "fixture.default.tee_hardware"
	cfgFund                    = "fixture.default.fund"
)

var (
	runtimeID    common.Namespace
	keymanagerID common.Namespace
)

// newDefaultFixture returns a default network fixture.
func newDefaultFixture() (*oasis.NetworkFixture, error) {
	var tee node.TEEHardware
	err := tee.FromString(viper.GetString(cfgTEEHardware))
	if err != nil {
		return nil, err
	}
	var mrSigner *sgx.MrSigner
	if tee == node.TEEHardwareIntelSGX {
		mrSigner = &sgx.FortanixDummyMrSigner
	}

	fixture := &oasis.NetworkFixture{
		TEE: oasis.TEEFixture{
			Hardware: tee,
			MrSigner: mrSigner,
		},
		Network: oasis.NetworkCfg{
			NodeBinary:             viper.GetString(cfgNodeBinary),
			RuntimeSGXLoaderBinary: viper.GetString(cfgRuntimeLoader),
			ConsensusTimeoutCommit: 1 * time.Second,
			EpochtimeMock:          viper.GetBool(cfgEpochtimeMock),
			HaltEpoch:              viper.GetUint64(cfgHaltEpoch),
			IAS: oasis.IASCfg{
				Mock: true,
			},
			DeterministicIdentities: viper.GetBool(cfgDeterministicIdentities),
			StakingGenesis: &staking.Genesis{
				TotalSupply: *quantity.NewFromUint64(10000000000),
				Ledger:      make(map[staking.Address]*staking.Account),
			},
		},
		Entities: []oasis.EntityCfg{
			{IsDebugTestEntity: true},
		},
		Validators: []oasis.ValidatorFixture{
			{Entity: 1},
		},
	}

	for i := 0; i < viper.GetInt(cfgNumEntities); i++ {
		fixture.Entities = append(fixture.Entities, oasis.EntityCfg{})
	}

	if viper.GetBool(cfgSetupRuntimes) {
		fixture.Runtimes = []oasis.RuntimeFixture{
			// Key manager runtime.
			{
				ID:         keymanagerID,
				Kind:       registry.KindKeyManager,
				Entity:     0,
				Keymanager: -1,
				Binaries:   viper.GetStringSlice(cfgKeymanagerBinary),
				AdmissionPolicy: registry.RuntimeAdmissionPolicy{
					AnyNode: &registry.AnyNodeRuntimeAdmissionPolicy{},
				},
			},
			// Compute runtime.
			{
				ID:         runtimeID,
				Kind:       registry.KindCompute,
				Entity:     0,
				Keymanager: 0,
				Binaries:   viper.GetStringSlice(cfgRuntimeBinary),
				Executor: registry.ExecutorParameters{
					GroupSize:       2,
					GroupBackupSize: 1,
					RoundTimeout:    20 * time.Second,
				},
				Merge: registry.MergeParameters{
					GroupSize:       2,
					GroupBackupSize: 1,
					RoundTimeout:    20 * time.Second,
				},
				TxnScheduler: registry.TxnSchedulerParameters{
					Algorithm:         registry.TxnSchedulerAlgorithmBatching,
					GroupSize:         2,
					MaxBatchSize:      1,
					MaxBatchSizeBytes: 16 * 1024 * 1024, // 16 MiB
					BatchFlushTimeout: 20 * time.Second,
				},
				Storage: registry.StorageParameters{
					GroupSize:               1,
					MinWriteReplication:     1,
					MaxApplyWriteLogEntries: 100_000,
					MaxApplyOps:             2,
					MaxMergeRoots:           8,
					MaxMergeOps:             2,
				},
				AdmissionPolicy: registry.RuntimeAdmissionPolicy{
					AnyNode: &registry.AnyNodeRuntimeAdmissionPolicy{},
				},
				GenesisStatePath: viper.GetString(cfgRuntimeGenesisState),
				GenesisRound:     0,
			},
		}
		fixture.KeymanagerPolicies = []oasis.KeymanagerPolicyFixture{
			{Runtime: 0, Serial: 1},
		}
		fixture.Keymanagers = []oasis.KeymanagerFixture{
			{Runtime: 0, Entity: 1},
		}
		fixture.StorageWorkers = []oasis.StorageWorkerFixture{
			{Backend: "badger", Entity: 1},
		}
		fixture.ComputeWorkers = []oasis.ComputeWorkerFixture{
			{Entity: 1, Runtimes: []int{1}},
			{Entity: 1, Runtimes: []int{1}},
			{Entity: 1, Runtimes: []int{1}},
		}
		fixture.Clients = []oasis.ClientFixture{{}}
	}

	for _, acc := range viper.GetStringSlice(cfgFund) {
		var addr api.Address
		if err := addr.UnmarshalText([]byte(acc)); err != nil {
			return nil, fmt.Errorf("Invalid fund address: %s, error: %w", acc, err)
		}
		fixture.Network.StakingGenesis.Ledger[addr] = &staking.Account{
			General: staking.GeneralAccount{
				Balance: *quantity.NewFromUint64(10000000000),
			},
		}

	}

	return fixture, nil
}

func init() {
	DefaultFixtureFlags.Bool(cfgDeterministicIdentities, false, "generate nodes with deterministic identities")
	DefaultFixtureFlags.Bool(cfgEpochtimeMock, false, "use mock epochtime")
	DefaultFixtureFlags.Bool(cfgSetupRuntimes, true, "initialize the network with runtimes and runtime nodes")
	DefaultFixtureFlags.Int(cfgNumEntities, 1, "number of (non debug) entities in genesis")
	DefaultFixtureFlags.String(cfgKeymanagerBinary, "simple-keymanager", "path to the keymanager runtime")
	DefaultFixtureFlags.String(cfgNodeBinary, "oasis-node", "path to the oasis-node binary")
	DefaultFixtureFlags.String(cfgRuntimeBinary, "simple-keyvalue", "path to the runtime binary")
	DefaultFixtureFlags.String(cfgRuntimeGenesisState, "", "path to the runtime genesis state")
	DefaultFixtureFlags.String(cfgRuntimeLoader, "oasis-core-runtime-loader", "path to the runtime loader")
	DefaultFixtureFlags.String(cfgTEEHardware, "", "TEE hardware to use")
	DefaultFixtureFlags.Uint64(cfgHaltEpoch, math.MaxUint64, "halt epoch height")
	DefaultFixtureFlags.StringSlice(cfgFund, nil, "specify account(s) that should be funded in genesis")

	_ = viper.BindPFlags(DefaultFixtureFlags)

	_ = runtimeID.UnmarshalHex("8000000000000000000000000000000000000000000000000000000000000000")
	_ = keymanagerID.UnmarshalHex("c000000000000000ffffffffffffffffffffffffffffffffffffffffffffffff")
}
