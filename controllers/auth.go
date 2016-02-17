package controllers

import (
	"strings"
	"violence/models"
	"violence/services"

	"github.com/astaxie/beego"
)

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : login.go
// Author      : ningzhong.zeng
// Revision    : 2015-11-14 16:43:19
// Description :
//****************************************************/

type AuthController struct {
	BaseController
}

func (this *AuthController) Login() {
	if this.Ctx.Input.IsGet() {
		this.TplNames = "auth/login.html"
	} else {
		user := &models.User{}
		if err := this.ParseForm(user); err != nil {
			this.ErrInputData("用户信息错误")
			return
		}
		// toolbox.Display("登录UserForm", user)
		valid := user.ValidUserLogin()
		b, _ := valid.Valid(user)
		if !b {
			beego.Info("Login UserForm error")
			this.Data["ErrorMap"] = valid.ErrorMap()
			this.Data["ErrorForm"] = user
			this.TplNames = "auth/login.html"
			return
		}

		var data map[string]string
		services.GetInstance().AuthService.WithContext(this.Ctx, func() {
			data = services.GetInstance().AuthService.Login(user)
		})

		for k, v := range data {
			if strings.Contains(k, "Error") {
				beego.Debug("登录失败", v)
				this.Data["UserForm"] = user
				this.Data["LoginMsg"] = "登录失败，" + v
				this.TplNames = "auth/login.html"
				return
			}
		}

		// var key string = "c"
		// services.GetInstance().CacheService.Put(key, "value", 10)
		// var v = services.GetInstance().CacheService.Get(key)
		// beego.Debug("v", v)
		beego.Debug("登录成功")
		this.DoSuccess("登录成功", "/admin/index")
	}

}

func (this *AuthController) Logout() {
	services.GetInstance().AuthService.WithContext(this.Ctx, func() {
		services.GetInstance().AuthService.Logout()
	})
	this.TplNames = "auth/login.html"
}
