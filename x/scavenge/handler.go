package scavenge

import (
	"fmt"
	"github.com/SiddarthVijay/scavenge/x/scavenge/types"
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler creates an sdk.Handler for all the scavenge type messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgCreateScavenge:
			return handleMsgCreateScavenge(ctx, k, msg)
		case MsgCommitSolution:
			return handleMsgCommitSolution(ctx, k, msg)
		case MsgRevealSolution:
			return handleMsgRevealSolution(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", ModuleName,  msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgCreateScavenge(ctx sdk.Context, k Keeper, msg MsgCreateScavenge) (*sdk.Result, error) {
	var scavenge = types.Scavenge{
		Creator:      msg.Creator,
		Description:  msg.Description,
		SolutionHash: msg.SolutionHash,
		Reward:       msg.Reward,
	}

	// Check if the scavenge exists already
	_, err := k.GetScavenge(ctx, msg.SolutionHash)
	if err == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Scavenge with that solution hash already exists")
	}

	// Send reward coins to the escrow account
	moduleAcct := sdk.AccAddress(crypto.AddressHash([]byte(types.ModuleName)))
	sdkError := k.CoinKeeper.SendCoins(ctx, scavenge.Creator, moduleAcct, scavenge.Reward)

	if sdkError != nil {
		return nil, sdkError
	}

	k.SetScavenge(ctx, scavenge)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.EventTypeCreateScavenge),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
			sdk.NewAttribute(types.AttributeDescription, msg.Description),
			sdk.NewAttribute(types.AttributeSolutionHash, msg.SolutionHash),
			sdk.NewAttribute(types.AttributeReward, msg.Reward.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handle<Action> does x
func handleMsgCommitSolution(ctx sdk.Context, k Keeper, msg MsgCommitSolution) (*sdk.Result, error) {
	var commit = types.Commit{
		Scavenger:             msg.Scavenger,
		SolutionHash:          msg.SolutionHash,
		SolutionScavengerHash: msg.SolutionScavengerHash,
	}

	_, err := k.GetCommit(ctx, msg.SolutionScavengerHash)
	if err == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Commit with that hash already exists")
	}

	k.SetCommit

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handle<Action> does x
func handleMsg<Action>(ctx sdk.Context, k Keeper, msg Msg<Action>) (*sdk.Result, error) {
	err := k.<Action>(ctx, msg.ValidatorAddr)
	if err != nil {
		return nil, err
	}

	// TODO: Define your msg events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
