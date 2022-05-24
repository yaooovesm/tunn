package config

import "errors"

//
// Admin
// @Description:
//
type Admin struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
	AdminUser
}

//
// AdminUser
// @Description:
//
type AdminUser struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

//
// Check
// @Description:
// @receiver u
// @param user
// @return error
//
func (u *AdminUser) Check(user AdminUser) error {
	if u.Password != user.Password {
		return errors.New("密码错误")
	}
	if u.User != user.User {
		return errors.New("用户错误")
	}
	return nil
}
