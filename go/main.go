package main

// var NatsUrl = "nats://192.168.50.2:4222"
import (
	"github.com/gogf/gf/frame/g"
	"nacos_go_demo/utils"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		glog.Error(err)
	}
	d := "natsmicro_dev"
	g2 := "natsmicro"
	res, err := configClient.GetConfig(vo.ConfigParam{
		DataId: d,
		Group:  g2})

	if err != nil {
		glog.Error(err)
	}
	glog.Info(res)
	reloadConf(res)
	configClient.ListenConfig(vo.ConfigParam{
		DataId: d,
		Group:  g2,
		OnChange: func(namespace, group, dataId, data string) {
			glog.Info("group:" + group + ", dataId:" + dataId + ", data:" + data)
			reloadConf(data)
		},
	})

	server := g.Server()
	server.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("哈喽世界！")
	})
	var servicePort uint64 = 8000
	server.SetPort(int(servicePort))
	go func() {
		server.Start()
	}()
	//closeChan := make(chan struct{}, 0) // Used for underlying server closing event notification.

	internalIp := ip.InternalIP()
	// namingClient
	success, _ := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          internalIp,
		Port:        servicePort,
		ServiceName: "demo.go",
		Weight:      10,
		ClusterName: "a",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	})
	glog.Info("namingClient.RegisterInstance: ", success)

	service, _ := namingClient.GetService(vo.GetServiceParam{
		ServiceName: "demo.go",
		Clusters:    []string{"a"},
	})
	glog.Info("\nservice", service)

	instances, _ := namingClient.SelectAllInstances(vo.SelectAllInstancesParam{
		ServiceName: "demo.go",
		Clusters:    []string{"a"},
	})
	glog.Info("\n instances", instances)

	instance, _ := namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: "demo.go",
		Clusters:    []string{"a"},
	})

	glog.Info("\n SelectOneHealthyInstance", instance)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		glog.Info("get a signal: ", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			success, _ := namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
				Ip:          internalIp,
				Port:        servicePort,
				ServiceName: "demo.go",
				Cluster:     "a",
				Ephemeral:   true,
			})
			glog.Info("DeregisterInstance", success)
			glog.Info("kratos-demo exit")
			time.Sleep(time.Second)
			server.Shutdown()

			return
		case syscall.SIGHUP:
		default:
			return
		}
	}


	//}()
}
