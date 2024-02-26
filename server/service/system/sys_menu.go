package system

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"jykj-cmbp-dev-platform/server/global"
	"jykj-cmbp-dev-platform/server/model/common/request"
	"jykj-cmbp-dev-platform/server/model/system"
	systemRsp "jykj-cmbp-dev-platform/server/model/system/response"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: getMenuTreeMap
//@description: 获取路由总树map
//@param: authorityId string
//@return: treeMap map[string][]system.SysMenu, err error

type MenuService struct{}

var MenuServiceApp = new(MenuService)

func (menuService *MenuService) getMenuTreeMap(menusName string, roleId string) (treeMap systemRsp.MenusList, err error) {
	//var allMenus []system.SysMenu
	//var baseMenu []system.SysBaseMenu
	//var btns []system.SysAuthorityBtn
	var allmenusList []systemRsp.MenusDetail
	var allbuttonList []systemRsp.ButtonDetail

	//treeMap = make(map[string][]system.SysMenu)
	//
	//var SysAuthorityMenus []system.SysAuthorityMenu

	var QUERY = global.CMBP_DB.Model(&system.Menus{})
	var lastMenus system.Menus

	if menusName != "" {
		err = QUERY.Where("menus_name = ?", menusName).First(&lastMenus).Error
		if err != nil {
			return systemRsp.MenusList{}, err
		}
	}

	var buttonList []system.Menus
	var menuList []system.Menus

	if lastMenus.LastMenu != "" {
		err = QUERY.Where("type = ?", 2).Where("last_menu = ?", lastMenus.LastMenu).Order("create_time").Find(&menuList).Error
		if err != nil {
			return systemRsp.MenusList{}, err
		}
		err = global.CMBP_DB.Model(&system.Menus{}).Where("type = ?", 3).Where("last_menu = ?", lastMenus.LastMenu).Order("create_time").Find(&buttonList).Error
		if err != nil {
			return systemRsp.MenusList{}, err
		}
	} else {
		err = QUERY.Where("type = ?", 2).Order("create_time").Find(&menuList).Error
		if err != nil {
			return systemRsp.MenusList{}, err
		}
		err = global.CMBP_DB.Model(&system.Menus{}).Where("type = ?", 3).Order("create_time").Find(&buttonList).Error
		if err != nil {
			return systemRsp.MenusList{}, err
		}
	}

	var roleName string
	err = global.CMBP_DB.Model(system.Roles{}).Where("id = ? ", roleId).Pluck("name", &roleName).Error
	if err != nil {
		return systemRsp.MenusList{}, err
	}

	for _, menu := range menuList {
		var roleList []string
		json.Unmarshal([]byte(menu.RoleList), &roleList)
		for _, role := range roleList {
			if roleName == role {
				menuDetail := systemRsp.MenusDetail{
					MenuID:      menu.ID,
					Type:        menu.Type,
					Name:        menu.Name,
					Url:         menu.Url,
					AssemblyUrl: menu.AssemblyUrl,
					Icon:        menu.Icon,
					IsRouting:   menu.IsRouting,
				}
				allmenusList = append(allmenusList, menuDetail)
			}
		}
	}

	for _, button := range buttonList {
		var roleList []string
		json.Unmarshal([]byte(button.RoleList), &roleList)
		for _, role := range roleList {
			if roleName == role {
				buttonDetail := systemRsp.ButtonDetail{
					MenuID:      button.ID,
					Type:        button.Type,
					Name:        button.Name,
					Url:         button.Url,
					AssemblyUrl: button.AssemblyUrl,
					Icon:        button.Icon,
					IsRouting:   button.IsRouting,
				}
				allbuttonList = append(allbuttonList, buttonDetail)
			}
		}
	}
	return systemRsp.MenusList{Button: allbuttonList, Menus: allmenusList}, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserMenu
//@description: 获取动态菜单树
//@param: authorityId string
//@return: menus []system.SysMenu, err error

func (menuService *MenuService) GetUserMenu(menusName string, roleId string) (menus systemRsp.MenusList, err error) {
	menuTree, err := menuService.getMenuTreeMap(menusName, roleId)
	//menus = menuTree["0"]
	//for i := 0; i < len(menus); i++ {
	//	err = menuService.getChildrenList(&menus[i], menuTree)
	//}
	return menuTree, err
}

func (menuService *MenuService) GetMenuTree(menu_id string, roleId string) (menus systemRsp.MenusList, err error) {
	menuTree, err := menuService.getMenuTreeMap(menu_id, roleId)
	//menus = menuTree["0"]
	//for i := 0; i < len(menus); i++ {
	//	err = menuService.getChildrenList(&menus[i], menuTree)
	//}
	return menuTree, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: getChildrenList
//@description: 获取子菜单
//@param: menu *model.SysMenu, treeMap map[string][]model.SysMenu
//@return: err error

func (menuService *MenuService) getChildrenList(menu *system.SysMenu, treeMap map[string][]system.SysMenu) (err error) {
	menu.Children = treeMap[menu.MenuId]
	for i := 0; i < len(menu.Children); i++ {
		err = menuService.getChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetInfoList
//@description: 获取路由分页
//@return: list interface{}, total int64,err error

func (menuService *MenuService) GetInfoList() (list interface{}, total int64, err error) {
	var menuList []system.SysBaseMenu
	treeMap, err := menuService.getBaseMenuTreeMap()
	menuList = treeMap["0"]
	for i := 0; i < len(menuList); i++ {
		err = menuService.getBaseChildrenList(&menuList[i], treeMap)
	}
	return menuList, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: getBaseChildrenList
//@description: 获取菜单的子菜单
//@param: menu *model.SysBaseMenu, treeMap map[string][]model.SysBaseMenu
//@return: err error

func (menuService *MenuService) getBaseChildrenList(menu *system.SysBaseMenu, treeMap map[string][]system.SysBaseMenu) (err error) {
	//menu.Children = treeMap[strconv.Itoa(int(menu.ID))]
	menu.Children = treeMap[menu.ID]
	for i := 0; i < len(menu.Children); i++ {
		err = menuService.getBaseChildrenList(&menu.Children[i], treeMap)
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: AddBaseMenu
//@description: 添加基础路由
//@param: menu model.SysBaseMenu
//@return: error

func (menuService *MenuService) AddBaseMenu(menu system.SysBaseMenu) error {
	if !errors.Is(global.CMBP_DB.Where("name = ?", menu.Name).First(&system.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	return global.CMBP_DB.Create(&menu).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: getBaseMenuTreeMap
//@description: 获取路由总树map
//@return: treeMap map[string][]system.SysBaseMenu, err error

func (menuService *MenuService) getBaseMenuTreeMap() (treeMap map[string][]system.SysBaseMenu, err error) {
	var allMenus []system.SysBaseMenu
	treeMap = make(map[string][]system.SysBaseMenu)
	err = global.CMBP_DB.Order("sort").Preload("MenuBtn").Preload("Parameters").Find(&allMenus).Error
	for _, v := range allMenus {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return treeMap, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetBaseMenuTree
//@description: 获取基础路由树
//@return: menus []system.SysBaseMenu, err error

func (menuService *MenuService) GetBaseMenuTree() (menus []system.SysBaseMenu, err error) {
	treeMap, err := menuService.getBaseMenuTreeMap()
	menus = treeMap["0"]
	for i := 0; i < len(menus); i++ {
		err = menuService.getBaseChildrenList(&menus[i], treeMap)
	}
	return menus, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: AddMenuAuthority
//@description: 为角色增加menu树
//@param: menus []model.SysBaseMenu, authorityId string
//@return: err error

func (menuService *MenuService) AddMenuAuthority(menus []system.SysBaseMenu, authorityId string) (err error) {
	var auth system.SysAuthority
	auth.AuthorityId = authorityId
	auth.SysBaseMenus = menus
	err = AuthorityServiceApp.SetMenuAuthority(&auth)
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetMenuAuthority
//@description: 查看当前角色树
//@param: info *request.GetAuthorityId
//@return: menus []system.SysMenu, err error

func (menuService *MenuService) GetMenuAuthority(info *request.GetAuthorityId) (menus []system.SysMenu, err error) {
	var baseMenu []system.SysBaseMenu
	var SysAuthorityMenus []system.SysAuthorityMenu
	err = global.CMBP_DB.Where("sys_authority_authority_id = ?", info.AuthorityId).Find(&SysAuthorityMenus).Error
	if err != nil {
		return
	}

	var MenuIds []string

	for i := range SysAuthorityMenus {
		MenuIds = append(MenuIds, SysAuthorityMenus[i].MenuId)
	}

	err = global.CMBP_DB.Where("id in (?) ", MenuIds).Order("sort").Find(&baseMenu).Error

	for i := range baseMenu {
		menus = append(menus, system.SysMenu{
			SysBaseMenu: baseMenu[i],
			AuthorityId: info.AuthorityId,
			MenuId:      baseMenu[i].ID,
			Parameters:  baseMenu[i].Parameters,
		})
	}
	return menus, err
}

// UserAuthorityDefaultRouter 用户角色默认路由检查
//
//	Author [SliverHorn](https://github.com/SliverHorn)
func (menuService *MenuService) UserAuthorityDefaultRouter(user *system.Users) {
	var menuIds []string
	err := global.CMBP_DB.Model(&system.RoleMenus{}).Joins("JOIN t_user_roles ON t_role_menus.role_id = t_user_roles.role_id").Where("t_user_roles.user_id = ?", user.ID).Pluck("t_role_menus.menu_id", &menuIds).Error
	if err != nil {
		return
	}
	//var am system.SysBaseMenu
	//err = global.CMBP_DB.First(&am, "name = ? and id in (?)", user.Authority.DefaultRouter, menuIds).Error
	//if errors.Is(err, gorm.ErrRecordNotFound) {
	//	user.Authority.DefaultRouter = "404"
	//}
}
