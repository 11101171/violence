package models

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : user.go
// Author      : ningzhong.zeng
// Revision    : 2015-11-14 17:34:46
// Description :
//****************************************************/

import (
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"strings"
	"violence/utils"

	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"

	_ "github.com/go-sql-driver/mysql"
)

type RoleType int

const (
	ROLE_SUPER_ADMIN RoleType = iota
	ROLE_ADMIN
	ROLE_GUEST
)

func Role(v string) RoleType {
	vi, _ := strconv.Atoi(v)
	switch vi {
	case 0:
		return ROLE_SUPER_ADMIN
	case 1:
		return ROLE_ADMIN
	case 2:
		return ROLE_GUEST
	default:
		return -1
	}
}

// username 用户名
type User struct {
	Id         string    `orm:"size(45);pk"`
	Username   string    `form:"Username",orm:"size(45);default("aaaa")"`
	Password   string    `form:"Password",orm:"size(45)"`
	Role       RoleType  `orm:"size(1)"`
	GmtCreated time.Time `orm:"auto_now_add;type(datetime)"`
	GmtUpdated time.Time `orm:"auto_now;type(datetime)"`
}

// Add user by this userForm
// type UserForm struct {
// Username string `form:"username"; valid:"Required,Min(7),Max(30)"`
// Password string `form:"password"; valid:"Required,Min(6),Max(30)"`
// }
func (this *User) TableName() string {
	return "v_user"
}

// create table index
func (this *User) TableIndex() [][]string {
	return [][]string{
		[]string{"id"},
		[]string{"username"},
	}
}

// create table unique
func (this *User) TableUnique() [][]string {
	return [][]string{
		[]string{"username"},
	}
}

var model *User = new(User)

func init() {

}

// Valid user edit
func (u *User) ValidUserEdit() *validation.Validation {
	valid := u.ValidUserAdd()
	valid.Required(u.Id, "UsernameId").Message("用户id不能为空")
	return valid
}

// Valid user add
func (u *User) ValidUserAdd() *validation.Validation {
	valid := validation.Validation{}

	u.Username = strings.TrimSpace(u.Username)
	u.Password = strings.TrimSpace(u.Password)

	valid.Required(u.Username, "UsernameReq").Message("用户名不能为空")
	valid.MinSize(u.Username, 6, "UsernameMixSize").Message("用户名不能小于6个字符")
	valid.MaxSize(u.Username, 10, "UsernameMaxSize").Message("用户姓名不能大于10个字符")

	valid.MinSize(u.Password, 6, "PasswordMinSize").Message("用户密码不能小于6为字符")
	valid.MaxSize(u.Password, 10, "PasswordMaxSize").Message("用户密码不能大于10个字符")

	// uLen := len(strings.TrimSpace(u.Username))
	// if uLen < 6 || uLen > 10 {
	// valid.SetError("Username", "用户不能为空，长度大于6小于10")
	// }

	// pLen := len(strings.TrimSpace(u.Password))
	// if pLen < 6 || pLen > 10 {
	// valid.SetError("Password", "用户密码不能为空，长度大于6小于10")
	// }

	return &valid
}

// Valid user login
func (u *User) ValidUserLogin() *validation.Validation {
	valid := u.ValidUserAdd()
	return valid
}

func EnPwd(password string) string {
	hash := md5.New()
	io.WriteString(hash, password)
	return fmt.Sprintf("%x", hash.Sum(nil))
}
func CheckPwd(password string, checkPassword string) bool {
	return EnPwd(checkPassword) == password
}

func (this *User) Login() map[string]string {
	data := make(map[string]string)
	user := &User{}
	err := orm.NewOrm().QueryTable(model).Filter("username", this.Username).One(user)
	if err != nil {
		beego.Debug("Login err: ", err)
		data["ErrorUserNotFount"] = "用户不存在"
		return data
	}
	this.Id = user.Id

	if !CheckPwd(user.Password, this.Password) {
		data["ErrorWrongPwd"] = "密码错误"
	}
	return data
}

func (this *User) SelectList() (users []*User) {
	beego.Error("user")
	qs := orm.NewOrm().QueryTable(model)
	qs.OrderBy("-GmtCreated").All(&users)
	return users
}

func (this *User) SelectOneById(id string) *User {
	user := &User{}
	qs := orm.NewOrm().QueryTable(model)
	qs.Filter("id", id).One(user)
	return user
}

func (this *User) InsertOne() error {
	this.Id = utils.Random().RandomString("user", 10)
	this.Password = EnPwd(this.Password)
	err := base.InsertOne(this)
	if err != nil {
		beego.Error("InsertOne Affected", "Err=>", err)
	}
	return err
}

func (this *User) UpdateOne() error {
	num, err := orm.NewOrm().QueryTable(model).Filter("id", this.Id).Update(orm.Params{
		"username": this.Username,
		"password": this.Password})
	if err != nil {
		beego.Error("UpdateOne Affected Num=>", num, " ,Err=>", err)
	}

	return err
}

func (this *User) DeleteOneById() error {
	num, err := orm.NewOrm().QueryTable(model).Filter("id", this.Id).Delete()
	if err != nil {
		beego.Error("DeleteOneById Affected Num=>", num, " ,Err=>", err)
	}
	return err
}

func (this *User) ToString() string {
	return fmt.Sprintf("%#v", this)
}
