package configuration

import (
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
)

var data = `
node :
  - name : node1
    ip : http://192.31.12.12
  - name : node2
    ip : http://192.31.12.13
  - name : node3
    ip : http://192.31.12.14
  - name : node4
    ip : http://192.31.12.15 

channel :
  - name : loopchain_default
    nodes : [node1, node2, node3, node4]
  - name : loopchain_default2
    nodes : [node1, node2]

prometheus :
  prometheusExternal : http://192.31.12.30:9090
  prometheusISAAC : http://localhost:9090
  queryPath: /api/v1/query
  crawlingInterval: 5
  nodeType: loopchain     # Can be used Node type (loopchain, goloop) 

blockchain:
  crawlingInterval: 5
  db:
    - type: sqlite3
      id: ""
      pass: ""
      database: ""
      path: data/crawling_data.db

etc :
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

authorizationAPI:
  thirdPartyUser: channels, nodes, blocks, txs			#channels, nodes, blocks, txs
`
var changeData = `
node :
  - name : node1
    ip : http://192.31.12.15
  - name : node2
    ip : http://192.31.12.16
  - name : node3
    ip : http://192.31.12.17
  - name : node4
    ip : http://192.31.12.18 

channel :
  - name : loopchain_default3
    nodes : [node1, node2, node3, node4]
  - name : loopchain_default4
    nodes : [node1, node2]        

prometheus :
  prometheusExternal : http://192.31.12.30:9090
  prometheusISAAC : http://localhost:9090
  queryPath: /api/v1/query
  crawlingInterval: 5
  nodeType: loopchain     # Can be used Node type (loopchain, goloop) 

blockchain:
  crawlingInterval: 10
  db:
    - type: sqlite3
      id: ""
      pass: ""
      database: ""
      path: data/crawling_data.db

etc :
  sessionTimeout: 60
  language: ko
  loglevel: 1
  db:
    - type: sqlite3
      id: ""
      pass: ""
      database: ""
      path: data/isaac.db
  loginLogoImagePath: images/iconloop.png

authorizationAPI:
  thirdPartyUser: channels, nodes, blocks, txs			#channels, nodes, blocks, txs
  
`

func TestYamlFileLoad(t *testing.T) {

	const confFilePath string = "myConf.yaml"
	InitConfigData(confFilePath)

	var node Nodes
	node = QueryNodeByNodeName("node0")

	if node.IP != "https://int-test-ctz.solidwallet.io" {
		t.Error("No such a node IP in node map.")
	}
	if node.Name != "node0" {
		t.Error("No such a node in node map.")
	}

}

func TestCreateYamlData(t *testing.T) {
	conf := Configuration{}

	err := yaml.Unmarshal([]byte(data), &conf)
	if err != nil {
		t.Error("Data Unmarshal Error ")
	}
	if (len(conf.Node) != 4) || (len(conf.Channel) != 2) {
		t.Error("Get Data Count Error")
	}
	if (conf.Node[0].Name != "node1") || (conf.Node[0].IP != "http://192.31.12.12") {
		t.Error("Get [node1] Data Error")
	}
	if (conf.Node[1].Name != "node2") || (conf.Node[1].IP != "http://192.31.12.13") {
		t.Error("Get [node2] Data Error")
	}
	if (conf.Node[2].Name != "node3") || (conf.Node[2].IP != "http://192.31.12.14") {
		t.Error("Get [node3] Data Error")
	}
	if (conf.Node[3].Name != "node4") || (conf.Node[3].IP != "http://192.31.12.15") {
		t.Error("Get [node4] Data Error")
	}
	if (conf.Channel[0].Name != "loopchain_default") || (conf.Channel[0].Nodes[0] != "node1") ||
		(conf.Channel[0].Nodes[1] != "node2") || (conf.Channel[0].Nodes[2] != "node3") ||
		(conf.Channel[0].Nodes[3] != "node4") {
		t.Error("Get [loopchain_default] Data Error")
	}
	if (conf.Channel[1].Name != "loopchain_default2") || (conf.Channel[1].Nodes[0] != "node1") ||
		(conf.Channel[1].Nodes[1] != "node2") {
		t.Error("Get [loopchain_default2] Data Error")
	}
	if conf.Prometheus.PrometheusExternal != "http://192.31.12.30:9090" {
		t.Error("Get [PrometheusExternal] Data Error")
	}
	if conf.Prometheus.PrometheusISAAC != "http://localhost:9090" {
		t.Error("Get [PrometheusISAAC] Data Error")
	}
	if conf.Prometheus.QueryPath != "/api/v1/query" {
		t.Error("Get [prometheus queryPath] Data Error")
	}
	if conf.Prometheus.CrawlingInterval != 5 {
		t.Error("Get [prometheus crawlingInterval] Data Error")
	}
	if conf.ETC.SessionTimeout != 30 {
		t.Error("Get [etc] Data Error")
	}
	if conf.ETC.LogLevel != 1 {
		t.Error("Cannot read Blockchain.crawlingInterval data")
	}
	if conf.Blockchain.CrawlingInterval != 5 {
		t.Error("Cannot read Blockchain.crawlingInterval data")
	}
	if conf.Blockchain.DB[0].DBType != "sqlite3" {
		t.Error("Cannot read DB[0].Target data")
	}
	if conf.ETC.DB[0].DBType != "sqlite3" {
		t.Error("Cannot read DB[0].DBType data")
	}

}

func TestCreateYamlFile(t *testing.T) {
	conf := Configuration{}

	err := yaml.Unmarshal([]byte(changeData), &conf)
	if err != nil {
		t.Error("Data Unmarshal Error ")
	}
	d, err := yaml.Marshal(&conf)
	if err != nil {
		t.Error("Data Marshal Error")
	}

	err = ioutil.WriteFile("changedYAML.yaml", d, 0644)
	if err != nil {
		t.Error("ioutil WriteFile Error")
	}

	dataChange, err := ioutil.ReadFile("changedYAML.yaml")
	if err != nil {
		t.Error("changed.yaml read Error")
	}
	confChange := &Configuration{}

	err = yaml.Unmarshal(dataChange, &confChange)
	if err != nil {
		t.Error("Data Unmarshal Error ")
	}

	if (len(confChange.Node) != 4) || (len(confChange.Channel) != 2) {
		t.Error("Get Data Count Erorr")
	}
	if (confChange.Node[0].Name != "node1") || (confChange.Node[0].IP != "http://192.31.12.15") {
		t.Error("Get [node1] Change Data Error")
	}
	if (confChange.Node[1].Name != "node2") || (confChange.Node[1].IP != "http://192.31.12.16") {
		t.Error("Get [node2] Change Data Error")
	}
	if (confChange.Node[2].Name != "node3") || (confChange.Node[2].IP != "http://192.31.12.17") {
		t.Error("Get [node3] Change Data Error")
	}
	if (confChange.Node[3].Name != "node4") || (confChange.Node[3].IP != "http://192.31.12.18") {
		t.Error("Get [node4] Change Data Error")
	}
	if (confChange.Channel[0].Name != "loopchain_default3") || (confChange.Channel[0].Nodes[0] != "node1") ||
		(confChange.Channel[0].Nodes[1] != "node2") || (confChange.Channel[0].Nodes[2] != "node3") ||
		(confChange.Channel[0].Nodes[3] != "node4") {
		t.Error("Get [loopchain_default] Change Data Error")
	}
	if (confChange.Channel[1].Name != "loopchain_default4") || (confChange.Channel[1].Nodes[0] != "node1") ||
		(confChange.Channel[1].Nodes[1] != "node2") {
		t.Error("Get [loopchain_default2] Change Data Error")
	}
	if conf.Prometheus.PrometheusExternal != "http://192.31.12.30:9090" {
		t.Error("Get [PrometheusExternal] Data Error")
	}
	if conf.Prometheus.PrometheusISAAC != "http://localhost:9090" {
		t.Error("Get [PrometheusISAAC] Data Error")
	}
	if confChange.Prometheus.QueryPath != "/api/v1/query" {
		t.Error("Get [prometheus queryPath] Change Data Error")
	}
	if confChange.Prometheus.CrawlingInterval != 5 {
		t.Error("Get [prometheus crawlingInterval] Change Data Error")
	}
	if confChange.ETC.SessionTimeout != 60 {
		t.Error("Get [etc] Change Data Error")
	}
	if confChange.ETC.LogLevel != 1 {
		t.Error("Cannot read Blockchain.crawlingInterval data")
	}

	if confChange.Blockchain.CrawlingInterval != 10 {
		t.Error("Cannot read Blockchain.crawlingInterval data")
	}
	if conf.Blockchain.DB[0].DBType != "sqlite3" {
		t.Error("Cannot read DB[0].Target data")
	}
	if conf.ETC.DB[0].DBType != "sqlite3" {
		t.Error("Cannot read DB[0].DBType data")
	}

}

// To-Do: Implement singletone test case.
