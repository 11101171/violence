package models

import (
	"testing"
	"violence/models"
)

func Test_Random(t *testing.T) {
	//t.Log(models.RandomString("t", 10))
}

func Test_Insert(t *testing.T) {
	user := &models.User{Id: "pktest", Username: "xliausername", Password: "pass"}
	if err := models.InsertOne(user); err == nil {
		t.Log("成功")
	} else {
		t.Error("插入失败")
	}
}

func Test_Select(t *testing.T) {
	user := &models.User{Id: "pktest"}
	if models.SelectOne(user) == nil {
		t.Log("查找成功", models.ToString(user))
	} else {
		t.Error("读取失败")
	}
}

func Test_Update(t *testing.T) {
	user := &models.User{Id: "pktest", Password: "haha"}
	if models.UpdateOne(user) == nil {
		t.Log("修改成功")
	} else {
		t.Error("修改失败")
	}
}

func Test_Delete(t *testing.T) {
	user := &models.User{Id: "pktest"}
	if models.DeleteOne(user) == nil {
		t.Log("删除成功")
	} else {
		t.Error("删除失败")
	}
}
