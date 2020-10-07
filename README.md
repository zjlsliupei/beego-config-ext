# beego-config-ext 为扩展beego-config模块
目前已经增加zookeeper扩展，基于beego config标准实现zookeeper接口
## 安装
```go
go get github.com/zjlsliupei/beego-config-ext
```

## 使用
```go
import (
    "github.com/astaxie/beego/config"
    _ "github.com/zjlsliupei/beego-config-ext/zookeeper"
)
// 参考beego模块实例config
c, err := config.NewConfig("zookeeper", `{"path":"/test","hosts":["localhost:2181"]`}`)
if err != nil {
    fmt.Println(err)
}
c.String("name") // 输出 ：liupei

