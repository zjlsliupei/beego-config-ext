package zookeeper

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/config"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"time"
)

type ZkConfig struct {
}

type ZkConfigContainer struct {
	data map[string]string
}

// Parse 获取zookeeper数据
// address格式：`{"path":"/test","hosts":["localhost:2181"]}`
func (z *ZkConfig) Parse(option string) (config.Configer, error) {
	options := gjson.Parse(option)
	path := options.Get("path").String()
	hosts := options.Get("hosts").Array()
	var _hosts []string
	for _, host := range hosts {
		_hosts = append(_hosts, host.String())
	}
	// 连接zookeeper
	conn, err := z.connect(_hosts)
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	// 获取所有一级路径下的值
	_data, err := z.getPathValues(conn, path)
	if err != nil {
		return nil, err
	}
	zc := &ZkConfigContainer{data: _data}
	return zc, nil
}

// ParseData 直接转化配置
func (z *ZkConfig) ParseData(data []byte) (config.Configer, error) {
	var _data map[string]string
	err := json.Unmarshal(data, &_data)
	if err != nil {
		return nil, err
	}
	zc := &ZkConfigContainer{data: _data}
	return zc, nil
}

// connect 连接zookeeper
func (z *ZkConfig) connect(hosts []string) (*zk.Conn, error) {
	//option := zk.WithEventCallback(func(event zk.Event) {
	//	fmt.Println("*******************")
	//	fmt.Println("path:", event.Path)
	//	fmt.Println("type:", event.Type.String())
	//	fmt.Println("state:", event.State.String())
	//	fmt.Println("-------------------")
	//})
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// getPathValues 获取路径下一级节点的值
func (z *ZkConfig) getPathValues(conn *zk.Conn, path string) (map[string]string, error) {
	keys, _, err := conn.Children(path)
	var data = make(map[string]string)
	if err != nil {
		return data, err
	}
	for _, v := range keys {
		value, _, err := conn.Get(path + "/" + v)
		if err == nil {
			data[v] = string(value)
		}
	}
	return data, nil
}

func (zc *ZkConfigContainer) Set(key, val string) error {
	zc.data[key] = val
	return nil
}

func (zc *ZkConfigContainer) String(key string) string {
	if val, ok := zc.data[key]; ok {
		return val
	}
	return ""
}

func (zc *ZkConfigContainer) Strings(key string) []string {
	if val, ok := zc.data[key]; ok {
		return strings.Split(val, ";")
	}
	return []string{}
}

func (zc *ZkConfigContainer) Int(key string) (int, error) {
	if val, ok := zc.data[key]; ok {
		return strconv.Atoi(val)
	}
	return 0, errors.New(key + " not exist")
}

func (zc *ZkConfigContainer) Int64(key string) (int64, error) {
	if val, ok := zc.data[key]; ok {
		return strconv.ParseInt(val, 10, 64)
	}
	return 0, errors.New(key + " not exist")
}

func (zc *ZkConfigContainer) Bool(key string) (bool, error) {
	if val, ok := zc.data[key]; ok {
		return config.ParseBool(val)
	}
	return false, errors.New(key + " not exist")
}

func (zc *ZkConfigContainer) Float(key string) (float64, error) {
	if val, ok := zc.data[key]; ok {
		return strconv.ParseFloat(val, 64)
	}
	return 0, errors.New(key + " not exist")
}

func (zc *ZkConfigContainer) DefaultString(key, defaultVal string) string {
	if val, ok := zc.data[key]; ok {
		return val
	}
	return defaultVal
}

func (zc *ZkConfigContainer) DefaultStrings(key string, defaultVal []string) []string {
	if val, ok := zc.data[key]; ok {
		return strings.Split(val, ";")
	}
	return defaultVal
}

func (zc *ZkConfigContainer) DefaultInt(key string, defaultVal int) int {
	if val, ok := zc.data[key]; ok {
		_val, _ := strconv.Atoi(val)
		return _val
	}
	return defaultVal
}

func (zc *ZkConfigContainer) DefaultInt64(key string, defaultVal int64) int64 {
	if val, ok := zc.data[key]; ok {
		_val, _ := strconv.ParseInt(val, 10, 64)
		return _val
	}
	return defaultVal
}

func (zc *ZkConfigContainer) DefaultBool(key string, defaultVal bool) bool {
	if val, ok := zc.data[key]; ok {
		_val, _ := config.ParseBool(val)
		return _val
	}
	return defaultVal
}

func (zc *ZkConfigContainer) DefaultFloat(key string, defaultVal float64) float64 {
	if val, ok := zc.data[key]; ok {
		_val, _ := strconv.ParseFloat(val, 64)
		return _val
	}
	return defaultVal
}

func (zc *ZkConfigContainer) DIY(key string) (interface{}, error) {
	if val, ok := zc.data[key]; ok {
		return val, nil
	}
	return nil, errors.New(key + " not exist")
}

func (zc *ZkConfigContainer) GetSection(section string) (map[string]string, error) {
	return nil, errors.New("GetSection function not support for zookeeper")
}

func (zc *ZkConfigContainer) SaveConfigFile(filename string) error {
	return nil
}

func init() {
	config.Register("zookeeper", &ZkConfig{})
}
