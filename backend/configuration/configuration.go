package configuration

import (
	"io/ioutil"
	"motherbear/backend/constants"
	"sync"

	"motherbear/backend/logger"

	"gopkg.in/yaml.v2"
)

// Configuration is map data YAML file format.
type Configuration struct {
	Node          []Nodes
	Channel       []Channels
	Prometheus    Prometheus
	ETC           Etc
	Blockchain    Blockchain
	Authorization Authorization
}

// Nodes configurations.
type Nodes struct {
	Name string `yaml:"name"`
	IP   string `yaml:"ip"`
}

// Channels configurations.
type Channels struct {
	Name  string   `yaml:"name"`
	Nodes []string `yaml:",flow"`
}

// Prometheus configurations.
type Prometheus struct {
	PrometheusExternal string `yaml:"prometheusExternal"`
	PrometheusISAAC    string `yaml:"prometheusISAAC"`
	QueryPath          string `yaml:"queryPath"`
	CrawlingInterval   int    `yaml:"crawlingInterval"`
	NodeType           string `yaml:"nodeType"`
	JobNameOfgoloop    string `yaml:"jobNameOfgoloop"`
}

// Blockchain configurations
type Blockchain struct {
	CrawlingInterval int        `yaml:"crawlingInterval"`
	DB               []DBConfig `yaml:"db"`
}

// Etc configurations
type Etc struct {
	SessionTimeout     int        `yaml:"sessionTimeout"`
	Language           string     `yaml:"language"`
	LogLevel           int        `yaml:"loglevel"`
	DB                 []DBConfig `yaml:"db"`
	LoginLogoImagePath string     `yaml:"loginLogoImagePath"`
}

// Authorization API configuration
type Authorization struct {
	ThirdPartyUserAPI []string `yaml:"thirdPartyUserAPI"`
}

// DB configurations
type DBConfig struct {
	DBType   string `yaml:"type"`
	Id       string `yaml:"id"`
	Pass     string `yaml:"pass"`
	Database string `yaml:"database"`
	DBPath   string `yaml:"path"`
}

var instance *Configuration
var once sync.Once
var currentFilePath string

// Map: node information by name.
var nodeMap map[string]Nodes

// We used the Singletone pattern.
func init() {
	once.Do(func() {
		instance = &Configuration{}
	})
}

// Conf returns the instance of configuration.
func Conf() *Configuration {
	return instance
}

// GetFilePath returns the current configuration file path.
func GetFilePath() string {
	return currentFilePath
}

// InitConfigData loads YAML configuration file.
func InitConfigData(filePath string) *Configuration {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		logger.Error("[InitConfigData]Config file read fail!")
		logger.Panicln(err)
	}

	err = yaml.Unmarshal(data, instance)
	if err != nil {
		logger.Error("[InitConfigData]yaml Unmarshal fail!")
		logger.Panicln(err)
	}

	currentFilePath = filePath

	// Build map for nodes by name.
	buildNodeMap()

	return instance
}

func buildNodeMap() {
	nodeMap = make(map[string]Nodes)
	for i, n := range instance.Node {
		nodeMap[instance.Node[i].Name] = n
	}
}

//QueryNodeByNodeName returns Nodes data by name in constant time.
func QueryNodeByNodeName(nodeName string) Nodes {
	return nodeMap[nodeName]
}

//QueryNodeType returns Node Type.
func QueryNodeType() string {
	if Conf().Prometheus.NodeType == constants.NodeType1 {
		return constants.NodeType1
	} else if Conf().Prometheus.NodeType == constants.NodeType2 {
		return constants.NodeType2
	} else {
		return constants.NodeUnknown
	}
}

//QueryPrometheusIP returns Prometheus IP + Query path.
func QueryPrometheusIP() string {
	return Conf().Prometheus.PrometheusExternal + Conf().Prometheus.QueryPath
}

// ChangeConfigFile updates the YAML configuration file. s
func ChangeConfigFile(filePath string, conf *Configuration) {
	changeData, err := yaml.Marshal(&conf)
	if err != nil {
		logger.Error("[ChangeConfigFile]yaml Marshal fail!")
		logger.Panicln(err)
	}
	err = ioutil.WriteFile(filePath, changeData, 0644)
	if err != nil {
		logger.Error("[ChangeConfigFile]Config file WriteFile fail!")
		logger.Panicln(err)
	}

	// Build map for nodes by name.
	buildNodeMap()

	return
}
