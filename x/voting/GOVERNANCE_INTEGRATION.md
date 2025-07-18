# Governance Integration Guide

This guide explains how to integrate the voting module's weighted voting system with the Cosmos SDK governance module.

## Overview

The voting module provides role-based voting multipliers that can be applied to governance votes. This allows different participants to have varying influence based on their role in the ecosystem.

## Integration Steps

### 1. Update Your App Configuration

In your `app.go`, after initializing the voting keeper and governance keeper:

```go
import (
    votingkeeper "cosmos-weighted-governance-sdk/x/voting/keeper"
    // ... other imports
)

// After creating keepers...

// Create the governance hooks wrapper
govHooks := votingkeeper.NewGovHooksWrapper(
    app.VotingKeeper,
    nil, // or existing gov hooks if you have any
)

// Set the hooks on the governance keeper
app.GovKeeper.SetHooks(govtypes.NewMultiGovHooks(govHooks))
```

### 2. Custom Tally Handler (Optional)

For full integration with weighted voting during tally:

```go
// Create custom tally handler
tallyHandler := votingkeeper.NewWeightedTallyHandler(
    app.VotingKeeper,
    app.StakingKeeper,
)

// Note: The current gov module doesn't support custom tally handlers directly
// You would need to modify the gov module or create a wrapper module
```

### 3. Governance Proposal Integration

To manage voter roles through governance:

1. Create governance proposals to add/update/delete voter roles
2. The proposal executor should have the governance module account address
3. All role management will go through the standard governance process

Example proposal:

```json
{
  "@type": "/cosmos.gov.v1beta1.MsgSubmitProposal",
  "content": {
    "@type": "/cosmos.gov.v1beta1.TextProposal",
    "title": "Add Core Contributor Role",
    "description": "Grant core contributor role with 2x voting multiplier to cosmos1..."
  },
  "initial_deposit": [{"denom": "stake", "amount": "1000000"}],
  "proposer": "cosmos1..."
}
```

## Current Limitations

1. **Tally Integration**: The current implementation emits events for tracking but doesn't directly modify vote tallying. Full integration would require:
   - Custom governance module that supports weighted tallying
   - Or modifications to the standard gov module

2. **Vote Storage**: Multipliers are calculated during voting but applying them during tally requires access to vote storage, which may need gov module modifications.

3. **Historical Votes**: Multipliers are applied based on current roles, not historical roles at vote time.

## Event Monitoring

The system emits the following events for monitoring:

- `vote_multiplier_applied`: Emitted when a vote is cast with a multiplier
  - `proposal_id`: The governance proposal ID
  - `voter`: The voter's address
  - `multiplier`: The applied multiplier

## Security Considerations

1. Only the governance module account can manage voter roles
2. Rate limiting prevents spam attacks
3. Maximum one role per address by default (configurable)
4. 5-minute cooldown between role creations

## Future Enhancements

1. **Full Tally Integration**: Modify gov module to support custom tally handlers
2. **Historical Multipliers**: Store multipliers at vote time for accurate historical tallying
3. **Delegation Support**: Apply multipliers to delegated voting power
4. **Vote Weight Display**: Show effective voting power in governance UI