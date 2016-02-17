# violence
api搭建 ssh登录操作(已经停止更新，转revel框架开发)

# Getting started
go get github.com/11101171/violence


# go start
##dev 方式启动:
bee run
##prod 方式启动:
./bin/run-9999.sh start

#配置
<pre><code>
appname = violence
httpport = 9099
runmode = dev
HttpServerTimeOut = 10

## 数据库
使用mysql 

[dev] 开发环境
mysqluser = "admin"
mysqlpass = "xxx"
mysqlurls = "127.0.0.1"
mysqldb   = "violence"

[test] 测试环境
mysqluser = "admin"
mysqlpass = "xxx"
mysqlurls = "127.0.0.1"
mysqldb   = "violence"

[prod] 生产环境
mysqluser = "admin"
mysqlpass = "xxx"
mysqlurls = "127.0.0.1"
mysqldb   = "violence"
</code></pre>
* 初始化登录用户和密码都是：admin1


#测试
####测试SSH 

####测试API 配置JSNO文件
<pre><code>
{
    "host":"http://www.baidu.com",
    "port":"80",
    "name":"rrr大搜车api",
    "lists":[
        {
            "folder":"css","sort":"1","api_params":[
                {
                    "id":"11",
                    "path":"/admin/server/encode",
                    "name":"测试",
                    "method":"GET",
                    "fileds":[
                            {"name":"username", "ftype":"text", "value":"zhangsan","lable":"用户名(不能小于3个字符)","placeholder":"如:张三"},
                             {"name":"uang", "ftype":"text", "value":"","lable":"xx用户名(不能小于3个字符)","placeholder":"如:你妈妈的"},
                            {"name":"passwrod", "ftype":"text", "value":"","lable":"密码","placeholder":"如:abc123"},
                            {"name":"sign", "ftype":"text", "value":"","lable":"秘钥","placeholder":"按加密按钮即可", "salt":"xxs", "ptype":"MD5","pway":"1"}
                     ]
                },
                 {
                    "id":"12",
                    "path":"/v",
                    "name":"认证接口",
                    "method":"POST",
                    "fileds":[
                            {"name":"bankCard", "ftype":"text", "value":"62261xxxxxxx","lable":"银行卡","placeholder":"如:6226197900918805"},
                             {"name":"name", "ftype":"text", "value":"张三","lable":"用户名(不能小于1个字符)","placeholder":"如:张三"},
                            {"name":"phone", "ftype":"text", "value":"1880444444","lable":"预存手机号","placeholder":"如:153000000"}
                     ]
                },
                {
                    "id":"13",
                    "path":"/xxx/v1/api",
                    "name":"验证接口",
                    "method":"POST",
                    "fileds":[
                            {"name":"server", "ftype":"text", "value":"verify","lable":"服务名","placeholder":""},
                             {"name":"password", "ftype":"text", "value":"123456","lable":"密码","placeholder":"如:张三"},
                            {"name":"sign", "ftype":"text", "value":"","lable":"秘钥","placeholder":"按加密按钮即可", "salt":"FFEF43Fexx", "ptype":"MD5","pway":"1"}
                     ]
                }
            ]
        },
        {
            "folder":"xxxxx","sort":"1","api_params":[]
        },
        {
            "folder":"yyyyyyy","sort":"1","api_params":[]
        }
    ]
}
</code></pre>




