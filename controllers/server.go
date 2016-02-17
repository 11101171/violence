package controllers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"violence/models"
	"violence/services"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

//****************************************************/

// File        : controllers/server.go
// Author      : ningzhong.zeng
// Revision    : 2016-01-7 13:14:55
// Description :
//****************************************************/

type ServerController struct {
	BaseController
}

var jsonContent = `
{
    "host":"http://www.baidu.com/",
    "port":"80",
	"name":"大搜车api",
    "lists":[
        {
            "folder":"css","sort":"1","api_params":[
            {
				"id":"iuiui",
                "path":"/admin/server/encode",
                "method":"GET",
                "fileds":[
                        {"name":"username", "ftype":"text", "value":"zhangsan","lable":"用户名(不能小于3个字符)","placeholder":"如:张三"},
                        {"name":"passwrod", "ftype":"text", "value":"","lable":"密码","placeholder":"如:abc123"},
                        {"name":"sign", "ftype":"text", "value":"","lable":"秘钥","placeholder":"按加密按钮即可", "salt":"xxs", "pftype":"MD5","pway":"1"}
                 ]
            }]
		},
		{
            "folder":"email","sort":"1","api_params":[]
        },
        {
            "folder":"img","sort":"1","api_params":[]
        }
    ]
}
`

func (this *ServerController) Index() {
	user := this.GetAuthUser()
	servers := new(models.Server).SelectListByUserId(user.Id)
	this.Data["Servers"] = servers
	serverId := this.Ctx.Input.Param(":serverId")
	server := &models.Server{}
	if serverId != "" {
		server = server.SelectOneById(serverId)
	} else {
		if len(servers) > 0 {
			server = servers[0]
		}
	}

	jsonBody := models.Api{}
	json.Unmarshal([]byte(server.Content), &jsonBody)
	beego.Debug("jsonbody", jsonBody.Name)
	this.Data["JsonBody"] = jsonBody
	this.Data["Server"] = server
	this.Data["Admin"] = true
	this.Layout = "layout/main.html"
	this.TplNames = "admin/server/index.html"
}
func (this *ServerController) Operate() {
	if this.Ctx.Input.IsGet() {
		serverId := this.Ctx.Input.Param(":serverId")
		beego.Debug("serverId=>", serverId)
		server := &models.Server{}
		if serverId != "" {
			server = server.SelectOneById(serverId)
		}
		this.Data["ServerForm"] = server
		this.Layout = "layout/main.html"
		this.TplNames = "admin/server/operate.html"
	} else {
		server := models.Server{}
		if err := this.ParseForm(&server); err != nil {
			this.ErrInputData("数据错误")
		}

		user := this.GetAuthUser()
		server.UserId = user.Id
		// beego.Debug("server=>",
		// fmt.Sprintf("%+v", server),
		// )
		var valid *validation.Validation
		if server.Id == "" {
			valid = server.ValidServerAdd()
		} else {
			valid = server.ValidServerEdit()
		}

		b, _ := valid.Valid(&server)
		if !b {
			this.Data["ErrorMap"] = valid.ErrorMap()
			this.Data["ErrorForm"] = server
			this.Layout = "layout/main.html"
			this.TplNames = "admin/server/operate.html"
			return
		}

		var message string
		if server.Id == "" {
			message = "新增API配置主题成功"
			if err := server.InsertOne(); err != nil {
				this.ErrDataBase("新增API配置主题失败")
				return
			}
		} else {
			message = "更新API配置主题成功"
			if err := server.UpdateOne(); err != nil {
				this.ErrDataBase("更新API配置主题失败")
				return
			}
		}
		services.GetInstance().CacheService.PutServerContent(server.Id, server.Content)
		this.DoSuccess(message, "/admin/server/index")

	}

}

// Delete server
func (this *ServerController) Delete() {
	serverId := this.Ctx.Input.Param(":serverId")
	if serverId == "" {
		this.ErrInputData("错误的server")
		return
	}
	server := &models.Server{Id: serverId}
	err := server.DeleteOneById()
	if err != nil {
		beego.Debug("err=>", err)
	}
	this.DoSuccess("删除成功", "/admin/server/index")
}
func (this *ServerController) Share() {
	serverId := this.Ctx.Input.Param(":serverId")
	server := &models.Server{}
	if serverId != "" {
		server = server.SelectOneById(serverId)
	} else {
		this.ErrInputData("数据错误，找不到")
	}
	jsonBody := models.Api{}
	json.Unmarshal([]byte(server.Content), &jsonBody)
	beego.Debug("jsonbody", jsonBody.Name)
	this.Data["JsonBody"] = jsonBody
	this.Data["Server"] = server

	this.Layout = "layout/main.html"
	this.TplNames = "admin/server/index.html"
}

var api = `
            {
                "host":"http://localhost:9099/admin/server/encode",
                "method":"GET",
                "fileds":[
						{"name":"username", "ftype":"text", "value":"zhangsan","lable":"用户名(不能小于3个字符)","placeholder":"如:张三"},
						{"name":"passwrod", "ftype":"text", "value":"","lable":"密码","placeholder":"如:abc123"},
						{"name":"sign", "ftype":"text", "value":"","lable":"秘钥","placeholder":"按加密按钮即可", "salt":"xxs", "pftype":"MD5","pway":"1"}
                 ]
            }
        `

func (this *ServerController) Active() {
	if this.Ctx.Input.IsGet() {
		serverId := this.Ctx.Input.Param(":serverId")
		apiParamId := this.Ctx.Input.Param(":apiParamId")

		beego.Debug(serverId, apiParamId)
		server := new(models.Server).SelectOneById(serverId)
		if server.Id == "" {
			this.ErrDataBase("数据库错误")
		}

		apiParam := new(models.Api).Get(server.Content, apiParamId)
		// json.Unmarshal([]byte(api), &apiParam)
		beego.Debug("host", apiParam.Host)
		this.Data["Api"] = apiParam
		this.Data["ServerId"] = serverId
		this.Data["ApiParamId"] = apiParamId
		this.TplNames = "admin/server/active.html"
	} else {
		apiParam := models.ApiParam{}
		if err := this.ParseForm(&apiParam); err != nil {
			this.Data["err"] = err
			this.RFailed()
			return
		}
		// beego.Debug("apiParam", apiParam)
		// transport := &httpclient.Transport{
		// ConnectTimeout:        1 * time.Second,
		// RequestTimeout:        10 * time.Second,
		// ResponseHeaderTimeout: 5 * time.Second,
		// }
		// defer transport.Close()
		var err = errors.New("")
		var values = this.Ctx.Request.Form
		values.Del("Host")
		values.Del("Method")
		var location = values["location[]"]
		values.Del("location[]")
		beego.Debug("values", values.Encode())
		// beego.Debug("location>", location)
		// var resp *http.Response
		// if apiParam.IsGet() {
		// resp, err = http.Get(apiParam.Host + "?" + values.Encode())
		// } else {
		// resp, err = http.PostForm(apiParam.Host, values)
		// }
		req, err := http.NewRequest(apiParam.Method, apiParam.Host, strings.NewReader(values.Encode()))
		if err != nil {
			beego.Error("http req err =>", err)
		}
		for i := 0; i < len(location); i = i + 1 {
			if value := values.Get(location[i]); value != "" {
				req.Header.Add(location[i], value)
			}
		}
		client := &http.Client{
			Timeout: 10 * time.Second,
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded; ")
		resp, err := client.Do(req)
		if err != nil {
			beego.Error("http resp err =>", err)
		}
		// client := &http.Client{Transport: transport}
		robots, _ := ioutil.ReadAll(resp.Body)
		beego.Debug("resp=>", string(robots))
		// this.Data["Api"] = apiParam
		// this.RSuccess()
		this.RSuccessWithData(
			map[string]interface{}{
				"resp": string(robots),
			},
		)
	}
}
func (this *ServerController) Encode() {
	serverId := this.Ctx.Input.Param(":serverId")
	apiParamId := this.Ctx.Input.Param(":apiParamId")

	beego.Debug(serverId, apiParamId)
	server := new(models.Server).SelectOneById(serverId)
	if server.Id == "" {
		this.ErrDataBase("数据库错误")
	}

	apiParam := new(models.Api).Get(server.Content, apiParamId)
	var values = this.Ctx.Request.Form
	values.Del("Host")
	values.Del("Method")
	// json.Unmarshal([]byte(api), &apiParam)
	// data := map[string]interface{}{"Jing": "xxx", "Pll": "pppp"}
	this.RSuccessWithData(apiParam.Encode(values))
}
func (this *ServerController) Info() {
	this.TplNames = "admin/server/info.html"
}
