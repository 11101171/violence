package main

import (
	"strings"
	"violence/controllers"
	"violence/services"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/validation"
)

func commonFilter(ctx *context.Context) {
	beego.Info("==============================================================")
	beego.Info("path=>", ctx.Request.URL.Path)
	beego.Info("params=>", ctx.Request.Form)
	// beego.Debug("postForm=>", ctx.Request.PostForm)
}
func alterFilter(ctx *context.Context) {
	// var localPort = beego.AppConfig.String("httpport")
}

func router() {
	authR := beego.NewNamespace("/auth",
		beego.NSCond(func(ctx *context.Context) bool {
			if ctx.Input.Domain() == "api.beego.me" {
				return true
			}
			commonFilter(ctx)
			return true
		}),
		// beego.NSBefore(func(ctx *context.Context) {
		// beego.Debug("path=>", ctx.Request.URL.Path)
		// services.etInstance().AuthService.WithContext(ctx, func() {
		// if services.GetInstance().AuthService.HasLogin() {
		// ctx.Redirect(200, "/admin/index")
		// }
		// })
		// }),
		beego.NSRouter("/login", &controllers.AuthController{}, "*:Login"),
		beego.NSRouter("/logout", &controllers.AuthController{}, "*:Logout"),
	)

	adminR := beego.NewNamespace("/admin",
		beego.NSBefore(func(ctx *context.Context) {
			commonFilter(ctx)
			services.GetInstance().AuthService.WithContext(ctx, func() {
				if !services.GetInstance().AuthService.HasLogin() {
					ctx.Redirect(302, "/auth/login")
				}
			})
		}),
		beego.NSRouter("/index", &controllers.AdminController{}, "*:Index"),
		beego.NSNamespace("/user",
			beego.NSGet("/info/:id", func(ctx *context.Context) {
				ctx.Output.Body([]byte("shopinfo"))
			}),
			beego.NSRouter("/list", &controllers.AdminController{}, "Get:UserList"),
			beego.NSRouter("/operate/?:userId", &controllers.AdminController{}, "*:UserOperate"),
			beego.NSRouter("/delete/:userId", &controllers.AdminController{}, "Get:UserDelete"),
		),
		beego.NSNamespace("/agent",
			beego.NSRouter("/list", &controllers.AdminController{}, "Get:AgentList"),
			beego.NSRouter("/operate/?:agentId", &controllers.AdminController{}, "*:AgentOperate"),
			beego.NSRouter("/delete/:agentId", &controllers.AdminController{}, "Get:AgentDelete"),
			beego.NSRouter("/:agentId/cmd/list", &controllers.AdminController{}, "Get:CmdList"),
			beego.NSRouter("/:agentId/cmd/operate", &controllers.AdminController{}, "*:CmdOperate"),
			beego.NSRouter("/:agentId/cmd/operate/?:cmdId", &controllers.AdminController{}, "*:CmdOperate"),
			beego.NSRouter("/cmd/delete/:cmdId", &controllers.AdminController{}, "Get:CmdDelete"),
		),
		beego.NSNamespace("/ssh",
			beego.NSRouter("/:agentId", &controllers.AdminController{}, "Get:SshIndex"),
			beego.NSRouter("/join/:agentId", &controllers.AdminController{}, "*:SshJoin"),
			beego.NSRouter("/list", &controllers.AdminController{}, "*:SshList"),
		),
		beego.NSNamespace("/server",
			beego.NSRouter("/index/?:serverId", &controllers.ServerController{}, "Get:Index"),
			beego.NSRouter("/operate/?:serverId", &controllers.ServerController{}, "*:Operate"),
			beego.NSRouter("/delete/:serverId", &controllers.ServerController{}, "Get:Delete"),
		),
	)

	guestR := beego.NewNamespace("/guest",
		beego.NSBefore(func(ctx *context.Context) {
			commonFilter(ctx)
		}),
		beego.NSNamespace("/server",
			beego.NSRouter("/share/:serverId", &controllers.ServerController{}, "Get:Share"),
			beego.NSRouter("/:serverId/active/:apiParamId", &controllers.ServerController{}, "*:Active"),
			beego.NSRouter("/:serverId/encode/:apiParamId", &controllers.ServerController{}, "Post:Encode"),
			beego.NSRouter("/info", &controllers.ServerController{}, "Get:Info"),
		),
	)

	//注册namespace
	beego.AddNamespace(authR, adminR, guestR)

}
func vaildError(filed string, errorMap map[string]*validation.ValidationError) interface{} {
	// beego.Info(str, errorMap)
	if errorMap == nil {
		return ""
	}
	str := strings.Split(filed, "-")
	for _, v := range errorMap {
		if strings.Contains(v.Key, str[0]) {

			switch {
			case str[1] == "tmpl":
				return v.Tmpl
			case str[1] == "name":
				return v.Name
			case str[1] == "value":
				return v.Value
			case str[1] == "message":
				return v.Message
			default:
				return ""
			}
		}

	}
	return ""
}
func main() {
	beego.AddFuncMap("VaildError", vaildError)
	beego.SetStaticPath("amaze", "static/amaze")
	beego.SetStaticPath("hp", "static/hp")
	beego.SetStaticPath("term", "static/term")
	beego.SetStaticPath("static", "static")
	beego.ErrorController(&controllers.ErrorController{})
	router()
	loggerConfig := `{
		"filename":"log/beego.log",
		"maxlines" : 1000,
		"maxsize"  : 10240,
		"rotate": true,
		"daily":true,
		"maxdays":10
	}`
	beego.SetLogger("file", loggerConfig)
	beego.Run()
}
