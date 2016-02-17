package services

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : services/sshService.go
// Author      : ningzhong.zeng
// Revision    : 2015-12-27 17:07:43
// Description :
//****************************************************/

import (
	"encoding/base64"
	"flag"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	// "io/ioutil"

	"violence/models"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/kr/pty"
)

var addrFlag, cmdFlag, staticFlag string

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1,
	WriteBufferSize: 1,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsPty struct {
	Cmd *exec.Cmd // pty builds on os.exec
	Pty *os.File  // a pty is simply an os.File
}

type SSHConfig struct {
	Id    string
	Agent *models.Agent
	WsPty *wsPty
}

type SSHService struct {
}

func init() {
	// go sshroom()
}

func (wp *wsPty) Start() {
	var err error
	args := flag.Args()
	wp.Cmd = exec.Command(cmdFlag, args...)
	wp.Pty, err = pty.Start(wp.Cmd)
	if err != nil {
		beego.Error("Failed to start command: %s\n", err)
	}
}

func (wp *wsPty) Stop() {
	wp.Pty.Close()
	wp.Cmd.Wait()
}

func (sshConfig *SSHConfig) Login() string {
	// beego.Debug(sshConfig.Agent.LoginName, sshConfig.Agent.Host)
	var err error
	sshConfig.WsPty.Cmd = exec.Command("ssh", "-o", "StrictHostKeyChecking=no", sshConfig.Agent.LoginName+"@"+sshConfig.Agent.Host)
	sshConfig.WsPty.Pty, err = pty.Start(sshConfig.WsPty.Cmd)
	if err != nil {
		beego.Error("Failed to start command: %s\n", err)
	}

	i := 0
	for {
		if i >= 10 {
			beego.Error("login error")
			sshConfig.WsPty.Pty.Close()
			return "login error"
		}
		buf := make([]byte, 1024)
		size, err := sshConfig.WsPty.Pty.Read(buf)
		if err != nil {
			beego.Error("login Read error")
		}

		if !strings.Contains(string([]byte(buf[:size])), "password") {
			i++
			continue
		}

		_, err = sshConfig.WsPty.Pty.Write([]byte(sshConfig.Agent.LoginPass + "\n"))
		if err != nil {
			beego.Error("login Write error")
		}
		return ""
	}
}

func NewSSHConfig(uname string) SSHConfig {
	agent := new(models.Agent).SelectOneById(uname)
	return SSHConfig{
		Agent: agent,
		Id:    "1",
	}
}

func (this *SSHService) Shell(w http.ResponseWriter, r *http.Request, sshConfig SSHConfig) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		beego.Error("Websocket upgrade failed: %s\n", err)
	}
	defer conn.Close()

	sshConfig.WsPty = &wsPty{}
	if sshConfig.Login() != "" {
		return
	}
	_, err = sshConfig.WsPty.Pty.Write([]byte(sshConfig.Agent.InitShellCmd + "\n"))
	if err != nil {
		beego.Error("InitShellCmd Error")
	}

	// copy everything from the pty master to the websocket
	// using base64 encoding for now due to limitations in term.js
	go func() {
		buf := make([]byte, 128)
		// TODO: more graceful exit on socket close / process exit
		for {
			n, err := sshConfig.WsPty.Pty.Read(buf)
			if err != nil {
				beego.Error("Failed to read from pty master: %s", err)
				return
			}

			out := make([]byte, base64.StdEncoding.EncodedLen(n))
			base64.StdEncoding.Encode(out, buf[0:n])

			err = conn.WriteMessage(websocket.TextMessage, out)

			if err != nil {
				beego.Error("Failed to send %d bytes on websocket: %s", n, err)
				return
			}
		}
	}()

	// read from the web socket, copying to the pty master
	// messages are expected to be text and base64 encoded
	for {
		mt, payload, err := conn.ReadMessage()
		if err != nil {
			if err != io.EOF {
				beego.Info("conn.ReadMessage failed: %s\n", err)
				return
			}
		}

		switch mt {
		case websocket.BinaryMessage:
			beego.Info("Ignoring binary message: %q\n", payload)
		case websocket.TextMessage:
			buf := make([]byte, base64.StdEncoding.DecodedLen(len(payload)))
			_, err := base64.StdEncoding.Decode(buf, payload)
			if err != nil {
				beego.Info("base64 decoding of payload failed: %s\n", err)
			}
			sshConfig.WsPty.Pty.Write(buf)
		default:
			beego.Info("Invalid message type %d\n", mt)
			return
		}
	}

	sshConfig.WsPty.Stop()
}
