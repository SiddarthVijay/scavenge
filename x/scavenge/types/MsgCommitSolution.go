package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCommitSolution{}

// MsgCommitSolution - struct for unjailing jailed validator
type MsgCommitSolution struct {
	Scavenger				sdk.AccAddress	`json:"scavenger" yaml:"scavenger"` 						// Address of the scavenger who is submitting this
	SolutionHash 			string			`json:"solutionHash" yaml:"solutionHash"` 					// Hash of the solution for verification
	SolutionScavengerHash	string			`json:"solutionScavengerHash" yaml:"solutionScavengerHash"`	// Need to check what this is
}

// NewMsgCommitSolution creates a new MsgCommitSolution instance
func NewMsgCommitSolution(scavenger sdk.AccAddress, solutionHash, solutionScavengerHash string) MsgCommitSolution {
	return MsgCommitSolution{
		Scavenger: scavenger,
		SolutionHash: solutionHash,
		SolutionScavengerHash: solutionScavengerHash,
	}
}

const CommitSolutionConst = "CommitSolution"

// nolint
func (msg MsgCommitSolution) Route() string { return RouterKey }
func (msg MsgCommitSolution) Type() string  { return CommitSolutionConst }
func (msg MsgCommitSolution) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Scavenger}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgCommitSolution) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgCommitSolution) ValidateBasic() error {
	if msg.Scavenger.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing creator address")
	}
	if msg.SolutionHash == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "solutionHash cannot be empty")
	}
	if msg.SolutionScavengerHash == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "solutionScavengerHash cannot be empty")
	}

	return nil
}
