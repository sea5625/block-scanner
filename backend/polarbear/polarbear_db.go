package polarbear

import (
	"encoding/json"
	"motherbear/backend/constants"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	util "motherbear/backend/utility"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Block is crawled block data.
type Block struct {
	gorm.Model
	Channel     string `gorm:"type:VARCHAR(64);not null;index"`
	BlockHeight int64  `gorm:"type:BIGINT;not null;index"`
	PeerID      string `gorm:"type:VARCHAR(512);not null"`
	Signature   string `gorm:"type:VARCHAR(512);not null"`
	Timestamp   time.Time
	BlockHash   string `gorm:"type:VARCHAR(512);not null;index"`
	Txs         []Tx   `gorm:"many2many:block_tx;"`
}

// Tx is transaction data.
type Tx struct {
	gorm.Model
	TxHash      string    `gorm:"type:VARCHAR(512);not null;index"`
	Status      string    `gorm:"type:VARCHAR(20);not null;index"`
	Channel     string    `gorm:"type:VARCHAR(64);not null;index"`
	BlockHeight int64     `gorm:"type:BIGINT;not null;index"`
	Timestamp   time.Time `gorm:"index"`
	From        string    `gorm:"type:VARCHAR(128);not null"`
	To          string    `gorm:"type:VARCHAR(128);not null"`
	Data        string    `gorm:"type:MEDIUMTEXT;not null"`
}

// Symptom is peer symptom data.
type Symptom struct {
	gorm.Model
	Channel     string `gorm:"type:varchar(40)"`
	Channel_PK  string `gorm:"type:varchar(40)"`
	Msg         string `gorm:"type:varchar(512)"`
	SymptomType string `gorm:"type:varchar(32)"` //Slow response, Unsync block
	Timestamp   time.Time
}

type TxSearch struct {
	Status      string
	BlockHeight int64
	From        time.Time
	To          time.Time
	FromAddress string
	ToAddress   string
	Data        string
}

var instance *gorm.DB
var once sync.Once

// Database returns the instance.
func Database() *gorm.DB {
	return instance
}

// InitDB create instance of DB.
func InitDB(dbType string, dbOption ...string) *gorm.DB {
	var err error

	switch dbType {
	case constants.DBTypeSqlite3: //sqlite3
		instance, err = gorm.Open(dbType, dbOption[0])
	case constants.DBTypeMysql: //mysql
		instance, err = gorm.Open(dbType, dbOption[0]+
			":"+dbOption[1]+"@"+dbOption[2]+
			"?charset=utf8mb4,utf8&parseTime=True&loc=Local")
	}

	if err != nil {
		panic(err)
	}
	if instance == nil {
		panic("Failed to create the handle")
	}

	// Create tables for Tx, and Block.
	initCreateTable()
	return instance
}

// CloseDBInstance delete instance of DB.
func CloseDBInstance() {
	err := Database().Close()
	if err != nil {
		logger.Panicln(err)
	}
	logger.Info("Delete polarbear instance of DB")
}

func initCreateTable() {
	if !instance.HasTable(&Tx{}) {
		instance.CreateTable(&Tx{})
	}
	if !instance.HasTable(&Block{}) {
		instance.CreateTable(&Block{})
	}
	if !instance.HasTable(&Symptom{}) {
		instance.CreateTable(&Symptom{})
	}
}

func convUnixTimeStampToTime(Timestamp int64) time.Time {
	tmpDecStr := strconv.FormatInt(Timestamp, 10)

	// UNIX time stamp in mili seconds timestamp.
	if len(tmpDecStr) == 16 {
		a, _ := strconv.ParseInt(tmpDecStr[0:10], 10, 64)
		b, _ := strconv.ParseInt(tmpDecStr[10:], 10, 64)
		return time.Unix(a, b)

		// UNIX time stamp in nano..
	} else if len(tmpDecStr) == 19 {
		return time.Unix(0, Timestamp)
	}

	// No handle for other cases.
	return time.Unix(0, 0)
}

func buildBlockRecordFromJSON(JSONData map[string]interface{}, URI string, channelName string, block *Block) {

	result := JSONData["result"].(map[string]interface{})
	height := int64(result["height"].(float64))

	// Put block data into table.
	block.BlockHash = util.AddHexHD(result["block_hash"].(string))
	block.Channel = channelName
	block.PeerID = result["peer_id"].(string)
	block.BlockHeight = height

	// Convert float64 to UNIX time.
	tmpDecInt := int64(result["time_stamp"].(float64))
	block.Timestamp = convUnixTimeStampToTime(tmpDecInt)

	// Singature
	block.Signature = result["signature"].(string)

	// Traversal confirmed TXs.
	confirmedTxList := result["confirmed_transaction_list"].([]interface{})

	logger.Infof("Tx count %d,  in the block %s", len(confirmedTxList), block.BlockHash)
	for _, t := range confirmedTxList {

		temp := t.(map[string]interface{})

		// Convert UNIX timestamp in  string  to Time structure.
		var tmpDecInt int64
		if reflect.TypeOf(temp["timestamp"]).Name() == "string" {
			tmpHexStr := temp["timestamp"].(string)
			if strings.Contains(tmpHexStr, "0x") {
				// Hex string to Int64.
				tmpDecInt, _ = strconv.ParseInt(tmpHexStr[2:], 16, 64)
			} else {
				// Decimal string to Int64.
				tmpDecInt, _ = strconv.ParseInt(tmpHexStr, 16, 64)
			}
		} else if reflect.TypeOf(temp["timestamp"]).Name() == "float64" {
			// Float64 timestamp to Int64.
			tmpDecInt = int64(temp["timestamp"].(float64))
		}
		txTimestamp := convUnixTimeStampToTime(tmpDecInt)

		// Parse data string. In some case no data in Tx.
		dataString := ""
		if val, ok := temp["data"]; ok {
			//typeName := reflect.TypeOf(val).String()
			//fmt.Println(typeName)
			if reflect.TypeOf(val).Name() == "string" {
				dataString = temp["data"].(string)
			} else {
				dataMap := val.(map[string]interface{})
				b, err := json.Marshal(dataMap)
				if err != nil {
					panic(err)
				}
				dataString = string(b)
			}
		}

		// Get the TX's status. If URI is "", then just set status as success because JSONData should be test data.
		var txHash string
		if val, ok := temp["txHash"]; ok {
			txHash = util.AddHexHD(val.(string))
		} else if val, ok := temp["tx_hash"]; ok {
			txHash = util.AddHexHD(val.(string))
		} else {
			logger.Error("No key for TxHash!! Block hash: ", result["block_hash"].(string))
		}

		// Set status by calling of JSON RPC.
		var status string
		if URI != "" {
			status, _ = getTxStatus(URI, channelName, txHash)
		} else {
			status = "Success"
		}

		block.Txs = append(
			block.Txs,
			Tx{
				TxHash:      txHash,
				Channel:     channelName,
				Status:      status,
				BlockHeight: height,
				Timestamp:   txTimestamp,
				From:        temp["from"].(string),
				To:          temp["to"].(string),
				Data:        dataString,
			},
		)
	}

}

// AddBlockRecordFromJSONResponse build block data from JSON data.
func AddBlockRecordFromJSONResponse(
	JSONData map[string]interface{},
	URI string,
	channelName string,
	block *Block) error {

	buildBlockRecordFromJSON(JSONData, URI, channelName, block)

	if Database().NewRecord(&block) {
		if err := Database().Save(&block).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetCurrentBlockHeightInDB is
func GetCurrentBlockHeightInDB(channelName string) int64 {

	// Check arguments.
	if channelName == "" {
		logger.Errorf("Arguments is wrong. channelName:%s", channelName)
		return -1
	}

	var block Block
	blockTable := Database().Preload("Txs").Model(&Block{}).Order("block_height desc")
	if err := blockTable.First(&block, &Block{Channel: channelName}).Error; err != nil {
		if err.Error() == "record not found" {
			logger.Info("No blocks in the DB. ")
			return 0
		} else {
			logger.Errorf("%s", err)
			return -1
		}

	} else {
		return block.BlockHeight
	}
}

// QueryBlocksInChannel queries blocks in channel.
func QueryBlocksInChannel(
	channelName string,
	limit int,
	offset int,
	out interface{}) (int64, error) {

	// Check arguments.
	if limit < 0 || offset < 0 || channelName == "" {
		logger.Errorf("Arguments is wrong. limit:%d, offset:%d, channelName:%s", limit, offset, channelName)
		return -1, isaacerror.SysErrFailToQueryBlockInChannel
	}

	//Query block data.
	var block Block
	Database().First(&block, 1)
	blockTable := Database().Preload("Txs").Model(&block).Order(
		"block_height desc")

	if err := blockTable.Offset(offset).Limit(limit).Find(out, &Block{
		Channel: channelName,
	}).Error; err != nil {
		return -1, isaacerror.SysErrFailToQueryBlockInChannel
	} else {
		var count int64
		var block2 Block
		Database().Model(&block2).Where(&Block{
			Channel: channelName,
		}).Count(&count)
		return count, nil
	}
}

// QueryBlockByHeightInChannel queries block by height in channel.
func QueryBlockByHeightInChannel(
	channelName string,
	height int64,
	out interface{}) error {

	// Check arguments.
	if height <= 0 || channelName == "" {
		logger.Errorf("Arguments is wrong. height:%d,  channelName:%s", height, channelName)
		return isaacerror.SysErrFailToQueryBlockInChannel
	}

	//Query block data.
	var block Block
	Database().First(&block, 1)
	blockTable := Database().Preload(
		"Txs").Model(
		&block).Order(
		"block_height desc")

	if err := blockTable.Find(out, &Block{
		Channel:     channelName,
		BlockHeight: height},
	).Error; err != nil {
		return isaacerror.SysErrFailToQueryBlocksInChannel
	} else {
		return nil
	}
}

// QueryBlockInChannelByHash queries block by height in channel.
func QueryBlockInChannelByHash(
	channelName string,
	blockHash string,
	out interface{}) error {

	//Query block data.
	var block Block
	Database().First(&block, 1)
	blockTable := Database().Preload("Txs").Model(&block).Order(
		"block_height desc")

	if err := blockTable.Find(out, &Block{
		Channel:   channelName,
		BlockHash: blockHash},
	).Error; err != nil {
		return isaacerror.SysErrFailToQueryBlockInChannel
	} else {
		return nil
	}
}

// QueryTxsInChannelByTime is
func QueryTxsInChannelByTime(
	channelName string,
	limit int,
	offset int,
	out interface{}) (int64, error) {

	//Query Tx data.
	var tx Tx
	Database().First(&tx, 1)
	txTable :=
		Database().Model(
			&tx).Offset(
			offset).Limit(
			limit).Order("block_height desc").Order("timestamp desc")

	if err := txTable.Find(out, &Tx{
		Channel: channelName,
	}).Error; err != nil {
		return -1, isaacerror.SysErrFailToQueryTxsInChannel
	} else {
		var count int64
		var tx2 Tx
		Database().Model(&tx2).Where(Tx{
			Channel: channelName,
		}).Count(&count)

		return count, nil
	}
}

// QueryTxsInChannelBySearch is
func QueryTxsInChannelBySearch(
	channelName string,
	limit int,
	offset int,
	search TxSearch,
	out interface{}) (int64, error) {

	//Query Tx data.
	var count int64

	// Search by channel name.
	txTable := Database().Where(Tx{Channel: channelName}).Model(&Tx{})

	// If exists search entity, search by search entity.
	//  entity : status, blockHeight, fromAddress, toAddress
	if search.BlockHeight != -1 {
		txTable = txTable.Where(Tx{BlockHeight: search.BlockHeight}).Model(&Tx{})
	}

	if search.Status != "" {
		txTable = txTable.Where(Tx{Status: search.Status}).Model(&Tx{})
	}

	if search.FromAddress != "" {
		txTable = txTable.Where(Tx{From: search.FromAddress}).Model(&Tx{})
	}

	if search.ToAddress != "" {
		txTable = txTable.Where(Tx{To: search.ToAddress}).Model(&Tx{})
	}

	// If exists 'from / to' entity, search by 'from / to'.
	if !search.From.IsZero() && !search.To.IsZero() {
		txTable = txTable.Where("timeStamp BETWEEN ? AND ?", search.From, search.To).Model(&Tx{})
	} else if search.From.IsZero() && search.To.IsZero() {
	} else {
		return -1, isaacerror.SysErrInvalidTimeSearchCondition
	}
	// If exists data entity, search by data entity.
	if search.Data != "" {
		txTable = txTable.Where("data LIKE ?", "%"+search.Data+"%").Model(&Tx{})
	}

	// Get count about Tx list.
	txTable.Count(&count)

	// Get Tx list.
	err := txTable.Offset(offset).Limit(limit).Order("block_height desc").Order("timestamp desc").Find(out, &Tx{}).Error
	if err != nil {
		return -1, isaacerror.SysErrFailToQueryTxsInChannel
	}

	return count, nil
}

// QueryTxInChannelByHash is
func QueryTxInChannelByHash(
	channelName string,
	txHash string,
	out interface{}) error {

	//Query block data.
	var tx Tx
	Database().First(&tx, 1)

	if err := Database().Model(&tx).Find(out, &Tx{
		Channel: channelName,
		TxHash:  txHash,
	}).Error; err != nil {
		return isaacerror.SysErrFailToQueryTxInChannel
	} else {
		return nil
	}

}


// AddPeerSymptom insert peer symptom data tables.
func AddPeerSymptom(
	channelName string,
	channelPK string,
	symptom string,
	msg string) error {

	peerSymptomMappingTB := &Symptom{
		Channel: channelName,
		Channel_PK: channelPK,
		Msg : msg,
		SymptomType: symptom,
		Timestamp: time.Now(),
	}

	if Database().NewRecord(&peerSymptomMappingTB) {
		if err := Database().Save(&peerSymptomMappingTB).Error; err != nil {
			return err
		}
	}

	return nil
}

// QueryPeerSymptomListTable return peersymptom table Info.
func QueryPeerSymptomListTable(
	limit int,
	offset int,
	from string,
	to string,
	channelPermission []string,
	out interface{}) (int64, error) {

	// Check arguments.
	if limit < 0 || offset < 0 {
		logger.Errorf("Arguments is wrong. limit:%d, offset:%d", limit, offset)
		return -1, isaacerror.SysErrFailToQueryPeerSymptom
	}

	//Query symptom data.
	var peerSymptomTB []Symptom
	var symptomTable *gorm.DB
	Database().First(&peerSymptomTB, 1)
	if from != "" && to != "" {
		fromTimer, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return -1, isaacerror.SysErrFailToParseTimeStringToTimeObject
		}
		toTimer, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return -1, isaacerror.SysErrFailToParseTimeStringToTimeObject
		}
		symptomTable = Database().Model(&peerSymptomTB).Order("Timestamp desc").
			Where("Timestamp BETWEEN ? AND ?", fromTimer, toTimer)
	} else {
		symptomTable = Database().Model(&peerSymptomTB).Order("Timestamp desc")
	}

	var count int64
	symptomTable.Where("Channel_PK IN (?)", channelPermission).
		Count(&count)

	if err := symptomTable.Offset(offset).Limit(limit).
		Where("Channel_PK IN (?)", channelPermission).
		Find(out, &Symptom{
	}).Error; err != nil {
		return -1, isaacerror.SysErrFailToQueryPeerSymptom
	}

	return count, nil
}