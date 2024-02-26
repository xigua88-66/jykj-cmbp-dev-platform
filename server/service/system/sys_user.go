package system

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
	"jykj-cmbp-dev-platform/server/global"
	"jykj-cmbp-dev-platform/server/model/common/request"
	"jykj-cmbp-dev-platform/server/model/system"
	"jykj-cmbp-dev-platform/server/utils"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: userInter system.SysUser, err error

type UserService struct{}

func (userService *UserService) Register(u system.SysUser) (userInter system.SysUser, err error) {
	var user system.SysUser
	if !errors.Is(global.CMBP_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return userInter, errors.New("用户名已注册")
	}
	// 否则 附加uuid 密码hash加密 注册
	u.Password = utils.BcryptHash(u.Password)
	u.UUID = uuid.Must(uuid.NewV4())
	err = global.CMBP_DB.Create(&u).Error
	return u, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.Users) (userInter *system.Users, err error) {
	if nil == global.CMBP_DB {
		return nil, fmt.Errorf("db not init")
	}

	var user system.Users
	err = global.CMBP_DB.Where("username = ?", u.Username).First(&user).Error
	fmt.Println("请求的密码是：", u.Password)
	fmt.Println("数据库的密码是：", user.Password, user.Username)

	// todo 登录公司平台

	if err == nil {
		if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
			return nil, errors.New("密码错误")
		}
		MenuServiceApp.UserAuthorityDefaultRouter(&user)

		if user.MineCode == "999999999" && !user.ExpireTime.IsZero() && time.Now().After(user.ExpireTime) && user.ExpireLoginNum > 0 {
			user.IsActive = false
			user.RootDisable = true
			err = global.CMBP_DB.Save(&user).Error
			if err != nil {
				return nil, errors.New("该角色修改权限时发生错误")
			}
			return nil, errors.New("该用户账号已过期，请联系管理员")
		} else if user.IsActive == false {
			return nil, errors.New("该用户已被管理员被禁用,请联系管理员")
		} else {
			if user.MineCode == "999999999" && !user.ExpireTime.IsZero() && time.Now().After(user.ExpireTime) {
				user.ExpireLoginNum = 1
				err = global.CMBP_DB.Save(&user).Error
				if err != nil {
					return nil, errors.New("该角色修改权限时发生错误")
				}
				return &user, errors.New("当前账号已经过期，这是您最后一次登录，请联系管理员")
			}
			return &user, err
		}
	}
	if strings.Contains(err.Error(), "record not found") {
		return nil, errors.New("用户不存在")
	}
	return nil, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ChangePassword
//@description: 修改用户密码
//@param: u *model.SysUser, newPassword string
//@return: userInter *model.SysUser,err error

func (userService *UserService) ChangePassword(u *system.Users, newPassword string) (userInter *system.Users, err error) {
	var user system.Users
	if err = global.CMBP_DB.Where("id = ?", u.ID).First(&user).Error; err != nil {
		return nil, err
	}
	if ok := utils.BcryptCheck(u.Password, user.Password); !ok {
		return nil, errors.New("原密码错误")
	}
	user.Password = utils.BcryptHash(newPassword)
	err = global.CMBP_DB.Save(&user).Error
	return &user, err

}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.CMBP_DB.Model(&system.SysUser{})
	var userList []system.SysUser
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Authority").Find(&userList).Error
	return userList, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserAuthority
//@description: 设置一个用户的权限
//@param: uuid uuid.UUID, authorityId string
//@return: err error

func (userService *UserService) SetUserAuthority(id string, authorityId string) (err error) {
	assignErr := global.CMBP_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", id, authorityId).First(&system.SysUserAuthority{}).Error
	if errors.Is(assignErr, gorm.ErrRecordNotFound) {
		return errors.New("该用户无此角色")
	}
	err = global.CMBP_DB.Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityId).Error
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserAuthorities
//@description: 设置一个用户的权限
//@param: id uint, authorityIds []string
//@return: err error

func (userService *UserService) SetUserAuthorities(id string, authorityIds []string) (err error) {
	return global.CMBP_DB.Transaction(func(tx *gorm.DB) error {
		TxErr := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		var useAuthority []system.SysUserAuthority
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, system.SysUserAuthority{
				SysUserId: id, SysAuthorityAuthorityId: v,
			})
		}
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}
		TxErr = tx.Where("id = ?", id).First(&system.SysUser{}).Update("authority_id", authorityIds[0]).Error
		if TxErr != nil {
			return TxErr
		}
		// 返回 nil 提交事务
		return nil
	})
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteUser
//@description: 删除用户
//@param: id float64
//@return: err error

func (userService *UserService) DeleteUser(id int) (err error) {
	return global.CMBP_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&system.SysUser{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&[]system.SysUserAuthority{}, "sys_user_id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser model.SysUser
//@return: err error, user model.SysUser

func (userService *UserService) SetUserInfo(req system.SysUser) error {
	return global.CMBP_DB.Model(&system.SysUser{}).
		Select("updated_at", "nick_name", "header_img", "phone", "email", "sideMode", "enable").
		Where("id=?", req.ID).
		Updates(map[string]interface{}{
			"updated_at": time.Now(),
			"nick_name":  req.NickName,
			"header_img": req.HeaderImg,
			"phone":      req.Phone,
			"email":      req.Email,
			"side_mode":  req.SideMode,
			"enable":     req.Enable,
		}).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetUserInfo
//@description: 设置用户信息
//@param: reqUser model.SysUser
//@return: err error, user model.SysUser

func (userService *UserService) SetSelfInfo(req system.SysUser) error {
	return global.CMBP_DB.Model(&system.SysUser{}).
		Where("id=?", req.ID).
		Updates(req).Error
}

//@author: [piexlmax](https://github.com/piexlmax)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user system.SysUser

func (userService *UserService) GetUserInfo(uid string, userId string) (user map[string]interface{}, err error) {
	var reqUser system.Users
	var roleObj system.Roles
	var role = ""

	if userId != "" {
		err = global.CMBP_DB.Model(system.Users{}).Where("id = ?", userId).First(&reqUser).Error
		global.CMBP_DB.Model(system.Roles{}).Joins("JOIN t_user_roles on t_roles_info.id=t_user_roles.role_id").Where("t_user_roles.user_id = ?", userId).First(&roleObj)
	} else {
		err = global.CMBP_DB.Model(system.Users{}).Where("id = ?", uid).First(&reqUser).Error
		global.CMBP_DB.Model(system.Roles{}).Joins("JOIN t_user_roles on t_roles_info.id=t_user_roles.role_id").Where("t_user_roles.user_id = ?", uid).First(&roleObj)
	}
	if err != nil {
		return nil, err
	}

	switch roleObj.Name {
	case "ROOT":
		role = "管理员"
	case "MODEL", "LABLE":
		role = "开发者"
	case "ADMIN":
		role = "企业用户"
	default:
		role = roleObj.RoleName
	}
	var mine system.MineRegistry
	err = global.CMBP_DB.Model(system.MineRegistry{}).Where("mine_code = ?", reqUser.MineCode).First(&mine).Error

	var expireTime string
	var lastDays = 99999
	if !reqUser.ExpireTime.IsZero() {
		expireTime = reqUser.ExpireTime.Format("2006-01-02 15:04:05") // 格式化时间或者以你需要的方式使用
		now := time.Now()
		duration := reqUser.ExpireTime.Sub(now)
		lastDays = int(duration.Hours() / 24)
	} else {
		expireTime = ""
	}
	rspUser := map[string]interface{}{
		"username":      reqUser.Username,
		"role":          role,
		"mine":          mine.MineShortname,
		"phone":         reqUser.Phone,
		"emial":         reqUser.Email,
		"expire_time":   expireTime,
		"last_days":     lastDays,
		"mine_code":     reqUser.MineCode,
		"ding_account":  reqUser.DingAccount,
		"register_time": reqUser.CreateTime.Format("2006-01-02 15:04:05"),
	}
	return rspUser, err

}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserById(id int) (user *system.SysUser, err error) {
	var u system.SysUser
	err = global.CMBP_DB.Where("id = ?", id).First(&u).Error
	return &u, err
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserByUuid(uuid string) (user *system.SysUser, err error) {
	var u system.SysUser
	if err = global.CMBP_DB.Where("uuid = ?", uuid).First(&u).Error; err != nil {
		return &u, errors.New("用户不存在")
	}
	return &u, nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: resetPassword
//@description: 修改用户密码
//@param: ID uint
//@return: err error

func (userService *UserService) ResetPassword(ID string) (err error) {
	err = global.CMBP_DB.Model(&system.SysUser{}).Where("id = ?", ID).Update("password", utils.BcryptHash("123456")).Error
	return err
}

//func (UserService *UserService) UpdateUserInfo(user system.Users) (err error) {
//	return global.CMBP_DB.Model(&system.Users{}).
//		Select("updated_at", "nick_name", "header_img", "phone", "email", "sideMode", "enable").
//		Where("id=?", req.ID).
//		Updates(map[string]interface{}{
//			"updated_at": time.Now(),
//			"nick_name":  req.NickName,
//			"header_img": req.HeaderImg,
//			"phone":      req.Phone,
//			"email":      req.Email,
//			"side_mode":  req.SideMode,
//			"enable":     req.Enable,
//		}).Error
//}
