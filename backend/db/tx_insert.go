package db

import (
	"motherbear/backend/configuration"
	"motherbear/backend/constants"
	"motherbear/backend/goloop/admin"
	"motherbear/backend/isaacerror"
	"motherbear/backend/logger"
	"time"
)

// InsertUserInfo insert user info and put it into other related tables.
func InsertUserInfo(userID string, firstName string, lastName string, password string, emailAdres string, mobilephoneNo string, userType string, channelPKID []string, permission []string) string {
	pk := getPKSource()
	pk = "PKID_" + pk

	tx := NewTransaction()
	defer tx.Close()

	hash, _ := HashPassword(password)

	var userInfoCreateTB USER_INFO_TB

	// Two user types are defined.
	//  - USER_COMMON : Common user.
	//  - THIRD_PARTY : User for 3rd party. Can use RESTFul API, GET /api/v1/channels.
	//   * Will not allow to change password, and other information.

	if userType == constants.DBUserTypeCodeThirdParty {
		userInfoCreateTB = USER_INFO_TB{
			PK:                               pk,
			USER_ID:                          userID,
			FIRST_NAME:                       firstName,
			LAST_NAME:                        lastName,
			PASSWORD:                         hash,
			EMAIL_ADRES:                      emailAdres,
			MOBILE_PHONE_NO:                  mobilephoneNo,
			TYPE_CODE:                        constants.DBUserTypeCodeThirdParty,
			CREATE_DATE:                      time.Now(),
			TIME_WHEN_USER_MODIFIED_PASSWORD: time.Now(),
		}
	} else {
		userInfoCreateTB = USER_INFO_TB{
			PK:              pk,
			USER_ID:         userID,
			FIRST_NAME:      firstName,
			LAST_NAME:       lastName,
			PASSWORD:        hash,
			EMAIL_ADRES:     emailAdres,
			MOBILE_PHONE_NO: mobilephoneNo,
			TYPE_CODE:       constants.DBUserTypeCodeCommon,
			CREATE_DATE:     time.Now(),
		}
	}

	if err := tx.db.Create(&userInfoCreateTB).Error; err != nil {
		logger.Error("InsertUserInfo, USER_INFO_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
		return ""
	}

	// If userType is THIRD_PARTY, do not execute.
	if userType != constants.DBUserTypeCodeThirdParty {
		var permissionInfoTB []PERMISSION_INFO_TB

		if err := tx.db.Find(&permissionInfoTB).Error; err != nil {
			logger.Error("InsertUserInfo, PERMISSION_INFO_TB Find failed!")
			logger.Errorf("%v+", err)
			tx.Fail()
			return ""
		}

		for _, permissionMap := range permissionInfoTB {
			var ckPermission = false
			for _, value := range permission {
				if value == permissionMap.PERMISSION_ALIAS {
					ckPermission = true
					break
				}
			}
			permissionUserMappingTB := &PERMISSION_USER_MAPPING_TB{
				PK:               pk,
				PERMISSION_ALIAS: permissionMap.PERMISSION_ALIAS,
				PERMISSION_CHECK: ckPermission,
			}
			if err := tx.db.Create(&permissionUserMappingTB).Error; err != nil {
				logger.Error("InsertUserInfo, PERMISSION_USER_MAPPING_TB insert failed!")
				logger.Errorf("%v+", err)
				tx.Fail()
				return ""
			}
		}

		for _, value := range channelPKID {
			chUserMappingTB := &CHANNEL_USER_MAPPING_TB{
				PK:         pk,
				CHANNEL_PK: value,
			}
			if err := tx.db.Create(&chUserMappingTB).Error; err != nil {
				logger.Error("InsertUserInfo, CHANNEL_USER_MAPPING_TB insert failed!")
				logger.Errorf("%v+", err)
				tx.Fail()
				return ""
			}
		}
	}

	return pk
}

// InsertUserPermissionChannels insert user permission channel and put it into other related tables.
func InsertUserPermissionChannels(pk string, chPKID ...string) {
	tx := NewTransaction()
	defer tx.Close()

	chUserMappingTB := &CHANNEL_USER_MAPPING_TB{}
	for _, value := range chPKID {
		chUserMappingTB = &CHANNEL_USER_MAPPING_TB{
			PK:         pk,
			CHANNEL_PK: value,
		}
		if err := tx.db.Create(&chUserMappingTB).Error; err != nil {
			logger.Error("InsertUserPermissionChannels, CHANNEL_USER_MAPPING_TB insert failed!")
			logger.Errorf("%v+", err)
			tx.Fail()
		}
	}
}

// InsertUserPermissions insert user permission and put it into other related tables.
func InsertUserPermissions(pk string, permissionAlias string, permissionCheck bool) {
	tx := NewTransaction()
	defer tx.Close()

	permissionUserMappingTB := &PERMISSION_USER_MAPPING_TB{
		PK:               pk,
		PERMISSION_ALIAS: permissionAlias,
		PERMISSION_CHECK: permissionCheck,
	}
	if err := tx.db.Create(&permissionUserMappingTB).Error; err != nil {
		logger.Error("InsertUserPermissions, PERMISSION_USER_MAPPING_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
	}
}

// InsertConfigurationChannel insert configuration channel and put it into other related tables.
func InsertConfigurationChannel(channelName string) string {
	pk := getPKSource()
	pk = "PKCH_" + pk

	tx := NewTransaction()
	defer tx.Close()

	idOfGoloop := ""
	// If using the goloop, set channel ID after execute '/chain' admin API on goloop.
	if configuration.QueryNodeType() == constants.NodeType2 {
		nodeIPList := make([]string, 0)

		for _, value := range configuration.Conf().Node {
			nodeIPList = append(nodeIPList, value.IP)
		}

		var err error
		idOfGoloop, err = admin.GetChannelIDByName(channelName, nodeIPList)
		if err != nil {
			logger.Error("InsertConfigurationChannel, CONFIGURATION_CHANNEL_TB insert failed!")
			logger.Errorf("%v+", err)
			tx.Fail()
		}
	}

	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{
		CHANNEL_PK:   pk,
		CHANNEL_NAME: channelName,
		CHANNEL_ID:   idOfGoloop,
	}
	if err := tx.db.Create(&configurationChannelTB).Error; err != nil {
		logger.Error("InsertConfigurationChannel, CONFIGURATION_CHANNEL_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
	}

	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{
		ALERT_METHOD:                1,
		ALERT_LEVEL:                 "MAJOR",
		MAX_TIME_SEC_FOR_UNSYNC:     constants.DBDefaultUnsyncBlockToleranceTime,
		MAX_TIME_SEC_FOR_RESPONSE:   constants.DBDefaultSlowResponseTime,
		MAX_UNSYNC_BLOCK_DIFFERENCE: constants.DBDefaultUnsyncBlockDifference,
		CHANNEL_PK:                  pk,
	}
	if err := tx.db.Create(&configurationAlertDataTB).Error; err != nil {
		logger.Error("InsertConfigurationChannel, CONFIGURATION_DATA_ALERT_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
	}

	configurationNotiAlertVisibilityTB := &CONFIGURATION_DATA_VISIBILITY_TB{
		CHECK_HOST_NAME:     true,
		CHECK_BLOCK_HEIGHT:  true,
		CHECK_RESPONSE_TIME: true,
		CHECK_NODE_IP:       true,
		CHECK_TRANSACTION:   true,
		CHECK_LEADER:        true,
		CHANNEL_PK:          pk,
	}
	if err := tx.db.Create(&configurationNotiAlertVisibilityTB).Error; err != nil {
		logger.Error("InsertConfigurationChannel, CONFIGURATION_DATA_VISIBILITY_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
	}

	userChannelPermissionTB := &CHANNEL_USER_MAPPING_TB{
		PK:         constants.AdminPK,
		CHANNEL_PK: pk,
	}
	if err := tx.db.Create(&userChannelPermissionTB).Error; err != nil {
		logger.Error("InsertConfigurationChannel, CHANNEL_USER_MAPPING_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
	}

	return pk
}

// InsertConfigurationChannel insert configuration channel and put it into other related tables and node-channel mapping table.
func InsertConfigurationChannelTBWithMapping(channelName string, nodePK ...string) string {
	pk := getPKSource()
	pk = "PKCH_" + pk

	tx := NewTransaction()
	defer tx.Close()

	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{
		CHANNEL_PK:   pk,
		CHANNEL_NAME: channelName,
	}
	if err := tx.db.Create(&configurationChannelTB).Error; err != nil {
		logger.Error("InsertConfigurationChannel, CONFIGURATION_CHANNEL_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
	}

	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{
		ALERT_METHOD:              1,
		ALERT_LEVEL:               "MAJOR",
		MAX_TIME_SEC_FOR_UNSYNC:   constants.DBDefaultUnsyncBlockToleranceTime,
		MAX_TIME_SEC_FOR_RESPONSE: constants.DBDefaultSlowResponseTime,
		CHANNEL_PK:                pk,
	}
	if err := tx.db.Create(&configurationAlertDataTB).Error; err != nil {
		logger.Error("InsertConfigurationChannel, CONFIGURATION_DATA_ALERT_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
	}

	configurationNotiAlertVisibilityTB := &CONFIGURATION_DATA_VISIBILITY_TB{
		CHECK_HOST_NAME:     true,
		CHECK_BLOCK_HEIGHT:  true,
		CHECK_RESPONSE_TIME: true,
		CHECK_NODE_IP:       true,
		CHECK_TRANSACTION:   true,
		CHECK_LEADER:        true,
		CHANNEL_PK:          pk,
	}
	if err := tx.db.Create(&configurationNotiAlertVisibilityTB).Error; err != nil {
		logger.Error("InsertConfigurationChannel, CONFIGURATION_DATA_VISIBILITY_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
	}

	for _, value := range nodePK {
		nodeChannelMappingTB := &NODE_CHANNEL_MAPPING_TB{
			NODE_PK:    value,
			CHANNEL_PK: pk,
		}
		if err := tx.db.Create(&nodeChannelMappingTB).Error; err != nil {
			logger.Error("InsertConfigurationChannel, NODE_CHANNEL_MAPPING_TB insert failed!")
			logger.Errorf("%v+", err)
			tx.Fail()
		}
	}

	userChannelPermissionTB := &CHANNEL_USER_MAPPING_TB{
		PK:         constants.AdminPK,
		CHANNEL_PK: pk,
	}
	if err := tx.db.Create(&userChannelPermissionTB).Error; err != nil {
		logger.Error("InsertConfigurationChannel, CHANNEL_USER_MAPPING_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
	}

	return pk
}

// InsertConfigurationNode insert node info and put it into other related tables.
func InsertConfigurationNode(nodeName string, nodeIP string) string {
	pk := getPKSource()
	pk = "PKND_" + pk

	nodeAddress := ""
	if configuration.QueryNodeType() == constants.NodeType2 {
		system, err := admin.GetSystem(nodeIP)
		if err != nil {
			logger.Error("InsertConfigurationNode, CONFIGURATION_NODE_TB insert failed!")
			logger.Errorf("%v+", err)
			return ""
		}
		if len(system.Setting.Address) == 0 {
			err := isaacerror.SysErrFailToGetSystemDataOfGoloop
			logger.Error("InsertConfigurationNode, CONFIGURATION_NODE_TB insert failed!")
			logger.Errorf("%v+", err)
			return ""
		}

		nodeAddress = system.Setting.Address
	}

	tx := NewTransaction()
	defer tx.Close()

	configurationNodeTB := &CONFIGURATION_NODE_TB{
		NODE_NAME:    nodeName,
		NODE_PK:      pk,
		NODE_IP:      nodeIP,
		NODE_ADDRESS: nodeAddress,
	}
	if err := tx.db.Create(&configurationNodeTB).Error; err != nil {
		logger.Error("InsertConfigurationNode, CONFIGURATION_NODE_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()

		return ""
	}

	return pk
}

// InsertChannelPermissionNode insert channel permission nodes and put it into other related tables.
func InsertChannelPermissionNode(channelPKID string, nodePKID string) {
	tx := NewTransaction()
	defer tx.Close()

	nodeChannelMappingTB := &NODE_CHANNEL_MAPPING_TB{
		NODE_PK:    nodePKID,
		CHANNEL_PK: channelPKID,
	}
	if err := tx.db.Create(&nodeChannelMappingTB).Error; err != nil {
		logger.Error("InsertChannelPermissionNode, NODE_CHANNEL_MAPPING_TB insert failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
	}
}
