package controllers

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : controllers/admin.go
// Author      : ningzhong.zeng
// Revision    : 2015-11-14 16:46:52
// Description :
//****************************************************/

import (
	"violence/models"
	"violence/services"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type AdminController struct {
	BaseController
}

func (this *AdminController) Index() {
	this.Layout = "layout/main.html"
	this.TplNames = "admin/index.html"
}

// 用户列表
func (this *AdminController) UserList() {
	users := new(models.User).SelectList()
	this.Data["Users"] = users
	this.Layout = "layout/main.html"
	this.TplNames = "admin/user/list.html"
}

// @router /user/operate/?:userId [get,post]
func (this *AdminController) UserOperate() {
	if this.Ctx.Input.Method() == "GET" {
		userId := this.Ctx.Input.Param(":userId")
		beego.Info("userId", userId)

		user := &models.User{}
		if userId != "" {
			user = user.SelectOneById(userId)
		}
		this.Data["UserForm"] = user
		this.Layout = "layout/main.html"
		this.TplNames = "admin/user/operate.html"
		return
	} else {
		user := &models.User{}
		if err := this.ParseForm(user); err != nil {
			this.Data["jons"] = err
			this.ServeJson()
			return
		}
		var valid *validation.Validation
		if user.Id == "" {
			valid = user.ValidUserAdd()
		} else {
			valid = user.ValidUserEdit()
		}

		b, _ := valid.Valid(user)
		if !b {
			this.Data["ErrorMap"] = valid.ErrorMap()
			this.Data["ErrorForm"] = user
			for k, v := range valid.ErrorMap() {
				beego.Info(k, " Error is ", "-", v.Value, "-", v.Field, "-", v.Tmpl, "-", v.Key, "=", v.Message, "-")
			}
			this.Layout = "layout/main.html"
			this.TplNames = "admin/user/operate.html"
			return
		}

		var message string
		if user.Id == "" {
			message = "新增用户成功"
			if err := user.InsertOne(); err != nil {
				this.ErrDataBase("新增用户失败")
				return
			}
		} else {
			message = "更新用户成功"
			if err := user.UpdateOne(); err != nil {
				this.ErrDataBase("更新用户失败")
				return
			}
		}

		this.DoSuccess(message, "/admin/user/list")
	}

}

// Delete user
func (this *AdminController) UserDelete() {
	userId := this.Ctx.Input.Param(":userId")
	if userId == "" {
		this.ErrInputData("错误的用户")
		return
	}
	user := &models.User{Id: userId}
	err := user.DeleteOneById()
	if err != nil {
		beego.Debug("err=>", err)
	}
	this.DoSuccess("删除成功", "/admin/user/list")
}

func (this *AdminController) SshList() {
	this.Layout = "layout/main_sidebar.html"
	this.TplNames = "admin/ssh/list.html"
}

// ssh
func (this *AdminController) AgentList() {
	user := this.GetAuthUser()
	agents := new(models.Agent).SelectListByUserId(user.Id)
	this.Data["Agents"] = agents
	this.Layout = "layout/main.html"
	this.TplNames = "admin/agent/list.html"
}

func (this *AdminController) AgentOperate() {
	if this.Ctx.Input.Method() == "GET" {
		agentId := this.Ctx.Input.Param(":agentId")
		beego.Info("agentId", agentId)

		agent := &models.Agent{}
		if agentId != "" {
			agent = agent.SelectOneById(agentId)
		}
		this.Data["AgentForm"] = agent
		this.Layout = "layout/main.html"
		this.TplNames = "admin/agent/operate.html"
	} else {
		agent := &models.Agent{}
		if err := this.ParseForm(agent); err != nil {
			this.Data["jons"] = err
			this.ServeJson()
			return
		}

		agent.UserId = this.GetAuthUser().Id
		var valid *validation.Validation
		if agent.Id == "" {
			valid = agent.ValidAgentAdd()
		} else {
			valid = agent.ValidAgentEdit()
		}

		b, _ := valid.Valid(agent)
		if !b {
			this.Data["ErrorMap"] = valid.ErrorMap()
			this.Data["ErrorForm"] = agent
			for k, v := range valid.ErrorMap() {
				beego.Info(k, " Error is ", "-", v.Value, "-", v.Field, "-", v.Tmpl, "-", v.Key, "=", v.Message, "-")
			}
			this.Layout = "layout/main.html"
			this.TplNames = "admin/agent/operate.html"
			return
		}

		var message string
		if agent.Id == "" {
			message = "新增ssh成功"
			if err := agent.InsertOne(); err != nil {
				this.ErrDataBase("新增ssh失败")
				return
			}
		} else {
			message = "更新ssh成功"
			if err := agent.UpdateOne(); err != nil {
				this.ErrDataBase("更新ssh失败")
				return
			}
		}

		this.DoSuccess(message, "/admin/agent/list")

	}
}
func (this *AdminController) AgentDelete() {
	agentId := this.Ctx.Input.Param(":agentId")
	if agentId == "" {
		this.ErrInputData("错误的用户")
		return
	}
	agent := &models.Agent{Id: agentId}
	err := agent.DeleteOneById()
	if err != nil {
		beego.Debug("err=>", err)
		this.DoFail("删除失败", "/admin/agent/list")
	}
	this.DoSuccess("删除成功", "/admin/agent/list")
}

// cmd
func (this *AdminController) CmdList() {
	agentId := this.Ctx.Input.Param(":agentId")
	cmds := new(models.Cmd).SelectListByAgentId(agentId)
	this.Data["Cmds"] = cmds
	this.Layout = "layout/main.html"
	this.TplNames = "admin/agent/list.html"
}
func (this *AdminController) CmdOperate() {
	if this.Ctx.Input.Method() == "GET" {
		cmdId := this.Ctx.Input.Param(":cmdId")
		agentId := this.Ctx.Input.Param(":agentId")
		beego.Info("cmdId", cmdId)

		cmd := &models.Cmd{}
		if cmdId != "" {
			cmd = cmd.SelectOneById(cmdId)
		}
		cmd.AgentId = agentId
		this.Data["CmdForm"] = cmd
		this.Layout = "layout/main.html"
		this.TplNames = "admin/cmd/operate.html"
	} else {
		cmd := &models.Cmd{}
		if err := this.ParseForm(cmd); err != nil {
			this.Data["jons"] = err
			this.ServeJson()
			return
		}
		var valid *validation.Validation
		if cmd.Id == "" {
			valid = cmd.ValidCmdAdd()
		} else {
			valid = cmd.ValidCmdEdit()
		}

		b, _ := valid.Valid(cmd)
		if !b {
			this.Data["ErrorMap"] = valid.ErrorMap()
			this.Data["ErrorForm"] = cmd
			this.Layout = "layout/main.html"
			this.TplNames = "admin/cmd/operate.html"
			return
		}

		var message string
		if cmd.Id == "" {
			message = "新增cmd成功"
			if err := cmd.InsertOne(); err != nil {
				this.ErrDataBase("新增cmd失败")
				return
			}
		} else {
			message = "更新cmd成功"
			if err := cmd.UpdateOne(); err != nil {
				this.ErrDataBase("更新ssh失败")
				return
			}
		}

		this.DoSuccess(message, "/admin/agent/list")

	}
}
func (this *AdminController) CmdDelete() {
	cmdId := this.Ctx.Input.Param(":cmdId")
	if cmdId == "" {
		this.ErrInputData("错误的cmd")
		return
	}
	cmd := &models.Cmd{Id: cmdId}
	err := cmd.DeleteOneById()
	if err != nil {
		beego.Debug("err=>", err)
		this.DoFail("删除失败", "/admin/agent/list")
	}

	this.DoSuccess("删除成功", "/admin/agent/list")
}

// ssh
func (this *AdminController) SshIndex() {
	agentId := this.Ctx.Input.Param(":agentId")
	beego.Info("agentId", agentId)

	agent := &models.Agent{}
	if agentId != "" {
		agent = agent.SelectOneById(agentId)
	}

	this.Data["Agent"] = agent
	this.Data["AuthUser"] = this.GetAuthUser()
	this.TplNames = "admin/ssh/index.html"
	this.Layout = "layout/main.html"
}

func (this *AdminController) SshJoin() {
	var agentId = this.Ctx.Input.Param(":agentId")
	var sshConfig = services.NewSSHConfig(agentId)
	if sshConfig.Agent.Id == "" {
		this.RFailed()
	}
	// Join .
	services.GetInstance().SSHService.Shell(this.Ctx.ResponseWriter, this.Ctx.Request, sshConfig)

	this.RSuccess()
}
