package global

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"carbon-go-relay/utils"

	"github.com/toolkits/file"
)

//HTTPConfig define http configs
type HTTPConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
}

//SocketConfig define socket configs
type SocketConfig struct {
	Enabled bool   `json:"enabled"`
	Listen  string `json:"listen"`
	Timeout int    `json:"timeout"`
}

//RelayPattern give address and relevant patterns
type RelayPattern struct {
	AliasName       string   `json:"alias_name"`
	Address         string   `json:"address"`
	Patterns        []string `json:"patterns"`
	MaxWorkerNumber int      `json:"max_worker_number"`
}

//RelayClusterStruct give relay clusters
type RelayClusterStruct struct {
	RelayClusterList []RelayPattern `json:"relay_cluster_list"`
}

//TotalConfig define global configs
type TotalConfig struct {
	Debug                  bool                `json:"debug"`
	MaxBrubeckLength       int                 `json:"max_brubeck_length"`
	SendBatchSize          int                 `json:"send_batch_size"`
	SendPeriodMilliseconds int                 `json:"send_period_milliseconds"`
	HTTP                   *HTTPConfig         `json:"http"`
	Socket                 *SocketConfig       `json:"socket"`
	RelayCluster           *RelayClusterStruct `json:"relay_cluster"`
	StatsGroup             []string            `json:"stats_group"`
}

var (
	//ConfigFile includes glboal configs
	ConfigFile string
	//Hostname give the hostname
	Hostname   string
	config     *TotalConfig
	configLock = new(sync.RWMutex)
)

//InitalizeGlobalConstants initalize some global constants
func InitalizeGlobalConstants() {
	getHostname()
}

//Config return global config
func Config() *TotalConfig {
	configLock.RLock()
	defer configLock.RUnlock()
	return config
}

//ParseConfig parse config file
func ParseConfig(cfg string) {
	if cfg == "" {
		utils.Zlog.Fatal("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		utils.Zlog.Fatal("config file:", cfg, "is not existent. maybe you need `mv cfg.example.json cfg.json`")
	}

	ConfigFile = cfg

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		utils.Zlog.Fatal("read config file:", cfg, "fail:", err)
	}

	var c TotalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		utils.Zlog.Fatal("parse config file:", cfg, "fail:", err)
	}

	configLock.Lock()
	defer configLock.Unlock()
	config = &c

	utils.Zlog.Info("g.ParseConfig ok, file ", cfg)
}

func getHostname() {
	tmpHostname, err := os.Hostname()
	if err != nil {
		utils.Zlog.Fatalf("can not get the hostname")
	}
	Hostname = strings.Replace(tmpHostname, ".", "_", -1)
}
