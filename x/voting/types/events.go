package types

// Event types
const (
	EventTypeVoterRoleCreated = "voter_role_created"
	EventTypeVoterRoleUpdated = "voter_role_updated"
	EventTypeVoterRoleDeleted = "voter_role_deleted"

	AttributeKeyRoleID      = "role_id"
	AttributeKeyAddress     = "address"
	AttributeKeyRole        = "role"
	AttributeKeyMultiplier  = "multiplier"
	AttributeKeyAddedBy     = "added_by"
	AttributeKeyAddedAt     = "added_at"
	AttributeKeyDeletedBy   = "deleted_by"
	AttributeKeyUpdatedBy   = "updated_by"
)