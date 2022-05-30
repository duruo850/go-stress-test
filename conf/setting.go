package conf

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// array 自定义数组参数
type array []string

// String string
func (a *array) String() string {
	return fmt.Sprint(*a)
}

// Set set
func (a *array) Set(s string) error {
	*a = append(*a, s)

	return nil
}

var (
	Concurrency     uint64  = 1       // 并发数
	TotalNumber     uint64  = 1       // 请求数(单个并发/协程)
	DebugStr                = "false" // 是否是debug
	RequestURL              = ""      // 压测的url 目前支持，http/https ws/wss
	Path                    = ""      // curl文件路径 http接口压测，自定义参数设置
	Verify                  = ""      // Verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json
	Headers         array             // 自定义头信息传递给服务器
	Body            = ""              // HTTP POST方式传送数据
	MaxCon          = 1               // 单个连接最大请求数
	Code            = 200             //成功状态码
	Http2           = false           // 是否开http2.0
	Keepalive       = false           // 是否开启长连接
	SocketKeepAlive = true            // 是否开启sccket级别的keepavlie
	ConnRetry       = 3               // 建立连接重试次数
	WriteData       = ""              // 写入的数据
	ConnectionMode  = 2               // 1:顺序建立长链接 2:并发建立长链接
)

func init() {
	// 配置文件参数
	conf, err := readConf("conf/setting.yaml")
	if err == nil {
		Concurrency = conf.StressTest.Concurrency
		TotalNumber = conf.StressTest.TotalNumber
		RequestURL = conf.StressTest.Url

		for k, v := range conf.StressTest.Headers {
			Headers = append(Headers, fmt.Sprintf("%s=%s", k, v))
		}

		WriteData = conf.StressTest.WriteData
		Keepalive = conf.StressTest.KeepAlive
		SocketKeepAlive = conf.StressTest.SocketKeepAlive
		ConnRetry = conf.StressTest.ConnRetry
		ConnectionMode = conf.StressTest.ConnectionMode
	}
	fmt.Println("setting.yaml", conf)

	// 启动参数
	flag.Uint64Var(&Concurrency, "c", Concurrency, "并发数")
	flag.Uint64Var(&TotalNumber, "n", TotalNumber, "请求数(单个并发/协程)")
	flag.StringVar(&DebugStr, "d", DebugStr, "调试模式")
	flag.StringVar(&RequestURL, "u", RequestURL, "压测地址")
	flag.StringVar(&Path, "p", Path, "curl文件路径")
	flag.StringVar(&Verify, "v", Verify, "验证方法 http 支持:statusCode、json webSocket支持:json")
	flag.Var(&Headers, "H", "自定义头信息传递给服务器 示例:-H 'Content-Type: application/json'")
	flag.StringVar(&Body, "data", Body, "HTTP POST方式传送数据")
	flag.IntVar(&MaxCon, "m", MaxCon, "单个host最大连接数")
	flag.IntVar(&Code, "Code", Code, "请求成功的状态码")
	flag.BoolVar(&Http2, "Http2", Http2, "是否开http2.0")
	flag.BoolVar(&Keepalive, "k", Keepalive, "是否开启长连接")
	flag.StringVar(&WriteData, "wd", WriteData, "写入的数据")
	// 解析参数
	flag.Parse()
}

type Conf struct {
	StressTest struct {
		Concurrency     uint64            `yaml:"concurrency"`
		TotalNumber     uint64            `yaml:"totalNumber"`
		Url             string            `yaml:"url"`
		Headers         map[string]string `yaml:"headers"`
		WriteData       string            `yaml:"writeData"`
		KeepAlive       bool              `yaml:"keep_alive"`
		SocketKeepAlive bool              `yaml:"socket_keep_alive"`
		ConnRetry       int               `yaml:"conn_retry"`
		ConnectionMode  int               `yaml:"connection_mode"`
	}
}

func readConf(filename string) (*Conf, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf Conf
	err = yaml.Unmarshal(buf, &conf)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}
	return &conf, nil
}
