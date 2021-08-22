package goloop_prom_crawler

import (
	"gopkg.in/go-playground/assert.v1"
	"gopkg.in/yaml.v2"
	"motherbear/backend/configuration"
	"motherbear/backend/db"
	"motherbear/backend/polarbear"
	"motherbear/backend/prometheus"
	"os"
	"testing"
)

var data = `
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

func TestGoloopPrometheus_Crawler(t *testing.T) {
	t.Skip("Skipping test CrawlingGoloopData. Should prepare prometheus server to test.")
	Setup()
	defer Teardown()

	goloop := GoloopPrometheus{}
	goloopData, err := goloop.Crawler()

	assert.NotEqual(t, nil, goloopData)
	assert.Equal(t, nil, err)
	assert.Equal(t, len(configuration.Conf().Channel), len(goloopData.PrometheusChannelData))
	assert.Equal(t, len(configuration.Conf().Node), len(goloopData.PrometheusChannelData[0].Nodes))
}

func TestAddJobName(t *testing.T) {
	t.Skip("Skipping test addJobName. Should prepare prometheus server to test.")
	Setup()
	defer Teardown()

	correctJobName := configuration.Conf().Prometheus.JobNameOfgoloop

	queryList := addJobName(query)

	for i, value := range queryList {
		assert.Equal(t, correctJobName+"_"+query[i], value)
	}
}

func TestRequestPrometheusQuery(t *testing.T) {
	t.Skip("Skipping test requestPrometheusQuery. Should prepare prometheus server to test.")

	Setup()
	defer Teardown()

	queryList := addJobName(query)
	prometheusIP := configuration.Conf().Prometheus.PrometheusISAAC + configuration.Conf().Prometheus.QueryPath
	prometheusData, _ := requestPrometheusQuery(prometheusIP, queryList)
	assert.Equal(t, "success", prometheusData.Status)
	assert.NotEqual(t, 0, len(prometheusData.Data.Result))
}

func TestIsNodeStatus(t *testing.T) {
	t.Skip("Skipping test Check NodeStatus. Should prepare prometheus server to test.")
	Setup()
	defer Teardown()

	goloop := GoloopPrometheus{}
	goloopData, err := goloop.Crawler()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(goloopData.PrometheusChannelData))

	for _, channel := range goloopData.PrometheusChannelData {
		channel.Status = ChannelNormal
		prevLastBlockHeight, lastBlockHeight := prometheus.IsLastBlockHeight(channel.Nodes)
		configurationAlertDataTB := &db.CONFIGURATION_DATA_ALERT_TB{}
		configurationAlertDataTB = db.GetAlertConfigInfoByName(channel.Name)

		for _, node := range channel.Nodes {
			if lastBlockHeight > node.BlockHeight && prevLastBlockHeight > node.PrevBlockHeight {
				node.Status = NodeUnSyncedBlock
			}
			if float64(configurationAlertDataTB.MAX_TIME_SEC_FOR_RESPONSE) < node.ResponseTimeInSec {
				if node.Status == NodeUnSyncedBlock {
					node.Status = NodeUnSyncBlockAndSlowResponse
				} else {
					node.Status = NodeSlowResponse
				}
			}
		}
	}
}

func TestIsChannelStatus(t *testing.T) {
	t.Skip("Skipping test Check NodeStatus. Should prepare prometheus server to test.")
	Setup()
	defer Teardown()

	goloop := GoloopPrometheus{}
	goloopData, err := goloop.Crawler()

	assert.Equal(t, nil, err)
	assert.NotEqual(t, 0, len(goloopData.PrometheusChannelData))

	for _, channel := range goloopData.PrometheusChannelData {
		channel.Status = ChannelNormal
		configurationAlertDataTB := &db.CONFIGURATION_DATA_ALERT_TB{}
		configurationAlertDataTB = db.GetAlertConfigInfoByName(channel.Name)

		for _, node := range channel.Nodes {
			node.Status = NodeUnSyncedBlock
			if float64(configurationAlertDataTB.MAX_TIME_SEC_FOR_RESPONSE) < node.ResponseTimeInSec {
				if node.Status == NodeUnSyncedBlock {
					node.Status = NodeUnSyncBlockAndSlowResponse
				} else {
					node.Status = NodeSlowResponse
				}
			}

			assert.Equal(t, node.Status, NodeUnSyncedBlock)

			if node.Status > 0 {
				channel.Status = ChannelAbnormal
			}
		}
		assert.Equal(t, channel.Status, ChannelAbnormal)
	}
}

func Setup() {
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
