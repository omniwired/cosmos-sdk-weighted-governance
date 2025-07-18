# Cosmos Weighted Governance SDK

A Cosmos SDK blockchain implementing two advanced features: weighted voting governance and vesting-aware staking restrictions. Built with Cosmos SDK v0.53.3 and Ignite CLI v29.2.0.

## Two Core Features

This project shows how to build advanced Cosmos SDK modules by implementing two features that tackle real blockchain governance challenges:

### Feature 1: Weighted Voting Governance
Role-based voting where different participants have different voting power multipliers

### Feature 2: Vesting-Aware Staking Restrictions
Prevents staking of unvested tokens to maintain proper token distribution schedules

## Architecture

### Weighted Voting Module (`voting`)

The weighted voting system lets different participants have varying influence based on their role and contribution to the network.

Key features include role-based multipliers where core contributors get 2x voting power, validators get 1.5x, strategic partners get 1.8x, and community members get 1x. The system allows dynamic role management through governance proposals, validates all inputs properly, and provides efficient queries for role lookups and statistics.

Technical implementation uses Collections for state management, Protocol Buffers for message serialization, custom keeper methods for cross-module queries, and AutoCLI for command-line interaction.

### Vesting-Aware Staking Module (`delegation`)

The vesting-aware staking system prevents people from staking tokens they don't technically own yet, making sure token distribution schedules work as intended.

Key features include automatic detection of vesting accounts, real-time eligibility checks against vesting schedules, detailed reporting with vesting status and amounts, and smooth integration with auth and bank modules.

Technical implementation uses interface-based design for vesting account abstraction, context-aware validation using block time, comprehensive error handling, and gRPC/REST API endpoints.

## Module Interactions

The system shows several important Cosmos SDK patterns:
- Inter-module communication with secure keeper-to-keeper interactions
- State management using Collections framework for efficient data storage
- Query optimization with paginated responses and indexed lookups
- Transaction validation with multi-layer validation and ante handlers

## Getting Started

### Prerequisites
- Go 1.21+
- Ignite CLI v29.2.0+
- Node.js 18+ (for frontend development)

### Quick Start

```bash
# Clone and setup
git clone https://github.com/yourusername/enhanced-governance-staking
cd enhanced-governance-staking

# Start development chain
ignite chain serve

# Build binary
go build -o build/cosmos-weighted-governance-sdkd ./cmd/cosmos-weighted-governance-sdkd
```

### Example Usage

```bash
# Query voting module
cosmos-weighted-governance-sdkd query voting list-voter-role
cosmos-weighted-governance-sdkd query voting get-voter-role 1

# Create voter role
cosmos-weighted-governance-sdkd tx voting create-voter-role cosmos1... validator 1.5 $(date +%s) cosmos1...

# Check staking eligibility
cosmos-weighted-governance-sdkd query delegation staking-eligibility cosmos1...

# Query module parameters
cosmos-weighted-governance-sdkd query voting params
cosmos-weighted-governance-sdkd query delegation params
```

## Development

### Testing
```bash
# Run all tests
go test ./...

# Test specific module
go test ./x/voting/...
go test ./x/delegation/...

# Integration tests
go test ./app/...
```

### Code Generation
```bash
# Generate protobuf files
ignite generate proto-go

# Generate OpenAPI docs
ignite generate openapi

# Generate TypeScript client
ignite generate ts-client
```

## Technical Implementation Details

### Security Features Demonstrated
- Input validation on all user-facing endpoints
- Role-based access control for administrative functions
- Protection against integer overflow in multiplier calculations
- Secure cross-module communication patterns

### Performance Optimizations
- Optimized state queries using Collections framework
- Indexed lookups for role-based operations
- Efficient pagination for large datasets
- Minimal computational overhead in validation logic

### Architectural Patterns
- Modular architecture supporting independent scaling
- Efficient state pruning strategies
- Configurable parameters for different network sizes
- Compatible with IBC for cross-chain governance

## A Word About Democracy... and Staking

You know, it's funny how we've recreated the same old power dynamics in blockchain governance. Core contributors get 2x voting weight because apparently writing code makes you twice as wise about network decisions. Validators get 1.5x because running servers is apparently worth more than just holding tokens. It's like we've invented digital aristocracy, but with better documentation and open-source transparency!

As for the vesting restrictions on staking - well, someone had to be the fun police and make sure people can't stake tokens they don't technically own yet. Think of it as blockchain parental controls: "No, you can't stake your allowance until it's actually yours!" But hey, at least the smart contracts are impartial enforcers, unlike that one friend who always changes the rules mid-game.

## Technical Specifications

### Dependencies
- **Cosmos SDK**: v0.53.3 (latest stable)
- **CometBFT**: v0.38.x (Tendermint consensus)
- **IBC**: v8.x (Inter-Blockchain Communication)
- **Collections**: Latest (state management framework)

### Module Specifications
- **State Management**: Collections-based with optimized indexing
- **Message Validation**: Multi-layer validation with custom ante handlers
- **Query Engine**: gRPC/REST with pagination support
- **CLI Integration**: AutoCLI with custom command extensions

## Project Purpose

This is a portfolio project I built to showcase my Cosmos SDK development skills. It shows experience with:
- Custom module development with complex state management
- Cross-module communication and keeper interactions
- Protocol buffer message design and validation
- CLI integration and query optimization
- Blockchain governance and staking mechanisms

Note: This is an educational/portfolio project, not intended for production use.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Resources

- [Cosmos SDK Documentation](https://docs.cosmos.network/)
- [Ignite CLI Docs](https://docs.ignite.com/)
- [Cosmos Developer Portal](https://cosmos.network/developers/)
- [IBC Protocol](https://ibcprotocol.org/)

---

*Built with ❤️ using Cosmos SDK and the power of modular blockchain architecture*