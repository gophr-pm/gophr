package main

// orderModulesByDeps arranges a list of modules in the order that puts the
// most depended upon thing first and the least depended upon thing last.
func orderModulesByDeps(
	modules map[string]*module,
	excludedModuleNames map[string]bool,
	reverse bool) []*module {
	// Put the all the modules in a availability map first.
	availableModules := map[string]*module{}
	for k, v := range modules {
		var excluded bool
		if excludedModuleNames != nil {
			_, excluded = excludedModuleNames[k]
		}

		// Only make available if not excluded.
		if !excluded {
			availableModules[k] = v
		}
	}

	// Insert all modules into ordered slice.
	orderedModules := []*module{}
	for _, m := range modules {
		insertIntoOrderedModules(m, &orderedModules, availableModules)
	}

	// Reverse the slice if necessary.
	if reverse {
		reverseOrderedModules := make([]*module, len(orderedModules))
		for i := len(orderedModules) - 1; i >= 0; i-- {
			reverseOrderedModules = append(reverseOrderedModules, orderedModules[i])
		}
		return reverseOrderedModules
	}

	return orderedModules
}

// insertOrderedModule insert module into an ordered slice.
func insertIntoOrderedModules(
	m *module,
	orderedModules *[]*module,
	availableModules map[string]*module) {
	// Exit immediately if not available.
	if _, mAvailable := availableModules[m.name]; !mAvailable {
		return
	}

	// Take care of all the deps first.
	if len(m.deps) > 0 {
		for _, depName := range m.deps {
			// Only act if the dep is available.
			if dep, available := availableModules[depName]; available {
				// Insert it into the ordered module.
				insertIntoOrderedModules(dep, orderedModules, availableModules)
			}
		}
	}

	// Now that the deps have been handled, insert into the slice.
	*orderedModules = append(*orderedModules, m)
	// Remove this module from the availability map.
	delete(availableModules, m.name)
}
