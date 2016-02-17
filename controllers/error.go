package controllers

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : controllers/error.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-6 22:35:36
// Description :
//****************************************************/

import "fmt"

type ErrorController struct {
	BaseController
	Message string
}

func (this *ErrorController) ErrorInside() {
	this.Data["Message"] = this.Message
	this.Layout = "layout/main.html"
	this.TplNames = "layout/error.html"
	return
}

func (this *ErrorController) Error404() {
	this.Layout = "layout/main.html"
	this.TplNames = "layout/404.html"
	return
}

func (this *ErrorController) Error403() {
	this.RanderError("您没有权限访问", "You Is Forbidden")
}
func (this *ErrorController) Error405() {
	this.RanderError("方法不允许", "Method Not Allow")
}

func (this *ErrorController) Error500() {
	this.RanderError("系统错误", "请联系站长")
}
func main() {
	fmt.Println("Start Main func()")
}
