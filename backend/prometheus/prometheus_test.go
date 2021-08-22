package prometheus

import (
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestSetPrometheusData(t *testing.T) {
	prometheus := PrometheusData{
		PrometheusChannelData: make([]PrometheusChannelData, 3),
		TimeStamp:             1,
		Status:                1,
	}

	prometheus.PrometheusChannelData[0].Name = "channel0"
	prometheus.PrometheusChannelData[1].Name = "channel1"
	prometheus.PrometheusChannelData[2].Name = "channel2"

	SetPrometheusData(&prometheus)

	assert.Equal(t, instance.PrometheusChannelData[0].Name, "channel0")
	assert.Equal(t, instance.PrometheusChannelData[1].Name, "channel1")
	assert.Equal(t, instance.PrometheusChannelData[2].Name, "channel2")
}

func TestGetPrometheusData(t *testing.T) {
	instance.PrometheusChannelData = make([]PrometheusChannelData, 3)
	instance.PrometheusChannelData[0].Name = "channel0"
	instance.PrometheusChannelData[1].Name = "channel1"
	instance.PrometheusChannelData[2].Name = "channel2"
	instance.TimeStamp = 1

	prometheus, _ := GetPrometheusData()

	assert.Equal(t, prometheus.PrometheusChannelData[0].Name, "channel0")
	assert.Equal(t, prometheus.PrometheusChannelData[1].Name, "channel1")
	assert.Equal(t, prometheus.PrometheusChannelData[2].Name, "channel2")
}

func TestGetPrometheusChannelData(t *testing.T) {
	instance.PrometheusChannelData = make([]PrometheusChannelData, 3)
	instance.PrometheusChannelData[0].Name = "channel0"
	instance.PrometheusChannelData[1].Name = "channel1"
	instance.PrometheusChannelData[2].Name = "channel2"
	instance.TimeStamp = 1
	instance.Status = CrawlingSuccess

	prometheusChannelData, _ := GetPrometheusChannelData("channel0")

	assert.Equal(t, prometheusChannelData.Name, "channel0")
}

func TestIsExistChannel(t *testing.T) {
	var loopChainData []PrometheusChannelData

	loopChainData = make([]PrometheusChannelData, 3)
	loopChainData[0].Name = "channel0"
	loopChainData[1].Name = "channel1"
	loopChainData[2].Name = "channel2"

	index, result := IsExistingChannel("channel1", loopChainData)
	assert.Equal(t, true, result)
	assert.Equal(t, 1, index)

	index, result = IsExistingChannel("channel3", loopChainData)
	assert.Equal(t, false, result)
	assert.Equal(t, 0, index)
}

func TestIsExistNode(t *testing.T) {
	var loopChainNodeData []PrometheusNodesData

	loopChainNodeData = make([]PrometheusNodesData, 3)
	loopChainNodeData[0].Name = "node0"
	loopChainNodeData[1].Name = "node1"
	loopChainNodeData[2].Name = "node2"

	index, result := IsExistingNode("node1", loopChainNodeData)
	assert.Equal(t, true, result)
	assert.Equal(t, 1, index)

	index, result = IsExistingNode("node3", loopChainNodeData)
	assert.Equal(t, false, result)
	assert.Equal(t, 0, index)
}

func TestIsLastBlockHeight(t *testing.T) {
	var loopChainNodeData []PrometheusNodesData

	loopChainNodeData = make([]PrometheusNodesData, 3)
	loopChainNodeData[0].Name = "node0"
	loopChainNodeData[0].BlockHeight = 2
	loopChainNodeData[0].PrevBlockHeight = 1
	loopChainNodeData[1].Name = "node1"
	loopChainNodeData[1].BlockHeight = 3
	loopChainNodeData[1].PrevBlockHeight = 2
	loopChainNodeData[2].Name = "node2"
	loopChainNodeData[2].BlockHeight = 4
	loopChainNodeData[2].PrevBlockHeight = 3

	prevBlockHeight, blockHeight := IsLastBlockHeight(loopChainNodeData)

	assert.Equal(t, uint64(4), blockHeight)
	assert.Equal(t, uint64(3), prevBlockHeight)
}
