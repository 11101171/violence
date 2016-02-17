package models

import (
	"fmt"
	"testing"
	"violence/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

func init() {
	fmt.Println("test init")
}
func Test_ValidUserAdd(t *testing.T) {
	user := models.User{Username: "zhangsan", Password: "password"}
	var valid *validation.Validation
	valid = user.ValidUserAdd()
	b, _ := valid.Valid(user)
	if !b {
		for k, v := range valid.ErrorMap() {
			beego.Info(k, " Error is ", "-", v.Value, "-", v.Field, "-", v.Tmpl, "-", v.Key, "=")
		}
	}

}

func Test_Valid(t *testing.T) {
	valid := validation.Validation{}
	user := &models.User{Username: "zhangsan", Password: "password"}
	b, err := valid.Valid(user)
	if err != nil {
		// handle error
		t.Log(err.Error())
		t.Error("Valid Error")
	}
	if !b {
		// validation does not pass
		// blabla...
		for _, err := range valid.Errors {
			t.Log(err.Field + "-" + err.String())
			t.Log(err.Key, err.Message)
		}

		t.Log("ppppp")
		for k, v := range valid.ErrorMap() {
			t.Log(k, "=", v)
		}
		t.Error("Valid Error")
	}

	t.Log("Valid not has errors")
}

func Test_InsertOne(t *testing.T) {
	user := &models.User{Username: "11zhangsan", Password: "lis"}
	err := user.InsertOne()
	t.Log("id=>", user.Id)
	if err != nil {
		t.Error("插入失败")
	} else {
		t.Log("插入成功")

	}
}

func Test_SelectList(t *testing.T) {
	users := new(models.User).SelectList()
	for num, user := range users {
		t.Log("获取数据", num, "=>", new(models.Base).BaseToString(user))
	}
}

func Test_SelectOneById(t *testing.T) {
	user := new(models.User).SelectOneById("pk_2")
	t.Log(user.ToString())
}
