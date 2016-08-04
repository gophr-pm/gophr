package main

import (
	"errors"

	"gopkg.in/urfave/cli.v1"
)

var (
	allModuleID   = "all"
	allModuleDeps = []string{}
)

type allModule struct{}

func (m *allModule) id() string {
	return allModuleID
}

func (m *allModule) deps() []string {
	return allModuleDeps
}

func (m *allModule) dockerfile() string {
	return ""
}

func (m *allModule) containerMetadata() ([]dockerPortMapping, []dockerLinkMapping, []dockerVolumeMapping) {
	return nil, nil, nil
}

func (m *allModule) build(c *cli.Context, shallow bool) error {
	var failedModuleIds []string

	// Build the list of module ids to detect if this call failed.
	traverseModules(func(m module) {
		err := m.build(c, true)
		if err != nil {
			failedModuleIds = append(failedModuleIds, m.id())
		}
	})

	// Fail if any of the sub-routines failed.
	if len(failedModuleIds) > 0 {
		return newFailedModulesError("build", failedModuleIds)
	}

	return nil
}

func (m *allModule) start(c *cli.Context, shallow bool) error {
	var failedModuleIds []string

	// Build the list of module ids to detect if this call failed.
	traverseModules(func(m module) {
		err := m.start(c, true)
		if err != nil {
			failedModuleIds = append(failedModuleIds, m.id())
		}
	})

	// Fail if any of the sub-routines failed.
	if len(failedModuleIds) > 0 {
		return newFailedModulesError("start", failedModuleIds)
	}

	return nil
}

func (m *allModule) stop(c *cli.Context, shallow bool) error {
	var failedModuleIds []string

	// Build the list of module ids to detect if this call failed.
	traverseModules(func(m module) {
		err := m.stop(c, true)
		if err != nil {
			failedModuleIds = append(failedModuleIds, m.id())
		}
	})

	// Fail if any of the sub-routines failed.
	if len(failedModuleIds) > 0 {
		return newFailedModulesError("stop", failedModuleIds)
	}

	return nil
}

func (m *allModule) log(c *cli.Context, shallow bool) error {
	return errors.New("Cannot log every module at once.")
}

func (m *allModule) ssh(c *cli.Context, shallow bool) error {
	return errors.New("Cannot ssh into every module at once.")
}

func (m *allModule) test(c *cli.Context, shallow bool) error {
	var failedModuleIds []string

	// Build the list of module ids to detect if this call failed.
	traverseModules(func(m module) {
		err := m.test(c, true)
		if err != nil {
			failedModuleIds = append(failedModuleIds, m.id())
		}
	})

	// Fail if any of the sub-routines failed.
	if len(failedModuleIds) > 0 {
		return newFailedModulesError("test", failedModuleIds)
	}

	return nil
}

func (m *allModule) restart(c *cli.Context, shallow bool) error {
	if err := m.stop(c, false); err != nil {
		return err
	} else if err = m.start(c, false); err != nil {
		return err
	}

	return nil
}
