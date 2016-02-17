package models

import (
	"time"
	"violence/utils"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : models/server.go
// Author      : ningzhong.zeng
// Revision    : 2016-01-12 11:54:04
// Description :
//****************************************************/

type Server struct {
	Id         string    `orm:"pk; size(45)"`
	UserId     string    `orm:"size(45); index"`
	Theme      string    `orm:"size(45)"`
	Content    string    `orm:"size(20000)"` // type(text)
	GmtCreated time.Time `orm:"auto_now_add;type(datetime)"`
	GmtUpdated time.Time `orm:"auto_now;type(datetime)"`
}

func (this *Server) TableName() string {
	return "a_server"
}

var serverModel = new(Server)

// Valid userForm edit
func (this *Server) ValidServerEdit() *validation.Validation {
	valid := this.ValidServerAdd()
	valid.Required(this.Id, "Id").Message("id不能为空")
	return valid
}

// Valid userForm add
func (this *Server) ValidServerAdd() *validation.Validation {
	valid := validation.Validation{}
	valid.Required(this.Content, "ContentReq").Message("内容不能为空")
	valid.MinSize(this.Content, 1, "ContentMixSize").Message("内容不能小于1个字符")
	valid.Required(this.Theme, "ThemeReq").Message("配置主题不能为空")
	valid.MinSize(this.Theme, 1, "ThemeMixSize").Message("配置主题不能小于1个字符")
	valid.Required(this.UserId, "UserIdReq").Message("请登录")
	valid.MinSize(this.UserId, 1, "UserIdMixSize").Message("请登录")
	return &valid
}

// Query
func (this *Server) SelectListByUserId(userId string) (servers []*Server) {
	qs := orm.NewOrm().QueryTable(serverModel).Filter("user_id", userId)
	qs.OrderBy("-GmtCreated").All(&servers)
	return servers
}
func init() {
	// orm.RegisterModel(new(Server))
}

// select
func (this *Server) SelectOneById(id string) *Server {
	server := Server{}
	qs := orm.NewOrm().QueryTable(serverModel)
	qs.Filter("id", id).One(&server)
	return &server
}

// Insert
func (this *Server) InsertOne() error {
	this.Id = utils.Random().RandomString("Server", 10)
	err := base.InsertOne(this)
	if err != nil {
		beego.Error("InsertOne Affected", "Err=>", err)
	}
	return err
}

// Update
func (this *Server) UpdateOne() error {
	num, err := orm.NewOrm().QueryTable(serverModel).Filter("id", this.Id).Update(orm.Params{
		"content": this.Content,
	},
	)
	if err != nil {
		beego.Error("UpdateOne Affected Num=>", num, " ,Err=>", err)
	}
	return err
}

// Dalete
func (this *Server) DeleteOneById() error {
	num, err := orm.NewOrm().QueryTable(serverModel).Filter("id", this.Id).Delete()
	if err != nil {
		beego.Error("DeleteOneById Affected Num=>", num, " ,Err=>", err)
	}
	return err
}
