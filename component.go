/*
------------------------------------------------------------------------------------------------------------------------
####### minikit ####### (c) 2020-2021 mls-361 ###################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package minikit

type (
	// Component AFAIRE.
	Component interface {
		Name() string
		Category() string
		Dependencies() []string
		IsBuilt() bool
		Initialize(m *Manager) error
		Build(m *Manager) error
		Close()
	}

	// Base AFAIRE.
	Base struct {
		name     string
		category string
		built    bool
	}
)

// NewBase AFAIRE.
func NewBase(name, category string) *Base {
	return &Base{
		name:     name,
		category: category,
	}
}

// Name AFAIRE.
func (cb *Base) Name() string {
	return cb.name
}

// Category AFAIRE.
func (cb *Base) Category() string {
	return cb.category
}

// Dependencies AFAIRE.
func (cb *Base) Dependencies() []string {
	return []string{}
}

// Built AFAIRE.
func (cb *Base) Built() {
	cb.built = true
}

// IsBuilt AFAIRE.
func (cb *Base) IsBuilt() bool {
	return cb.built
}

// Initialize AFAIRE.
func (cb *Base) Initialize(_ *Manager) error {
	return nil
}

// Build AFAIRE.
func (cb *Base) Build(_ *Manager) error {
	cb.Built()
	return nil
}

// Close AFAIRE.
func (cb *Base) Close() {}

/*
######################################################################################################## @(°_°)@ #######
*/
