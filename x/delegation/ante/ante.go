package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"cosmos-weighted-governance-sdk/x/delegation/keeper"
)

// NewAnteHandler creates a new ante handler with vesting delegation checks
func NewAnteHandler(dk keeper.Keeper, originalAnteHandler sdk.AnteHandler) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewVestingDelegationDecorator(dk),
		// The original ante handler should be run after our custom decorators
		// This is a terminal decorator that runs the original ante handler
		terminalDecorator{originalAnteHandler},
	)
}

// terminalDecorator is used to run the original ante handler after our custom decorators
type terminalDecorator struct {
	originalHandler sdk.AnteHandler
}

func (t terminalDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Ignore the next handler and run our original handler
	return t.originalHandler(ctx, tx, simulate)
}