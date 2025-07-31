package keeper

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"cosmos-weighted-governance-sdk/x/voting/types"

	"cosmossdk.io/collections"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateVoterRole(ctx context.Context, msg *types.MsgCreateVoterRole) (*types.MsgCreateVoterRoleResponse, error) {
	creator, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid creator address: %s", err))
	}

	// only gov can do this
	if !bytes.Equal(k.GetAuthority(), creator) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "only governance account can create voter roles; expected %s, got %s", expectedAuthorityStr, msg.Creator)
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get module params")
	}

	// rate limiting check
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	currentTime := sdkCtx.BlockTime().Unix()
	
	lastCreationTime, err := k.LastRoleCreationTime.Get(ctx)
	if err == nil && params.RoleCreationCooldown > 0 {
		timeSinceLastCreation := currentTime - lastCreationTime
		if timeSinceLastCreation < int64(params.RoleCreationCooldown) {
			return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
				"role creation is rate limited: %d seconds remaining",
				int64(params.RoleCreationCooldown)-timeSinceLastCreation)
		}
	}

	if err := k.ValidateVoterRole(msg.Address, msg.Role, msg.Multiplier); err != nil {
		return nil, err
	}

	// check if they already have a role
	if k.HasVoterRole(ctx, msg.Address) {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("address %s already has a voter role", msg.Address))
	}
	
	// max roles check
	roleCount := k.CountRolesForAddress(ctx, msg.Address)
	if roleCount >= params.MaxVoterRolesPerAddress {
		return nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
			"address %s already has maximum number of roles (%d)",
			msg.Address, params.MaxVoterRolesPerAddress)
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

	// update last creation time
	if err = k.LastRoleCreationTime.Set(ctx, currentTime); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to update last creation time")
	}

	// TODO: maybe add metrics here for role creation tracking?
	
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVoterRoleCreated,
			sdk.NewAttribute(types.AttributeKeyRoleID, fmt.Sprintf("%d", nextId)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
			sdk.NewAttribute(types.AttributeKeyRole, msg.Role),
			sdk.NewAttribute(types.AttributeKeyMultiplier, msg.Multiplier),
			sdk.NewAttribute(types.AttributeKeyAddedBy, msg.AddedBy),
			sdk.NewAttribute(types.AttributeKeyAddedAt, fmt.Sprintf("%d", msg.AddedAt)),
		),
	)

	return &types.MsgCreateVoterRoleResponse{
		Id: nextId,
	}, nil
}

func (k msgServer) UpdateVoterRole(ctx context.Context, msg *types.MsgUpdateVoterRole) (*types.MsgUpdateVoterRoleResponse, error) {
	creator, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid creator address: %s", err))
	}

	// same check as create - only gov
	if !bytes.Equal(k.GetAuthority(), creator) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "only governance account can update voter roles; expected %s, got %s", expectedAuthorityStr, msg.Creator)
	}

	if err := k.ValidateVoterRole(msg.Address, msg.Role, msg.Multiplier); err != nil {
		return nil, err
	}

	// make sure it exists first
	_, err = k.VoterRole.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get voterRole")
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

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVoterRoleUpdated,
			sdk.NewAttribute(types.AttributeKeyRoleID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(types.AttributeKeyAddress, msg.Address),
			sdk.NewAttribute(types.AttributeKeyRole, msg.Role),
			sdk.NewAttribute(types.AttributeKeyMultiplier, msg.Multiplier),
			sdk.NewAttribute(types.AttributeKeyUpdatedBy, msg.Creator),
		),
	)

	return &types.MsgUpdateVoterRoleResponse{}, nil
}

func (k msgServer) DeleteVoterRole(ctx context.Context, msg *types.MsgDeleteVoterRole) (*types.MsgDeleteVoterRoleResponse, error) {
	creator, err := k.addressCodec.StringToBytes(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err))
	}

	// Check if creator is the governance module account
	if !bytes.Equal(k.GetAuthority(), creator) {
		expectedAuthorityStr, _ := k.addressCodec.BytesToString(k.GetAuthority())
		return nil, errorsmod.Wrapf(types.ErrInvalidSigner, "only governance account can delete voter roles; expected %s, got %s", expectedAuthorityStr, msg.Creator)
	}

	// Checks that the element exists
	val, err := k.VoterRole.Get(ctx, msg.Id)
	if err != nil {
		if errors.Is(err, collections.ErrNotFound) {
			return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
		}

		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to get voterRole")
	}

	// No need to check if msg creator matches val.Creator since only governance can delete

	if err := k.VoterRole.Remove(ctx, msg.Id); err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrLogic, "failed to delete voterRole")
	}

	// Emit event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeVoterRoleDeleted,
			sdk.NewAttribute(types.AttributeKeyRoleID, fmt.Sprintf("%d", msg.Id)),
			sdk.NewAttribute(types.AttributeKeyAddress, val.Address),
			sdk.NewAttribute(types.AttributeKeyDeletedBy, msg.Creator),
		),
	)

	return &types.MsgDeleteVoterRoleResponse{}, nil
}
