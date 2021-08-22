package admin

import (
	"gopkg.in/go-playground/assert.v1"
	"gopkg.in/yaml.v2"
	"motherbear/backend/configuration"
	"os"
	"testing"
)

var data = `
node :
  - name : node0
    ip : http://13.125.120.242:9080
  - name : node1
    ip : http://15.164.163.252:9080
  - name : node2
    ip : http://13.209.77.121:9080
  - name : node3
    ip : http://13.125.244.31:9080
`

const confFilePath string = "testConfiguration.yaml"

func TestGetChains(t *testing.T) {
	t.Skip("Skipping test GetChains. Should prepare the goloop to test.")
	Setup()
	defer Teardown()

	ip := configuration.Conf().Node[0].IP

	chains, err := GetChains(ip)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(chains), 3)
}

func TestGetChainsAllTry(t *testing.T) {
	t.Skip("Skipping test GetChainsAllTry. Should prepare the goloop to test.")
	Setup()
	defer Teardown()

	ipList := make([]string, 0)
	for _, value := range configuration.Conf().Node {
		ipList = append(ipList, value.IP)
	}

	chains, err := GetChainsAllTry(ipList)

	assert.Equal(t, err, nil)
	assert.Equal(t, len(chains), 3)
}

func TestGetChannelIDByName(t *testing.T) {
	t.Skip("Skipping test GetChannelIDByName. Should prepare the goloop to test.")
	Setup()
	defer Teardown()

	ipList := make([]string, 0)
	for _, value := range configuration.Conf().Node {
		ipList = append(ipList, value.IP)
	}

	channelName := "822027"
	correctChannelID := "0x822027"

	channelID, err := GetChannelIDByName(channelName, ipList)

	assert.Equal(t, err, nil)
	assert.Equal(t, channelID, correctChannelID)
}

func TestGetSystem(t *testing.T) {
	t.Skip("Skipping test GetSystem. Should prepare the goloop to test.")
	Setup()
	defer Teardown()

	for _, value := range configuration.Conf().Node {
		system, _ := GetSystem(value.IP)

		assert.NotEqual(t, 0, len(system.BuildVersion))
	}
}

func Setup() {
	conf := configuration.Configuration{}

	err := yaml.Unmarshal([]byte(data), &conf)
	if err == nil {
		configuration.ChangeConfigFile(confFilePath, &conf)
		configuration.InitConfigData(confFilePath)
	}
}

func Teardown() {
	_ = os.Remove(confFilePath)
}
