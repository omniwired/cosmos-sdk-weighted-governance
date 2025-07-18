package types

import (
	"fmt"

	host "github.com/cosmos/ibc-go/v10/modules/core/24-host"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
		PortId: PortID, VoterRoleList: []VoterRole{}}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	voterRoleIdMap := make(map[uint64]bool)
	voterRoleCount := gs.GetVoterRoleCount()
	for _, elem := range gs.VoterRoleList {
		if _, ok := voterRoleIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for voterRole")
		}
		if elem.Id >= voterRoleCount {
			return fmt.Errorf("voterRole id should be lower or equal than the last id")
		}
		voterRoleIdMap[elem.Id] = true
	}

	return gs.Params.Validate()
}
