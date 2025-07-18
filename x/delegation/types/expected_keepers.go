package types

import (
	"context"
	"time"

	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AuthKeeper defines the expected interface for the Auth module.
type AuthKeeper interface {
	AddressCodec() address.Codec
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI // only used for simulation
	// Methods imported from account should be defined here
}

// VestingAccount defines the expected interface for vesting accounts.
type VestingAccount interface {
	sdk.AccountI
	GetVestedCoins(time.Time) sdk.Coins
	GetVestingCoins(time.Time) sdk.Coins
	GetOriginalVesting() sdk.Coins
	GetDelegatedFree() sdk.Coins
	GetDelegatedVesting() sdk.Coins
}

// BankKeeper defines the expected interface for the Bank module.
type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

// ParamSubspace defines the expected Subspace interface for parameters.
type ParamSubspace interface {
	Get(context.Context, []byte, interface{})
	Set(context.Context, []byte, interface{})
}
