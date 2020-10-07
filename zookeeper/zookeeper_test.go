package zookeeper

import (
	"encoding/json"
	"github.com/astaxie/beego/config"
	"github.com/samuel/go-zookeeper/zk"
	"testing"
	"time"
)

var (
	path  = "/test"
	hosts = []string{"localhost:2181"}
)

func TestZk(t *testing.T) {
	// 写入测试数据
	conn, _, err := zk.Connect(hosts, time.Second*5)
	defer conn.Close()
	if err != nil {
		t.Error(err)
	}
	key := path + "/name"
	val := "liupei"
	conn.Create(key, []byte(val), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))

	// 测试获取
	_hosts, _ := json.Marshal(hosts)
	c, err := config.NewConfig("zookeeper", `{"path":"`+path+`","hosts":`+string(_hosts)+`}`)
	if err != nil {
		t.Error(err)
	}
	if c.String("name") != val {
		t.Error(key + " value is invalid")
	}
}
