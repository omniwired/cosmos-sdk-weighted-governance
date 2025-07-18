package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	"cosmos-weighted-governance-sdk/x/voting/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	// Address capable of executing a MsgUpdateParams message.
	// Typically, this should be the x/gov module account.
	authority []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	Port collections.Item[string]

	ibcKeeperFn  func() *ibckeeper.Keeper
	VoterRoleSeq collections.Sequence
	VoterRole    collections.Map[uint64, types.VoterRole]
	// LastRoleCreationTime tracks the last time a role was created (for rate limiting)
	LastRoleCreationTime collections.Item[int64]
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	ibcKeeperFn func() *ibckeeper.Keeper,

) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,

		ibcKeeperFn:          ibcKeeperFn,
		Port:                 collections.NewItem(sb, types.PortKey, "port", collections.StringValue),
		Params:               collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		VoterRole:            collections.NewMap(sb, types.VoterRoleKey, "voterRole", collections.Uint64Key, codec.CollValue[types.VoterRole](cdc)),
		VoterRoleSeq:         collections.NewSequence(sb, types.VoterRoleCountKey, "voterRoleSequence"),
		LastRoleCreationTime: collections.NewItem(sb, collections.NewPrefix([]byte("last_role_creation")), "lastRoleCreation", collections.Int64Value),
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() []byte {
	return k.authority
}
