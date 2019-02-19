package modules

import (
	"golang-fave/engine/wrapper"
)

func (this *Modules) RegisterModule_Some() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "some",
		Name:   "Some Module",
		Order:  1,
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		// Back-end
		return this.getSidebarModules(wrap), "Some", "Some Sidebar"
	})
}

func (this *Modules) RegisterModule_More() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "more",
		Name:   "More Module",
		Order:  2,
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		// Back-end
		return this.getSidebarModules(wrap), "More", "More Sidebar"
	})
}

func (this *Modules) RegisterModule_System() *Module {
	return this.newModule(MInfo{
		WantDB: true,
		Mount:  "system",
		Name:   "System Module",
		Order:  800,
		System: true,
	}, nil, func(wrap *wrapper.Wrapper) (string, string, string) {
		// Back-end
		return this.getSidebarModules(wrap), "System", "System Sidebar"
	})
}
