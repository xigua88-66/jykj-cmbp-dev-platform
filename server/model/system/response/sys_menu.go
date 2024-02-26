package response

import "jykj-cmbp-dev-platform/server/model/system"

type SysMenusResponse struct {
	Menus []system.SysMenu `json:"menus"`
}

type SysBaseMenusResponse struct {
	Menus []system.SysBaseMenu `json:"menus"`
}

type SysBaseMenuResponse struct {
	Menu system.SysBaseMenu `json:"menu"`
}

type MenusList struct {
	Button []ButtonDetail `json:"button"`
	Menus  []MenusDetail  `json:"menus"`
}

type MenusDetail struct {
	MenuID      string `json:"menu_id"`
	Type        int    `json:"type"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	AssemblyUrl string `json:"assembly_url"`
	Icon        string `json:"icon"`
	IsRouting   int    `json:"is_routing"`
}

type ButtonDetail struct {
	MenuID      string `json:"menu_id"`
	Type        int    `json:"type"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	AssemblyUrl string `json:"assembly_url"`
	Icon        string `json:"icon"`
	IsRouting   int    `json:"is_routing"`
}
