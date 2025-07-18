package voting

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"cosmos-weighted-governance-sdk/testutil/sample"
	votingsimulation "cosmos-weighted-governance-sdk/x/voting/simulation"
	"cosmos-weighted-governance-sdk/x/voting/types"
)

// GenerateGenesisState creates a randomized GenState of the module.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	votingGenesis := types.GenesisState{
		Params:        types.DefaultParams(),
		PortId:        types.PortID,
		VoterRoleList: []types.VoterRole{{Id: 0, Creator: sample.AccAddress()}, {Id: 1, Creator: sample.AccAddress()}}, VoterRoleCount: 2,
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&votingGenesis)
}

// RegisterStoreDecoder registers a decoder.
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)
	const (
		opWeightMsgCreateVoterRole          = "op_weight_msg_voting"
		defaultWeightMsgCreateVoterRole int = 100
	)

	var weightMsgCreateVoterRole int
	simState.AppParams.GetOrGenerate(opWeightMsgCreateVoterRole, &weightMsgCreateVoterRole, nil,
		func(_ *rand.Rand) {
			weightMsgCreateVoterRole = defaultWeightMsgCreateVoterRole
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCreateVoterRole,
		votingsimulation.SimulateMsgCreateVoterRole(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgUpdateVoterRole          = "op_weight_msg_voting"
		defaultWeightMsgUpdateVoterRole int = 100
	)

	var weightMsgUpdateVoterRole int
	simState.AppParams.GetOrGenerate(opWeightMsgUpdateVoterRole, &weightMsgUpdateVoterRole, nil,
		func(_ *rand.Rand) {
			weightMsgUpdateVoterRole = defaultWeightMsgUpdateVoterRole
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgUpdateVoterRole,
		votingsimulation.SimulateMsgUpdateVoterRole(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))
	const (
		opWeightMsgDeleteVoterRole          = "op_weight_msg_voting"
		defaultWeightMsgDeleteVoterRole int = 100
	)

	var weightMsgDeleteVoterRole int
	simState.AppParams.GetOrGenerate(opWeightMsgDeleteVoterRole, &weightMsgDeleteVoterRole, nil,
		func(_ *rand.Rand) {
			weightMsgDeleteVoterRole = defaultWeightMsgDeleteVoterRole
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgDeleteVoterRole,
		votingsimulation.SimulateMsgDeleteVoterRole(am.authKeeper, am.bankKeeper, am.keeper, simState.TxConfig),
	))

	return operations
}

// ProposalMsgs returns msgs used for governance proposals for simulations.
func (am AppModule) ProposalMsgs(simState module.SimulationState) []simtypes.WeightedProposalMsg {
	return []simtypes.WeightedProposalMsg{}
}
