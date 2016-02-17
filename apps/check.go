package main

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : apps/check.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-25 16:44:46
// Description :
//****************************************************/

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/httplib"
)

func main() {
	fmt.Println("Start Main func()")
	req := httplib.Get("http://localhost:9099/admin/index").SetTimeout(100*time.Second, 30*time.Second)
	req.Header("Accept-Encoding", "gzip,deflate,sdch")
	req.Header("Host", "localhost")
	req.Header("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.57 Safari/537.36")
	req.Param("username", "astaxie")
	req.Param("password", "123456")
	str, err := req.String()
	if err != nil {
		fmt.Println("err=>", err)
	}
	fmt.Sprintln("reps=>", str)
}
