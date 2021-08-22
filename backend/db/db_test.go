package db

import (
	"motherbear/backend/constants"
	"motherbear/backend/utility"
	"os"
	"strings"
	"testing"
	"time"

	"motherbear/backend/configuration"

	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/assert.v1"
)

func TestInitCreateTable(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	table := []string{
		"USER_INFO_TB",
		"CHANNEL_USER_MAPPING_TB",
		"PERMISSION_USER_MAPPING_TB",
		"PERMISSION_INFO_TB",
		"ROLE_GROUP_TB",
		"CHANNEL_GROUP_MAPPING_TB",
		"PERMISSION_GROUP_MAPPING_TB",
		"CONFIGURATION_CHANNEL_TB",
		"CONFIGURATION_MASTER_TB",
		"CONFIGURATION_NODE_TB",
		"NODE_CHANNEL_MAPPING_TB",
		"CONFIGURATION_DATA_ALERT_TB",
	}

	for i := 0; i < len(table); i++ {
		if !DBgorm().HasTable(strings.ToLower(table[i]) + "s") {
			t.Error("Error Not Created " + table[i] + " Table")
		}
	}

	userInfoTB := &USER_INFO_TB{}

	if DBgorm().HasTable(userInfoTB) {
		var count int

		userInfoTB := &USER_INFO_TB{}

		DBgorm().First(&userInfoTB)

		DBgorm().Model(&USER_INFO_TB{}).Where("PK = ?", "PKID_0000000000000000").Count(&count)

		if count != 1 {
			t.Error("Error Insert initial user data")
		}
	}
}

func Setup(path string) *gorm.DB {

	// Initialize yaml settings
	if _, err := os.Stat("./config"); os.IsNotExist(err) {
		os.MkdirAll("./config", os.ModePerm)
	}

	// Read the configuration file.
	configuration.InitConfigData("config/configuration.yaml")

	database := InitDB("sqlite3", path)
	userInfoTB := &USER_INFO_TB{
		PK:                               "PKID_0000000000000001",
		USER_ID:                          "test1",
		FIRST_NAME:                       "test1First",
		LAST_NAME:                        "test1Last",
		PASSWORD:                         "$2a$10$XToqmLMc7XTwwhuOokpFaOBpqUBEMQuoEyS1b7EU5eQES4ijBGxIu",
		EMAIL_ADRES:                      "test@test.com",
		MOBILE_PHONE_NO:                  "010-111-1111",
		TYPE_CODE:                        constants.DBUserTypeCodeCommon,
		TIME_WHEN_USER_MODIFIED_PASSWORD: time.Now(),
		CREATE_DATE:                      time.Now(),
		LATEST_LOGIN_DATE:                time.Now(),
	}
	database.Create(&userInfoTB)
	channelUserMappingTB := &CHANNEL_USER_MAPPING_TB{
		PK:         "PKID_0000000000000001",
		CHANNEL_PK: "PKCH_0000000000000000",
	}
	channelUserMappingTB2 := &CHANNEL_USER_MAPPING_TB{
		PK:         "PKID_0000000000000001",
		CHANNEL_PK: "PKCH_0000000000000001",
	}
	database.Create(&channelUserMappingTB)
	database.Create(&channelUserMappingTB2)
	permissionUserMappingTB := &PERMISSION_USER_MAPPING_TB{
		PK:               "PKID_0000000000000001",
		PERMISSION_ALIAS: constants.DBUserPermissionNode,
		PERMISSION_CHECK: true,
	}
	permissionUserMappingTB2 := &PERMISSION_USER_MAPPING_TB{
		PK:               "PKID_0000000000000001",
		PERMISSION_ALIAS: constants.DBUserPermissionMonitoringLog,
		PERMISSION_CHECK: false,
	}
	database.Create(&permissionUserMappingTB)
	database.Create(&permissionUserMappingTB2)
	roleGroupTB := &ROLE_GROUP_TB{
		GROUP_PK:         "PKGR_0000000000000001",
		GROUP_NAME_ALIAS: "System CH1",
	}
	database.Create(&roleGroupTB)

	channelGroupMappingTB := &CHANNEL_GROUP_MAPPING_TB{
		GROUP_PK:   "PKGR_0000000000000001",
		CHANNEL_PK: "PKCH_0000000000000000",
	}
	channelGroupMappingTB2 := &CHANNEL_GROUP_MAPPING_TB{
		GROUP_PK:   "PKGR_0000000000000001",
		CHANNEL_PK: "PKCH_0000000000000001",
	}
	database.Create(&channelGroupMappingTB)
	database.Create(&channelGroupMappingTB2)

	permissionGroupMappingTB := &PERMISSION_GROUP_MAPPING_TB{
		GROUP_PK:         "PKGR_0000000000000001",
		PERMISSION_INDEX: 1,
		PERMISSION_CHECK: false,
	}
	permissionGroupMappingTB2 := &PERMISSION_GROUP_MAPPING_TB{
		GROUP_PK:         "PKGR_0000000000000001",
		PERMISSION_INDEX: 2,
		PERMISSION_CHECK: true,
	}
	database.Create(&permissionGroupMappingTB)
	database.Create(&permissionGroupMappingTB2)

	return database
}

func Teardown(path string) {
	os.Remove(path)
}

func TestGetNodePKToDataMap(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	nodeTB := GetConfigurationNodeTable()
	nodePKToData := GetNodePKToDataMap()

	for _, value := range nodeTB {
		assert.Equal(t, value, nodePKToData[value.NODE_PK])
	}
}

func TestGetNodeNameToPKMap(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	nodeTB := GetConfigurationNodeTable()
	nodeNameToPK := GetNodeNameToPKMap()

	for _, value := range nodeTB {
		assert.Equal(t, value.NODE_PK, nodeNameToPK[value.NODE_NAME])
	}
}

func TestGetChannelPKToNameMap(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	channelTB := GetConfigurationChannelTable()
	channelPKToData := GetChannelPKToNameMap()

	for _, value := range channelTB {
		assert.Equal(t, value.CHANNEL_NAME, channelPKToData[value.CHANNEL_PK])
	}
}

func TestGetChannelNameToPKMap(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	channelTB := GetConfigurationChannelTable()
	channelPKToData := GetChannelNameToPKMap()

	for _, value := range channelTB {
		assert.Equal(t, value.CHANNEL_PK, channelPKToData[value.CHANNEL_NAME])
	}
}

func TestGetChannelIDToNameMap(t *testing.T) {
	t.Skip("Skipping test GetChannelIDToNameMap. Should prepare configuration.yaml of the goloop to test.")
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	channelIDToName := GetChannelIDToNameMap()

	channelLen := len(configuration.Conf().Channel)

	assert.Equal(t, channelLen, len(channelIDToName))

}

func TestConvertChannelIDToName(t *testing.T) {
	t.Skip("Skipping test ConvertChannelIDToName. Should prepare configuration.yaml of the goloop to test.")

	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	correctChannelName := "1"
	channelID := "0x6b738e"

	channelName, _ := ConvertChannelIDToName(channelID)

	assert.Equal(t, correctChannelName, channelName)
}

func TestIsExistChannelListInDB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	channelTB1 := &CONFIGURATION_CHANNEL_TB{
		CHANNEL_PK:   "PKCH_0000000000000000",
		CHANNEL_NAME: "channel1",
	}
	channelTB2 := &CONFIGURATION_CHANNEL_TB{
		CHANNEL_PK:   "PKCH_0000000000000001",
		CHANNEL_NAME: "channel2",
	}
	DBgorm().Create(&channelTB1)
	DBgorm().Create(&channelTB2)

	channelList1 := []string{"PKCH_0000000000000000", "PKCH_0000000000000001"}
	channelList2 := []string{"PKCH_0000000000000003"}

	assert.Equal(t, true, IsExistChannelListInDB(channelList1))
	assert.Equal(t, false, IsExistChannelListInDB(channelList2))
}

func TestIsExistPermissionListInDB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	permissionTB := GetPermissionTable()
	var permissionList1 []string
	for _, value := range permissionTB {
		permissionList1 = append(permissionList1, value.PERMISSION_ALIAS)
	}

	permissionList2 := []string{"asdf"}

	assert.Equal(t, true, IsExistPermissionListInDB(permissionList1))
	assert.Equal(t, false, IsExistPermissionListInDB(permissionList2))
}

func TestGetUserInfo(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var userInfoTB []USER_INFO_TB
	userInfoTB = GetUserInfo()

	assert.Equal(t, 2, len(userInfoTB))
	assert.Equal(t, userInfoTB[0].USER_ID, "admin")
	assert.Equal(t, userInfoTB[0].EMAIL_ADRES, "admin@admin.com")
	assert.Equal(t, userInfoTB[0].MOBILE_PHONE_NO, "010-000-0000")
	assert.Equal(t, userInfoTB[0].TYPE_CODE, constants.DBUserTypeCodeAdmin)
}

func TestGetUserInfoByUserID(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	userInfoTB := &USER_INFO_TB{}
	userInfoTB = GetUserInfoByUserID("admin")

	assert.Equal(t, userInfoTB.USER_ID, "admin")
	assert.Equal(t, userInfoTB.EMAIL_ADRES, "admin@admin.com")
	assert.Equal(t, userInfoTB.MOBILE_PHONE_NO, "010-000-0000")
	assert.Equal(t, userInfoTB.TYPE_CODE, constants.DBUserTypeCodeAdmin)
}

func TestGetUserPKID(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	pk := GetUserPK("admin")

	assert.Equal(t, pk, "PKID_0000000000000000")
}

func TestGetUserPermissionChannels(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var channelUserMappingTB []CHANNEL_USER_MAPPING_TB
	channelUserMappingTB = GetUserPermissionChannels("PKID_0000000000000001")

	assert.Equal(t, channelUserMappingTB[0].CHANNEL_PK, "PKCH_0000000000000000")
	assert.Equal(t, channelUserMappingTB[1].CHANNEL_PK, "PKCH_0000000000000001")
}

func TestGetUserPermissionChannelsTable(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var channelUserMappingTB []CHANNEL_USER_MAPPING_TB
	channelUserMappingTB = GetUserPermissionChannelsTable()

	assert.Equal(t, channelUserMappingTB[2].CHANNEL_PK, "PKCH_0000000000000000")
	assert.Equal(t, channelUserMappingTB[3].CHANNEL_PK, "PKCH_0000000000000001")
}

func TestGetUserPermissionTable(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var permissionUserMappingTB []PERMISSION_USER_MAPPING_TB
	permissionUserMappingTB = GetUserPermissionTable()

	for _, permissionMap := range permissionUserMappingTB {
		if permissionMap.PK == "PKID_0000000000000001" && permissionMap.PERMISSION_ALIAS == constants.DBUserPermissionNode {
			assert.Equal(t, permissionMap.PERMISSION_CHECK, true)
		}
		if permissionMap.PK == "PKID_0000000000000001" && permissionMap.PERMISSION_ALIAS == constants.DBUserPermissionMonitoringLog {
			assert.Equal(t, permissionMap.PERMISSION_CHECK, false)
		}

	}
}

func TestGetUserPermissions(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var permissionUserMappingTB []PERMISSION_USER_MAPPING_TB
	permissionUserMappingTB = GetUserPermissions("PKID_0000000000000001")

	for _, permissionMap := range permissionUserMappingTB {
		if permissionMap.PERMISSION_ALIAS == constants.DBUserPermissionNode {
			assert.Equal(t, permissionMap.PERMISSION_CHECK, true)
		}
		if permissionMap.PERMISSION_ALIAS == constants.DBUserPermissionMonitoringLog {
			assert.Equal(t, permissionMap.PERMISSION_CHECK, false)
		}

	}
}

func TestGetPermissionTable(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var permissionInfoTB []PERMISSION_INFO_TB

	permissionInfoTB = GetPermissionTable()

	assert.Equal(t, permissionInfoTB[0].PERMISSION_ALIAS, constants.DBUserPermissionNode)
}

func TestGetPermissionInfo(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	permissionInfoTB := &PERMISSION_INFO_TB{}

	permissionInfoTB = GetPermissionInfo(1)

	assert.Equal(t, permissionInfoTB.PERMISSION_ALIAS, constants.DBUserPermissionNode)
}

func TestGetRoleGroupInfo(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	roleGroupTB := &ROLE_GROUP_TB{}
	roleGroupTB = GetRoleGroupInfo("PKGR_0000000000000001")

	assert.Equal(t, roleGroupTB.GROUP_NAME_ALIAS, "System CH1")
}

func TestGetRoleGroupTable(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var roleGroupTB []ROLE_GROUP_TB
	roleGroupTB = GetRoleGroupTable()

	assert.Equal(t, roleGroupTB[0].GROUP_PK, "PKGR_0000000000000001")
	assert.Equal(t, roleGroupTB[0].GROUP_NAME_ALIAS, "System CH1")
}

func TestGetGroupPermissionChannels(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var channelGroupMappingTB []CHANNEL_GROUP_MAPPING_TB
	channelGroupMappingTB = GetGroupPermissionChannels("PKGR_0000000000000001")

	assert.Equal(t, channelGroupMappingTB[0].CHANNEL_PK, "PKCH_0000000000000000")
	assert.Equal(t, channelGroupMappingTB[1].CHANNEL_PK, "PKCH_0000000000000001")
}

func TestGetGroupPermissionChannelsTable(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var channelGroupMappingTB []CHANNEL_GROUP_MAPPING_TB
	channelGroupMappingTB = GetGroupPermissionChannelsTable()

	assert.Equal(t, channelGroupMappingTB[0].CHANNEL_PK, "PKCH_0000000000000000")
	assert.Equal(t, channelGroupMappingTB[1].CHANNEL_PK, "PKCH_0000000000000001")
}

func TestGetGroupPermissions(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var permissionGroupMappingTB []PERMISSION_GROUP_MAPPING_TB
	permissionGroupMappingTB = GetGroupPermissions("PKGR_0000000000000001")

	assert.Equal(t, permissionGroupMappingTB[0].PERMISSION_CHECK, false)
	assert.Equal(t, permissionGroupMappingTB[1].PERMISSION_CHECK, true)
}

func TestGetAlertConfigInfoByPK(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var configurationChannelTB []CONFIGURATION_CHANNEL_TB
	configurationChannelTB = GetConfigurationChannelTable()

	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{}
	configurationAlertDataTB = GetAlertConfigInfoByPK(configurationChannelTB[0].CHANNEL_PK)

	assert.Equal(t, configurationAlertDataTB.MAX_TIME_SEC_FOR_UNSYNC, constants.DBDefaultUnsyncBlockToleranceTime)
}

func TestGetAlertConfigTable(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var configurationAlertDataTB []CONFIGURATION_DATA_ALERT_TB
	configurationAlertDataTB = GetAlertConfigTable()

	assert.Equal(t, configurationAlertDataTB[0].MAX_TIME_SEC_FOR_UNSYNC, constants.DBDefaultUnsyncBlockToleranceTime)
}

func TestGetConfigurationChannelTable(t *testing.T) {
	// dbpath := "mbears.db"
	// database := Setup(dbpath)
	// defer Teardown(dbpath)
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var configurationChannelTB []CONFIGURATION_CHANNEL_TB
	configurationChannelTB = GetConfigurationChannelTable()

	assert.Equal(t, configurationChannelTB[0].CHANNEL_NAME, "loopchain_default")
}

func TestGetConfigurationChannelInfo(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var configurationChannelTB []CONFIGURATION_CHANNEL_TB
	configurationChannelTB = GetConfigurationChannelTable()

	configuration_channel_ch_tb := &CONFIGURATION_CHANNEL_TB{}
	configuration_channel_ch_tb = GetConfigurationChannelInfo(configurationChannelTB[0].CHANNEL_PK)

	assert.Equal(t, configuration_channel_ch_tb.CHANNEL_NAME, "loopchain_default")
}

func TestGetConfigurationChannelInfoByList(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var configurationChannelTB []CONFIGURATION_CHANNEL_TB
	configurationChannelTB = GetConfigurationChannelTable()

	channelList := []string{configurationChannelTB[0].CHANNEL_PK, configurationChannelTB[1].CHANNEL_PK}

	var configuration_channel_ch_tb []CONFIGURATION_CHANNEL_TB
	configuration_channel_ch_tb = GetConfigurationChannelInfoByList(channelList)

	assert.Equal(t, len(channelList), len(configuration_channel_ch_tb))
	for _, value := range configurationChannelTB {
		assert.Equal(t, true, utility.IsExistValueInList(value.CHANNEL_PK, channelList))
	}
}

func TestGetConfigurationChannelInfoByName(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var configurationChannelTB []CONFIGURATION_CHANNEL_TB
	configurationChannelTB = GetConfigurationChannelTable()

	configuration_channel_ch_tb := &CONFIGURATION_CHANNEL_TB{}
	configuration_channel_ch_tb = GetConfigurationChannelInfoByName(configurationChannelTB[0].CHANNEL_NAME)

	assert.Equal(t, configuration_channel_ch_tb.CHANNEL_NAME, "loopchain_default")
}

func TestGetConfigurationMasterTable(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	configurationMasterTB := &CONFIGURATION_MASTER_TB{}
	configurationMasterTB = GetConfigurationMasterTable()

	assert.Equal(t, configurationMasterTB.SESSION_TIMEOUT, configuration.Conf().ETC.SessionTimeout)
}

func TestGetConfigurationNodeInfoByNodeID(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	configurationNodeTB := &CONFIGURATION_NODE_TB{}
	configurationNodeTB = GetConfigurationNodeInfoByNodeName("node1")

	assert.Equal(t, configurationNodeTB.NODE_NAME, "node1")
}

func TestGetConfigurationNodeTable(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var configurationNodeTB []CONFIGURATION_NODE_TB
	configurationNodeTB = GetConfigurationNodeTable()

	assert.Equal(t, configurationNodeTB[0].NODE_NAME, "node0")
}

func TestGetChannelPermissionNodesTable(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertChannelPermissionNode("PKCH_0000000000000001", "PKND_0000000000000001")
	InsertChannelPermissionNode("PKCH_0000000000000001", "PKND_0000000000000002")

	var nodeChannelMappingTB []NODE_CHANNEL_MAPPING_TB
	nodeChannelMappingTB = GetChannelPermissionNodesTable()

	assert.Equal(t, nodeChannelMappingTB[5].CHANNEL_PK, "PKCH_0000000000000001")
	assert.Equal(t, nodeChannelMappingTB[5].NODE_PK, "PKND_0000000000000001")
	assert.Equal(t, nodeChannelMappingTB[6].CHANNEL_PK, "PKCH_0000000000000001")
	assert.Equal(t, nodeChannelMappingTB[6].NODE_PK, "PKND_0000000000000002")
}

func TestGetGetChannelPermissionNodes(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertChannelPermissionNode("PKCH_0000000000000001", "PKND_0000000000000001")
	InsertChannelPermissionNode("PKCH_0000000000000001", "PKND_0000000000000002")

	var nodeChannelMappingTB []NODE_CHANNEL_MAPPING_TB
	nodeChannelMappingTB = GetChannelPermissionNodes("PKCH_0000000000000001")

	assert.Equal(t, nodeChannelMappingTB[0].NODE_PK, "PKND_0000000000000001")
	assert.Equal(t, nodeChannelMappingTB[1].NODE_PK, "PKND_0000000000000002")
}

func TestInsertUserInfoTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	channel := []string{"PKCH_0000000000000001", "PKCH_0000000000000002"}
	permission := []string{constants.DBUserPermissionNode, constants.DBUserPermissionMonitoringLog}

	InsertUserInfo("test2", "test2", "test2", "password", "test2@test.com", "010-2222-2222", constants.DBUserTypeCodeCommon, channel, permission)
	InsertUserInfo("test3", "test3", "test3", "password", "test2@test.com", "010-2222-2222", constants.DBUserTypeCodeCommon, channel, permission)

	userInfoTB := &USER_INFO_TB{}
	userInfoTB = GetUserInfoByUserID("test3")

	//	assert.Equal(t, userInfoTB.PK, "PKID_0000000000000020")
	assert.Equal(t, userInfoTB.EMAIL_ADRES, "test2@test.com")
	assert.Equal(t, userInfoTB.MOBILE_PHONE_NO, "010-2222-2222")
	assert.Equal(t, userInfoTB.TYPE_CODE, constants.DBUserTypeCodeCommon)

	var channelUserMappingTB []CHANNEL_USER_MAPPING_TB
	channelUserMappingTB = GetUserPermissionChannels(userInfoTB.PK)

	assert.Equal(t, channelUserMappingTB[0].CHANNEL_PK, "PKCH_0000000000000001")
	assert.Equal(t, channelUserMappingTB[1].CHANNEL_PK, "PKCH_0000000000000002")

	var permissionUserMappingTB []PERMISSION_USER_MAPPING_TB
	permissionUserMappingTB = GetUserPermissions(userInfoTB.PK)

	assert.Equal(t, permissionUserMappingTB[0].PERMISSION_ALIAS, constants.DBUserPermissionMonitoringLog)
	assert.Equal(t, permissionUserMappingTB[0].PERMISSION_CHECK, true)
}

func TestInsertUserPermissionChannelsTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	channel := []string{"PKCH_0000000000000001", "PKCH_0000000000000002"}
	permission := []string{constants.DBUserPermissionNode, constants.DBUserPermissionMonitoringLog}

	InsertUserInfo("test2", "test2", "test2", "password", "test2@test.com", "010-2222-2222", constants.DBUserTypeCodeCommon, channel, permission)

	userInfoTB := &USER_INFO_TB{}
	userInfoTB = GetUserInfoByUserID("test2")

	InsertUserPermissionChannels(userInfoTB.PK, "PKCH_0000000000000003", "PKCH_0000000000000004")

	var channelUserMappingTB []CHANNEL_USER_MAPPING_TB
	channelUserMappingTB = GetUserPermissionChannels(userInfoTB.PK)

	assert.Equal(t, channelUserMappingTB[0].CHANNEL_PK, "PKCH_0000000000000001")
	assert.Equal(t, channelUserMappingTB[1].CHANNEL_PK, "PKCH_0000000000000002")
	assert.Equal(t, channelUserMappingTB[2].CHANNEL_PK, "PKCH_0000000000000003")
	assert.Equal(t, channelUserMappingTB[3].CHANNEL_PK, "PKCH_0000000000000004")
}

func TestInsertConfigurationChannelTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertConfigurationChannel("ch3")

	var configuration_channel_ch_tb []CONFIGURATION_CHANNEL_TB
	configuration_channel_ch_tb = GetConfigurationChannelTable()

	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{}
	configurationAlertDataTB = GetAlertConfigInfoByPK(configuration_channel_ch_tb[2].CHANNEL_PK)

	//assert.Equal(t, configuration_channel_ch_tb[1].CHANNEL_PK, "PKID_0000000000000020")
	assert.Equal(t, configuration_channel_ch_tb[2].CHANNEL_NAME, "ch3")
	assert.Equal(t, configurationAlertDataTB.MAX_TIME_SEC_FOR_UNSYNC, constants.DBDefaultUnsyncBlockToleranceTime)
}

// Test to insert channel when nodeType is goloop.
func TestInsertConfigurationChannelTBgoloop(t *testing.T) {
	t.Skip("Skipping test InsertConfigurationChannel. Should prepare the goloop to test.")
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertConfigurationChannel("3")

	var configuration_channel_ch_tb []CONFIGURATION_CHANNEL_TB
	configuration_channel_ch_tb = GetConfigurationChannelTable()

	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{}
	configurationAlertDataTB = GetAlertConfigInfoByPK(configuration_channel_ch_tb[2].CHANNEL_PK)

	//assert.Equal(t, configuration_channel_ch_tb[1].CHANNEL_PK, "PKID_0000000000000020")
	assert.Equal(t, configuration_channel_ch_tb[2].CHANNEL_NAME, "3")
	assert.Equal(t, configuration_channel_ch_tb[2].CHANNEL_ID, "0x53cdd6")
	assert.Equal(t, configurationAlertDataTB.MAX_TIME_SEC_FOR_UNSYNC, constants.DBDefaultUnsyncBlockToleranceTime)
}

func TestInsertConfigurationChannelTBWithMapping(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertConfigurationChannelTBWithMapping("ch3", "PKND_0000000000000001", "PKND_0000000000000002")

	var configuration_channel_ch_tb []CONFIGURATION_CHANNEL_TB
	configuration_channel_ch_tb = GetConfigurationChannelTable()

	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{}
	configurationAlertDataTB = GetAlertConfigInfoByPK(configuration_channel_ch_tb[2].CHANNEL_PK)

	nodeChannelMappingTB := GetChannelPermissionNodes(configuration_channel_ch_tb[2].CHANNEL_PK)

	//assert.Equal(t, configuration_channel_ch_tb[1].CHANNEL_PK, "PKID_0000000000000020")
	assert.Equal(t, configuration_channel_ch_tb[2].CHANNEL_NAME, "ch3")
	assert.Equal(t, configurationAlertDataTB.MAX_TIME_SEC_FOR_UNSYNC, constants.DBDefaultUnsyncBlockToleranceTime)
	assert.Equal(t, nodeChannelMappingTB[0].NODE_PK, "PKND_0000000000000001")
	assert.Equal(t, nodeChannelMappingTB[1].NODE_PK, "PKND_0000000000000002")
}

func TestInsertConfigurationNodeTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertConfigurationNode("node3", "127.0.0.2")

	var configurationNodeTB []CONFIGURATION_NODE_TB
	configurationNodeTB = GetConfigurationNodeTable()

	assert.Equal(t, configurationNodeTB[3].NODE_NAME, "node3")
}

func TestInsertChannelPermissionNodesTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertChannelPermissionNode("PKCH_0000000000000001", "PKND_0000000000000003")

	var nodeChannelMappingTB []NODE_CHANNEL_MAPPING_TB
	nodeChannelMappingTB = GetChannelPermissionNodes("PKCH_0000000000000001")

	assert.Equal(t, nodeChannelMappingTB[0].NODE_PK, "PKND_0000000000000003")
}

func TestDeleteUserInfoTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	channel := []string{"PKCH_0000000000000001", "PKCH_0000000000000002"}
	permission := []string{constants.DBUserPermissionNode, constants.DBUserPermissionMonitoringLog}

	InsertUserInfo("test3", "test3", "test3", "password", "test2@test.com", "010-2222-2222", constants.DBUserTypeCodeCommon, channel, permission)

	userInfoTB := &USER_INFO_TB{}
	userInfoTB = GetUserInfoByUserID("test3")

	user_info_ch_tb := &USER_INFO_TB{}
	user_info_ch_tb = GetUserInfoByPK(userInfoTB.PK)

	DeleteUserInfo(user_info_ch_tb.PK)

	user_info_d_tb := &USER_INFO_TB{}
	user_info_d_tb = GetUserInfoByUserID("test3")

	assert.Equal(t, user_info_d_tb.PK, "")
}

func TestDeleteConfigurationChannelTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertConfigurationChannel("ch2")

	var configuration_channel_ch_tb []CONFIGURATION_CHANNEL_TB
	configuration_channel_ch_tb = GetConfigurationChannelTable()

	DeleteConfigurationChannel(configuration_channel_ch_tb[1].CHANNEL_PK)

	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{}
	configurationChannelTB = GetConfigurationChannelInfo(configuration_channel_ch_tb[1].CHANNEL_PK)

	assert.Equal(t, configurationChannelTB.CHANNEL_PK, "")
}

func TestDeleteNodeInfoTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertConfigurationNode("node03", "127.0.0.2")
	var configurationNodeTB []CONFIGURATION_NODE_TB
	configurationNodeTB = GetConfigurationNodeTable()

	DeleteConfigurationNodeInfo(configurationNodeTB[2].NODE_PK)

	configuration_node_ch_tb := &CONFIGURATION_NODE_TB{}
	configuration_node_ch_tb = GetConfigurationNodeInfoByNodePK(configurationNodeTB[2].NODE_PK)

	assert.Equal(t, configuration_node_ch_tb.NODE_NAME, "")
}

func TestUpdateConfigurationChannelInfoTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertConfigurationChannel("ch3")

	var configuration_channel_ch_tb []CONFIGURATION_CHANNEL_TB
	configuration_channel_ch_tb = GetConfigurationChannelTable()

	UpdateConfigurationChannelInfo(configuration_channel_ch_tb[2].CHANNEL_PK, "ch4")

	configuration_channel_ch1_tb := &CONFIGURATION_CHANNEL_TB{}
	configuration_channel_ch1_tb = GetConfigurationChannelInfo(configuration_channel_ch_tb[2].CHANNEL_PK)

	assert.Equal(t, configuration_channel_ch1_tb.CHANNEL_NAME, "ch4")
}

func TestUpdateUser(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	changeUserID := "test10"
	channelTB := GetConfigurationChannelTable()
	permissionTB := GetPermissionTable()
	channel := []string{channelTB[0].CHANNEL_PK}
	permission := []string{permissionTB[0].PERMISSION_ALIAS}

	UpdateUser("PKID_0000000000000001", changeUserID, "test1First", "test1Last", "qwer1234", " test@test.com", "010-111-1111", constants.DBUserTypeCodeCommon, channel, permission)

	userTB := GetUserInfoByPK("PKID_0000000000000001")
	channelMappingTB := GetUserPermissionChannels("PKID_0000000000000001")
	permissionMappingTB := GetUserPermissions("PKID_0000000000000001")
	var permissionList []string
	for _, value := range permissionMappingTB {
		if value.PERMISSION_CHECK {
			permissionList = append(permissionList, value.PERMISSION_ALIAS)
		}
	}

	assert.Equal(t, changeUserID, userTB.USER_ID)
	assert.Equal(t, 1, len(channelMappingTB))
	assert.Equal(t, channel[0], channelMappingTB[0].CHANNEL_PK)
	assert.Equal(t, 1, len(permissionList))
	assert.Equal(t, permission[0], permissionList[0])
}

func TestUpdateUserPermissionChannelsTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	channel := []string{"PKCH_0000000000000001", "PKCH_0000000000000002"}
	permission := []string{constants.DBUserPermissionNode, constants.DBUserPermissionMonitoringLog}

	InsertUserInfo("test2", "test2", "test3", "password", "test2@test.com", "010-2222-2222", constants.DBUserTypeCodeCommon, channel, permission)
	userInfoTB := &USER_INFO_TB{}
	userInfoTB = GetUserInfoByUserID("test2")

	InsertUserPermissionChannels(userInfoTB.PK, "PKCH_0000000000000003", "PKCH_0000000000000004")

	UpdateUserPermissionChannels(userInfoTB.PK, "PKCH_0000000000000005", "PKCH_0000000000000006")

	var channelUserMappingTB []CHANNEL_USER_MAPPING_TB
	channelUserMappingTB = GetUserPermissionChannels(userInfoTB.PK)

	assert.Equal(t, channelUserMappingTB[0].CHANNEL_PK, "PKCH_0000000000000005")
	assert.Equal(t, channelUserMappingTB[1].CHANNEL_PK, "PKCH_0000000000000006")
}

func TestUpdateUserPermissions(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	channel := []string{"PKCH_0000000000000001", "PKCH_0000000000000002"}
	permission := []string{constants.DBUserPermissionNode, constants.DBUserPermissionMonitoringLog}

	InsertUserInfo("test2", "test2", "test2", "password", "test2@test.com", "010-2222-2222", constants.DBUserTypeCodeCommon, channel, permission)
	userInfoTB := &USER_INFO_TB{}
	userInfoTB = GetUserInfoByUserID("test2")

	permission2 := []string{constants.DBUserPermissionNode}
	UpdateUserPermissions(userInfoTB.PK, permission2)

	var permissionUserMappingTB []PERMISSION_USER_MAPPING_TB
	permissionUserMappingTB = GetUserPermissions(userInfoTB.PK)

	assert.Equal(t, permissionUserMappingTB[0].PERMISSION_ALIAS, constants.DBUserPermissionMonitoringLog)
	assert.Equal(t, permissionUserMappingTB[0].PERMISSION_CHECK, false)
	assert.Equal(t, permissionUserMappingTB[1].PERMISSION_ALIAS, constants.DBUserPermissionNode)
	assert.Equal(t, permissionUserMappingTB[1].PERMISSION_CHECK, true)
}

func TestUpdateChannelPermissionNodeTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	channel := []string{"PKCH_0000000000000001", "PKCH_0000000000000002"}
	permission := []string{constants.DBUserPermissionNode, constants.DBUserPermissionMonitoringLog}

	InsertUserInfo("test2", "test2", "test2", "password", "test2@test.com", "010-2222-2222", constants.DBUserTypeCodeCommon, channel, permission)
	var channel_info_tb []CONFIGURATION_CHANNEL_TB
	channel_info_tb = GetConfigurationChannelTable()

	UpdateChannelPermissionNode(channel_info_tb[1].CHANNEL_PK, "PKND_0000000000000003", "PKND_0000000000000004")

	var nodeChannelMappingTB []NODE_CHANNEL_MAPPING_TB
	nodeChannelMappingTB = GetChannelPermissionNodes(channel_info_tb[1].CHANNEL_PK)

	assert.Equal(t, nodeChannelMappingTB[0].NODE_PK, "PKND_0000000000000003")
	assert.Equal(t, nodeChannelMappingTB[1].NODE_PK, "PKND_0000000000000004")
}

func TestUpdateConfigurationNodeInfoTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertConfigurationNode("node03", "127.0.0.2")

	var configurationNodeTB []CONFIGURATION_NODE_TB
	configurationNodeTB = GetConfigurationNodeTable()

	UpdateConfigurationNodeInfo(configurationNodeTB[2].NODE_PK, "node05", "127.0.0.0")

	configuration_node_ch_tb := &CONFIGURATION_NODE_TB{}
	configuration_node_ch_tb = GetConfigurationNodeInfoByNodePK(configurationNodeTB[2].NODE_PK)

	assert.Equal(t, configuration_node_ch_tb.NODE_NAME, "node05")
	assert.Equal(t, configuration_node_ch_tb.NODE_IP, "127.0.0.0")
}

func TestUpdateConfigurationAlertTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var configurationAlertDataTB []CONFIGURATION_DATA_ALERT_TB
	configurationAlertDataTB = GetAlertConfigTable()

	UpdateConfigurationAlertInfo(configurationAlertDataTB[0].CHANNEL_PK, 1, 60, 60)

	configuration_alert_ch_data_tb := &CONFIGURATION_DATA_ALERT_TB{}
	configuration_alert_ch_data_tb = GetAlertConfigInfoByPK(configurationAlertDataTB[0].CHANNEL_PK)

	assert.Equal(t, configuration_alert_ch_data_tb.MAX_TIME_SEC_FOR_RESPONSE, 60)
}

func TestUpdateConfigurationVisibilityTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	var configurationVisibilityDataTB []CONFIGURATION_DATA_VISIBILITY_TB
	configurationVisibilityDataTB = GetVisibilityConfigTable()

	UpdateConfigurationVisibilityInfo(configurationVisibilityDataTB[0].CHANNEL_PK, true, true, true, false, false, false)

	configuration_visibility_ch_data_tb := &CONFIGURATION_DATA_VISIBILITY_TB{}
	configuration_visibility_ch_data_tb = GetVisibilityConfigInfo(configurationVisibilityDataTB[0].CHANNEL_PK)

	assert.Equal(t, configuration_visibility_ch_data_tb.CHECK_NODE_IP, false)
}

func TestUpdateConfigurationMasterTB(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	UpdateConfigurationMasterInfo(60, "English")

	configurationMasterTB := &CONFIGURATION_MASTER_TB{}
	configurationMasterTB = GetConfigurationMasterTable()

	assert.Equal(t, configurationMasterTB.SET_LANGUAGE, "English")
}

func TestTransactionUpdateSuccessProcess(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	tx := NewTransaction()

	configurationMasterTB := &CONFIGURATION_MASTER_TB{}
	tx.db.Model(&configurationMasterTB).Updates(map[string]interface{}{
		"SESSION_TIMEOUT": 60, "SET_LANGUAGE": "ko"})

	tx.db.Model(&configurationMasterTB).Updates(map[string]interface{}{
		"SESSION_TIMEOUT": 120, "SET_LANGUAGE": "en"})

	//success transaction -> Commit
	tx.db.Commit()

	cofiguration_master_ch_tb := &CONFIGURATION_MASTER_TB{}
	cofiguration_master_ch_tb = GetConfigurationMasterTable()

	assert.Equal(t, cofiguration_master_ch_tb.SET_LANGUAGE, "en")
	assert.Equal(t, cofiguration_master_ch_tb.SESSION_TIMEOUT, 120)
}

func TestTransactionUpdateFailProcess(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	tx := NewTransaction()

	configurationMasterTB := &CONFIGURATION_MASTER_TB{}
	tx.db.Model(&configurationMasterTB).Updates(map[string]interface{}{
		"SESSION_TIMEOUT": 60, "SET_LANGUAGE": "en"})

	//failed transaction -> Rollback
	tx.db.Rollback()

	cofiguration_master_ch_tb := &CONFIGURATION_MASTER_TB{}
	cofiguration_master_ch_tb = GetConfigurationMasterTable()

	assert.Equal(t, cofiguration_master_ch_tb.SET_LANGUAGE, "ko")
}

func TestTransactionDeleteSuccessProcess(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertConfigurationChannel("ch2")

	var configuration_channel_ch_tb []CONFIGURATION_CHANNEL_TB
	configuration_channel_ch_tb = GetConfigurationChannelTable()

	tx := NewTransaction()

	tx.db.Delete(CONFIGURATION_CHANNEL_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[2].CHANNEL_PK)
	tx.db.Delete(CONFIGURATION_DATA_ALERT_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[2].CHANNEL_PK)
	tx.db.Delete(CONFIGURATION_DATA_VISIBILITY_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[2].CHANNEL_PK)
	tx.db.Delete(NODE_CHANNEL_MAPPING_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[2].CHANNEL_PK)
	tx.db.Delete(CHANNEL_USER_MAPPING_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[2].CHANNEL_PK)
	tx.db.Delete(CHANNEL_GROUP_MAPPING_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[2].CHANNEL_PK)

	//success transaction -> Commit
	tx.db.Commit()

	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{}
	configurationChannelTB = GetConfigurationChannelInfo(configuration_channel_ch_tb[2].CHANNEL_PK)

	assert.Equal(t, configurationChannelTB.CHANNEL_NAME, "")
}

func TestTransactionDeleteFailProcess(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	InsertConfigurationChannel("ch2")

	var configuration_channel_ch_tb []CONFIGURATION_CHANNEL_TB
	configuration_channel_ch_tb = GetConfigurationChannelTable()

	tx := NewTransaction()

	tx.db.Delete(CONFIGURATION_CHANNEL_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[1].CHANNEL_PK)
	tx.db.Delete(CONFIGURATION_DATA_ALERT_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[1].CHANNEL_PK)
	tx.db.Delete(CONFIGURATION_DATA_VISIBILITY_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[1].CHANNEL_PK)
	tx.db.Delete(NODE_CHANNEL_MAPPING_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[1].CHANNEL_PK)
	tx.db.Delete(CHANNEL_USER_MAPPING_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[1].CHANNEL_PK)
	tx.db.Delete(CHANNEL_GROUP_MAPPING_TB{}, "CHANNEL_PK = ?", configuration_channel_ch_tb[1].CHANNEL_PK)

	//failed transaction -> Rollback
	tx.db.Rollback()

	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{}
	configurationChannelTB = GetConfigurationChannelInfo(configuration_channel_ch_tb[1].CHANNEL_PK)

	assert.Equal(t, configurationChannelTB.CHANNEL_NAME, "loopchain_default2")
}

func TestTransactionInsertSuccessProcess(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	tx := NewTransaction()

	var pk = "PKCH_" + "0000000000000011"
	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{
		CHANNEL_PK:   pk,
		CHANNEL_NAME: "ch3",
	}
	tx.db.Create(&configurationChannelTB)

	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{
		ALERT_METHOD:              1,
		ALERT_LEVEL:               "MAJOR",
		MAX_TIME_SEC_FOR_UNSYNC:   10,
		MAX_TIME_SEC_FOR_RESPONSE: 5,
		CHANNEL_PK:                pk,
	}
	tx.db.Create(&configurationAlertDataTB)

	configurationNotiAlertVisibilityTB := &CONFIGURATION_DATA_VISIBILITY_TB{
		CHECK_HOST_NAME:     true,
		CHECK_BLOCK_HEIGHT:  true,
		CHECK_RESPONSE_TIME: true,
		CHECK_NODE_IP:       true,
		CHECK_TRANSACTION:   true,
		CHECK_LEADER:        true,
		CHANNEL_PK:          pk,
	}
	tx.db.Create(&configurationNotiAlertVisibilityTB)

	//success transaction -> Commit
	tx.db.Commit()

	var configuration_channel_ch_tb []CONFIGURATION_CHANNEL_TB
	configuration_channel_ch_tb = GetConfigurationChannelTable()

	configuration_noti_ch_data_tb := &CONFIGURATION_DATA_ALERT_TB{}
	configuration_noti_ch_data_tb = GetAlertConfigInfoByPK(configuration_channel_ch_tb[2].CHANNEL_PK)

	assert.Equal(t, configuration_channel_ch_tb[2].CHANNEL_NAME, "ch3")
	assert.Equal(t, configuration_noti_ch_data_tb.MAX_TIME_SEC_FOR_UNSYNC, 10)
}

func TestTransactionInsertFailProcess(t *testing.T) {
	dbpath := ":memory:"
	Setup(dbpath)
	defer Teardown(dbpath)

	tx := NewTransaction()

	var pk = "PKCH_" + "0000000000000011"
	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{
		CHANNEL_PK:   pk,
		CHANNEL_NAME: "ch3",
	}
	tx.db.Create(&configurationChannelTB)

	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{
		ALERT_METHOD:              1,
		ALERT_LEVEL:               "MAJOR",
		MAX_TIME_SEC_FOR_UNSYNC:   10,
		MAX_TIME_SEC_FOR_RESPONSE: 5,
		CHANNEL_PK:                pk,
	}
	tx.db.Create(&configurationAlertDataTB)

	configurationNotiAlertVisibilityTB := &CONFIGURATION_DATA_VISIBILITY_TB{
		CHECK_HOST_NAME:     true,
		CHECK_BLOCK_HEIGHT:  true,
		CHECK_RESPONSE_TIME: true,
		CHECK_NODE_IP:       true,
		CHECK_TRANSACTION:   true,
		CHECK_LEADER:        true,
		CHANNEL_PK:          pk,
	}
	tx.db.Create(&configurationNotiAlertVisibilityTB)

	//success transaction -> Commit
	tx.db.Rollback()

	var count int
	DBgorm().Model(&CONFIGURATION_CHANNEL_TB{}).Where("CHANNEL_PK = ?", "PKCH_0000000000000011").Count(&count)

	assert.Equal(t, count, 0)
}

func TestBcryptHash(t *testing.T) {
	password := "admin123"
	hash, _ := HashPassword(password)

	match := CheckPasswordHash(password, "$2a$10$XToqmLMc7XTwwhuOokpFaOBpqUBEMQuoEyS1b7EU5eQES4ijBGxIu")
	match2 := CheckPasswordHash(password, hash)

	assert.Equal(t, match, true)
	assert.Equal(t, match2, true)
}
