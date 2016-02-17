package controllers

import (
	"violence/models"
	"violence/services"

	"github.com/astaxie/beego"
)

const (
	errInputData = "输入数据有错误"
	errDataBase  = "数据库操作错误"
)

type ValidForm interface {
	Valid()
}

type Err interface {
	ErrInputData(title string, message string)
	ErrDataBase(title string, message string)
	RSuccess()
	RFailed()
}

type BaseController struct {
	beego.Controller
	Err
}

// Err
func (this *BaseController) ErrInputData(message string) {
	this.RanderError(errInputData, message)
}
func (this *BaseController) ErrDataBase(message string) {
	this.RanderError(errDataBase, message)
}
func (this *BaseController) RanderError(title string, message string) {
	this.Data["Message"] = message
	this.Data["Title"] = title
	this.Layout = "layout/main.html"
	this.TplNames = "layout/error.html"
	return
}

var suc = map[string]interface{}{
	"r_code": 1,
	"r_msg":  "操作成功",
}
var fai = map[string]interface{}{
	"r_code": 0,
	"r_msg":  "操作失败",
}

func init() {

}

// Suc
func (this *BaseController) RSuccess() {
	this.Data["json"] = suc
	this.ServeJson()
}
func (this *BaseController) RSuccessWithData(data map[string]interface{}) {
	suc["data"] = data
	this.Data["json"] = suc
	this.ServeJson()
}

// Fai
func (this *BaseController) RFailed() {
	this.Data["json"] = fai
	this.ServeJson()
}
func (this *BaseController) RFailedWithData(data map[string]interface{}) {
	fai["data"] = data
	this.Data["json"] = fai
	this.ServeJson()
}

// redirect
func (this *BaseController) DoSuccess(message string, url string) {
	this.Data["RedirectUrl"] = url
	this.Data["Message"] = message
	this.Data["St"] = true
	this.TplNames = "layout/redirect.html"
	this.Layout = "layout/main.html"
}
func (this *BaseController) DoFail(message string, url string) {
	this.Data["RedirectUrl"] = url
	this.Data["Message"] = message
	this.Data["St"] = false
	this.TplNames = "layout/redirect.html"
	this.Layout = "layout/main.html"
}

// auth
func (this *BaseController) GetAuthUser() *models.User {
	user := &models.User{}
	services.GetInstance().AuthService.WithContext(this.Ctx, func() {
		user = services.GetInstance().AuthService.GetAuthUser()
	})
	return user
}

// 公用日志方法
// func (this *BaseController) Info(v ...interface{}) {
// services.GetInstance().LogService.Info(v)
// }

// func (this *BaseController) Error(v ...interface{}) {
// beego.Error(v)
// }

// func (this *BaseController) Debug(v ...interface{}) {
// beego.Debug(v)
// }

//
