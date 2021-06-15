## httpcase - 接口自动化测试工具

### 介绍
httpcases是一款golang开发的、通过编写接口测试脚本执行接口自动化测试的工具。

### 功能特性
* 测试脚本简单，容易编写，学习成本低，无需编程基础；
* 支持函数功能，轻松处理测试接口返回数据；
* 函数功能丰富，包含字符串处理、数字处理、文件处理、断言函数等；
* 丰富的断言函数支持；
* 支持JavaScript编写扩展插件；


### 安装

#### 1. 直接下载

[点此下载](https://github.com/roseboy/httpcase/releases)

#### 2.brew 安装

```sh
brew install roseboy/tap/httpcase
```

#### 3. 源码安装
```sh
//1. 下载源码
git clone https://gitee.com/roseboy/httpcase.git

//2. 编译源码
//编译环境 go version go1.13.5 darwin/amd64
cd httpcase
go build -o hc main.go

//国内网络不好请设置代理
go env -w GOPROXY=https://goproxy.cn,direct
```


### 开始使用

#### 简单示例

保存以下代码为文件 apitest.hc 
```
//声明一个测试用例,名字叫做 ApiTest
@ApiTest
//使用POST方法请求一个url
POST http://localhost:8080/test
//设置请求头
user-agent: httpcasev1.0
x-requested-with: XMLHttpRequest
cookie: username=admin; name=value;

//设置请求体
{
    "name":"王二丫",
    "sex":"男",
    "age":18
}

//断言函数
//判断http状态是否等于200
assert ${res.status} == 200
//判断返回结果的name字段是否等于tom
assertEq ${res.body.name} tom
//以字符串格式打印返回结果
print ${res.text}

```

运行以下命令执行测试用例
```sh
hc run apitest.hc
```

#### 一个完整的增删改查用例


保存以下代码为文件 httpcase_helloworld.hc

```
//fileName httpcase_helloworld.http

!env test
!envSet test apiUrl http://127.0.0.1:8000
!envSet dev apiUrl http://localhost:80
//全局token
!header token 123456

@添加用户
POST ${apiUrl}/user

{
    "name":"王二丫",
    "sex":"男",
    "age":18
}

assert ${res.status} == 200
print ${res.body.id}
set userId ${res.body.id}
print ${res.text}

@查询用户
GET ${apiUrl}/user/${userId}

assert ${res.status} == 200
assertEq ${res.body.name} 王二丫
print ${res.text}

@修改用户
POST ${apiUrl}/user/${userId}

{
    "id":${userId},
    "name":"王二丫",
    "sex":"男",
    "age":20
}

assert ${res.status} == 200
print ${res.text}

@查询用户是否修改成功
GET ${apiUrl}/user/${userId}

assert ${res.status} == 200
assertEq ${res.body.age} 20
print ${res.text}

@删除用户
DELETE ${apiUrl}/user/${userId}

assert ${res.status} == 200
print ${res.text}

@查询用户是否删除成功
GET ${apiUrl}/user/${userId}

assert ${res.status} == 404
print ${res.text}

```

命令行执行以下命令，启动一个demo接口服务

```sh
hc demo -p 8000
```

打开新的命令行，执行以下命令，运行测试用例

```sh
hc run httpcase_helloworld.hc
```

执行完成后，控制台打印如下执行结果，同时生成 httpcase_helloworld_report_xxxxxxxx.html的测试报告。
```sh
---------------------------------------------------------------------
  Test Result (Total:6, Pass:6, Fail:0, Skip:0, Duration:1ms)
---------------------------------------------------------------------
[1] [Pass] 添加用户 POST http://127.0.0.1:8000/user
[2] [Pass] 查询用户 GET http://127.0.0.1:8000/user/1622822023
[3] [Pass] 修改用户 POST http://127.0.0.1:8000/user/1622822023
[4] [Pass] 查询用户是否修改成功 GET http://127.0.0.1:8000/user/1622822023
[5] [Pass] 删除用户 DELETE http://127.0.0.1:8000/user/1622822023
[6] [Pass] 查询用户是否删除成功 GET http://127.0.0.1:8000/user/1622822023

```

### 使用说明

#### 发送POST请求
```
//TODO:
```

#### 断言
```
//TODO:
```

#### 自定义header、body等请求参数
```
//TODO:
```

#### 函数调用
```
//TODO:
```

#### 全局函数、全局变量、设置变量、使用变量
```
//TODO:
```

#### 请求之前执行函数
```
//TODO:
```

#### 用例文件包含
```
//TODO:
```

#### 加载外部json数据
```
//TODO:
```

#### 编写JavaScript插件，调用
```
//TODO:
```

#### 测试结果回调
```
//TODO:
```

### 命令行参数

|  参数   | 作用  |
|  ----  | ----  |
| run  | 运行测试用例 |
| demo | 启动一个示例接口服务 |
| version | 打印版本信息 |
| help | 显示帮助 |


### 请求响应结构体(${res})

```json
{
    "cookie":{ //Cookies
        "JSESSIONID":"d6f775bb0765885473b0cba3a5fa9c12",
        "_xsrf":"AWgoSoiwRqFW1p431145342bOJ1X8ZfeQ",
        ...
    },
    "header":{ //响应header
        "Cache-Control":"no-cache, no-store, must-revalidate, private, max-age=0",
        "Content-Type":"application/json; charset=utf-8",
        "Date":"Sun, 04 Apr 2021 13:24:40 GMT",
        "Etag":"W/"8937a8d575c57c91d8bcbc5f43850e0cf2f95d06"",
        "Pragma":"no-cache",
        "Referrer-Policy":"no-referrer-when-downgrade",
        "Server":"CLOUD ELB 1.0.0",
        ...
    },
    "length":100, //响应体长度
    "protocol":"HTTP/2.0", //协议版本
    "status":200, //http响应状态码
    "time":318, //请求耗时，单位：ms
    "text":"", //响应内容文本形式
    "body":{} //响应内容json对象形式
}
```

### 测试结果回调
#### 回调参数说明
```
TODO:
```
#### 使用JavaScript处理回调参数
```
TODO:
```

### 函数列表
#### 公共函数
|  函数   | 描述  | 参数 |
|  ----  | ----  | ----  |
| set  |设置一个变量 |<u>**set**</u> name value|
| env |设置环境|<u>**env**</u> prod|
| envSet |设置一个环境变量|<u>**envSet**</u> prod name value|
| callback |测试完成后回调接口|<u>**callback**</u> url|
| callbackWithFunction |测试完成后回调接口，并且使用js函数处理回调请求体|<u>**callbackWithFunction**</u> funName url|
| loadData |加载json数据，相当于批量set变量|<u>**loadData**</u> jsPath|
| import |导入js插件|<u>**import**</u> name jsPath|
| header |设置请求头|<u>**header**</u> key value|
| body |设置请求体|<u>**body**</u> jsonBody|
| param |设置请求参数|<u>**param**</u> name value|
| file |设置上传文件|<u>**file**</u> field path|
| allowRedirect |设置允许重定向|<u>**allowRedirect**</u> false|
| print |打印|<u>**print**</u> value|
| sleep |延时（毫秒）|<u>**sleep**</u> time|
| while |循环条件判断|<u>**while**</u> val1 opt val2|

#### 断言函数
|  函数   | 描述  | 参数 |
|  ----  | ----  | ----  |
| assert |断言|<u>**assert**</u> val1 opt val2|
| assertEq |断言-相等|<u>**assertEq**</u> val1 val2|
| assertNe |断言-不等|<u>**assertNe**</u> val1 val2|
| assertLt |断言-小于|<u>**assertLt**</u> val1 val2|
| assertGt |断言-大于|<u>**assertGt**</u> val1 val2|
| assertLe |断言-小于等于|<u>**assertLe**</u> val1 val2|
| assertGe |断言-大于等于|<u>**assertGe**</u> val1 val2|
| assertContains |断言-包含|<u>**assertContains**</u> str subStr|
| assertMatch |断言-正则匹配|<u>**assertMatch**</u> str regxStr|
| assertBefore |断言-在之前|<u>**assertBefore**</u> date1 date2|
| assertAfter |断言-在之后|<u>**assertAfter**</u> date1 date2|
| assertEmpty |断言-为空|<u>**assertEmpty**</u> str|
| assertNotEmpty |断言-不为空|<u>**assertNotEmpty**</u> str|

#### 字符串函数
|  函数   | 描述  | 参数 |
|  ----  | ----  | ----  |
| len |字符串长度|<u>**len**</u> str|
| replace |字符串替换|<u>**replace**</u> str old new|
| toLower |字符串转小写|<u>**toLower**</u> str|
| toUpper |字符串转大些|<u>**toUpper**</u> str|
| trim |去除左右两边不可见字符|<u>**trim**</u> str|
| trimLeft |去除左边不可见字符|<u>**trimLeft**</u> str|
| trimRight |去除右边不可见字符|<u>**trimRight**</u> str|
| match |字符串是否匹配|<u>**match**</u> str regxStr|
| indexOf |字符串查找|<u>**indexOf**</u> str substr|
| subStr |根据字符索引截取子字符串|<u>**subStr**</u> str indexStr indexStr2|
| subStr2 |根据字符截取子字符串|<u>**subStr2**</u> str beginStr endStr|
| concat |连接两个字符串|<u>**concat**</u> str1 str2|

#### 数字函数
|  函数   | 描述  | 参数 |
|  ----  | ----  | ----  |
| add |加|<u>**add**</u> num1 num2|
| sub |减|<u>**sub**</u> num1 num2|
| multiply |乘|<u>**multiply**</u> num1 num2|
| divide |除|<u>**divide**</u> num1 num2|
| mod |取余|<u>**mod**</u> num1 num2|

#### 文件函数
|  函数   | 描述  | 参数 |
|  ----  | ----  | ----  |
| readFile |读取文件|<u>**readFile**</u> filePath|
| writeFile |写入文件|<u>**writeFile**</u> filePath text|
| appendFile |追加写入文件|<u>**appendFile**</u> filePath text|


### 协议声明
[MulanPSL2](https://license.coscl.org.cn/MulanPSL2/)



