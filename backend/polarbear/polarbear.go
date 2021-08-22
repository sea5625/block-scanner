package polarbear

import (
	"encoding/json"
	"math/rand"
	"motherbear/backend/configuration"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"motherbear/backend/utility"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/jasonlvhit/gocron"
)

// BlockHeaderData is
type BlockHeaderData struct {
	JsonRpc string          `json:"jsonrpc" `
	Id      int64           `json:"id"`
	Result  BlockResultData `json:"result"`
}

// BlockResultData is
type BlockResultData struct {
	Version                  string   `json:"version" `
	PrevBlockHash            string   `json:"prev_block_hash" `
	MerkleTreeRootHash       string   `json:"merkle_tree_root_hash" `
	Timestamp                int64    `json:"time_stamp" `
	ConfirmedTransactionList []TxData `json:"confirmed_transaction_list" `
	BlockHash                string   `json:"block_hash" `
	Height                   int64    `json:"height" `
	PeerID                   string   `json:"peer_id" `
	Signature                string   `json:"signature" `
}

// TxData is
type TxData struct {
	Version   string `json:"version" `
	From      string `json:"from"`
	To        string `json:"to"`
	StepLimit string `json:"stepLimit"`
	Timestamp string `json:"timestamp" `
	NID       string `json:"nid"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
	Status    string `json:"status"`
	TxHash    string `json:"txHash"`
	DataType  string `json:"dataType"`
	Data      string `json:"data"`
}

type status int

const (
	crawling = iota
	readyToStarCrawling
)

var crawlingStatus status
var scheduler *gocron.Scheduler

// Init is initilizing polarbear module.
func Init(dbType string, dbOption ...string) error {

	// Init DB file.
	if InitDB(dbType, dbOption...) == nil {
		return isaacerror.SysErrFailToInitDBForPolarbear
	}

	crawlingStatus = readyToStarCrawling
	return nil
}

// BeginToCrawl crawl the data from the loopchain node.
func BeginToCrawl() *gocron.Scheduler {
	job := func() {
		cronJobForEveryChannel(configuration.Conf())
	}

	// Set interval and begin crawling.
	var interval uint64
	interval = uint64(configuration.Conf().Blockchain.CrawlingInterval)

	scheduler = gocron.NewScheduler()
	scheduler.Every(interval).Seconds().Do(job)
	scheduler.Start()

	return scheduler
}

func StopToCrawl() {
	logger.Info("Stop crawling block data from loopchain.")
	scheduler.Clear()
}

func unitCrawlAndStoreBlock(nodeIP string, channelName string, height int64) error {

	// Request block data from node.
	var blockData map[string]interface{}
	if err := getBlockByHeight(
		&blockData,
		nodeIP,
		channelName,
		height); err != nil {
		logger.Fatalln(err)
		return err
	}

	// Put log with block hash.
	if blockData["result"] == nil {
		logger.Errorf("Fail to get the block %d in %s", height, channelName)
		return isaacerror.SysErrFailToGetBlockData
	}

	result := blockData["result"].(map[string]interface{})
	logger.Infof("Begin to crawl %d block in %s. %s",
		height, channelName,
		utility.AddHexHD(result["block_hash"].(string)))

	// Add block data into DB.
	var blockRecord Block
	if err := AddBlockRecordFromJSONResponse(
		blockData,
		nodeIP,
		channelName,
		&blockRecord); err != nil {
		return err

	}

	// Put log with hash.
	logger.Infof("End to crawl %d block in %s. %s",
		height, channelName,
		utility.AddHexHD(result["block_hash"].(string)))
	return nil
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func crawlAndStoreBlock(numCPU int, nodeIPList []string, channelName string,
	beginBlockheight int64, endBlockHeight int64) error {
	logger.Infof("Crawl block from %d to %d in %s.", beginBlockheight, endBlockHeight, channelName)

	stepToCrawl := int64(2*numCPU + 1)
	logger.Infof("Step to crawl = %d", stepToCrawl)

	diffBlockHeight := endBlockHeight - beginBlockheight

	// Do linear crawling if diffBloekcHeight is small because more efficient.
	if diffBlockHeight < stepToCrawl {
		for h := beginBlockheight; h <= endBlockHeight; h++ {
			var  randomIndex int = 0
			if   len(nodeIPList) != 1{
				randomIndex = random(0, len(nodeIPList)-1)
			}
			err := unitCrawlAndStoreBlock(nodeIPList[randomIndex], channelName, h)
			if err != nil {
				return err
			}
		}

	} else {
		// Do pararell crawling.
		var countCrawledBlock int64 = 0
		remains := endBlockHeight % stepToCrawl

		for h := beginBlockheight; h <= endBlockHeight-remains; h += stepToCrawl {
			// Pararellized code block
			var wg sync.WaitGroup
			for i := int64(0); i < stepToCrawl; i++ {
				wg.Add(1)
				countCrawledBlock++
				go func(nodeIPList []string, channelName string, height int64) {
					defer wg.Done()
					randomIndex := random(0, len(nodeIPList)-1)
					err := unitCrawlAndStoreBlock(nodeIPList[randomIndex], channelName, height)
					if err != nil {
						logger.Errorf("%s", err)
					}
				}(nodeIPList, channelName, h+i)
			}
			wg.Wait()
		}

		// If there are some remain block height, do linear crawling.
		if remains > 0 {
			begin := endBlockHeight - remains + 1
			logger.Infof("Remains to crawl from %d to %d", begin, endBlockHeight)
			for h := begin; h <= endBlockHeight; h++ {
				randomIndex := random(0, len(nodeIPList)-1)
				err := unitCrawlAndStoreBlock(nodeIPList[randomIndex], channelName, h)
				if err != nil {
					return err
				}
			}
		}

	}

	return nil
}

func crawlBlockchain(channelName string) error {

	// 	// Get node list.
	nodeIPList := []string{}
	for _, c := range configuration.Conf().Channel {
		if c.Name == channelName {
			for _, n := range c.Nodes {
				nodeIPList = append(nodeIPList, configuration.QueryNodeByNodeName(n).IP)
			}
		}
	}

	// Check the block height of channel and start to crawl if it needs.
	blockHeightInDB := GetCurrentBlockHeightInDB(channelName)
	blockHeight, _ := getLastBlockHeight(nodeIPList[0], channelName)

	//	If block height in local is lower then block height online, then  start to crawl.
	numCPU := runtime.NumCPU()

	if blockHeightInDB < blockHeight {
		return crawlAndStoreBlock(numCPU, nodeIPList, channelName, blockHeightInDB+1, blockHeight)
	} else {
		logger.Debugf("Don't need to crawl in %s", channelName)
		return nil
	}
}

func cronJobForEveryChannel(conf *configuration.Configuration) {

	if crawlingStatus == readyToStarCrawling {
		crawlingStatus = crawling

		for _, c := range conf.Channel {
			err := crawlBlockchain(c.Name)
			if err != nil {
				panic(err)
			}
		}

		crawlingStatus = readyToStarCrawling
	} else {
		logger.Infof("Crawling process is running..")
	}

}

func generateRandString(n int) string {
	var letterRunes = []rune("abcdef0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateWalletID() string {
	return "hx" + generateRandString(40)
}

func generateBlockTxHash() string {
	return "0x" + generateRandString(68)
}

// GenerateTestBlockJSONData data generate pseudo block data for testing..
func GenerateTestBlockJSONData(countOfTx int, height int64, data *map[string]interface{}) ([]byte, error) {
	blockHeader := BlockHeaderData{
		JsonRpc: "2.0",
		Id:      3123,
		Result: BlockResultData{
			Version:                  "0.1a",
			PrevBlockHash:            generateRandString(40),
			MerkleTreeRootHash:       generateRandString(40),
			Timestamp:                time.Now().UnixNano() / 1000,
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
				Timestamp: strconv.FormatInt(time.Now().UnixNano()/1000, 10),
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
