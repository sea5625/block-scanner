package polarbear

import (
	"encoding/json"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"strconv"
	"time"

	"bytes"

	"github.com/ybbus/jsonrpc"
)

func generalJSONRPCReq(out interface{}, URI string,
	channelName string, method string, params ...interface{}) error {
	apiURI := URI + "/api/v3"

	if channelName != "" {
		if channelName != "default" {
			apiURI = apiURI + "/" + channelName
		}
	}

	rpcClient := jsonrpc.NewClient(apiURI)

	var err error
	if len(params) == 0 {
		err = rpcClient.CallFor(out, method)
		return err

	} else {
		res, err := rpcClient.Call(method, params)
		out = res
		return err
	}

}

func getLastBlockHeight(URI string, channelName string) (int64, error) {
	var body map[string]interface{}
	err := generalJSONRPCReq(&body, URI, channelName, "icx_getLastBlock")
	if err == nil {
		height := int64(body["height"].(float64))
		return height, err
	} else {
		return -1, err
	}
}


const countOfTrial = 10
const durationInMiliSec = 500

func getTxStatus(URI string, channelName string, txHash string) (string, error) {

	var txStatus string
	var err error
	err = isaacerror.SysUnknownError
	for i := 0; i< countOfTrial; i++ {
		txStatus, err =_getTxStatus(URI, channelName, txHash)
		if err == nil {
			return  txStatus, err
		}else{
			time.Sleep( durationInMiliSec * time.Microsecond)
		}
	}

	return "", err
}


func _getTxStatus(URI string, channelName string, txHash string) (string, error) {

	apiURI := URI + "/api/v3"
	if channelName != "" {
		if channelName != "default" {
			apiURI = apiURI + "/" + channelName
		}
	}

	// Try to call JSON RPC.
	rpcClient := jsonrpc.NewClient(apiURI)
	res, err := rpcClient.Call(
		"icx_getTransactionResult",
		map[string]interface{}{
			"txHash": txHash,
		})

	// If there are no response, try again.
	if err != nil {
		res, err = rpcClient.Call(
			"icx_getTransactionResult",
			map[string]interface{}{
				"txHash": txHash,
			})

		if err != nil {
			logger.Fatalln("Fail  to request tx status : ", txHash)
			return "", err
		}
	}

	// Convert res => byte
	resBytes := new(bytes.Buffer)
	err = json.NewEncoder(resBytes).Encode(res)
	if err != nil {
		return "", err
	}

	// Convert byte to map[string]interface{}
	var f interface{}
	err = json.Unmarshal(resBytes.Bytes(), &f)
	if err != nil {
		return "", err
	}

	// Convert byte to map[string]interface{}
	response := f.(map[string]interface{})
	result := response["result"].(map[string]interface{})

	// Parse status information.
	var txStatus string
	if result["status"].(string) == "0x1" {
		txStatus = "Success"
	} else {
		txStatus = "Failure"
	}

	return txStatus, err
}



func getBlockByHeight(out *map[string]interface{}, URI string,
	channelName string, height int64) error {

	HexHeightInString := "0x" +
		strconv.FormatInt(height, 16)

	apiURI := URI + "/api/v3"

	if channelName != "" {
		if channelName != "default" {
			apiURI = apiURI + "/" + channelName
		}
	}

	rpcClient := jsonrpc.NewClient(apiURI)
	res, err := rpcClient.Call(
		"icx_getBlockByHeight",
		map[string]interface{}{
			"height": HexHeightInString,
		})

	if err != nil {
		logger.Fatalln("Fail  to request block by height : ", height)
		return err
	}

	// Convert res => byte
	resBytes := new(bytes.Buffer)
	err = json.NewEncoder(resBytes).Encode(res)
	if err != nil {
		logger.Fatalln("Fail  to encode response.  ", res)
		return err
	}

	// Convert byte to map[string]interface{}
	var f interface{}
	err = json.Unmarshal(resBytes.Bytes(), &f)
	if err != nil {
		logger.Fatalln("Fail  to decode byte to JSON data.  ", resBytes.Bytes())
		return err
	}

	*out = f.(map[string]interface{})
	return err
}

func getTxByHash(channelName string, URI string, txHash string) error {
	return nil
}
