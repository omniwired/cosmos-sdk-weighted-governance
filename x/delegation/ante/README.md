# Delegation Module Ante Handler

This ante handler enforces vesting restrictions on staking operations.

## Integration

To integrate the vesting delegation ante handler into your application, add the following to your `app.go`:

```go
// In your app initialization, after setting up keepers:

// Get the existing ante handler
existingAnteHandler := authante.NewAnteHandler(
    authante.HandlerOptions{
        AccountKeeper:   app.AccountKeeper,
        BankKeeper:      app.BankKeeper,
        SignModeHandler: txConfig.SignModeHandler(),
        // ... other options
    },
)

// Wrap it with our vesting delegation decorator
anteHandler := delegationante.NewAnteHandler(
    app.DelegationKeeper,
    existingAnteHandler,
)

// Set the ante handler
app.SetAnteHandler(anteHandler)
```

## How It Works

The ante handler intercepts staking transactions and validates:

1. **MsgDelegate**: Ensures delegators can only stake vested tokens
2. **MsgBeginRedelegate**: Validates that redelegation amount doesn't exceed vested tokens
3. **MsgUndelegate**: Currently allows all undelegations (can be restricted if needed)

For vesting accounts, the handler checks:
- Current vesting schedule
- Amount of vested vs unvested tokens
- Requested delegation amount

If a vesting account attempts to delegate more than their vested amount, the transaction is rejected.