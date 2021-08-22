package db

import (
	"math/rand"
	"motherbear/backend/constants"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"strconv"
	"strings"
	"sync"

	"time"

	. "motherbear/backend/configuration"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	_ "github.com/jinzhu/gorm/dialects/mysql"

	"golang.org/x/crypto/bcrypt"
)

type USER_INFO_TB struct {
	PK                               string `gorm:"type:varchar(40);primary_key"`
	USER_ID                          string `gorm:"type:varchar(40);not null"`
	FIRST_NAME                       string `gorm:"type:varchar(40);not null"`
	LAST_NAME                        string `gorm:"type:varchar(40);not null"`
	PASSWORD                         string `gorm:"type:varchar(80);not null"`
	EMAIL_ADRES                      string `gorm:"type:varchar(80);not null"`
	MOBILE_PHONE_NO                  string `gorm:"type:varchar(40);not null"`
	TYPE_CODE                        string `gorm:"type:varchar(20);not null"`
	TIME_WHEN_USER_MODIFIED_PASSWORD time.Time
	CREATE_DATE                      time.Time
	LATEST_LOGIN_DATE                time.Time
}

type CHANNEL_USER_MAPPING_TB struct {
	PK         string `gorm:"type:varchar(40);primary_key;not null"`
	CHANNEL_PK string `gorm:"type:varchar(40);primary_key;not null"`
}

type PERMISSION_USER_MAPPING_TB struct {
	PK               string `gorm:"type:varchar(40);primary_key;not null"`
	PERMISSION_ALIAS string `gorm:"type:varchar(120);primary_key;not null"`
	PERMISSION_CHECK bool
}

type PERMISSION_INFO_TB struct {
	PERMISSION_INDEX int    `gorm:"AUTO_INCREMENT;primary_key;not null"`
	PERMISSION_ALIAS string `gorm:"type:varchar(120);primary_key;not null"`
}

type CONFIGURATION_CHANNEL_TB struct {
	CHANNEL_NAME string `gorm:"type:varchar(40)"`
	CHANNEL_PK   string `gorm:"type:varchar(40);primary_key;not null"`
	CHANNEL_ID   string `gorm:"type:varchar(40)"` // It is the channel ID of goloop. Can be use only on the goloop.
}

type CONFIGURATION_MASTER_TB struct {
	PROMETHEUS_IP   string `gorm:"type:varchar(256)"`
	EXPORTER_IP     string `gorm:"type:varchar(256)"`
	SESSION_TIMEOUT int
	SET_LANGUAGE    string `gorm:"type:varchar(16)"`
	PK_SOURCE       int
}

type CONFIGURATION_NODE_TB struct {
	NODE_NAME    string `gorm:"type:varchar(40)"`
	NODE_PK      string `gorm:"type:varchar(40);primary_key;not null"`
	NODE_IP      string `gorm:"type:varchar(40)"`
	NODE_ADDRESS string `gorm:"type:varchar(50)"`
}

type NODE_CHANNEL_MAPPING_TB struct {
	CHANNEL_PK string `gorm:"type:varchar(40);primary_key;not null"`
	NODE_PK    string `gorm:"type:varchar(40);primary_key;not null"`
}

type CONFIGURATION_DATA_ALERT_TB struct {
	ALERT_METHOD                int    //1: ALARM, 2:SMS, 3:EMAIL
	ALERT_LEVEL                 string `gorm:"type:varchar(5)"`
	MAX_TIME_SEC_FOR_UNSYNC     int
	MAX_UNSYNC_BLOCK_DIFFERENCE int
	MAX_TIME_SEC_FOR_RESPONSE   int
	CHANNEL_PK                  string `gorm:"type:varchar(40);primary_key;not null"`
}

type CONFIGURATION_DATA_VISIBILITY_TB struct {
	CHECK_HOST_NAME     bool
	CHECK_BLOCK_HEIGHT  bool
	CHECK_RESPONSE_TIME bool
	CHECK_NODE_IP       bool
	CHECK_TRANSACTION   bool
	CHECK_LEADER        bool
	CHANNEL_PK          string `gorm:"type:varchar(40);primary_key;not null"`
}

// Not currently in use. Implemented for future use
type ROLE_GROUP_TB struct {
	GROUP_PK         string `gorm:"type:varchar(40);primary_key;not null"`
	GROUP_NAME_ALIAS string `gorm:"type:varchar(40)"`
}

// Not currently in use. Implemented for future use
type CHANNEL_GROUP_MAPPING_TB struct {
	GROUP_PK   string `gorm:"type:varchar(40);primary_key;not null"`
	CHANNEL_PK string `gorm:"type:varchar(40);primary_key;not null"`
}

// Not currently in use. Implemented for future use
type PERMISSION_GROUP_MAPPING_TB struct {
	GROUP_PK         string `gorm:"type:varchar(40);primary_key;not null"`
	PERMISSION_INDEX int    `gorm:"type:int;primary_key;not null"`
	PERMISSION_CHECK bool
}

var instance *gorm.DB
var once sync.Once
var seededRand *rand.Rand

func init() {
	once.Do(func() {
		seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
}

// DB returns the instance.
func DBgorm() *gorm.DB {
	return instance
}

// Transaction instance struct
type Transaction struct {
	once     sync.Once
	rollback bool
	db       *gorm.DB
}

func (t *Transaction) Close() {
	t.once.Do(func() {
		if t.rollback {
			t.db.Rollback()
		} else {
			t.db.Commit()
		}
	})
}

func (t *Transaction) Fail() {
	t.rollback = true
}

func NewTransaction() (tx *Transaction) {
	transaction := &Transaction{}
	transaction.db = DBgorm().Begin()
	return transaction
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
			"?charset=utf8&parseTime=True&loc=Local")
	}

	if err != nil {
		panic(err)
	}
	if instance == nil {
		panic("Failed to create the handle")
	}

	// Create tables for ISAAC
	InitCreateTable()
	return instance
}

// CloseDBInstance delete instance of DB.
func CloseDBInstance() {
	err := DBgorm().Close()
	if err != nil {
		logger.Panicln(err)
	}
	logger.Info("Delete ISAAC instance of DB")
}

// InitCreateTable create init table.
func InitCreateTable() {
	checkFirstRun := false
	if !DBgorm().HasTable(&PERMISSION_INFO_TB{}) {
		permissionInfoTB := &PERMISSION_INFO_TB{
			PERMISSION_ALIAS: constants.DBUserPermissionNode,
		}
		permissionInfoTB1 := &PERMISSION_INFO_TB{
			PERMISSION_ALIAS: constants.DBUserPermissionMonitoringLog,
		}

		DBgorm().CreateTable(&permissionInfoTB)
		DBgorm().Create(&permissionInfoTB)
		DBgorm().Create(&permissionInfoTB1)
	}

	if !DBgorm().HasTable(&USER_INFO_TB{}) {
		userInfoTB := &USER_INFO_TB{
			PK:              constants.AdminPK,                                              // User pk id
			USER_ID:         constants.AdminID,                                              // User ID
			FIRST_NAME:      "firstName",                                                    // User first name
			LAST_NAME:       "lastName",                                                     // User last name
			PASSWORD:        "$2a$10$XToqmLMc7XTwwhuOokpFaOBpqUBEMQuoEyS1b7EU5eQES4ijBGxIu", // bcrypt(“admin123”)
			EMAIL_ADRES:     "admin@admin.com",                                              // Email
			MOBILE_PHONE_NO: "010-000-0000",                                                 // Mobile phone number
			TYPE_CODE:       constants.DBUserTypeCodeAdmin,                                  // Type code
			CREATE_DATE:     time.Now(),                                                     // create time
		}

		DBgorm().CreateTable(&userInfoTB)
		DBgorm().Create(&userInfoTB)

		checkFirstRun = true
	}

	if !DBgorm().HasTable(&CHANNEL_USER_MAPPING_TB{}) {
		DBgorm().CreateTable(&CHANNEL_USER_MAPPING_TB{})
	}

	if !DBgorm().HasTable(&PERMISSION_USER_MAPPING_TB{}) {
		DBgorm().CreateTable(&PERMISSION_USER_MAPPING_TB{})

		permissionList := GetPermissionTable()
		for _, value := range permissionList {
			InsertUserPermissions(constants.AdminPK, value.PERMISSION_ALIAS, true)
		}
	}

	if !DBgorm().HasTable(&PERMISSION_INFO_TB{}) {
		DBgorm().CreateTable(&PERMISSION_INFO_TB{})
	}

	if !DBgorm().HasTable(&ROLE_GROUP_TB{}) {
		DBgorm().CreateTable(&ROLE_GROUP_TB{})
	}

	if !DBgorm().HasTable(&CHANNEL_GROUP_MAPPING_TB{}) {
		DBgorm().CreateTable(&CHANNEL_GROUP_MAPPING_TB{})
	}

	if !DBgorm().HasTable(&PERMISSION_GROUP_MAPPING_TB{}) {
		DBgorm().CreateTable(&PERMISSION_GROUP_MAPPING_TB{})
	}

	if !DBgorm().HasTable(&CONFIGURATION_CHANNEL_TB{}) {
		DBgorm().CreateTable(&CONFIGURATION_CHANNEL_TB{})
	}

	if !DBgorm().HasTable(&CONFIGURATION_MASTER_TB{}) {
		configurationMasterTB := &CONFIGURATION_MASTER_TB{
			SESSION_TIMEOUT: 30,
			SET_LANGUAGE:    "ko",
			PK_SOURCE:       8,
		}
		DBgorm().CreateTable(&CONFIGURATION_MASTER_TB{})
		DBgorm().Create(&configurationMasterTB)
	}

	if !DBgorm().HasTable(&CONFIGURATION_NODE_TB{}) {
		DBgorm().CreateTable(&CONFIGURATION_NODE_TB{})
	}

	if !DBgorm().HasTable(&NODE_CHANNEL_MAPPING_TB{}) {
		DBgorm().CreateTable(&NODE_CHANNEL_MAPPING_TB{})
	}

	if !DBgorm().HasTable(&CONFIGURATION_DATA_ALERT_TB{}) {
		DBgorm().CreateTable(&CONFIGURATION_DATA_ALERT_TB{})
	}

	if !DBgorm().HasTable(&CONFIGURATION_DATA_VISIBILITY_TB{}) {
		DBgorm().CreateTable(&CONFIGURATION_DATA_VISIBILITY_TB{})
	}

	if checkFirstRun == true {
		for _, e := range Conf().Node {
			InsertConfigurationNode(e.Name, e.IP)
		}

		for _, ch := range Conf().Channel {
			pk := InsertConfigurationChannel(ch.Name)
			for _, nd := range ch.Nodes {
				configurationNodeTB := &CONFIGURATION_NODE_TB{}
				configurationNodeTB = GetConfigurationNodeInfoByNodeName(nd)
				InsertChannelPermissionNode(pk, configurationNodeTB.NODE_PK)
			}
		}
		UpdateConfigurationMasterInfo(Conf().ETC.SessionTimeout, Conf().ETC.Language)
	}
}

// leftPad2Len left padding
func leftPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	var retStr = strings.Repeat(padStr, padCountInt) + s
	return retStr[(len(retStr) - overallLen):]
}

// StringWithCharset return rand string
func stringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// getPKSource return pk source.
func getPKSource() string {
	configurationMasterTB := &CONFIGURATION_MASTER_TB{}
	DBgorm().First(&configurationMasterTB)
	var pkSource int = configurationMasterTB.PK_SOURCE
	if pkSource > 0 {
		pkSource = pkSource + 1
	} else {
		pkSource = 8
	}

	pkSourceStr := strconv.Itoa(pkSource)
	pkRanStr := stringWithCharset(8, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	pkSourceRet := leftPad2Len(pkSourceStr, "0", 8)
	pkSourceRet = pkRanStr + pkSourceRet

	DBgorm().Model(&configurationMasterTB).Update("PK_SOURCE", pkSource)

	return pkSourceRet
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetNodePKToDataMap() map[string]CONFIGURATION_NODE_TB {
	var nodeTB []CONFIGURATION_NODE_TB
	nodeTB = GetConfigurationNodeTable()
	nodePKToData := make(map[string]CONFIGURATION_NODE_TB)
	for _, value := range nodeTB {
		nodePKToData[value.NODE_PK] = value
	}

	return nodePKToData
}

func GetNodeNameToPKMap() map[string]string {
	var nodeTB []CONFIGURATION_NODE_TB
	nodeTB = GetConfigurationNodeTable()
	nodeNameToPK := make(map[string]string)
	for _, value := range nodeTB {
		nodeNameToPK[value.NODE_NAME] = value.NODE_PK
	}

	return nodeNameToPK
}

func GetChannelPKToNameMap() map[string]string {
	var channelTB []CONFIGURATION_CHANNEL_TB
	channelTB = GetConfigurationChannelTable()
	channelPKToName := make(map[string]string)
	for _, value := range channelTB {
		channelPKToName[value.CHANNEL_PK] = value.CHANNEL_NAME
	}

	return channelPKToName
}

func GetChannelNameToPKMap() map[string]string {
	var channelTB []CONFIGURATION_CHANNEL_TB
	channelTB = GetConfigurationChannelTable()
	channelNameToPK := make(map[string]string)
	for _, value := range channelTB {
		channelNameToPK[value.CHANNEL_NAME] = value.CHANNEL_PK
	}

	return channelNameToPK
}

func GetChannelIDToNameMap() map[string]string {
	var channelTB []CONFIGURATION_CHANNEL_TB
	channelTB = GetConfigurationChannelTable()
	channelIDToName := make(map[string]string)
	for _, value := range channelTB {
		channelIDToName[value.CHANNEL_ID] = value.CHANNEL_NAME
	}

	return channelIDToName
}

func ConvertPKCHtoChannelName(channelName string) (string, error) {
	if strings.Contains(channelName, "PKCH") {
		pkChannelString := channelName
		channelInfo := GetConfigurationChannelInfo(pkChannelString)
		if channelInfo != nil {
			channelName = channelInfo.CHANNEL_NAME
		} else {
			return "", isaacerror.SysErrNoChannelInDB
		}
	}
	return channelName, nil
}

func ConvertChannelIDToName(channelID string) (string, error) {
	channelName := ""

	channelTB := GetConfigurationChannelTable()
	for _, channel := range channelTB {
		if channel.CHANNEL_ID == channelID {
			channelName = channel.CHANNEL_NAME
			break
		}
	}

	if channelName == "" {
		err := isaacerror.SysErrNoChannelInISAAC
		logger.Error(err.Error())
		return "", err
	}

	return channelName, nil
}

func IsExistChannelListInDB(channelIDList []string) bool {
	channelListTB := GetConfigurationChannelTable()

	for _, channelID := range channelIDList {
		check := false
		for _, channelTB := range channelListTB {
			if channelID == channelTB.CHANNEL_PK {
				check = true
				break
			}
		}
		if !check {
			return false
		}
	}

	return true
}

func IsExistPermissionListInDB(permissionList []string) bool {
	permissionListTB := GetPermissionTable()

	for _, permission := range permissionList {
		check := false
		for _, permissionTB := range permissionListTB {
			if permission == permissionTB.PERMISSION_ALIAS {
				check = true
				break
			}
		}
		if !check {
			return false
		}
	}

	return true
}
