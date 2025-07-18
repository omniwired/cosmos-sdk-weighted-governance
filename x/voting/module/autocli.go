package voting

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"cosmos-weighted-governance-sdk/x/voting/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "ListVoterRole",
					Use:       "list-voter-role",
					Short:     "List all VoterRole",
				},
				{
					RpcMethod:      "GetVoterRole",
					Use:            "get-voter-role [id]",
					Short:          "Gets a VoterRole by id",
					Alias:          []string{"show-voter-role"},
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateVoterRole",
					Use:            "create-voter-role [address] [role] [multiplier] [added-at] [added-by]",
					Short:          "Create VoterRole",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "address"}, {ProtoField: "role"}, {ProtoField: "multiplier"}, {ProtoField: "added_at"}, {ProtoField: "added_by"}},
				},
				{
					RpcMethod:      "UpdateVoterRole",
					Use:            "update-voter-role [id] [address] [role] [multiplier] [added-at] [added-by]",
					Short:          "Update VoterRole",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}, {ProtoField: "address"}, {ProtoField: "role"}, {ProtoField: "multiplier"}, {ProtoField: "added_at"}, {ProtoField: "added_by"}},
				},
				{
					RpcMethod:      "DeleteVoterRole",
					Use:            "delete-voter-role [id]",
					Short:          "Delete VoterRole",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "id"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
