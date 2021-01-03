/*
------------------------------------------------------------------------------------------------------------------------
####### minikit ####### (c) 2020-2021 mls-361 ###################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package minikit

import (
	"fmt"

	"github.com/mls-361/failure"
)

const (
	// CategoryRun AFAIRE.
	CategoryRunner = "runner"
)

type (
	// Application AFAIRE.
	Application interface {
		Debug() int
	}

	// Manager AFAIRE.
	Manager struct {
		application Application
		components  map[string]Component
		closeList   []Component
	}

	// Runner AFAIRE.
	Runner interface {
		Run(m *Manager) error
	}
)

// NewManager AFAIRE.
func NewManager(application Application) *Manager {
	return &Manager{
		application: application,
		components:  make(map[string]Component),
		closeList:   make([]Component, 0),
	}
}

func (m *Manager) appDebug() int {
	if m.application == nil {
		return -1
	}

	return m.application.Debug()
}

// AddComponents AFAIRE.
func (m *Manager) AddComponents(cList ...Component) error {
	for _, c := range cList {
		category := c.Category()

		d, ok := m.components[category]
		if ok {
			return failure.New(nil).
				Set("category", category).
				Set("component", d.Name()).
				Msg("a component of this category already exists") /////////////////////////////////////////////////////
		}

		m.components[category] = c
	}

	return nil
}

// GetComponent AFAIRE.
func (m *Manager) GetComponent(category string, mustExist bool) (Component, error) {
	c, ok := m.components[category]
	if !ok {
		if mustExist {
			return nil, failure.New(nil).
				Set("category", category).
				Msg("no component of this category exists") ////////////////////////////////////////////////////////////
		}

		return nil, nil
	}

	return c, nil
}

// Components AFAIRE.
func (m *Manager) Components() []Component {
	components := make([]Component, 0)

	for _, c := range m.components {
		components = append(components, c)
	}

	return components
}

// InitializeComponents AFAIRE.
func (m *Manager) InitializeComponents() error {
	for _, c := range m.Components() {
		if c.IsBuilt() {
			continue
		}

		if err := c.Initialize(m); err != nil {
			return failure.New(err).Set("component", c.Name()).Msg("initialization error") /////////////////////////////
		}
	}

	return nil
}

// CloseComponents AFAIRE.
func (m *Manager) CloseComponents() {
	last := len(m.closeList) - 1

	for i := range m.closeList {
		c := m.closeList[last-i]
		c.Close()

		if m.appDebug() > 1 {
			fmt.Printf("=== Component: %s ==> CLOSED\n", c.Name()) //:::::::::::::::::::::::::::::::::::::::::::::::::::
		}
	}
}

func (m *Manager) recursiveBuild(snitch map[string]bool, c Component) error {
	snitch[c.Category()] = false

	if m.appDebug() > 1 {
		fmt.Printf("=== Component: %s ==> TO BUILD\n", c.Name()) //:::::::::::::::::::::::::::::::::::::::::::::::::::::
	}

	for _, cc := range c.Dependencies() {
		d, err := m.GetComponent(cc, true)
		if err != nil {
			return err
		}

		if d.IsBuilt() {
			continue
		}

		done, ok := snitch[cc]
		if ok {
			if done {
				continue
			}

			return failure.New(nil).
				Set("component1", c.Name()).
				Set("component2", d.Name()).
				Msg("these two components are interdependent") /////////////////////////////////////////////////////////
		}

		if err := m.recursiveBuild(snitch, d); err != nil {
			return err
		}
	}

	if err := c.Build(m); err != nil {
		return failure.New(err).Set("component", c.Name()).Msg("build error") //////////////////////////////////////////
	}

	m.closeList = append(m.closeList, c)
	snitch[c.Category()] = true

	if m.appDebug() > 1 {
		fmt.Printf("=== Component: %s ==> BUILT\n", c.Name()) //::::::::::::::::::::::::::::::::::::::::::::::::::::::::
	}

	return nil
}

// BuildComponent AFAIRE.
func (m *Manager) BuildComponent(category string) error {
	c, err := m.GetComponent(category, true)
	if err != nil {
		return err
	}

	if c.IsBuilt() {
		return nil
	}

	return m.recursiveBuild(make(map[string]bool), c)
}

// BuildComponents AFAIRE.
func (m *Manager) BuildComponents() error {
	snitch := make(map[string]bool)

	for _, c := range m.Components() {
		if c.IsBuilt() {
			continue
		}

		done, ok := snitch[c.Category()]
		if ok {
			if !done {
				return failure.New(nil).
					Set("component", c.Name()).
					Msg("this error should not occur") /////////////////////////////////////////////////////////////////
			}

			continue
		}

		if err := m.recursiveBuild(snitch, c); err != nil {
			return err
		}
	}

	return nil
}

// InterfaceError AFAIRE.
func (m *Manager) InterfaceError(c Component) error {
	return failure.New(nil).
		Set("name", c.Name()).
		Set("category", c.Category()).
		Msg("this component does not implement the interface of its category") /////////////////////////////////////////
}

// Run AFAIRE.
func (m *Manager) Run() error {
	c, err := m.GetComponent(CategoryRunner, true)
	if err != nil {
		return err
	}

	cr, ok := c.(Runner)
	if !ok {
		return m.InterfaceError(c) /////////////////////////////////////////////////////////////////////////////////////
	}

	return cr.Run(m)
}

/*
######################################################################################################## @(°_°)@ #######
*/
