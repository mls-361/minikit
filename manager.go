/*
------------------------------------------------------------------------------------------------------------------------
####### minikit ####### (c) 2020-2021 mls-361 ###################################################### MIT License #######
------------------------------------------------------------------------------------------------------------------------
*/

package minikit

import (
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/mls-361/failure"
)

const (
	// CategoryRun AFAIRE.
	CategoryRunner = "runner"
)

type (
	// Application AFAIRE.
	Application interface {
		Name() string
		Debug() int
	}

	// Manager AFAIRE.
	Manager struct {
		application Application
		components  map[string]Component
		closeList   []Component
	}

	// PluginCb AFAIRE.
	PluginCb func(pSym plugin.Symbol) error

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
				Set("component", d.Description()).
				Msg("a component of this category already exists") /////////////////////////////////////////////////////
		}

		m.components[category] = c
	}

	return nil
}

// AddPlugin AFAIRE.
func (m *Manager) AddPlugin(path, symName string, callback PluginCb) error {
	p, err := plugin.Open(path)
	if err != nil {
		return failure.New(err).
			Set("plugin", path).
			Msg("impossible to open this plugin") //////////////////////////////////////////////////////////////////////
	}

	pSym, err := p.Lookup(symName)
	if err != nil {
		return failure.New(err).
			Set("plugin", path).
			Set("symbol", symName).
			Msg("this plugin does not export this symbol") /////////////////////////////////////////////////////////////
	}

	if err := callback(pSym); err != nil {
		return failure.New(err).
			Set("plugin", path).
			Msg("plugin error") ////////////////////////////////////////////////////////////////////////////////////////
	}

	return nil
}

// AddPlugins AFAIRE.
func (m *Manager) AddPlugins(dirname, symName string, callback PluginCb) error {
	paths, err := filepath.Glob(filepath.Join(dirname, fmt.Sprintf("%s-*.so", m.application.Name())))
	if err != nil {
		return err
	}

	for _, path := range paths {
		if err := m.AddPlugin(path, symName, callback); err != nil {
			return err
		}
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
			return failure.New(err).Set("component", c.Description()).Msg("initialization error") //////////////////////
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
			fmt.Printf("=== Component: %s ==> CLOSED\n", c.Description()) //::::::::::::::::::::::::::::::::::::::::::::
		}
	}
}

func (m *Manager) recursiveBuild(snitch map[string]bool, c Component) error {
	snitch[c.Category()] = false

	if m.appDebug() > 1 {
		fmt.Printf("=== Component: %s ==> TO BUILD\n", c.Description()) //::::::::::::::::::::::::::::::::::::::::::::::
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
				Set("component1", c.Description()).
				Set("component2", d.Description()).
				Msg("these two components are interdependent") /////////////////////////////////////////////////////////
		}

		if err := m.recursiveBuild(snitch, d); err != nil {
			return err
		}
	}

	if err := c.Build(m); err != nil {
		return failure.New(err).Set("component", c.Description()).Msg("build error") ///////////////////////////////////
	}

	c.Built()

	m.closeList = append(m.closeList, c)
	snitch[c.Category()] = true

	if m.appDebug() > 1 {
		fmt.Printf("=== Component: %s ==> BUILT\n", c.Description()) //:::::::::::::::::::::::::::::::::::::::::::::::::
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
					Set("component", c.Description()).
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
		Set("name", c.Description()).
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
