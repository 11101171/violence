package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"sort"
	"strings"

	"github.com/astaxie/beego"
)

//****************************************************/
//Copyright(c) 2015 Tencent, all rights reserved
// File        : models/api.go
// Author      : ningzhong.zeng
// Revision    : 2016-01-10 17:41:17
// Description :
//****************************************************/

type Api struct {
	Host  string      `json:"host"`
	Port  string      `json:"port"`
	Name  string      `json:"name"`
	Lists []ApiFolder `json:"lists"`
}

type ApiFolder struct {
	Folder    string     `json:"folder"`
	Sort      string     `json:"sort"`
	ApiParams []ApiParam `json:"api_params"`
}

type ApiParam struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Host   string
	Path   string     `json:"path"`
	Method string     `json:"method"`
	Fileds []ApiFiled `json:"fileds"`
}

// 按照 Person.Age 从大到小排序
type ApiFiledSlice []ApiFiled

func (a ApiFiledSlice) Len() int { // 重写 Len() 方法
	return len(a)
}
func (a ApiFiledSlice) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}
func (a ApiFiledSlice) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	return a[j].Name > a[i].Name
	// return a[j].Age < a[i].Age
}

type ApiFiled struct {
	Name        string `json:"name`
	Ftype       string `json:"ftype"`
	Value       string `json:"value"`
	Lable       string `json:"lable"`
	Placeholder string `json:"placeholder"`
	Location    string `json:"location"`
	Salt        string `json:"salt"`
	Ptype       string `json:"ptype"`
	Pway        string `json:"pway"`
}

// 获取api参数
func (api *Api) Get(content string, apiParamId string) ApiParam {
	json.Unmarshal([]byte(content), &api)
	for _, apiFolder := range api.Lists {
		for _, apiParam := range apiFolder.ApiParams {
			if apiParam.Id == apiParamId {
				beego.Debug("api.Port", api.Port)
				if api.Port != "" && api.Port != "80" {
					apiParam.Host = api.Host + ":" + api.Port
				} else {
					apiParam.Host = api.Host
				}
				return apiParam
			}
		}
	}
	return ApiParam{}
}
func (apiParam *ApiParam) IsGet() bool {
	if apiParam.Method == "Get" || apiParam.Method == "GET" {
		return true
	} else {
		return false
	}
}
func (apiParam *ApiParam) ToUrlValues() url.Values {
	v := url.Values{}
	for _, filed := range apiParam.Fileds {
		v.Set(filed.Name, filed.Value)
	}
	return v
}
func (apiParam *ApiParam) ReplaceValue(values url.Values) {
	for k, v := range values {
		// beego.Debug("k=>", k, "v=>", v)
		for i := 0; i < len(apiParam.Fileds); i = i + 1 {
			if k == apiParam.Fileds[i].Name {
				apiParam.Fileds[i].Value = v[0]
			}
		}
		// for _, filed := range apiParam.Fileds {
		// beego.Debug("filed", filed.Name)
		// if filed.Name == k {
		// beego.Debug("replace")
		// filed.Value = v[0]
		// continue
		// }
		// }
	}
	beego.Debug("fileds+>", apiParam.Fileds)
}

func (apiParam *ApiParam) Encode(values url.Values) map[string]interface{} {
	apiParam.ReplaceValue(values)
	sort.Sort(ApiFiledSlice(apiParam.Fileds))
	data := make(map[string]interface{}, len(apiParam.Fileds))
	for _, filed := range apiParam.Fileds {
		switch {
		case ((filed.Ptype == "md5" || filed.Ptype == "MD5") && filed.Pway == "1"):
			data[filed.Name] = apiParam.Md5(filed.Salt)
		default:

		}
	}
	beego.Debug("encode", data)
	return data
}
func (apiParam *ApiParam) Md5(keyvalue string) string {
	var params string = ""
	for _, filed := range apiParam.Fileds {
		switch {
		case filed.Pway != "1" && filed.Location != "header":
			params = params + "&" + filed.Name + "=" + filed.Value
		default:
		}
	}
	params = strings.TrimLeft(params, "&")
	params = params + "key=" + keyvalue

	beego.Debug("params", params)
	cipherStr := md5.Sum([]byte(params))
	// fmt.Print(cipherStr)
	// fmt.Print(hex.EncodeToString(cipherStr))
	return strings.ToUpper(hex.EncodeToString(cipherStr[:]))
}
