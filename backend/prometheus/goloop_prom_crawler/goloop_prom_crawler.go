package goloop_prom_crawler

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
	"strings"
	"time"
)

type GoloopPrometheus struct {
}

type prometheusResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string                   `json:"resultType"`
		Result     []goloopPrometheusResult `json:"result"`
	} `json:"data"`
}

// goloopPrometheusResult ...
type goloopPrometheusResult struct {
	Metric struct {
		Name     string `json:"__name__"`
		Channel  string `json:"channel"`
		Hostname string `json:"hostname"`
		Instance string `json:"instance"`
		Job      string `json:"job"`
	} `json:"metric"`
	Values [][]interface{} `json:"values"`
}

var query []string = []string{
	"consensus_height",          // "BLOCKHEIGHT"
	"txpool_user_remove_sum",    // "COUNTOFTX"
	"txpool_user_drop_sum",      // "COUNTOFUNCONFIRMEDTX"
	"consensus_height_duration", // "RESPONSETIMEINSEC"
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

func (c GoloopPrometheus) Crawler() (*prometheus.PrometheusData, error) {

	ip := configuration.Conf().Prometheus.PrometheusISAAC

	queryPath := configuration.Conf().Prometheus.QueryPath

	queryList := addJobName(query)

	prometheusData, err := requestPrometheusQuery(ip+queryPath, queryList)
	if err != nil {
		logger.Error(isaacerror.SysErrFailToGetPrometheusDataFromGL.Error())
		return nil, err
	}

	prometheusResult := prometheusData.Data.Result

	var channelData []prometheus.PrometheusChannelData
	channelData = make([]prometheus.PrometheusChannelData, 0)

	// This value is prometheus crawling time
	//var crawlingIntervalInSec int
	var unixTime float64
	timeStamp := ""

	channelIDToName := db.GetChannelIDToNameMap()
	nodeTB := db.GetConfigurationNodeTable()

	// Convert prometheus result to our response data format.
	for _, value := range prometheusResult {
		// Get channel name by channel ID.
		channelName := channelIDToName["0x"+value.Metric.Channel]
		if channelName == "" {
			err := isaacerror.SysErrNoChannelInISAAC
			logger.Error(err.Error())
			return nil, err
		}

		channelIndex, isExistingChannel := prometheus.IsExistingChannel(channelName, channelData)
		if isExistingChannel == false {
			channelIndex = len(channelData)

			var data prometheus.PrometheusChannelData
			data.Name = channelName
			data.Status = ChannelNormal
			channelData = append(channelData, data)
		}

		// Get node name by node address.
		nodeName := ""
		for _, nodeData := range nodeTB {
			index := strings.Index(nodeData.NODE_ADDRESS, value.Metric.Hostname)
			if index == 2 {
				nodeName = nodeData.NODE_NAME
			}
		}
		if nodeName == "" {
			err := isaacerror.SysErrNoNodeInDB
			logger.Error(err.Error())
			return nil, err
		}

		nodeIndex, isExistingNode := prometheus.IsExistingNode(nodeName, channelData[channelIndex].Nodes)
		if isExistingNode == false {
			nodeIndex = len(channelData[channelIndex].Nodes)

			var data prometheus.PrometheusNodesData
			data.Name = nodeName
			data.Status = NodeNormal
			data.UnSyncBlockHoldInSec = 0
			channelData[channelIndex].Nodes = append(channelData[channelIndex].Nodes, data)
		} else {
			unSyncBlockHoldInSec, _ := prometheus.GetUnSyncBlockHoldInSec(channelData[channelIndex].Name,
				channelData[channelIndex].Nodes[nodeIndex].Name)
			channelData[channelIndex].Nodes[nodeIndex].UnSyncBlockHoldInSec = unSyncBlockHoldInSec
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
		channelData[channelIndex].Nodes[nodeIndex].TimeStamp = timeStamp

		// IsLeader is only 0 on goloop.
		// Todo: if seed then validate.
		channelData[channelIndex].Nodes[nodeIndex].IsLeader = 0

		// Parse metrics.
		switch value.Metric.Name {
		case queryList[0]:
			valueUInt64, err := strconv.ParseUint(value.Values[len(value.Values)-1][1].(string), 10, 64)
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToInt.Error())
				return nil, err
			}
			channelData[channelIndex].Nodes[nodeIndex].BlockHeight = valueUInt64

			valueUInt64Prev, err := strconv.ParseUint(value.Values[0][1].(string), 10, 64)
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToInt.Error())
				return nil, err
			}
			channelData[channelIndex].Nodes[nodeIndex].PrevBlockHeight = valueUInt64Prev
		case queryList[1]:
			valueUInt64, err := strconv.ParseUint(value.Values[len(value.Values)-1][1].(string), 10, 64)
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToInt.Error())
				return nil, err
			}
			channelData[channelIndex].Nodes[nodeIndex].CountOfTX = valueUInt64
		case queryList[2]:
			valueUInt64, err := strconv.ParseUint(value.Values[len(value.Values)-1][1].(string), 10, 64)
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToInt.Error())
				return nil, err
			}
			channelData[channelIndex].Nodes[nodeIndex].CountOfUnconfirmedTX = valueUInt64
		case queryList[3]:
			valueUInt64, err := strconv.ParseUint(value.Values[len(value.Values)-1][1].(string), 10, 64)
			if err != nil {
				logger.Error(isaacerror.SysErrFailConvertingStringToFloat.Error())
				return nil, err
			}

			// Convert msec to sec.
			valueFloat := float64(valueUInt64) / float64(1000)

			channelData[channelIndex].Nodes[nodeIndex].ResponseTimeInSec = valueFloat
		}
	}

	if len(channelData) == 0 {
		err := isaacerror.SysErrFailToGetPrometheusDataFromGL
		logger.Error(err.Error())
		return nil, err
	}

	for chIDX, channel := range channelData {
		_, lastBlockHeight := prometheus.IsLastBlockHeight(channel.Nodes)

		// Read configuration value.
		configurationAlertDataTB := db.GetAlertConfigInfoByName(channel.Name)

		// Check the status of all nodes in every channel.
		for ndIDX, node := range channel.Nodes {
			if node.BlockHeight == node.PrevBlockHeight ||
				lastBlockHeight - node.BlockHeight > constants.DBDefaultUnsyncBlockDifference {
				if channelData[chIDX].Nodes[ndIDX].UnSyncBlockHoldInSec == 0 {
					channelData[chIDX].Nodes[ndIDX].UnSyncBlockHoldInSec = CrawlingRangeTimeSec
				} else {
					channelData[chIDX].Nodes[ndIDX].UnSyncBlockHoldInSec +=
						configuration.Conf().Prometheus.CrawlingInterval + 1
				}
				if channelData[chIDX].Nodes[ndIDX].UnSyncBlockHoldInSec >
					configurationAlertDataTB.MAX_TIME_SEC_FOR_UNSYNC {
					channelPK := db.GetChannelPK(channel.Name)
					channelData[chIDX].Nodes[ndIDX].Status = NodeUnSyncedBlock
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
				channelData[chIDX].Nodes[ndIDX].UnSyncBlockHoldInSec = 0
			}

			if float64(configurationAlertDataTB.MAX_TIME_SEC_FOR_RESPONSE) < node.ResponseTimeInSec {
				if channelData[chIDX].Nodes[ndIDX].Status == NodeUnSyncedBlock {
					channelData[chIDX].Nodes[ndIDX].Status = NodeUnSyncBlockAndSlowResponse
				} else {
					channelData[chIDX].Nodes[ndIDX].Status = NodeSlowResponse
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

			if channelData[chIDX].Nodes[ndIDX].Status > 0 {
				channelData[chIDX].Status = ChannelAbnormal
			}
			logger.Debugf("[ch:%s][node:%s][status:%d] Succeed to crawling node data.",
				channelData[chIDX].Name,
				channelData[chIDX].Nodes[ndIDX].Name,
				channelData[chIDX].Nodes[ndIDX].Status)
		}
		logger.Debugf("[ch:%s][status:%d] Succeed to crawling channel data.",
			channelData[chIDX].Name,
			channelData[chIDX].Status)
	}
	var goloopData prometheus.PrometheusData
	goloopData.PrometheusChannelData = channelData
	goloopData.TimeStamp = unixTime
	goloopData.Status = prometheus.CrawlingSuccess
	return &goloopData, nil
}

func addJobName(queryList []string) []string {
	addQueryList := make([]string, len(queryList))

	jobName := configuration.Conf().Prometheus.JobNameOfgoloop

	for i, value := range queryList {
		addQueryList[i] = jobName + "_" + value
	}

	return addQueryList
}

// requestPrometheusQuery is function to query data from prometheus.
func requestPrometheusQuery(url string, queryList []string) (*prometheusResult, error) {
	var output prometheusResult

	queryString := "{__name__=~\""
	for i, value := range queryList {
		queryString += value
		if i != (len(queryList) - 1) {
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
