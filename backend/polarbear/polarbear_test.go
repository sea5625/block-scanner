package polarbear

import (
	conf "motherbear/backend/configuration"
	"testing"
)

func TestJSONRPCGetBlockHeight(t *testing.T) {

	endPoint := "https://test-ctz.solidwallet.io"

	height, err := getLastBlockHeight(endPoint, "")

	if err != nil || height == 0 {
		t.Error("Fail to get the height of blockchain. ", err)
	}
}

func TestGetTXStatus(t *testing.T) {
    t.Skip("Skipping test Check TX Status.")
	endPoint := "https://test-ctz.solidwallet.io"
	txHash := "0x32f32c4021d8c35e76bc55e8ccde917c9a9626dc835113f2f10b21fa43bd91ac"

	status, err := getTxStatus(endPoint, "", txHash)

	if err != nil && status != "Success" {
		t.Error("Fail to get Tx status.", status, err)
	}
}

func TestJSONRPCGetBlockByHeight(t *testing.T) {

	endPoint := "https://test-ctz.solidwallet.io"

	blockHeightToGet := int64(40020)
	var blockData map[string]interface{}
	err := getBlockByHeight(&blockData, endPoint, "", blockHeightToGet)

	if err != nil {
		t.Error("Fail to get the block by height. ", err)

	} else {

		var blockRecord Block
		buildBlockRecordFromJSON(blockData, endPoint, "", &blockRecord)

		if blockRecord.BlockHeight != int64(blockHeightToGet) {
			t.Error("Wrong block height data. ",
				blockRecord.BlockHeight, blockHeightToGet)
		}

	}

}

func TestCrawlBlockchain(t *testing.T) {

	// Read the configuration file.
	conf.InitConfigData("test_conf.yaml")

	// Init DB for polarbear
	err := Init("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
	}

	// Run crawling.
	cronJobForEveryChannel(conf.Conf())

	// Check the crawled the data.
	channelName := "loopchain_default1"
	height := GetCurrentBlockHeightInDB(channelName)
	if height < 0 {
		t.Error("Fail to crawl the block data.")
	}

	// Try to query block data.
	for h := int64(1); h <= height; h++ {

		var block Block
		err := QueryBlockByHeightInChannel(channelName, h, &block)
		if h != block.BlockHeight {
			t.Error("Fail to query block by height.")
		}

		if err != nil {
			t.Error(err)
		}

		var blockDub Block
		_ = QueryBlockInChannelByHash(channelName, block.BlockHash, &blockDub)
		if blockDub.BlockHash != block.BlockHash {
			t.Error("Fail to query block by block hash.")
		}

		if len(block.Txs) != 0 {
			txHashToQuery := block.Txs[0].TxHash

			var queriedTx Tx
			error := QueryTxInChannelByHash(channelName, txHashToQuery, &queriedTx)

			if error != nil {
				t.Error("Internal error!")
			}

			if queriedTx.TxHash != txHashToQuery {
				t.Error("Fail to query Tx.")
			}

		}

	}

}

func TestCrawlTestNetBlockchain(t *testing.T) {
	// go test allow to finish the unit test under 6 seconds.
	t.Skip("It takes very long time. We'll skip this test in general unit test. ")

	// Read the configuration file.
	conf.InitConfigData("testnet_conf.yaml")

	// Init DB for polarbear
	err := Init("sqlite3", ":memory:")
	if err != nil {
		t.Error(err)
	}

	// Run crawling.
	cronJobForEveryChannel(conf.Conf())
}
