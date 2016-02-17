package models

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : models/agent.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-27 19:25:47
// Description :
//****************************************************/

import (
	"strings"
	"time"
	"violence/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Agent struct {
	Id           string    `orm:"pk; size(45)"`
	Host         string    `orm:"size(225)"`
	Port         string    `orm:"size(6)"`
	LoginName    string    `orm:"size(225)"`
	LoginPass    string    `orm:"size(225)"`
	InitShellCmd string    `orm:"size(225)"; default("")`
	UserId       string    `orm:"size(225)"; index`
	Memo         string    `orm:"size(225)"`
	GmtCreated   time.Time `orm:"auto_now_add;type(datetime)"`
	GmtUpdated   time.Time `orm:"auto_now;type(datetime)"`
}

func (this *Agent) TableName() string {
	return "s_agent"
}

func (this *Agent) TableIndex() [][]string {
	return [][]string{
		[]string{"id"},
		[]string{"user_id"},
	}
}

var agentModel *Agent = new(Agent)

func init() {

}

// Valid userForm edit
func (a *Agent) ValidAgentEdit() *validation.Validation {
	valid := a.ValidAgentAdd()
	valid.Required(a.Id, "Id").Message("id不能为空")
	return valid
}

// Valid userForm add
func (a *Agent) ValidAgentAdd() *validation.Validation {
	valid := validation.Validation{}

	a.Host = strings.TrimSpace(a.Host)
	a.Port = strings.TrimSpace(a.Port)
	a.LoginName = strings.TrimSpace(a.LoginName)
	a.LoginPass = strings.TrimSpace(a.LoginPass)

	valid.Required(a.Host, "HostReq").Message("远程地址不能为空")
	valid.MinSize(a.Host, 10, "HostMixSize").Message("远程地址不能小于10个字符")
	valid.MaxSize(a.Host, 225, "HostMaxSize").Message("远程地址不能大于225个字符")

	valid.Required(a.Port, "PortReq").Message("端口不能为空")
	valid.MinSize(a.Port, 2, "PortMixSize").Message("端口不能小于2个字符")
	valid.MaxSize(a.Port, 6, "PortMaxSize").Message("端口不能大于6个字符")

	valid.Required(a.LoginName, "LoginNameReq").Message("登录用户名不能为空")
	valid.MinSize(a.LoginName, 1, "LoginNameMixSize").Message("登录用户名不能小于1个字符")
	valid.MaxSize(a.LoginName, 225, "LoginNameMaxSize").Message("登录用户名不能大于225个字符")

	valid.Required(a.LoginPass, "LoginPassReq").Message("登录密码不能为空")
	valid.MinSize(a.LoginPass, 1, "LoginPassMixSize").Message("登录密码不能小于1个字符")
	valid.MaxSize(a.LoginPass, 225, "LoginPassMaxSize").Message("登录密码不能大于225个字符")

	return &valid
}

func (a *Agent) SelectListByUserId(userId string) (agents []*Agent) {
	qs := orm.NewOrm().QueryTable(agentModel)
	qs.Filter("user_id", userId)
	qs.OrderBy("-GmtCreated").All(&agents)
	return agents
}

func (a *Agent) GetCmds() []*Cmd {
	return new(Cmd).SelectListByAgentId(a.Id)
}

func (a *Agent) SelectOneById(id string) *Agent {
	agent := &Agent{}
	qs := orm.NewOrm().QueryTable(agentModel)
	qs.Filter("id", id).One(agent)
	return agent
}

func (a *Agent) InsertOne() error {
	a.Id = utils.Random().RandomString("agent", 10)
	err := base.InsertOne(a)
	if err != nil {
		beego.Error("InsertOne Affected", "Err=>", err)
	}
	return err
}

func (a *Agent) UpdateOne() error {
	num, err := orm.NewOrm().QueryTable(agentModel).Filter("id", a.Id).Update(orm.Params{
		"host":           a.Host,
		"port":           a.Port,
		"login_name":     a.LoginName,
		"login_pass":     a.LoginPass,
		"init_shell_cmd": a.InitShellCmd,
		"memo":           a.Memo,
	},
	)
	if err != nil {
		beego.Error("UpdateOne Affected Num=>", num, " ,Err=>", err)
	}
	return err
}

func (a *Agent) DeleteOneById() error {
	num, err := orm.NewOrm().QueryTable(agentModel).Filter("id", a.Id).Delete()
	if err != nil {
		beego.Error("DeleteOneById Affected Num=>", num, " ,Err=>", err)
	}
	return err
}
