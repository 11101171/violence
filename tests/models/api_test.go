package models

import (
	"encoding/json"
	"fmt"
	"testing"
	"violence/models"
)

var api = `
            {
                "host":"http://www.baidu.com",
                "method":"GET",
                "fileds":[
						{"name":"username", "ftype":"text", "value":""},
						{"name":"passwrod", "ftype":"text", "value":""},
						{"name":"sign", "ftype":"text", "value":"", "salt":"xxs", "pftype":"MD5","pway":"1"}
                 ]
            }
        `

func Test_Json(t *testing.T) {

	t.Log(api)
	// api = `{"host": "http://www.baidu.com", "method":"GET","fileds":[{"name":"username","ftype":"text","value":"xxxx","salt":"oooo","ptyle":"xxxx","pway":"llll"}]}`
	apiParam := models.ApiParam{}
	json.Unmarshal([]byte(api), &apiParam)
	t.Log("JsonObj",
		fmt.Sprintf("%+v", apiParam),
	)
}
