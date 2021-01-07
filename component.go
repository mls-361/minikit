/*
------------------------------------------------------------------------------------------------------------------------
####### minikit ####### (c) 2020-2021 mls-361 ###################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package minikit

import "flag"

type (
	// Component AFAIRE.
	Component interface {
		Description() string
		Category() string
		Dependencies() []string
		Configure(fs *flag.FlagSet)
		Initialize(m *Manager) error
		IsBuilt() bool
		Build(m *Manager) error
		Built()
		Close()
	}

	// Base AFAIRE.
	Base struct {
		description string
		category    string
		built       bool
	}
)

// NewBase AFAIRE.
func NewBase(description, category string) *Base {
	return &Base{
		description: description,
		category:    category,
	}
}

// Name AFAIRE.
func (cb *Base) Description() string {
	return cb.description
}

// Category AFAIRE.
func (cb *Base) Category() string {
	return cb.category
}

// Dependencies AFAIRE.
func (cb *Base) Dependencies() []string {
	return []string{}
}

// Configure AFAIRE.
func (cb *Base) Configure(_ *flag.FlagSet) {}

// Initialize AFAIRE.
func (cb *Base) Initialize(_ *Manager) error {
	return nil
}

// IsBuilt AFAIRE.
func (cb *Base) IsBuilt() bool {
	return cb.built
}

// Build AFAIRE.
func (cb *Base) Build(_ *Manager) error {
	return nil
}

// Built AFAIRE.
func (cb *Base) Built() {
	cb.built = true
}

// Close AFAIRE.
func (cb *Base) Close() {}

/*
######################################################################################################## @(°_°)@ #######
*/
