package types

import "fmt"

const (
	// DefaultMaxVoterRolesPerAddress is the default maximum number of roles per address
	DefaultMaxVoterRolesPerAddress uint32 = 1
	
	// DefaultRoleCreationCooldown is the default cooldown period in seconds
	DefaultRoleCreationCooldown uint32 = 300 // 5 minutes
)

// NewParams creates a new Params instance.
func NewParams(maxRolesPerAddress, cooldown uint32) Params {
	return Params{
		MaxVoterRolesPerAddress: maxRolesPerAddress,
		RoleCreationCooldown:    cooldown,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultMaxVoterRolesPerAddress, DefaultRoleCreationCooldown)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.MaxVoterRolesPerAddress == 0 {
		return fmt.Errorf("max voter roles per address must be greater than 0")
	}
	
	// Cooldown can be 0 to disable the feature
	return nil
}
