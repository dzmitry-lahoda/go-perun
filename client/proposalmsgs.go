// Copyright 2019 - See NOTICE file for copyright holders.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"hash"
	"io"

	"golang.org/x/crypto/sha3"

	"github.com/pkg/errors"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/log"
	perunio "perun.network/go-perun/pkg/io"
	"perun.network/go-perun/wallet"
	"perun.network/go-perun/wire"
)

func init() {
	wire.RegisterDecoder(wire.LedgerChannelProposal,
		func(r io.Reader) (wire.Msg, error) {
			m := LedgerChannelProposal{}
			return &m, m.Decode(r)
		})
	wire.RegisterDecoder(wire.LedgerChannelProposalAcc,
		func(r io.Reader) (wire.Msg, error) {
			var m LedgerChannelProposalAcc
			return &m, m.Decode(r)
		})
	wire.RegisterDecoder(wire.SubChannelProposal,
		func(r io.Reader) (wire.Msg, error) {
			m := SubChannelProposal{}
			return &m, m.Decode(r)
		})
	wire.RegisterDecoder(wire.SubChannelProposalAcc,
		func(r io.Reader) (wire.Msg, error) {
			var m SubChannelProposalAcc
			return &m, m.Decode(r)
		})
	wire.RegisterDecoder(wire.ChannelProposalRej,
		func(r io.Reader) (wire.Msg, error) {
			var m ChannelProposalRej
			return &m, m.Decode(r)
		})
}

func newHasher() hash.Hash { return sha3.New256() }

// ProposalID uniquely identifies the channel proposal as  specified by the
// Channel Proposal Protocol (CPP).
type ProposalID = [32]byte

// NonceShare is used to cooperatively calculate a channel's nonce.
type NonceShare = [32]byte

type (
	// ChannelProposal is the interface that describes all channel proposal
	// message types.
	ChannelProposal interface {
		wire.Msg
		perunio.Decoder

		// Base returns the channel proposal's common values.
		Base() *BaseChannelProposal

		// Matches checks whether an accept message is of the correct type. This
		// does not check any contents of the accept message, only its type.
		Matches(ChannelProposalAccept) bool

		// Valid checks whether a channel proposal is valid.
		Valid() error

		// ProposalID calculates the proposal's unique identifier.
		ProposalID() ProposalID
	}

	// BaseChannelProposal contains all data necessary to propose a new
	// channel to a given set of peers. It is also sent over the wire.
	//
	// BaseChannelProposal implements the channel proposal messages from the
	// Multi-Party Channel Proposal Protocol (MPCPP).
	BaseChannelProposal struct {
		ChallengeDuration uint64              // Dispute challenge duration.
		NonceShare        NonceShare          // Proposer's channel nonce share.
		App               channel.App         // App definition, or nil.
		InitData          channel.Data        // Initial App data.
		InitBals          *channel.Allocation // Initial balances.
	}

	// LedgerChannelProposal is a channel proposal for ledger channels.
	LedgerChannelProposal struct {
		BaseChannelProposal
		Participant wallet.Address // Proposer's address in the channel.
		Peers       []wire.Address // Participants' wire addresses.
	}

	// SubChannelProposal is a channel proposal for subchannels.
	SubChannelProposal struct {
		BaseChannelProposal
		Parent channel.ID
	}
)

// proposalPeers returns the wire addresses of a proposed channel's
// participants.
func (c *Client) proposalPeers(p ChannelProposal) (peers []wire.Address) {
	switch prop := p.(type) {
	case *LedgerChannelProposal:
		peers = prop.Peers
	case *SubChannelProposal:
		ch, ok := c.channels.Get(prop.Parent)
		if !ok {
			c.log.Panic("ProposalPeers: invalid parent channel")
		}
		peers = ch.Peers()
	default:
		c.log.Panicf("ProposalPeers: unhandled proposal type %T")
	}
	return
}

// makeBaseChannelProposal creates a BaseChannelProposal and applies the supplied
// options. For more information, see ProposalOpts.
func makeBaseChannelProposal(
	challengeDuration uint64,
	initBals *channel.Allocation,
	opts ...ProposalOpts,
) BaseChannelProposal {
	opt := union(opts...)

	return BaseChannelProposal{
		ChallengeDuration: challengeDuration,
		NonceShare:        opt.nonce(),
		App:               opt.App(),
		InitData:          opt.AppData(),
		InitBals:          initBals,
	}
}

// Base returns the channel proposal's common values.
func (p *BaseChannelProposal) Base() *BaseChannelProposal {
	return p
}

// NumPeers returns the number of peers in a channel.
func (p BaseChannelProposal) NumPeers() int {
	return len(p.InitBals.Balances[0])
}

// Encode encodes the BaseChannelProposal into an io.Writer.
func (p BaseChannelProposal) Encode(w io.Writer) error {
	if w == nil {
		return errors.New("writer must not be nil")
	}

	if err := perunio.Encode(w, p.ChallengeDuration, p.NonceShare); err != nil {
		return err
	}

	return perunio.Encode(w, OptAppAndDataEnc{p.App, p.InitData}, p.InitBals)
}

// OptAppAndDataEnc makes an optional pair of App definition and Data encodable.
type OptAppAndDataEnc struct {
	channel.App
	channel.Data
}

// Encode encodes an optional pair of App definition and Data.
func (o OptAppAndDataEnc) Encode(w io.Writer) error {
	return perunio.Encode(w, channel.OptAppEnc{App: o.App}, o.Data)
}

// OptAppAndDataDec makes an optional pair of App definition and Data decodable.
type OptAppAndDataDec struct {
	App  *channel.App
	Data *channel.Data
}

// Decode decodes an optional pair of App definition and Data.
func (o OptAppAndDataDec) Decode(r io.Reader) (err error) {
	if err = perunio.Decode(r, channel.OptAppDec{App: o.App}); err != nil {
		return err
	}
	*o.Data, err = (*o.App).DecodeData(r)
	return err
}

// Decode decodes a BaseChannelProposal from an io.Reader.
func (p *BaseChannelProposal) Decode(r io.Reader) (err error) {
	if r == nil {
		return errors.New("reader must not be nil")
	}

	if err := perunio.Decode(r, &p.ChallengeDuration, &p.NonceShare); err != nil {
		return err
	}

	if p.InitBals == nil {
		p.InitBals = new(channel.Allocation)
	}

	return perunio.Decode(r, OptAppAndDataDec{&p.App, &p.InitData}, p.InitBals)
}

// Valid checks that the channel proposal is valid:
// * InitBals must not be nil
// * ValidateProposalParameters returns nil
// * InitBals are valid
// * No locked sub-allocations
// * non-zero ChallengeDuration.
func (p *BaseChannelProposal) Valid() error {
	if p.InitBals == nil {
		return errors.New("invalid nil fields")
	} else if err := channel.ValidateProposalParameters(
		p.ChallengeDuration, p.NumPeers(), p.App); err != nil {
		return errors.WithMessage(err, "invalid channel parameters")
	} else if err := p.InitBals.Valid(); err != nil {
		return err
	} else if len(p.InitBals.Locked) != 0 {
		return errors.New("initial allocation cannot have locked funds")
	}
	return nil
}

// Accept constructs an accept message that belongs to a proposal message. It
// should be used instead of manually constructing an accept message.
func (p LedgerChannelProposal) Accept(
	participant wallet.Address,
	nonceShare ProposalOpts,
) *LedgerChannelProposalAcc {
	if !nonceShare.isNonce() {
		log.WithField("proposal", p.ProposalID()).
			Panic("LedgerChannelProposal.Accept: nonceShare has no configured nonce")
	}
	return &LedgerChannelProposalAcc{
		BaseChannelProposalAcc: makeBaseChannelProposalAcc(
			p.ProposalID(), nonceShare.nonce()),
		Participant: participant,
	}
}

// Matches requires that the accept message is a LedgerChannelAcc message.
func (LedgerChannelProposal) Matches(acc ChannelProposalAccept) bool {
	_, ok := acc.(*LedgerChannelProposalAcc)
	return ok
}

// NewLedgerChannelProposal creates a ledger channel proposal and applies the
// supplied options. For more information, see ProposalOpts.
func NewLedgerChannelProposal(
	challengeDuration uint64,
	participant wallet.Address,
	initBals *channel.Allocation,
	peers []wire.Address,
	opts ...ProposalOpts,
) *LedgerChannelProposal {
	return &LedgerChannelProposal{
		BaseChannelProposal: makeBaseChannelProposal(
			challengeDuration,
			initBals,
			opts...),
		Participant: participant,
		Peers:       peers}
}

// Type returns wire.LedgerChannelProposal.
func (LedgerChannelProposal) Type() wire.Type {
	return wire.LedgerChannelProposal
}

// ProposalID returns the identifier of this channel proposal request.
func (p LedgerChannelProposal) ProposalID() (propID ProposalID) {
	hasher := newHasher()
	if err := perunio.Encode(hasher,
		p.Base(),
		p.Participant,
		wire.Addresses(p.Peers)); err != nil {
		log.Panicf("proposal ID nonce encoding: %v", err)
	}

	copy(propID[:], hasher.Sum(nil))
	return
}

// Encode encodes a ledger channel proposal.
func (p LedgerChannelProposal) Encode(w io.Writer) error {
	if err := p.assertValidNumParts(); err != nil {
		return err
	}
	return perunio.Encode(w,
		p.BaseChannelProposal,
		p.Participant,
		wire.AddressesWithLen(p.Peers))
}

// Decode decodes a ledger channel proposal.
func (p *LedgerChannelProposal) Decode(r io.Reader) error {
	err := perunio.Decode(r,
		&p.BaseChannelProposal,
		wallet.AddressDec{Addr: &p.Participant},
		(*wire.AddressesWithLen)(&p.Peers))
	if err != nil {
		return err
	}

	return p.assertValidNumParts()
}

func (p LedgerChannelProposal) assertValidNumParts() error {
	if len(p.Peers) < 2 || len(p.Peers) > channel.MaxNumParts {
		return errors.Errorf("expected 2-%d participants, got %d",
			channel.MaxNumParts, len(p.Peers))
	}
	return nil
}

// Valid checks whether the participant address is nil.
func (p LedgerChannelProposal) Valid() error {
	if err := p.BaseChannelProposal.Valid(); err != nil {
		return err
	}
	if p.Participant == nil {
		return errors.New("invalid nil participant")
	}
	return nil
}

// NewSubChannelProposal creates a subchannel proposal and applies the
// supplied options. For more information, see ProposalOpts.
func NewSubChannelProposal(
	parent channel.ID,
	challengeDuration uint64,
	initBals *channel.Allocation,
	opts ...ProposalOpts,
) *SubChannelProposal {
	return &SubChannelProposal{
		BaseChannelProposal: makeBaseChannelProposal(
			challengeDuration,
			initBals,
			opts...),
		Parent: parent,
	}
}

// ProposalID returns the identifier of this channel proposal request.
func (p SubChannelProposal) ProposalID() (propID ProposalID) {
	hasher := newHasher()
	if err := perunio.Encode(hasher,
		p.Base(),
		p.Parent); err != nil {
		log.Panicf("proposal ID nonce encoding: %v", err)
	}

	copy(propID[:], hasher.Sum(nil))
	return
}

// Encode encodes the SubChannelProposal into an io.Writer.
func (p SubChannelProposal) Encode(w io.Writer) error {
	return perunio.Encode(w, p.BaseChannelProposal, p.Parent)
}

// Decode decodes a SubChannelProposal from an io.Reader.
func (p *SubChannelProposal) Decode(r io.Reader) error {
	return perunio.Decode(r, &p.BaseChannelProposal, &p.Parent)
}

// Type returns wire.SubChannelProposal.
func (SubChannelProposal) Type() wire.Type {
	return wire.SubChannelProposal
}

// Accept constructs an accept message that belongs to a proposal message. It
// should be used instead of manually constructing an accept message.
func (p SubChannelProposal) Accept(
	nonceShare ProposalOpts,
) *SubChannelProposalAcc {
	propID := p.ProposalID()
	if !nonceShare.isNonce() {
		log.WithField("proposal", propID).
			Panic("SubChannelProposal.Accept: nonceShare has no configured nonce")
	}
	return &SubChannelProposalAcc{
		BaseChannelProposalAcc: makeBaseChannelProposalAcc(
			propID, nonceShare.nonce()),
	}
}

// Matches requires that the accept message is a sub channel proposal accept
// message.
func (SubChannelProposal) Matches(acc ChannelProposalAccept) bool {
	_, ok := acc.(*SubChannelProposalAcc)
	return ok
}

type (
	// ChannelProposalAccept is the generic interface for channel proposal
	// accept messages.
	ChannelProposalAccept interface {
		wire.Msg
		Base() *BaseChannelProposalAcc
	}

	// BaseChannelProposalAcc contains all data for a response to a channel proposal
	// message. The ProposalID must correspond to the channel proposal request one
	// wishes to respond to. Participant should be a participant address just
	// for this channel instantiation.
	//
	// The type implements the channel proposal response messages from the
	// Multi-Party Channel Proposal Protocol (MPCPP).
	BaseChannelProposalAcc struct {
		ProposalID ProposalID // Proposal session ID we're answering.
		NonceShare NonceShare // Responder's channel nonce share.
	}

	// LedgerChannelProposalAcc is the accept message type corresponding to
	// ledger channel proposals. ParticipantAdd is recommended to be unique for
	// each channel instantiation.
	LedgerChannelProposalAcc struct {
		BaseChannelProposalAcc
		Participant wallet.Address // Responder's participant address.
	}

	// SubChannelProposalAcc is the accept message type corresponding to sub
	// channel proposals.
	SubChannelProposalAcc struct {
		BaseChannelProposalAcc
	}
)

func makeBaseChannelProposalAcc(
	proposalID ProposalID,
	nonceShare NonceShare,
) BaseChannelProposalAcc {
	return BaseChannelProposalAcc{
		ProposalID: proposalID,
		NonceShare: nonceShare,
	}
}

// Encode encodes a BaseChannelProposalAcc.
func (acc BaseChannelProposalAcc) Encode(w io.Writer) error {
	return perunio.Encode(w,
		acc.ProposalID,
		acc.NonceShare)
}

// Decode decodes a BaseChannelProposalAcc.
func (acc *BaseChannelProposalAcc) Decode(r io.Reader) error {
	return perunio.Decode(r,
		&acc.ProposalID,
		&acc.NonceShare)
}

// Type returns wire.ChannelProposalAcc.
func (LedgerChannelProposalAcc) Type() wire.Type {
	return wire.LedgerChannelProposalAcc
}

// Base returns the common proposal accept values.
func (acc *LedgerChannelProposalAcc) Base() *BaseChannelProposalAcc {
	return &acc.BaseChannelProposalAcc
}

// Encode encodes the LedgerChannelProposalAcc into an io.Writer.
func (acc LedgerChannelProposalAcc) Encode(w io.Writer) error {
	return perunio.Encode(w,
		acc.BaseChannelProposalAcc,
		acc.Participant)
}

// Decode decodes a LedgerChannelProposalAcc from an io.Reader.
func (acc *LedgerChannelProposalAcc) Decode(r io.Reader) error {
	return perunio.Decode(r,
		&acc.BaseChannelProposalAcc,
		wallet.AddressDec{Addr: &acc.Participant})
}

// Type returns wire.SubChannelProposalAcc.
func (SubChannelProposalAcc) Type() wire.Type {
	return wire.SubChannelProposalAcc
}

// Base returns the common proposal accept values.
func (acc *SubChannelProposalAcc) Base() *BaseChannelProposalAcc {
	return &acc.BaseChannelProposalAcc
}

// Encode encodes the SubChannelProposalAcc into an io.Writer.
func (acc SubChannelProposalAcc) Encode(w io.Writer) error {
	return perunio.Encode(w, acc.BaseChannelProposalAcc)
}

// Decode decodes a SubChannelProposalAcc from an io.Reader.
func (acc *SubChannelProposalAcc) Decode(r io.Reader) error {
	return perunio.Decode(r, &acc.BaseChannelProposalAcc)
}

// ChannelProposalRej is used to reject a ChannelProposalReq.
// An optional reason for the rejection can be set.
//
// The message is one of two possible responses in the
// Multi-Party Channel Proposal Protocol (MPCPP).
type ChannelProposalRej struct {
	ProposalID ProposalID // The channel proposal to reject.
	Reason     string     // The rejection reason.
}

// Type returns wire.ChannelProposalRej.
func (ChannelProposalRej) Type() wire.Type {
	return wire.ChannelProposalRej
}

// Encode encodes a ChannelProposalRej into an io.Writer.
func (rej ChannelProposalRej) Encode(w io.Writer) error {
	return perunio.Encode(w, rej.ProposalID, rej.Reason)
}

// Decode decodes a ChannelProposalRej from an io.Reader.
func (rej *ChannelProposalRej) Decode(r io.Reader) error {
	return perunio.Decode(r, &rej.ProposalID, &rej.Reason)
}
