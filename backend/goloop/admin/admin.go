package admin

import (
	"encoding/json"
	"io/ioutil"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"net/http"
)

type System struct {
	BuildVersion string        `json:"buildVersion"`
	BuildTag     string        `json:"buildTag"`
	Setting      SystemSetting `json:"setting"`
	Config       SystemConfig  `json:"config"`
}

type SystemSetting struct {
	Address   string `json:"address"`
	P2p       string `json:"p2p"`
	P2pListen string `json:"p2pListen"`
	RpcAddr   string `json:"rpcAddr"`
	RpcDump   bool   `json:"rpcDump"`
}

type SystemConfig struct {
	EeInstances       int    `json:"eeInstances"`
	RpcDefaultChannel string `json:"rpcDefaultChannel"`
	RpcIncludeDebug   bool   `json:"rpcIncludeDebug"`
}

type Chain struct {
	Nid       string `json:"nid"`
	Channel   string `json:"channel"`
	State     string `json:"state"`
	Height    int    `json:"height"`
	LastError string `json:"lastError"`
}

const GOLOOP_ADMIN_API_PATH = "/admin"
const GOLOOP_ADMIN_CHAIN_API_PATH = GOLOOP_ADMIN_API_PATH + "/chain"
const GOLOOP_ADMIN_SYSTEM_API_PATH = GOLOOP_ADMIN_API_PATH + "/system"

func GetChains(nodeIP string) ([]Chain, error) {

	var chains []Chain

	// Do request '/chain' API to node.
	uri := nodeIP + GOLOOP_ADMIN_CHAIN_API_PATH

	res, err := http.Get(uri)
	if err != nil {
		logger.Error(isaacerror.SysErrFailToConnectionNodeOfGoloop.Error())
		return nil, err
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(isaacerror.SysErrFailToReadBodyNodeOfGoloop.Error())
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		logger.Error(isaacerror.SysErrFailToReadBodyNodeOfGoloop.Error())
		return nil, err
	}

	err = json.Unmarshal(resData, &chains)
	if err != nil {
		logger.Error(isaacerror.SysErrFailToUnmarshalChainsDataOfGoloop.Error())
		return nil, err
	}

	return chains, nil
}

func GetChainsAllTry(nodeIPList []string) ([]Chain, error) {
	chains := make([]Chain, 0)

	for _, value := range nodeIPList {
		var err error
		chains, err = GetChains(value)
		if err == nil && len(chains) > 0 {
			break
		} else if err != nil {
			logger.Debug(err.Error())
		} else {
			logger.Debug(isaacerror.SysErrFailToGetChainsDataOfGoloop.Error())
		}
	}

	if len(chains) <= 0 {
		err := isaacerror.SysErrFailToGetChainsDataOfGoloop
		logger.Error(err.Error())
		return nil, err
	}

	return chains, nil
}

func GetChannelIDByName(channelName string, nodeIPList []string) (string, error) {

	chains, err := GetChainsAllTry(nodeIPList)
	if err != nil {
		return "", err
	}

	idOfGoloop := ""

	for _, value := range chains {
		if value.Channel == channelName {
			idOfGoloop = value.Nid
			break
		}
	}

	if idOfGoloop == "" {
		err := isaacerror.SysErrNotExistChannelNameInTheGoloop
		logger.Error(err.Error())
		return "", err
	}

	return idOfGoloop, nil
}

func GetSystem(nodeIP string) (System, error) {
	var system System

	// Do request '/system' API to node.
	uri := nodeIP + GOLOOP_ADMIN_SYSTEM_API_PATH

	res, err := http.Get(uri)
	if err != nil {
		logger.Error(isaacerror.SysErrFailToConnectionNodeOfGoloop.Error())
		return system, err
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(isaacerror.SysErrFailToReadBodyNodeOfGoloop.Error())
		return system, err
	}

	err = res.Body.Close()
	if err != nil {
		logger.Error(isaacerror.SysErrFailToReadBodyNodeOfGoloop.Error())
		return system, err
	}

	err = json.Unmarshal(resData, &system)
	if err != nil {
		logger.Error(isaacerror.SysErrFailToUnmarshalSystemsDataOfGoloop.Error())
		return system, err
	}

	return system, nil
}
