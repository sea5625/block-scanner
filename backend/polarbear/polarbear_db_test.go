package polarbear

import (
	"encoding/json"
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Setup(path string) *gorm.DB {
	//database := InitDB("mysql", "isaac_index2", "helloworld123", "tcp(127.0.0.1:3306)/index2")
	database := InitDB("sqlite3", path)
	return database
}

func Teardown(path string) {
	_ = os.Remove(path)
	defer Database().Close()
}

func DBSetup(channelName string, count int) {
	for i := 0; i < count; i++ {
		var tempBlock map[string]interface{}
		_, err := buildTestBlockJSONData(100, int64(i), &tempBlock)
		if err != nil {
			fmt.Println("error buildTestBlockJSONData")
		}
		var blockRecord Block
		if err := AddBlockRecordFromJSONResponse(
			tempBlock,
			"",
			channelName,
			&blockRecord); err != nil {
			fmt.Println("error AddBlockRecordFromJSONResponse")
		}
	}
}

func buildTestBlockJSONData(countOfTx int, height int64, data *map[string]interface{}) ([]byte, error) {
	blockHeader := BlockHeaderData{
		JsonRpc: "2.0",
		Id:      3123,
		Result: BlockResultData{
			Version:                  "0.1a",
			PrevBlockHash:            generateRandString(40),
			MerkleTreeRootHash:       generateRandString(40),
			Timestamp:                time.Now().Unix(),
			BlockHash:                generateBlockTxHash(),
			Height:                   height,
			PeerID:                   generateWalletID(),
			Signature:                generateRandString(80),
			ConfirmedTransactionList: []TxData{},
		},
	}

	// Add the pseudo TX data.
	for i := 0; i < countOfTx; i++ {
		blockHeader.Result.ConfirmedTransactionList = append(
			blockHeader.Result.ConfirmedTransactionList,
			TxData{
				Version:   "0x3",
				From:      generateWalletID(),
				To:        generateWalletID(),
				StepLimit: "0x12345",
				NID:       "0x5",
				Nonce:     generateRandString(52),
				Signature: generateRandString(32),
				TxHash:    generateBlockTxHash(),
				Timestamp: strconv.FormatInt(time.Now().Unix(), 10),
				DataType:  "call",
				Data: string([]byte(`
							{
							"method": "transfer",
							"params": {
								"to": "hxab2d8215eab14bc6bdd8bfb2c8151257032ecd8b",
								"value": "0x1"
							}
						}`),
				),
			},
		)
	}

	byteData, err := json.Marshal(blockHeader)

	// Try to unmarshal JSON data.
	if err := json.Unmarshal(byteData, &data); err != nil {
		return []byte(``), err
	}

	return byteData, err
}

func TestPolarbearTestBlock(t *testing.T) {

	// Basic test configuration.
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	// Create new pseudo block data.
	var data map[string]interface{}
	_, err := buildTestBlockJSONData(3, 0, &data)
	if err != nil {
		t.Error("Fail to build block JSON data.")
	}

	// Try to build Block table data from JSON data.
	var blockRecord Block
	if err := AddBlockRecordFromJSONResponse(
		data,
		"",
		"myChannel",
		&blockRecord); err != nil {
		t.Error("Fail to convert JSON to Block table.", err)
	}

	// Test data.
	if len(blockRecord.Txs) != 3 {
		t.Error("Fail to get TX data.")
	}

	// Parse data in Tx.
	var dataInTx map[string]interface{}
	if err := json.Unmarshal(
		[]byte(blockRecord.Txs[0].Data),
		&dataInTx); err != nil {
		t.Error(err)
	}

	if dataInTx["method"].(string) != "transfer" {
		t.Error("Fail to parse data in Tx.")
		t.Error(dataInTx)
	}

}

// Test cases for AddBlock().
func TestAddBlock(t *testing.T) {

	// Basic test configuration.
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	// Generate data for Ch0 channel.
	channelName := "Ch0"
	for i := 0; i < 6; i++ {
		var tempBlock map[string]interface{}
		_, err := buildTestBlockJSONData(10, int64(i), &tempBlock)
		if err != nil {
			t.Error("Fail to generate block data. :", i)
		}
		var blockRecord Block
		if err := AddBlockRecordFromJSONResponse(
			tempBlock,
			"",
			channelName,
			&blockRecord); err != nil {
			t.Error("Fail to add tx record.", err)
		}
	}

	// Generate data for Ch1 channel.
	channelName = "Ch1"
	for i := 0; i < 100; i++ {
		var tempBlock map[string]interface{}
		_, err := buildTestBlockJSONData(10, int64(i), &tempBlock)
		if err != nil {
			t.Error("Fail to generate block data. :", i)
		}
		var blockRecord Block
		if err := AddBlockRecordFromJSONResponse(
			tempBlock,
			"",
			channelName,
			&blockRecord); err != nil {
			t.Error("Fail to add tx record.", err)
		}
	}

	// Query block of Ch0 channel.
	var blocksInCh0 []Block
	count, err := QueryBlocksInChannel("Ch0", 10, 0, &blocksInCh0)
	if err != nil || count > int64(len(blocksInCh0)) {
		t.Error("Fail to query blocks. ", err)
	}

	// Verify blocks.
	if len(blocksInCh0) != 6 {
		t.Error("Fail to query blocks in Ch0")
	}

	if len(blocksInCh0[4].Txs) == 0 {
		t.Error("Fail to query Txs of blocks in Ch0")
	} else {
		if blocksInCh0[0].BlockHeight < blocksInCh0[1].BlockHeight {
			t.Error("Fail to sort by block height in descending. ")
		}
	}

	// Query block of Ch1 channel.
	var blocksInCh1 []Block
	count, err = QueryBlocksInChannel("Ch1", 20, 3, &blocksInCh1)
	if err != nil || count < int64(len(blocksInCh1)) {
		t.Error("Fail to query blocks. ", err)
	}
	// Verify blocks.
	if len(blocksInCh1) == 0 {
		t.Error("Fail to query blocks in Ch1")
	}

	if len(blocksInCh1[9].Txs) == 0 {
		t.Error("Fail to query Txs of blocks in Ch1")
	} else {
		if blocksInCh1[0].BlockHeight < blocksInCh1[1].BlockHeight {
			t.Error("Fail to sort by block height in descending. ")
		}
	}

	//  Get the height of crawled blockchain .
	currentHeightOfCH1 := GetCurrentBlockHeightInDB("Ch1")
	if currentHeightOfCH1 < 0 {
		t.Error("Fail to get the block height of Ch1. ", currentHeightOfCH1)
	}

	// Query block of bad channel name.
	var blocksInBadCh []Block
	count, err = QueryBlocksInChannel("", 10, 0, &blocksInBadCh)
	if err != nil && count > 0 {
		t.Error("Fail to get the blocks in wrong channel. ", err)
	}

}

func TestPolarbearTxDataParse(t *testing.T) {

	// Basic test configuration.
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	// Block data
	byteBlock := []byte(`
		{
			"jsonrpc": "2.0",
			"id": 1234,
			"result": {
				"version": "0.1a", 
				"prev_block_hash": "48757af881f76c858890fb41934bee228ad50a71707154a482826c39b8560d4b",
				"merkle_tree_root_hash": "fabc1884932cf52f657475b6d62adcbce5661754ff1a9d50f13f0c49c7d48c0c",
				"time_stamp": 1516498781094429,
				"confirmed_transaction_list": [ 
					{
						"version": "0x3",
						"from": "hxbe258ceb872e08851f1f59694dac2558708ece11",
						"to": "cxb0776ee37f5b45bfaea8cff1d8232fbb6122ec32",
						"value": "0xde0b6b3a7640000",
						"stepLimit": "0x12345",
						"timestamp": "0x563a6cf330136",
						"nid": "0x3",
						"nonce": "0x1",
						"signature": "VAia7YZ2Ji6igKWzjR2YsGa2m53nKPrfK7uXYW78QLE+ATehAVZPC40szvAiA6NEU5gCYB4c4qaQzqDh2ugcHgA=",
						"txHash": "0xb903239f8543d04b5dc1ba6579132b143087c68db1b2168786408fcbce568238",
						"dataType": "call",
						"data": {
							"method": "transfer",
							"params": {
								"to": "hxab2d8215eab14bc6bdd8bfb2c8151257032ecd8b",
								"value": "0x1"
							}
						}
					}
				],
				"block_hash": "1fcf7c34dc875681761bdaa5d75d770e78e8166b5c4f06c226c53300cbe85f57",
				"height": 3,
				"peer_id": "e07212ee-fe4b-11e7-8c7b-acbc32865d5f",
				"signature": "MEQCICT8mTIL6pRwMWsJjSBHcl4QYiSgG8+0H3U32+05mO9HAiBOhIfBdHNm71WpAZYwJWwQbPVVXFJ8clXGKT3ScDWcvw=="
			}
		}
`)
	var data map[string]interface{}

	if err := json.Unmarshal(byteBlock, &data); err != nil {
		panic(err)
	}
}

func TestGetPeerSymptom(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	fromTimer, _ := time.Parse(time.RFC3339, "2019-08-13T10:41:05+09:00")

	peerSymptomMappingTB := &Symptom{
		Channel:     "loopchain_default",
		Channel_PK:  "PKCH_000000000000101010",
		Msg:         "[node4]response time slowly [25.106827] sec",
		SymptomType: "Slow response",
		Timestamp:   fromTimer,
	}

	if Database().NewRecord(&peerSymptomMappingTB) {
		Database().Save(&peerSymptomMappingTB)
	}

	var peer_symptom []Symptom
	var channel []string

	channel = append(channel, "PKCH_000000000000101010")

	count, _ := QueryPeerSymptomListTable(10, 0, "2019-08-13T10:41:04+09:00", "2019-08-13T10:51:04+09:00", channel, &peer_symptom)

	assert.Equal(t, count, int64(1))
	assert.Equal(t, peer_symptom[0].Channel, "loopchain_default")
}

func TestQueryBlock(t *testing.T) {
	t.Skip("Skipping test QueryBlocksInChannel. Should prepare polarbear data(DB) to test.")

	// Basic test configuration.
	dbpath := "performanceTest.db"
	Setup(dbpath)
	channelName := "1"

	// Query block of Ch0 channel.
	var blocksInCh0 []Block
	count, err := QueryBlocksInChannel(channelName, 10, 0, &blocksInCh0)
	if err != nil {
		t.Error("Fail to query blocks. ", err)
	}

	t.Log("block count =", count, len(blocksInCh0))
}

func TestQueryTx(t *testing.T) {
	t.Skip("Skipping test QueryTxsInChannelBySearch. Should prepare polarbear data(DB) to test.")

	// Basic test configuration.
	dbpath := "../../data/crawling_data.db"
	Setup(dbpath)
	channelName := "1"

	// Query tx of Ch0 channel.
	txSearch := TxSearch{
		Status:      "",
		BlockHeight: -1,
		From:        time.Time{},
		To:          time.Time{},
		FromAddress: "",
		ToAddress:   "",
		Data:        "",
	}

	var txInCh []Tx
	count, err := QueryTxsInChannelBySearch(channelName, 10, 0, txSearch, &txInCh)
	if err != nil {
		t.Error("Fail to query tx. ", err)
	}

	t.Log("tx count =", count, len(txInCh))
}

func TestPerformanceAboutQueryBlockByHeight(t *testing.T) {
	t.Skip("Skipping test performance about QueryBlockByHeightInChannel. Should prepare polarbear data(DB) to test.")

	// Basic test configuration.
	dbpath := "performanceTest.db"
	Setup(dbpath)
	channelName := "1"

	// Query block of Ch0 channel.
	var blocksInCh0 []Block
	count, err := QueryBlocksInChannel(channelName, 10, 0, &blocksInCh0)
	if err != nil {
		t.Error("Fail to query blocks. ", err)
	}

	t.Log("count =", count, len(blocksInCh0))

	loopCount := 100

	if count != 0 {
		var blocksByHeight []Block
		var elapsedTimeSum time.Duration
		for i := 0; i < loopCount; i++ {
			randomNum := rand.Intn(int(count))

			_, err := QueryBlocksInChannel(channelName, 10, randomNum, &blocksInCh0)
			if err != nil {
				t.Error("Fail to query blocks. ", err)
			}

			startTime := time.Now()
			err = QueryBlockByHeightInChannel(channelName, blocksInCh0[0].BlockHeight, &blocksByHeight)
			if err != nil {
				t.Error("Fail to query blocks. ", err)
			}
			elapsedTime := time.Since(startTime)

			elapsedTimeSum += elapsedTime
		}
		elapsedTimeAvr := elapsedTimeSum / time.Duration(loopCount)

		t.Log("add elapsed time about query block by height (count:", loopCount, ") :", elapsedTimeSum)
		t.Log("average elapsed time about query block by height (count:", loopCount, ") :", elapsedTimeAvr)
	}
}

func TestPerformanceAboutQueryBlockByHash(t *testing.T) {
	t.Skip("Skipping test performance about QueryBlockInChannelByHash. Should prepare polarbear data(DB) to test.")

	// Basic test configuration.
	dbpath := "performanceTest.db"
	Setup(dbpath)
	channelName := "1"

	// Query block of Ch0 channel.
	var blocksInCh0 []Block
	count, err := QueryBlocksInChannel(channelName, 10, 0, &blocksInCh0)
	if err != nil {
		t.Error("Fail to query blocks. ", err)
	}

	t.Log("count =", count, len(blocksInCh0))

	loopCount := 100

	if count != 0 {
		var blocksByHash []Block
		var elapsedTimeSum time.Duration
		for i := 0; i < loopCount; i++ {
			randomNum := rand.Intn(int(count))

			_, err := QueryBlocksInChannel(channelName, 10, randomNum, &blocksInCh0)
			if err != nil {
				t.Error("Fail to query blocks. ", err)
			}

			startTime := time.Now()
			err = QueryBlockInChannelByHash(channelName, blocksInCh0[0].BlockHash, &blocksByHash)
			if err != nil {
				t.Error("Fail to query blocks. ", err)
			}
			elapsedTime := time.Since(startTime)

			elapsedTimeSum += elapsedTime
		}
		elapsedTimeAvr := elapsedTimeSum / time.Duration(loopCount)

		t.Log("add elapsed time about query block by hash (count:", loopCount, ") :", elapsedTimeSum)
		t.Log("average elapsed time about query block by hash (count:", loopCount, ") :", elapsedTimeAvr)
	}
}

func TestPerformanceAboutQueryTxByHash(t *testing.T) {
	t.Skip("Skipping test performance about QueryTxInChannelByHash. Should prepare polarbear data(DB) to test.")

	// Basic test configuration.
	dbpath := "performanceTest.db"
	Setup(dbpath)
	channelName := "1"

	// Query tx of Ch0 channel.
	txSearch := TxSearch{
		Status:      "",
		BlockHeight: -1,
		From:        time.Time{},
		To:          time.Time{},
		FromAddress: "",
		ToAddress:   "",
	}

	var txInCh []Tx
	count, err := QueryTxsInChannelBySearch(channelName, 10, 0, txSearch, &txInCh)
	if err != nil {
		t.Error("Fail to query tx. ", err)
	}

	t.Log("count =", count, len(txInCh))

	loopCount := 100

	if count != 0 {
		var txByHash []Tx
		var elapsedTimeSum time.Duration
		for i := 0; i < loopCount; i++ {
			randomNum := rand.Intn(int(count))

			_, err := QueryTxsInChannelBySearch(channelName, 10, randomNum, txSearch, &txInCh)
			if err != nil {
				t.Error("Fail to query tx. ", err)
			}

			startTime := time.Now()
			err = QueryTxInChannelByHash(channelName, txInCh[0].TxHash, &txByHash)
			if err != nil {
				t.Error("Fail to query tx. ", err)
			}
			elapsedTime := time.Since(startTime)

			elapsedTimeSum += elapsedTime
		}
		elapsedTimeAvr := elapsedTimeSum / time.Duration(loopCount)

		t.Log("add elapsed time about query tx by hash (count:", loopCount, ") :", elapsedTimeSum)
		t.Log("average elapsed time about query tx by hash (count:", loopCount, ") :", elapsedTimeAvr)
	}
}

func TestPerformanceAboutQueryTxByBlockHeight(t *testing.T) {
	t.Skip("Skipping test performance about QueryTxsInChannelByBlockHeight. Should prepare polarbear data(DB) to test.")

	// Basic test configuration.
	dbpath := "performanceTest.db"
	Setup(dbpath)
	channelName := "1"

	// Query tx of Ch0 channel.
	txSearch := TxSearch{
		Status:      "",
		BlockHeight: -1,
		From:        time.Time{},
		To:          time.Time{},
		FromAddress: "",
		ToAddress:   "",
		Data:        "",
	}

	var txInCh []Tx
	count, err := QueryTxsInChannelBySearch(channelName, 10, 0, txSearch, &txInCh)
	if err != nil {
		t.Error("Fail to query tx. ", err)
	}

	t.Log("count =", count, len(txInCh))

	loopCount := 100

	if count != 0 {
		var txByHash []Tx
		var elapsedTimeSum time.Duration
		for i := 0; i < loopCount; i++ {
			randomNum := rand.Intn(int(count))

			_, err := QueryTxsInChannelBySearch(channelName, 10, randomNum, txSearch, &txInCh)
			if err != nil {
				t.Error("Fail to query tx. ", err)
			}

			txSearch = TxSearch{
				Status:      "",
				BlockHeight: txInCh[0].BlockHeight,
				From:        time.Time{},
				To:          time.Time{},
				FromAddress: "",
				ToAddress:   "",
				Data:        "",
			}

			startTime := time.Now()
			count, err = QueryTxsInChannelBySearch(channelName, 10, 0, txSearch, &txByHash)
			if err != nil {
				t.Error("Fail to query tx. ", err)
			}
			elapsedTime := time.Since(startTime)

			elapsedTimeSum += elapsedTime
		}
		elapsedTimeAvr := elapsedTimeSum / time.Duration(loopCount)

		t.Log("add elapsed time about query tx by height (count:", loopCount, ") :", elapsedTimeSum)
		t.Log("average elapsed time about query tx by height (count:", loopCount, ") :", elapsedTimeAvr)
	}
}

func TestPerformanceAboutQueryTxByStatus(t *testing.T) {
	t.Skip("Skipping test performance about QueryTxsInChannelByBlockStatus. Should prepare polarbear data(DB) to test.")

	// Basic test configuration.
	dbpath := "performanceTest.db"
	Setup(dbpath)
	channelName := "1"

	// Query tx of Ch0 channel.
	txSearch := TxSearch{
		Status:      "",
		BlockHeight: -1,
		From:        time.Time{},
		To:          time.Time{},
		FromAddress: "",
		ToAddress:   "",
		Data:        "",
	}

	var txInCh []Tx
	count, err := QueryTxsInChannelBySearch(channelName, 10, 0, txSearch, &txInCh)
	if err != nil {
		t.Error("Fail to query tx. ", err)
	}

	t.Log("count =", count, len(txInCh))

	loopCount := 100

	if count != 0 {
		var txByHash []Tx
		var elapsedTimeSum time.Duration
		for i := 0; i < loopCount; i++ {
			txSearch = TxSearch{
				Status:      "Success",
				BlockHeight: -1,
				From:        time.Time{},
				To:          time.Time{},
				FromAddress: "",
				ToAddress:   "",
				Data:        "",
			}

			startTime := time.Now()
			count, err = QueryTxsInChannelBySearch(channelName, 10, 0, txSearch, &txByHash)
			if err != nil {
				t.Error("Fail to query tx. ", err)
			}
			elapsedTime := time.Since(startTime)

			elapsedTimeSum += elapsedTime
		}
		elapsedTimeAvr := elapsedTimeSum / time.Duration(loopCount)

		t.Log("add elapsed time about query tx by status (count:", loopCount, ") :", elapsedTimeSum)
		t.Log("average elapsed time about query tx by status (count:", loopCount, ") :", elapsedTimeAvr)
	}
}

func TestPerformanceAboutQueryTxByDate(t *testing.T) {
	t.Skip("Skipping test performance about QueryTxsInChannelByDate. Should prepare polarbear data(DB) to test.")

	// Basic test configuration.
	dbpath := "performanceTest.db"
	Setup(dbpath)
	channelName := "1"

	// Query tx of Ch0 channel.
	txSearch := TxSearch{
		Status:      "",
		BlockHeight: -1,
		From:        time.Time{},
		To:          time.Time{},
		FromAddress: "",
		ToAddress:   "",
		Data:        "",
	}

	var txInCh []Tx
	count, err := QueryTxsInChannelBySearch(channelName, 10, 0, txSearch, &txInCh)
	if err != nil {
		t.Error("Fail to query tx. ", err)
	}

	t.Log("count =", count, len(txInCh))

	loopCount := 100

	if count != 0 {
		var txByHash []Tx
		var elapsedTimeSum time.Duration
		for i := 0; i < loopCount; i++ {
			txSearch = TxSearch{
				Status:      "",
				BlockHeight: -1,
				From:        time.Now().UTC().Add(-time.Minute * 30),
				To:          time.Now().UTC(),
				FromAddress: "",
				ToAddress:   "",
				Data:        "",
			}

			startTime := time.Now()
			count, err = QueryTxsInChannelBySearch(channelName, 10, 0, txSearch, &txByHash)
			if err != nil {
				t.Error("Fail to query tx. ", err)
			}
			elapsedTime := time.Since(startTime)

			elapsedTimeSum += elapsedTime
		}
		elapsedTimeAvr := elapsedTimeSum / time.Duration(loopCount)

		t.Log("add elapsed time about query tx by date (count:", loopCount, ") :", elapsedTimeSum)
		t.Log("average elapsed time about query tx by date (count:", loopCount, ") :", elapsedTimeAvr)
	}
}

func TestPerformanceAboutQuery(t *testing.T) {
	t.Skip("Skipping test performance about PerformanceAboutQuery. Should prepare polarbear data(DB) to test.")

	TestPerformanceAboutQueryBlockByHash(t)
	TestPerformanceAboutQueryBlockByHeight(t)
	TestPerformanceAboutQueryTxByHash(t)
	TestPerformanceAboutQueryTxByBlockHeight(t)
	TestPerformanceAboutQueryTxByStatus(t)
	TestPerformanceAboutQueryTxByDate(t)
}
