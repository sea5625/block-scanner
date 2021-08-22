package prometheus

import (
	"motherbear/backend/isaacerror"
	"sync"
)

// PrometheusData ...
type PrometheusData struct {
	PrometheusChannelData []PrometheusChannelData
	TimeStamp             float64
	Status                int
}

// PrometheusChannelData ...
type PrometheusChannelData struct {
	Name   string
	Status int
	Nodes  []PrometheusNodesData
}

// PrometheusNodesData ...
type PrometheusNodesData struct {
	Name                 string
	BlockHeight          uint64
	PrevBlockHeight      uint64 //BlockHeight for alerting confirmation
	CountOfTX            uint64
	CountOfUnconfirmedTX uint64
	ResponseTimeInSec    float64
	IsLeader             int
	TimeStamp            string
	Status               int
	UnSyncBlockHoldInSec int
}

// Node and channel status.
const (
	CrawlingSuccess int = iota // Succeed to crawl the data from the prometheus.
	CrawlingFail               // Failed to crawl the data from the prometheus.
	WarmingUpCrawling
)

var instance PrometheusData
var once sync.Once

// We used the Singleton pattern.
func init() {
	once.Do(func() {
		instance = PrometheusData{}
	})
}

func SetPrometheusData(prometheusData *PrometheusData) {
	instance = *prometheusData
}

func GetPrometheusData() (PrometheusData, error) {
	var prometheusData PrometheusData
	if instance.TimeStamp == 0 {
		err := isaacerror.SysErrFailToGetPrometheusData
		return prometheusData, err
	}
	prometheusData = instance

	return prometheusData, nil
}

func GetPrometheusChannelData(channelName string) (PrometheusChannelData, error) {
	var prometheusChannelData PrometheusChannelData

	prometheusData, err := GetPrometheusData()
	if err != nil {
		return prometheusChannelData, err
	}
	if prometheusData.Status != CrawlingSuccess {
		err := isaacerror.SysErrFailToGetPrometheusData
		return prometheusChannelData, err
	}

	for i, value := range prometheusData.PrometheusChannelData {
		if value.Name == channelName {
			prometheusChannelData = prometheusData.PrometheusChannelData[i]
		}
	}

	return prometheusChannelData, nil
}

func IsExistingChannel(name string, prometheusData []PrometheusChannelData) (int, bool) {
	for i, value := range prometheusData {
		if name == value.Name {
			return i, true
		}
	}

	return 0, false
}

func IsExistingNode(name string, nodesData []PrometheusNodesData) (int, bool) {
	for i, value := range nodesData {
		if name == value.Name {
			return i, true
		}
	}

	return 0, false
}

func IsLastBlockHeight(nodesData []PrometheusNodesData) (uint64, uint64) {
	var lastBlockHeight, prevLastBlockHeight uint64
	lastBlockHeight = 0
	prevLastBlockHeight = 0
	for _, value := range nodesData {
		if lastBlockHeight < value.BlockHeight {
			lastBlockHeight = value.BlockHeight
		}
		if prevLastBlockHeight < value.PrevBlockHeight {
			prevLastBlockHeight = value.PrevBlockHeight
		}
	}
	return prevLastBlockHeight, lastBlockHeight
}

func GetUnSyncBlockHoldInSec(channelName string, nodeName string) (int, error) {
	prometheusData, err := GetPrometheusData()
	if err != nil {
		return 0, err
	}
	if prometheusData.Status != CrawlingSuccess {
		err := isaacerror.SysErrFailToGetPrometheusDataFromLC
		return 0, err
	}

	for _, value := range prometheusData.PrometheusChannelData {
		if value.Name == channelName {
			for _, node := range value.Nodes {
				if nodeName == node.Name {
					return node.UnSyncBlockHoldInSec, nil
				}
			}
		}
	}

	return 0, nil
}