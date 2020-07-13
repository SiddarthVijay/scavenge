package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// verify interface at compile time
var _ sdk.Msg = &MsgCreateScavenge{}

// MsgCreateScavenge - struct for unjailing jailed validator
type MsgCreateScavenge struct {
	Creator			sdk.AccAddress `json:"creator" yaml:"creator"` 				// Address of the scavenge creator
	Description		string			`json:"description" yaml:"description"` 	// The question or the description of the challenge
	SolutionHash 	string 			`json:"solutionHash" yaml:"solutionHash"` 	// Hash of the solution for verification
	Reward 			sdk.Coins 		`json:"reward" yaml:"reward"` 				// Reward for the person who solves the scavenge first
}

// NewMsgCreateScavenge creates a new MsgCreateScavenge instance
func NewMsgCreateScavenge(creator sdk.AccAddress, description, solutionHash string, reward sdk.Coins) MsgCreateScavenge {
	return MsgCreateScavenge{
		Creator: creator,
		Description: description,
		SolutionHash: solutionHash,
		Reward: reward,
	}
}

const CreateScavengeConst = "CreateScavenge"

// nolint
func (msg MsgCreateScavenge) Route() string { return RouterKey }
func (msg MsgCreateScavenge) Type() string  { return CreateScavengeConst }
func (msg MsgCreateScavenge) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgCreateScavenge) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgCreateScavenge) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing creator address")
	}
	if msg.SolutionHash == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "solutionHash cannot be empty")
	}

	return nil
}
