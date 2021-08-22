package loopchain_prom_crawler

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

const dbPath string = ":memory:"
const confFilePath string = "testConfiguration.yaml"

var testIP string = "http://localhost:9090/api/v1/query"
var testQuery []string = []string{
	"block_height",         // "BLOCKHEIGHT"
	"tx_count",             // "COUNTOFTX"
	"unconfirmed_tx_count", // "COUNTOFUNCONFIRMEDTX"
	"is_leader",            // "ISLEADER"
	"response_time",        // "RESPONSETIMEINSEC"
}

func TestCrawlingLoopchainData(t *testing.T) {
	t.Skip("Skipping test CrawlingLoopChainData. Should prepare prometheus server to test.")
	Setup()
	defer Teardown()

	loopchain := LoopChainPrometheus{}
	loopChainData, err := loopchain.Crawler()

	assert.NotEqual(t, nil, loopChainData)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(loopChainData.PrometheusChannelData))
	assert.Equal(t, 5, len(loopChainData.PrometheusChannelData[0].Nodes))
}

func TestRequestPrometheusQuery(t *testing.T) {
	t.Skip("Skipping test requestPrometheusQuery. Should prepare prometheus server to test.")

	prometheusData, _ := requestPrometheusQuery(testIP, testQuery)
	assert.Equal(t, "success", prometheusData.Status)
}

func TestIsNodeStatus(t *testing.T) {
	t.Skip("Skipping test Check NodeStatus. Should prepare prometheus server to test.")
	Setup()
	defer Teardown()

	loopchain := LoopChainPrometheus{}
	loopChainData, err := loopchain.Crawler()

	assert.NotEqual(t, 0, len(loopChainData.PrometheusChannelData))
	assert.Equal(t, nil, err)

	for _, channel := range loopChainData.PrometheusChannelData {
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

			assert.Equal(t, node.Status, NodeNormal)
		}
	}
}

func TestIsChannelStatus(t *testing.T) {
	t.Skip("Skipping test Check NodeStatus. Should prepare prometheus server to test.")
	Setup()
	defer Teardown()

	loopchain := LoopChainPrometheus{}
	loopChainData, err := loopchain.Crawler()

	assert.NotEqual(t, 0, len(loopChainData.PrometheusChannelData))
	assert.Equal(t, nil, err)

	for _, channel := range loopChainData.PrometheusChannelData {
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
