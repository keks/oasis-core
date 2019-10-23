package state

import (
	"errors"
	"time"

	"github.com/tendermint/iavl"

	"github.com/oasislabs/oasis-core/go/common/cbor"
	"github.com/oasislabs/oasis-core/go/common/crypto/signature"
	"github.com/oasislabs/oasis-core/go/common/keyformat"
	registry "github.com/oasislabs/oasis-core/go/registry/api"
	"github.com/oasislabs/oasis-core/go/roothash/api/block"
	"github.com/oasislabs/oasis-core/go/tendermint/abci"
)

var (
	// runtimeKeyFmt is the key format used for per-runtime roothash state.
	//
	// Value is CBOR-serialized runtime state.
	runtimeKeyFmt = keyformat.New(0x20, &signature.MapKey{})
	// roundTimeoutKeyFmt is the key format used for round timeout value.
	//
	// Value is CBOR-serialized round timeout.
	roundTimeoutKeyFmt = keyformat.New(0x21)

	_ cbor.Marshaler   = (*RuntimeState)(nil)
	_ cbor.Unmarshaler = (*RuntimeState)(nil)
)

type RuntimeState struct {
	Runtime      *registry.Runtime `json:"runtime"`
	CurrentBlock *block.Block      `json:"current_block"`
	GenesisBlock *block.Block      `json:"genesis_block"`
	Round        *Round            `json:"round"`
	Timer        abci.Timer        `json:"timer"`
}

func (s *RuntimeState) MarshalCBOR() []byte {
	return cbor.Marshal(s)
}

func (s *RuntimeState) UnmarshalCBOR(data []byte) error {
	return cbor.Unmarshal(data, s)
}

type ImmutableState struct {
	*abci.ImmutableState
}

func NewImmutableState(state *abci.ApplicationState, version int64) (*ImmutableState, error) {
	inner, err := abci.NewImmutableState(state, version)
	if err != nil {
		return nil, err
	}

	return &ImmutableState{inner}, nil
}

func (s *ImmutableState) RuntimeState(id signature.PublicKey) (*RuntimeState, error) {
	_, raw := s.Snapshot.Get(runtimeKeyFmt.Encode(&id))
	if raw == nil {
		return nil, nil
	}

	var state RuntimeState
	err := state.UnmarshalCBOR(raw)
	return &state, err
}

func (s *ImmutableState) Runtimes() []*RuntimeState {
	var runtimes []*RuntimeState
	s.Snapshot.IterateRange(
		runtimeKeyFmt.Encode(),
		nil,
		true,
		func(key, value []byte) bool {
			if !runtimeKeyFmt.Decode(key) {
				return true
			}

			var state RuntimeState
			cbor.MustUnmarshal(value, &state)

			runtimes = append(runtimes, &state)
			return false
		},
	)

	return runtimes
}

func (s *ImmutableState) RoundTimeout() (time.Duration, error) {
	_, raw := s.Snapshot.Get(roundTimeoutKeyFmt.Encode())
	if raw == nil {
		panic(errors.New("tendermint/roothash: expected round timeout to be present in app state"))
	}

	var rt time.Duration
	err := cbor.Unmarshal(raw, &rt)
	return rt, err
}

type MutableState struct {
	*ImmutableState

	tree *iavl.MutableTree
}

func NewMutableState(tree *iavl.MutableTree) *MutableState {
	inner := &abci.ImmutableState{Snapshot: tree.ImmutableTree}

	return &MutableState{
		ImmutableState: &ImmutableState{inner},
		tree:           tree,
	}
}

func (s *MutableState) SetRuntimeState(state *RuntimeState) {
	s.tree.Set(runtimeKeyFmt.Encode(&state.Runtime.ID), state.MarshalCBOR())
}

func (s *MutableState) SetRoundTimeout(rt time.Duration) {
	s.tree.Set(roundTimeoutKeyFmt.Encode(), cbor.Marshal(rt))
}
