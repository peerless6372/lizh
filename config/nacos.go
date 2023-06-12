// @Title
// @Description
// @Author lizh
// @Date   2023/6/11  5:25
package config

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type NacosClicfg struct {
	TimeoutMs   uint64 // 请求Nacos服务端的超时时间，默认是10000ms
	NamespaceId string // ACM的命名空间Id
	Endpoint    string // 当使用ACM时，需要该配置. https://help.aliyun.com/document_detail/130146.html
	RegionId    string // ACM&KMS的regionId，用于配置中心的鉴权
	AccessKey   string // ACM&KMS的AccessKey，用于配置中心的鉴权
	SecretKey   string // ACM&KMS的SecretKey，用于配置中心的鉴权
	OpenKMS     bool   // 是否开启kms，默认不开启，kms可以参考文档 https://help.aliyun.com/product/28933.html
	// 同时DataId必须以"cipher-"作为前缀才会启动加解密逻辑
	CacheDir             string // 缓存service信息的目录，默认是当前运行目录
	UpdateThreadNum      int    // 监听service变化的并发数，默认20
	NotLoadCacheAtStart  bool   // 在启动的时候不读取缓存在CacheDir的service信息
	UpdateCacheWhenEmpty bool   // 当service返回的实例列表为空时，不更新缓存，用于推空保护
	Username             string // Nacos服务端的API鉴权Username
	Password             string // Nacos服务端的API鉴权Password
	LogDir               string // 日志存储路径
	RotateTime           string // 日志轮转周期，比如：30m, 1h, 24h, 默认是24h
	MaxAge               int64  // 日志最大文件数，默认3
	LogLevel             string // 日志默认级别，值必须是：debug,info,warn,error，默认值是info
}
type NacosServercfg struct {
	ContextPath string // Nacos的ContextPath，默认/nacos，在2.0中不需要设置
	IpAddr      string // Nacos的服务地址
	Port        uint64 // Nacos的服务端口
	Scheme      string // Nacos的服务地址前缀，默认http，在2.0中不需要设置
	GrpcPort    uint64 // Nacos的 grpc 服务端口, 默认为 服务端口+1000, 不是必填
}

// 创建nacos客户端
func NewNacosCli(nacosCcfg NacosClicfg, nacosScfg NacosServercfg) (namingClient naming_client.INamingClient, configClient config_client.IConfigClient, err error) {
	//// 创建clientConfig
	//clientConfig := constant.ClientConfig{
	//	NamespaceId:         "e525eafa-f7d7-4029-83d9-008937f9d468", // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
	//	TimeoutMs:           5000,
	//	NotLoadCacheAtStart: true,
	//	LogDir:              "/tmp/nacos/log",
	//	CacheDir:            "/tmp/nacos/cache",
	//	LogLevel:            "debug",
	//}

	// 创建clientConfig的另一种方式
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(nacosCcfg.NamespaceId), //当namespace是public时，此处填空字符串。
		constant.WithTimeoutMs(nacosCcfg.TimeoutMs),
		constant.WithNotLoadCacheAtStart(nacosCcfg.NotLoadCacheAtStart),
		constant.WithLogDir(nacosCcfg.LogDir),
		constant.WithCacheDir(nacosCcfg.CacheDir),
		constant.WithLogLevel(nacosCcfg.LogLevel),
	)

	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      nacosScfg.IpAddr,
			ContextPath: nacosScfg.ContextPath,
			Port:        nacosScfg.Port,
			Scheme:      nacosScfg.Scheme,
		},
		//{
		//	IpAddr:      "console2.nacos.io",
		//	ContextPath: "/nacos",
		//	Port:        80,
		//	Scheme:      "http",
		//},
	}

	//// 创建serverConfig的另一种方式
	//serverConfigs := []constant.ServerConfig{
	//	*constant.NewServerConfig(
	//		"console1.nacos.io",
	//		80,
	//		constant.WithScheme("http"),
	//		constant.WithContextPath("/nacos"),
	//	),
	//	*constant.NewServerConfig(
	//		"console2.nacos.io",
	//		80,
	//		constant.WithScheme("http"),
	//		constant.WithContextPath("/nacos"),
	//	),
	//}

	//// 创建服务发现客户端
	//_, _ := clients.CreateNamingClient(map[string]interface{}{
	//	"serverConfigs": serverConfigs,
	//	"clientConfig":  clientConfig,
	//})

	//// 创建动态配置客户端
	//_, _ := clients.CreateConfigClient(map[string]interface{}{
	//	"serverConfigs": serverConfigs,
	//	"clientConfig":  clientConfig,
	//})

	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return namingClient, configClient, nil
}
