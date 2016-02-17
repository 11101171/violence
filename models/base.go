package models

//*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : .go
// Author      : ningzhong.zeng
// Revision    : 2015-11-14 17:41:51
// Description :
//*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base*Base/

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Interface interface {
	InsertOne(obj interface{}) error
	SelectOne(obj interface{}) error
	UpdateOne(obj interface{}) error
	DeleteOne(obj interface{}) error
	ToString(obj interface{}) string
}

var base Base = Base{}

type Base struct {
	Interface
}

func init() {
	orm.Debug = true

	orm.RegisterModel(
		new(User),
		new(Agent),
		new(Cmd),
		new(Server),
	)

	var runmode = beego.AppConfig.String("runmode")
	var mysqluser = beego.AppConfig.String(runmode + "::mysqluser")
	var mysqlpass = beego.AppConfig.String(runmode + "::mysqlpass")
	var mysqlurls = beego.AppConfig.String(runmode + "::mysqlurls")
	var mysqldb = beego.AppConfig.String(runmode + "::mysqldb")
	var mysqlport = beego.AppConfig.String(runmode + "::mysqlport")
	if mysqlport == "" {
		mysqlport = "3306"
	}
	var url = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		mysqluser,
		mysqlpass,
		mysqlurls,
		mysqlport,
		mysqldb,
	)

	fmt.Println("========================")
	fmt.Println("启动模式:" + runmode)
	fmt.Println("链接地址:" + url)
	fmt.Println("========================")
	orm.RegisterDataBase("default", "mysql", url, 30)
	// orm.RegisterDataBase("default", "mysql", "admin:zhong@/violence?charset=utf8", 30)
	errRunSync := orm.RunSyncdb("default", false, true)
	if errRunSync != nil {
		panic(errRunSync)
	}
	user := User{
		Username: "admin1",
		Password: "admin1",
		Role:     ROLE_SUPER_ADMIN,
	}
	user.InsertOne()

}

func (b *Base) InsertOne(obj interface{}) error {
	_, err := orm.NewOrm().Insert(obj)
	return err
}

func (b *Base) SelectOne(obj interface{}) error {
	return orm.NewOrm().Read(obj)
}

func (b *Base) UpdateOne(obj interface{}) error {
	_, err := orm.NewOrm().Update(obj)
	return err
}

func (b *Base) DeleteOne(obj interface{}) error {
	_, err := orm.NewOrm().Delete(obj)
	return err
}

func (b *Base) ToString(obj interface{}) string {
	return fmt.Sprintf("%#v\n", obj)
}

// func RandomString(pixff string, strlen int) string {
// // rand.Seed(time.Now().UTC().UnixNano())
// result := make([]byte, strlen)
// for i := 0; i < strlen; i++ {
// result[i] = alphanum[rand.Intn(len(alphanum))]
// }
// return time.Now().Format("20151212010203") + "-" + pixff + string(result)
// }
