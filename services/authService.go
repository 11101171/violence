package services

import (
	"fmt"
	"strings"
	"violence/models"
	"violence/utils"
)

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : services/AuthService.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-23 18:13:14
// Description :
//****************************************************/

// "violence/util/desUtil"

const (
	authUsername string = "001"
	authPassword string = "002"

	fUsername string = "9ok37yHx[n@!1"
	fPassword string = "8Uj3#2!*ds5nmz"
)

var (
	seq       string = "="
	c_timeout int64  = 60 * 60 * 24
)

type AuthService struct {
	BaseService
}

// 是否已经登录
func (as *AuthService) HasLogin() bool {
	var result bool = false
	as.Factory().CookieService.WithContext(as.ctx, func() {
		var usernameCookie = as.Factory().CookieService.Get(authUsername)
		if usernameCookie != "" {
			var usernameBody = utils.DES().Decrypt(usernameCookie)
			if strings.HasPrefix(usernameBody, fUsername) {
				as.Factory().CookieService.Set(authUsername, usernameCookie)
				result = true
			}
		}
	})
	return result
}

// 获取登录者信息
func (as *AuthService) GetAuthUser() *models.User {
	user := &models.User{}
	as.Factory().CookieService.WithContext(as.ctx, func() {
		var usernameCookie = as.Factory().CookieService.Get(authUsername)
		if usernameCookie != "" {
			var usernameBody = utils.DES().Decrypt(usernameCookie)
			if strings.HasPrefix(usernameBody, fUsername) {
				authStr := strings.Split(usernameBody, seq)
				user.Id = authStr[1]
				user.Role = models.Role(authStr[2])
				user.Username = authStr[3]
			}
		}
	})
	return user
}

// 获取用户id
func (as *AuthService) Id() string {
	user := as.GetAuthUser()
	if user.Id != "" {
		return user.Id
	}
	return ""
}

// 获取用户姓名
func (as *AuthService) Username() string {
	user := as.GetAuthUser()
	if user.Username != "" {
		return user.Username
	}
	return ""
}

// 获取用户权限
func (as *AuthService) Role() models.RoleType {
	user := as.GetAuthUser()
	return user.Role
}

// 登录
func (as *AuthService) Login(user *models.User) map[string]string {
	data := user.Login()
	for k, _ := range data {
		if strings.Contains(k, "Error") {
			return data
		}
	}

	as.Factory().CookieService.WithContext(as.ctx, func() {
		var rolestr = fmt.Sprintf("%s", user.Role)
		var value = utils.DES().Encrypt(fUsername + seq + user.Id + seq + rolestr + seq + user.Username)
		as.Factory().CookieService.SetWithT(authUsername, value, c_timeout)
	})
	return data
}

// 登出
func (as *AuthService) Logout() {
	as.Factory().CookieService.WithContext(as.ctx, func() {
		as.Factory().CookieService.SetWithT(authUsername, "0000000", 1)
	})
}
