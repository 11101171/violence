package models

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : models/cmd.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-30 00:30:17
// Description :
//****************************************************/

import (
	"time"
	"violence/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Cmd struct {
	Id         string    `orm:"pk; size(45)"`
	Shell      string    `orm:"size(225)"`
	AgentId    string    `orm:"size(45)"`
	GmtCreated time.Time `orm:"auto_now_add;type(datetime)"`
	GmtUpdated time.Time `orm:"auto_now;type(datetime)"`
}

func (this *Cmd) TableName() string {
	return "s_cmd"
}

func (this *Cmd) TableIndex() [][]string {
	return [][]string{
		[]string{"id"},
	}
}

var cmdModel *Cmd = new(Cmd)

func init() {

}

// Valid userForm edit
func (this *Cmd) ValidCmdEdit() *validation.Validation {
	valid := this.ValidCmdAdd()
	valid.Required(this.Id, "Id").Message("id不能为空")
	return valid
}

// Valid cmdForm add
func (this *Cmd) ValidCmdAdd() *validation.Validation {
	valid := validation.Validation{}

	valid.Required(this.Shell, "ShellReq").Message("脚本不能为空")
	valid.MinSize(this.Shell, 1, "ShellMixSize").Message("脚本不能小于1个字符")
	valid.MaxSize(this.Shell, 225, "ShellMaxSize").Message("脚本不能大于225个字符")

	valid.Required(this.AgentId, "AgentIdReq").Message("主机不能为空")

	return &valid
}

func (this *Cmd) SelectListByAgentId(agentId string) (cmds []*Cmd) {
	qs := orm.NewOrm().QueryTable(cmdModel).Filter("agent_id", agentId)
	qs.OrderBy("-GmtCreated").All(&cmds)
	return cmds
}

func (this *Cmd) SelectOneById(id string) *Cmd {
	cmd := &Cmd{}
	qs := orm.NewOrm().QueryTable(cmdModel)
	qs.Filter("id", id).One(cmd)
	return cmd
}

func (this *Cmd) InsertOne() error {
	this.Id = utils.Random().RandomString("cmd", 10)
	err := base.InsertOne(this)
	if err != nil {
		beego.Error("InsertOne Affected ", "Err=>", err)
	}
	return err
}

func (this *Cmd) UpdateOne() error {
	num, err := orm.NewOrm().QueryTable(cmdModel).Filter("id", this.Id).Update(orm.Params{
		"Shell": this.Shell,
	},
	)
	if err != nil {
		beego.Error("UpdateOne Affected Num=>", num, " ,Err=>", err)
	}
	return err
}

func (this *Cmd) DeleteOneById() error {
	num, err := orm.NewOrm().QueryTable(cmdModel).Filter("id", this.Id).Delete()
	if err != nil {
		beego.Error("DeleteOneById Affected Num=>", num, " ,Err=>", err)
	}
	return err
}
