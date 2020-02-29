package main

// var NatsUrl = "nats://192.168.50.2:4222"
import (
	"github.com/gogf/gf/os/glog"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
)

var NatsUrl = "nats://192.168.50.2:4222,nats://192.168.50.1:4221,nats://192.168.50.1:4222,nats://192.168.50.1:4223"

// var NatsUrl = "nats://192.168.50.2:4222,nats://192.168.50.1:4221"

type T struct {
	Num int      `yaml:"num"`
	A   string   `yaml:"natsmicro_dev"`
	X   []string `yaml:"natsmicro_devnatsmicro_dev"`
	B   struct {
		RenamedC int   `yaml:"c"`
		D        []int `yaml:",flow"`
	}
}

var t = &T{}

func reloadConf(res string) {
	err := yaml.Unmarshal([]byte(res), t)
	if err != nil {
		glog.Error("error: %v", err)
	}
	glog.Info("\n\n", t)
}
func main() {
	// 可以没有，采用默认值
	clientConfig := constant.ClientConfig{
		TimeoutMs:      1 * 1000,
		ListenInterval: 2 * 1000,
		BeatInterval:   1 * 1000,
		LogDir:         "./logs/nacos",
		CacheDir:       "./nacos",
	}

	// 至少一个
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      "localhost",
			ContextPath: "/nacos",
			Port:        8848,
		},
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		glog.Error(err)
	}
	d := "natsmicro_dev"
	g := "natsmicro"
	res, err := configClient.GetConfig(vo.ConfigParam{
		DataId: d,
		Group:  g})

	if err != nil {
		glog.Error(err)
	}
	glog.Info(res)
	reloadConf(res)
	configClient.ListenConfig(vo.ConfigParam{
		DataId: d,
		Group:  g,
		OnChange: func(namespace, group, dataId, data string) {
			glog.Info("group:" + group + ", dataId:" + dataId + ", data:" + data)
			reloadConf(data)
		},
	})
	select {}
}
