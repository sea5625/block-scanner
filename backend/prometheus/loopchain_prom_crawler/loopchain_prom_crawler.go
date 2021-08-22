package loopchain_prom_crawler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"motherbear/backend/db"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"motherbear/backend/polarbear"
	"motherbear/backend/prometheus"
	"net/http"
	"strconv"
	"time"
)

type LoopChainPrometheus struct {
}

type prometheusResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string                      `json:"resultType"`
		Result     []loopchainPrometheusResult `json:"result"`
	} `json:"data"`
}

// loopchainPrometheusResult ...
type loopchainPrometheusResult struct {
	Metric struct {
		Name     string `json:"__name__"`
		Alias    string `json:"alias"`
		Channel  string `json:"channel"`
		Instance string `json:"instance"`
		Job      string `json:"job"`
	} `json:"metric"`
	Values [][]interface{} `json:"values"`
}

var query []string = []string{
	"block_height",         // "BLOCKHEIGHT"
	"tx_count",             // "COUNTOFTX"
	"unconfirmed_tx_count", // "COUNTOFUNCONFIRMEDTX"
	"is_leader",            // "ISLEADER"
	"response_time",        // "RESPONSETIMEINSEC"
}

// Node and channel status.
const (
	NodeNormal                     int = iota // Node is normal
	NodeUnSyncedBlock                         // Node block is unsync.
	NodeSlowResponse                          // Node responses slowly.
	NodeUnSyncBlockAndSlowResponse            // Node block is unsync + Node responses slowly
	ChannelNormal                             // Channel is normal.
	ChannelAbnormal                           // Channel is abnormal.
)

const CrawlingRangeTimeSec int = 60 // Time to collect at once in prometheus, Sec

func (c LoopChainPrometheus) Crawler() (*prometheus.PrometheusData, error) {
	ip := configuration.Conf().Prometheus.PrometheusISAAC

	queryPath := configuration.Conf().Prometheus.QueryPath

	prometheusData, err := requestPrometheusQuery(ip+queryPath, query)
	if err != nil {
		logger.Error(isaacerror.SysErrFailToGetPrometheusDataFromLC.Error())
		return nil, err
	}

	prometheusResult := prometheusData.Data.Result

	var loopchainChannelData []prometheus.PrometheusChannelData
	loopchainChannelData = make([]prometheus.PrometheusChannelData, 0)

	// This value is prometheus crawling time
	var unixTime float64
	timeStamp := ""

	// Convert prometheus result to our response data format.
	for _, value := range prometheusResult {
		channelIndex, isExistingChannel := prometheus.IsExistingChannel(value.Metric.Channel, loopchainChannelData)
		if isExistingChannel == false {
			channelIndex = len(loopchainChannelData)

			var data prometheus.PrometheusChannelData
			data.Name = value.Metric.Channel
			data.Status = ChannelNormal
			loopchainChannelData = append(loopchainChannelData, data)
		}

		nodeIndex, isExistingNode:= prometheus.IsExistingNode(value.Metric.Alias, loopchainChannelData[channelIndex].Nodes)
		if isExistingNode == false {
			nodeIndex = len(loopchainChannelData[channelIndex].Nodes)

			var data prometheus.PrometheusNodesData
			data.Name = value.Metric.Alias
			data.Status = NodeNormal
			data.UnSyncBlockHoldInSec = 0
			loopchainChannelData[channelIndex].Nodes = append(loopchainChannelData[channelIndex].Nodes, data)
		} else {
			unSyncBlockHoldInSec, _ := prometheus.GetUnSyncBlockHoldInSec(loopchainChannelData[channelIndex].Name,
				loopchainChannelData[channelIndex].Nodes[nodeIndex].Name )
			loopchainChannelData[channelIndex].Nodes[nodeIndex].UnSyncBlockHoldInSec = unSyncBlockHoldInSec
		}

		//  Choose one timestamp of metric in metrics.
		if timeStamp == "" {
			//  The last time stamp in metrics.
			unixTime = value.Values[len(value.Values)-1][0].(float64)
			location, err := time.LoadLocation("UTC")
			if err != nil {
				logger.Error(isaacerror.SysErrFailToLoadTimeLocation.Error())
				return nil, err
			}

			// Set crawling interval to begin crawling.
			//var crawlingInterval float64

			// Got one metric only. prometheusResult is warming.
			if len(value.Values) == 1 {
				return nil, nil // We decide to interprete 'warming state' if return (nil, nil).
			}
			//else {
			//	// Got many metrics. So calculate the interval for crawling and begin to crawl the data.
			//	crawlingInterval =
			//		value.Values[len(value.Values)-1][0].(float64) -
			//			value.Values[len(value.Values)-2][0].(float64)
			//}
			//crawlingIntervalInSec = int(math.Round(crawlingInterval))

			// Set time stamp from current metric.
			timeStamp = time.Unix(int64(unixTime), 0).In(location).Format(time.RFC3339)
		}

		// Use this time stamp for this metrics.
		loopchainChannelData[channelIndex].Nodes[nodeIndex].TimeStamp = timeStamp

		// Parse metrics.
		switch value.Metric.Name {
		case "block_height":
			valueUInt64, err := strconv.ParseUint(value.Values[len(value.Values)-1][1].(string), 10, 64)
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToInt.Error())
				return nil, err
			}
			loopchainChannelData[channelIndex].Nodes[nodeIndex].BlockHeight = valueUInt64

			valueUInt64Prev, err := strconv.ParseUint(value.Values[0][1].(string), 10, 64)
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToInt.Error())
				return nil, err
			}
			loopchainChannelData[channelIndex].Nodes[nodeIndex].PrevBlockHeight = valueUInt64Prev
		case "tx_count":
			valueUInt64, err := strconv.ParseUint(value.Values[len(value.Values)-1][1].(string), 10, 64)
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToInt.Error())
				return nil, err
			}
			loopchainChannelData[channelIndex].Nodes[nodeIndex].CountOfTX = valueUInt64
		case "unconfirmed_tx_count":
			valueUInt64, err := strconv.ParseUint(value.Values[len(value.Values)-1][1].(string), 10, 64)
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToInt.Error())
				return nil, err
			}
			loopchainChannelData[channelIndex].Nodes[nodeIndex].CountOfUnconfirmedTX = valueUInt64
		case "is_leader":
			valueInt, err := strconv.Atoi(value.Values[len(value.Values)-1][1].(string))
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToInt.Error())
				return nil, err
			}
			loopchainChannelData[channelIndex].Nodes[nodeIndex].IsLeader = valueInt
		case "response_time":
			valueFloat, err := strconv.ParseFloat(value.Values[len(value.Values)-1][1].(string), 64)
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToFloat.Error())
				return nil, err
			}
			loopchainChannelData[channelIndex].Nodes[nodeIndex].ResponseTimeInSec = valueFloat
		}

	}

	if len(loopchainChannelData) == 0 {
		err := isaacerror.SysErrFailToGetPrometheusDataFromLC
		logger.Error(err.Error())
		return nil, err
	}

	for chIDX, channel := range loopchainChannelData {
		prevLastBlockHeight, lastBlockHeight := prometheus.IsLastBlockHeight(channel.Nodes)

		// Read configuration value.
		configurationAlertDataTB := db.GetAlertConfigInfoByName(channel.Name)

		// Check the status of all nodes in every channel.
		for ndIDX, node := range channel.Nodes {
			if lastBlockHeight > node.BlockHeight && prevLastBlockHeight > node.PrevBlockHeight {
				if loopchainChannelData[chIDX].Nodes[ndIDX].UnSyncBlockHoldInSec == 0 {
					loopchainChannelData[chIDX].Nodes[ndIDX].UnSyncBlockHoldInSec =	CrawlingRangeTimeSec
				} else {
					loopchainChannelData[chIDX].Nodes[ndIDX].UnSyncBlockHoldInSec +=
						configuration.Conf().Prometheus.CrawlingInterval + 1
				}
				if loopchainChannelData[chIDX].Nodes[ndIDX].UnSyncBlockHoldInSec >
					configurationAlertDataTB.MAX_TIME_SEC_FOR_UNSYNC {
					channelPK := db.GetChannelPK(channel.Name)
					loopchainChannelData[chIDX].Nodes[ndIDX].Status = NodeUnSyncedBlock
					msg := fmt.Sprintf("[%s] block height [%d] is unsync [%d]",
						node.Name, uint64(node.BlockHeight), uint64(lastBlockHeight))
					logger.Symptomf(channel.Name, constants.UnsyncBlock, msg)
					err := polarbear.AddPeerSymptom(channel.Name, channelPK, constants.UnsyncBlock, msg)
					if err != nil {
						logger.Error("AddPeerSymptom, Symptom insert UnsyncBlock failed!")
						logger.Errorf("%v+", err)
					}
				}
			} else {
				loopchainChannelData[chIDX].Nodes[ndIDX].UnSyncBlockHoldInSec = 0
			}
			if float64(configurationAlertDataTB.MAX_TIME_SEC_FOR_RESPONSE) < node.ResponseTimeInSec {
				if loopchainChannelData[chIDX].Nodes[ndIDX].Status == NodeUnSyncedBlock {
					loopchainChannelData[chIDX].Nodes[ndIDX].Status = NodeUnSyncBlockAndSlowResponse
				} else {
					loopchainChannelData[chIDX].Nodes[ndIDX].Status = NodeSlowResponse
				}
				channelPK := db.GetChannelPK(channel.Name)
				msg := fmt.Sprintf("[%s]response time slowly [%f] sec",
					node.Name, node.ResponseTimeInSec)
				logger.Symptomf(channel.Name, constants.SlowResponse, msg)
				err := polarbear.AddPeerSymptom(channel.Name, channelPK, constants.SlowResponse, msg)
				if err != nil {
					logger.Error("AddPeerSymptom, Symptom insert SlowResponse failed!")
					logger.Errorf("%v+", err)
				}
			}

			if loopchainChannelData[chIDX].Nodes[ndIDX].Status > 0 {
				loopchainChannelData[chIDX].Status = ChannelAbnormal
			}
			logger.Debugf("[ch:%s][node:%s][status:%d] Succeed to crawling node data.",
				loopchainChannelData[chIDX].Name,
				loopchainChannelData[chIDX].Nodes[ndIDX].Name,
				loopchainChannelData[chIDX].Nodes[ndIDX].Status)
		}
		logger.Debugf("[ch:%s][status:%d] Succeed to crawling channel data.",
			loopchainChannelData[chIDX].Name,
			loopchainChannelData[chIDX].Status)
	}
	var loopchainData prometheus.PrometheusData
	loopchainData.PrometheusChannelData = loopchainChannelData
	loopchainData.TimeStamp = unixTime
	loopchainData.Status = prometheus.CrawlingSuccess
	return &loopchainData, nil
}

// requestPrometheusQuery is function to query data from prometheus.
func requestPrometheusQuery(url string, query []string) (*prometheusResult, error) {
	var output prometheusResult

	queryString := "{__name__=~\""
	for i, value := range query {
		queryString += value
		if i != (len(query) - 1) {
			queryString += "|"
		}
	}
	queryString += "\"}[" + strconv.Itoa(CrawlingRangeTimeSec) + "s]"

	uri := url + "?query=" + queryString

	res, err := http.Get(uri)
	if err != nil {
		logger.Error(err.Error())
		errorMessage := isaacerror.SysErrFailToConnectionPrometheus
		logger.Error(errorMessage.Error())
		return nil, errorMessage
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(err.Error())
		errorMessage := isaacerror.SysErrFailToReadBodyPrometheus
		logger.Error(errorMessage.Error())
		return nil, errorMessage
	}

	err = res.Body.Close()
	if err != nil {
		logger.Error(err.Error())
		errorMessage := isaacerror.SysErrFailToReadBodyClosePrometheus
		logger.Error(errorMessage.Error())
		return nil, errorMessage
	}

	err = json.Unmarshal(resData, &output)
	if err != nil {
		logger.Error(err.Error())
		errorMessage := isaacerror.SysErrFailToUnmarshalPrometheusData
		logger.Error(errorMessage.Error())
		return nil, errorMessage
	}

	if output.Status != "success" {
		err := isaacerror.SysErrFailToGetPrometheusData
		logger.Error(err.Error())
		return nil, err
	}

	return &output, nil
}
