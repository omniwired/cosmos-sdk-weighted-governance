package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmos-weighted-governance-sdk/x/voting/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateVoterRole(ctx context.Context, msg *types.MsgCreateVoterRole) (*types.MsgCreateVoterRoleResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid creator address: %s", err))
	}

	// Validate voter role parameters
	if err := k.ValidateVoterRole(msg.Address, msg.Role, msg.Multiplier); err != nil {
		return nil, err
	}

	// Check if address already has a voter role
	if k.HasVoterRole(ctx, msg.Address) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("address %s already has a voter role", msg.Address))
	}

	nextId, err := k.VoterRoleSeq.Next(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "failed to get next id")
	}

	var voterRole = types.VoterRole{
		Id:         nextId,
		Creator:    msg.Creator,
		Address:    msg.Address,
		Role:       msg.Role,
		Multiplier: msg.Multiplier,
		AddedAt:    msg.AddedAt,
		AddedBy:    msg.AddedBy,
	}

	if err = k.VoterRole.Set(
		ctx,
		nextId,
		voterRole,
	); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to set voterRole")
	}

	return &types.MsgCreateVoterRoleResponse{
		Id: nextId,
	}, nil
}

func (k msgServer) UpdateVoterRole(ctx context.Context, msg *types.MsgUpdateVoterRole) (*types.MsgUpdateVoterRoleResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid creator address: %s", err))
	}

	// Validate voter role parameters
	if err := k.ValidateVoterRole(msg.Address, msg.Role, msg.Multiplier); err != nil {
		return nil, err
	}

	// Checks that the element exists
	val, err := k.VoterRole.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get voterRole")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var voterRole = types.VoterRole{
		Creator:    msg.Creator,
		Id:         msg.Id,
		Address:    msg.Address,
		Role:       msg.Role,
		Multiplier: msg.Multiplier,
		AddedAt:    msg.AddedAt,
		AddedBy:    msg.AddedBy,
	}

	if err := k.VoterRole.Set(ctx, msg.Id, voterRole); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update voterRole")
	}

	return &types.MsgUpdateVoterRoleResponse{}, nil
}

func (k msgServer) DeleteVoterRole(ctx context.Context, msg *types.MsgDeleteVoterRole) (*types.MsgDeleteVoterRoleResponse, error) {
	if _, err := k.addressCodec.StringToBytes(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Checks that the element exists
	val, err := k.VoterRole.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get voterRole")
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	if err := k.VoterRole.Remove(ctx, msg.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete voterRole")
	}

	return &types.MsgDeleteVoterRoleResponse{}, nil
}
