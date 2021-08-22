package prom_crawler

import (
	"fmt"
	"github.com/jasonlvhit/gocron"
	"motherbear/backend/configuration"
	"motherbear/backend/db"
	"motherbear/backend/polarbear"
	"motherbear/backend/prometheus"
	"os"
	"strconv"
	"testing"
	"time"

	"gopkg.in/yaml.v2"
)

var dataForLoopChain = `
node :
  - name : node0
    ip : http://34.97.23.118:9000
  - name : node1
    ip : http://34.97.202.145:9000
  - name : node2
    ip : http://34.97.177.55:9000
  - name : node3
    ip : http://34.97.16.82:9000
  - name : node4
    ip : http://34.97.219.137:9000

channel :
  - name : "loopchain_default1"
    nodes : [node0, node1, node2, node3, node4]
  - name : "loopchain_default2"
    nodes : [node0, node1, node2, node3, node4]
  - name : "loopchain_default3"
    nodes : [node0, node1, node2, node3, node4]

prometheus:
  prometheusExternal : http://localhost:9090
  prometheusISAAC : http://localhost:9090
  queryPath : /api/v1/query
  crawlingInterval : 5
  nodeType : loopchain     # Can be used Node type (loopchain, goloop)
  jobNameOfgoloop : goloop  # If using the goloop, should set the job name of prometheus.

blockchain:
  crawlingInterval: 10
  db:
    - type: sqlite3
      id: ""
      pass: ""
      database: ""    # For remote MySQL server, use tcp($SERVER_IP:3306):$DB_NAME.
      path: data/crawling_data.db

etc:
  sessionTimeout: 30
  language: ko
  loglevel: 1
  db:
    - type: sqlite3
      id: ""
      pass: ""
      database: ""
      path: data/isaac.db
  loginLogoImagePath: images/iconloop.png

authorization:
  thirdPartyUserAPI: [channels]			# Can be used API list is channels, nodes, blocks, txs.
`

var dataForGoloop = `
node:
  - name: node0
    ip: http://13.125.120.242:9080
  - name: node1
    ip: http://15.164.163.252:9080
  - name: node2
    ip: http://13.209.77.121:9080
  - name: node3
    ip: http://13.125.244.31:9080

channel:
  - name: "5d46aa"
    nodes: [node0, node1, node2, node3]
  - name: "822027"
    nodes: [node0, node1, node2, node3]
  - name: "2b127a"
    nodes: [node0, node1, node2, node3]

prometheus:
  prometheusExternal: http://13.125.120.242:9090
  prometheusISAAC: http://13.125.120.242:9090
  queryPath: /api/v1/query
  crawlingInterval: 5
  nodeType: goloop     # Can be used Node type (loopchain, goloop)
  jobNameOfgoloop: goloop  # If using the goloop, should set the job name of prometheus.

blockchain:
  crawlingInterval: 10
  db:
    - type: sqlite3
      id: ""
      pass: ""
      database: ""    # For remote MySQL server, use tcp($SERVER_IP:3306):$DB_NAME.
      path: data/crawling_data.db

etc:
  sessionTimeout: 30
  language: ko
  loglevel: 1
  db:
    - type: sqlite3
      id: ""
      pass: ""
      database: ""
      path: data/isaac.db
  loginLogoImagePath: images/iconloop.png

authorization:
  thirdPartyUserAPI: [channels]            # Can be used API list is channels, nodes, blocks, txs.
`

const dbPath string = ":memory:"
const confFilePath string = "testConfiguration.yaml"

func TestCrawlingFromLoopChain(t *testing.T) {
	t.Skip("Skipping test RequestPrometheusQuery. Should prepare prometheus server to test.")
	Setup(dataForLoopChain)
	defer Teardown()
	BeginToCrawl()

	job := func() {
		//floatString := strconv.FormatFloat(value.TimeStamp, 'f', 3, 64)
		location, _ := time.LoadLocation("UTC")
		loopchainInstance, _ := prometheus.GetPrometheusData()
		timeStamp := time.Unix(int64(loopchainInstance.TimeStamp), 0).In(location).Format(time.RFC3339)
		fmt.Print(timeStamp + " " + strconv.Itoa(loopchainInstance.Status) + " ")
		fmt.Println("")
	}

	interval := uint64(configuration.Conf().Prometheus.CrawlingInterval)
	gocron.Every(interval).Seconds().Do(job)
	<-gocron.Start()
}

func TestCrawlingFromGoloop(t *testing.T) {
	t.Skip("Skipping test RequestPrometheusQuery. Should prepare prometheus server to test.")
	Setup(dataForGoloop)
	defer Teardown()
	BeginToCrawl()

	job := func() {
		//floatString := strconv.FormatFloat(value.TimeStamp, 'f', 3, 64)
		location, _ := time.LoadLocation("UTC")
		goloopInstance, _ := prometheus.GetPrometheusData()
		timeStamp := time.Unix(int64(goloopInstance.TimeStamp), 0).In(location).Format(time.RFC3339)
		fmt.Print(timeStamp + " " + strconv.Itoa(goloopInstance.Status) + " ")
		fmt.Println("")
	}

	interval := uint64(configuration.Conf().Prometheus.CrawlingInterval)
	gocron.Every(interval).Seconds().Do(job)
	<-gocron.Start()
}

func Setup(data string) {
	conf := configuration.Configuration{}

	err := yaml.Unmarshal([]byte(data), &conf)
	if err == nil {
		configuration.ChangeConfigFile(confFilePath, &conf)
		configuration.InitConfigData(confFilePath)
	}

	db.InitDB("sqlite3", dbPath)
	db.InitCreateTable()

	polarbear.InitDB("sqlite3", dbPath)
}

func Teardown() {
	_ = os.Remove(confFilePath)
}
