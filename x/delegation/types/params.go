package types

const (
	// DefaultStakeDenom is the default denomination for staking
	DefaultStakeDenom = "stake"
)

// NewParams creates a new Params instance.
func NewParams(stakeDenom string) Params {
	return Params{
		StakeDenom: stakeDenom,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(DefaultStakeDenom)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if p.StakeDenom == "" {
		return ErrInvalidStakeDenom
	}
	return nil
}
