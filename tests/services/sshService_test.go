package services

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : services/sshService_test.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-27 17:51:47
// Description :
//****************************************************/

import (
	"fmt"
	"testing"
	"violence/services"
)

func Test_Login(t *testing.T) {
	sshConfig := &services.SSHConfig{
		Host:         "120.26.78.245",
		Port:         "22",
		LoginName:    "souche",
		LoginPass:    "pwd!Z@X@souche2013",
		InitShellCmd: "whoiam;cd /home;",
	}

	t.Log("登录")
	sshService := &services.SSHService{}
	sshService.LoginShell(sshConfig)

}

func main() {
	fmt.Println("Start Main func()")
}
