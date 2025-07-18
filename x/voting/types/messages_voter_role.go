package types

func NewMsgCreateVoterRole(creator string, address string, role string, multiplier string, addedAt int64, addedBy string) *MsgCreateVoterRole {
	return &MsgCreateVoterRole{
		Creator:    creator,
		Address:    address,
		Role:       role,
		Multiplier: multiplier,
		AddedAt:    addedAt,
		AddedBy:    addedBy,
	}
}

func NewMsgUpdateVoterRole(creator string, id uint64, address string, role string, multiplier string, addedAt int64, addedBy string) *MsgUpdateVoterRole {
	return &MsgUpdateVoterRole{
		Id:         id,
		Creator:    creator,
		Address:    address,
		Role:       role,
		Multiplier: multiplier,
		AddedAt:    addedAt,
		AddedBy:    addedBy,
	}
}

func NewMsgDeleteVoterRole(creator string, id uint64) *MsgDeleteVoterRole {
	return &MsgDeleteVoterRole{
		Id:      id,
		Creator: creator,
	}
}
