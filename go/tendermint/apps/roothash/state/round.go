package state

import (
	"errors"
	"time"

	"github.com/oasislabs/oasis-core/go/common/cbor"
	"github.com/oasislabs/oasis-core/go/roothash/api/block"
	"github.com/oasislabs/oasis-core/go/roothash/api/commitment"
)

var (
	_ cbor.Marshaler   = (*Round)(nil)
	_ cbor.Unmarshaler = (*Round)(nil)
)

// Round is a roothash round.
type Round struct {
	ComputePool *commitment.MultiPool `json:"compute_pool"`
	MergePool   *commitment.Pool      `json:"merge_pool"`

	CurrentBlock *block.Block `json:"current_block"`
	Finalized    bool         `json:"finalized"`
}

func (r *Round) Reset() {
	r.ComputePool.ResetCommitments()
	r.MergePool.ResetCommitments()
	r.Finalized = false
}

func (r *Round) GetNextTimeout() (timeout time.Time) {
	timeout = r.ComputePool.GetNextTimeout()
	if timeout.IsZero() || (!r.MergePool.NextTimeout.IsZero() && r.MergePool.NextTimeout.Before(timeout)) {
		timeout = r.MergePool.NextTimeout
	}
	return
}

func (r *Round) AddComputeCommitment(commitment *commitment.ComputeCommitment, sv commitment.SignatureVerifier) (*commitment.Pool, error) {
	if r.Finalized {
		return nil, errors.New("tendermint/roothash: round is already finalized, can't commit")
	}
	return r.ComputePool.AddComputeCommitment(r.CurrentBlock, sv, commitment)
}

func (r *Round) AddMergeCommitment(commitment *commitment.MergeCommitment, sv commitment.SignatureVerifier) error {
	if r.Finalized {
		return errors.New("tendermint/roothash: round is already finalized, can't commit")
	}
	return r.MergePool.AddMergeCommitment(r.CurrentBlock, sv, commitment, r.ComputePool)
}

func (r *Round) Transition(blk *block.Block) {
	r.CurrentBlock = blk
	r.Reset()
}

// MarshalCBOR serializes the type into a CBOR byte vector.
func (r *Round) MarshalCBOR() []byte {
	return cbor.Marshal(r)
}

// UnmarshalCBOR deserializes a CBOR byte vector into given type.
func (r *Round) UnmarshalCBOR(data []byte) error {
	return cbor.Unmarshal(data, r)
}

func NewRound(
	computePool *commitment.MultiPool,
	mergePool *commitment.Pool,
	blk *block.Block,
) *Round {
	r := &Round{
		CurrentBlock: blk,
		ComputePool:  computePool,
		MergePool:    mergePool,
	}
	r.Reset()

	return r
}