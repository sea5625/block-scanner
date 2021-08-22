package db

import (
	"motherbear/backend/isaacerror"
	. "motherbear/backend/logger"
	"motherbear/backend/utility"
	"time"
)

// UpdateConfigurationChannelInfo update channel name and update it other related tables.
func UpdateConfigurationChannelInfo(pk string, channelName string) error {
	tx := NewTransaction()
	defer tx.Close()

	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{}
	if err := tx.db.Where("CHANNEL_PK = ?", pk).First(&configurationChannelTB).Error; err != nil {
		Logger().Error("UpdateConfigurationChannelInfo, CONFIGURATION_CHANNEL_TB Where failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
		return err
	}

	if err := tx.db.Model(&configurationChannelTB).Update("CHANNEL_NAME", channelName).Error; err != nil {
		Logger().Error("UpdateConfigurationChannelInfo, CONFIGURATION_CHANNEL_TB Update failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
		return err
	}

	return nil
}

// UpdateChannelPermissionNode update channel permission nodes and update it other related tables.
func UpdateChannelPermissionNode(channelPKID string, nodePKID ...string) error {
	tx := NewTransaction()
	defer tx.Close()

	if err := tx.db.Delete(NODE_CHANNEL_MAPPING_TB{}, "CHANNEL_PK = ?", channelPKID).Error; err != nil {
		Logger().Error("UpdateChannelPermissionNode, NODE_CHANNEL_MAPPING_TB delete failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
		return err
	}

	chNodeMappingTB := &NODE_CHANNEL_MAPPING_TB{}
	for _, value := range nodePKID {
		chNodeMappingTB = &NODE_CHANNEL_MAPPING_TB{
			CHANNEL_PK: channelPKID,
			NODE_PK:    value,
		}
		if err := tx.db.Create(&chNodeMappingTB).Error; err != nil {
			Logger().Error("UpdateChannelPermissionNode, NODE_CHANNEL_MAPPING_TB Create failed!")
			Logger().Errorf("%v+", err)
			tx.Fail()
			return err
		}
	}

	return nil
}

// UpdateConfigurationChannelInfo update channel name and update it other related tables and node-channel mapping data.
func UpdateConfigurationChannelInfoWithMapping(pk string, channelName string, nodePKID ...string) error {
	tx := NewTransaction()
	defer tx.Close()

	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{}
	if err := tx.db.Where("CHANNEL_PK = ?", pk).First(&configurationChannelTB).Error; err != nil {
		Logger().Error("UpdateConfigurationChannelInfo, CONFIGURATION_CHANNEL_TB Where failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
		return err
	}

	if err := tx.db.Model(&configurationChannelTB).Update("CHANNEL_NAME", channelName).Error; err != nil {
		Logger().Error("UpdateConfigurationChannelInfo, CONFIGURATION_CHANNEL_TB Update failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
		return err
	}

	if err := tx.db.Delete(NODE_CHANNEL_MAPPING_TB{}, "CHANNEL_PK = ?", pk).Error; err != nil {
		Logger().Error("UpdateChannelPermissionNode, NODE_CHANNEL_MAPPING_TB delete failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
		return err
	}

	chNodeMappingTB := &NODE_CHANNEL_MAPPING_TB{}
	for _, value := range nodePKID {
		chNodeMappingTB = &NODE_CHANNEL_MAPPING_TB{
			CHANNEL_PK: pk,
			NODE_PK:    value,
		}
		if err := tx.db.Create(&chNodeMappingTB).Error; err != nil {
			Logger().Error("UpdateChannelPermissionNode, NODE_CHANNEL_MAPPING_TB Create failed!")
			Logger().Errorf("%v+", err)
			tx.Fail()
			return err
		}
	}

	return nil
}

func UpdateUser(pk string, userID string, firstName string, lastName string, password string, emailAdres string, mobilephoneNo string, typeCode string, channelPKID []string, permission []string) error {

	// Check if valid channel.
	if !IsExistChannelListInDB(channelPKID) {
		err := isaacerror.SysErrSelectedChannelThatNotExistInDB
		Logger().Error("UpdateUser, " + err.Error())
		Logger().Errorf("%v+", err)
		return err
	}

	// Check if valid permission.
	if !IsExistPermissionListInDB(permission) {
		err := isaacerror.SysErrSelectedPermissionThatNotExistInDB
		Logger().Error("UpdateUser, " + err.Error())
		Logger().Errorf("%v+", err)
		return err
	}

	channelMappingTB := GetUserPermissionChannels(pk)
	permissionMappingTB := GetUserPermissions(pk)

	tx := NewTransaction()
	defer tx.Close()

	// Update user table.
	userTB := &USER_INFO_TB{}
	if err := tx.db.Where("PK = ?", pk).First(&userTB).Error; err != nil {
		Logger().Error("UpdateUser, USER_INFO_TB Where failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
		return err
	}

	userTB.USER_ID = userID
	userTB.FIRST_NAME = firstName
	userTB.LAST_NAME = lastName
	userTB.EMAIL_ADRES = emailAdres
	userTB.MOBILE_PHONE_NO = mobilephoneNo
	// Hash password.
	if password != "" {
		hash, _ := HashPassword(password)
		userTB.PASSWORD = hash
		userTB.TIME_WHEN_USER_MODIFIED_PASSWORD = time.Now()
	}

	if err := tx.db.Model(&userTB).Updates(&userTB).Error; err != nil {
		Logger().Error("UpdateUser, USER_INFO_TB Update failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
		return err
	}

	// Update channel user mapping table when channel mapping data is changed.
	var channelListInDB []string
	for _, value := range channelMappingTB {
		channelListInDB = append(channelListInDB, value.CHANNEL_PK)
	}
	if !utility.IsIdenticalSlice(channelPKID, channelListInDB) {
		if err := tx.db.Delete(CHANNEL_USER_MAPPING_TB{}, "PK = ?", pk).Error; err != nil {
			Logger().Error("UpdateUser, CHANNEL_USER_MAPPING_TB delete failed!")
			Logger().Errorf("%v+", err)
			tx.Fail()
		}

		for _, value := range channelPKID {
			chUserMappingTB := &CHANNEL_USER_MAPPING_TB{
				PK:         pk,
				CHANNEL_PK: value,
			}
			if err := tx.db.Create(&chUserMappingTB).Error; err != nil {
				Logger().Error("UpdateUser, CHANNEL_USER_MAPPING_TB insert failed!")
				Logger().Errorf("%v+", err)
				tx.Fail()
				return err
			}
		}
	}

	// Update permission user mapping table when permission mapping data is changed.
	for _, value := range permissionMappingTB {
		ckPermission := false

		if utility.IsExistValueInList(value.PERMISSION_ALIAS, permission) {
			ckPermission = true
		}

		// Update permission user mapping column when DB data and input permission data are different.
		if ckPermission != value.PERMISSION_CHECK {
			if err := tx.db.Model(&value).Update("PERMISSION_CHECK", ckPermission).Error; err != nil {
				Logger().Error("UpdateUser, PERMISSION_USER_MAPPING_TB Update failed!")
				Logger().Errorf("%v+", err)
				tx.Fail()
				return err
			}
			break
		}
	}

	// Create permission data if permission data is not in DB.
	var permissionListInDB []string
	for _, value := range permissionMappingTB {
		permissionListInDB = append(permissionListInDB, value.PERMISSION_ALIAS)
	}
	for _, value := range permission {
		if !utility.IsExistValueInList(value, permissionListInDB) {
			permissionUserMappingTB := &PERMISSION_USER_MAPPING_TB{
				PK:               pk,
				PERMISSION_ALIAS: value,
				PERMISSION_CHECK: true,
			}

			if err := tx.db.Create(&permissionUserMappingTB).Error; err != nil {
				Logger().Error("UpdateUser, PERMISSION_USER_MAPPING_TB insert failed!")
				Logger().Errorf("%v+", err)
				tx.Fail()
				return err
			}
		}
	}

	return nil
}

// UpdateUserPermissionChannels update user permission channel and update it other related tables.
func UpdateUserPermissionChannels(pk string, chPKID ...string) {
	tx := NewTransaction()
	defer tx.Close()

	if err := tx.db.Delete(CHANNEL_USER_MAPPING_TB{}, "PK = ?", pk).Error; err != nil {
		Logger().Error("UpdateUserPermissionChannels, CHANNEL_USER_MAPPING_TB delete failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
	}

	chUserMappingTB := &CHANNEL_USER_MAPPING_TB{}
	for _, value := range chPKID {
		chUserMappingTB = &CHANNEL_USER_MAPPING_TB{
			PK:         pk,
			CHANNEL_PK: value,
		}
		if err := tx.db.Create(&chUserMappingTB).Error; err != nil {
			Logger().Error("UpdateUserPermissionChannels, CHANNEL_USER_MAPPING_TB Create failed!")
			Logger().Errorf("%v+", err)
			tx.Fail()
		}
	}
}

// UpdateUserPermissionEach update user permission and update it other related tables.
func UpdateUserPermissionEach(pk string, permission string, permissionCheck bool) {
	tx := NewTransaction()
	defer tx.Close()

	var permissionUserMappingTB []PERMISSION_USER_MAPPING_TB
	permissionUserMappingTB = GetUserPermissions(pk)

	if err := tx.db.Where("PK = ?", pk).Find(&permissionUserMappingTB).Error; err != nil {
		Logger().Error("UpdateUserPermissionEach, PERMISSION_USER_MAPPING_TB Where failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
	}

	if err := tx.db.Model(&permissionUserMappingTB).Where("PERMISSION_ALIAS = ?", permission).Update("PERMISSION_CHECK", permissionCheck).Error; err != nil {
		Logger().Error("UpdateUserPermissionEach, PERMISSION_USER_MAPPING_TB Update failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
	}
}

// UpdateUserPermissions update user permission and update it other related tables.
func UpdateUserPermissions(pk string, permission []string) {
	tx := NewTransaction()
	defer tx.Close()

	if err := tx.db.Delete(PERMISSION_USER_MAPPING_TB{}, "PK = ?", pk).Error; err != nil {
		Logger().Error("UpdateUserPermissions, PERMISSION_USER_MAPPING_TB delete failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
	}

	var permissionInfoTB []PERMISSION_INFO_TB
	if err := tx.db.Find(&permissionInfoTB).Error; err != nil {
		Logger().Error("UpdateUserPermissions, PERMISSION_INFO_TB Find failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
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
			Logger().Error("UpdateUserPermissions, PERMISSION_USER_MAPPING_TB insert failed!")
			Logger().Errorf("%v+", err)
			tx.Fail()
		}

	}
}

// UpdateConfigurationNodeInfo update node name, ip and update it other related tables.
func UpdateConfigurationNodeInfo(pk string, nodeName string, nodeIP string) error {
	tx := NewTransaction()
	defer tx.Close()

	configurationNodeTB := &CONFIGURATION_NODE_TB{}
	if err := tx.db.Where("NODE_PK = ?", pk).First(&configurationNodeTB).Error; err != nil {
		Logger().Error("UpdateConfigurationNodeInfo, CONFIGURATION_NODE_TB Where failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()

		return err
	}

	if err := tx.db.Model(&configurationNodeTB).Updates(map[string]interface{}{
		"NODE_NAME": nodeName, "NODE_IP": nodeIP}).Error; err != nil {
		Logger().Error("UpdateConfigurationNodeInfo, CONFIGURATION_NODE_TB Update failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()

		return err
	}

	return nil
}

// UpdateConfigurationAlertInfo update alarm alert field and update it other related tables.
func UpdateConfigurationAlertInfo(pk string, notiMethod int, maxTimeUnsync int, maxTimeSlowRes int) {
	tx := NewTransaction()
	defer tx.Close()

	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{}
	if err := tx.db.Where("CHANNEL_PK = ?", pk).First(&configurationAlertDataTB).Error; err != nil {
		Logger().Error("UpdateConfigurationAlertInfo, CONFIGURATION_DATA_ALERT_TB Where failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
	}

	if err := tx.db.Model(&configurationAlertDataTB).Updates(map[string]interface{}{
		"NOTI_METHOD": notiMethod, "NOTI_LEVEL": "Major", "MAX_TIME_SEC_FOR_UNSYNC": maxTimeUnsync,
		"MAX_TIME_SEC_FOR_RESPONSE": maxTimeSlowRes}).Error; err != nil {
		Logger().Error("UpdateConfigurationAlertInfo, CONFIGURATION_DATA_ALERT_TB Update failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
	}
}

// UpdateConfigurationVisibilityInfo update peer visibility field and update it other related tables.
func UpdateConfigurationVisibilityInfo(pk string, checkHost bool, checkBH bool, checkRT bool,
	checkIP bool, checkTX bool, checkLeader bool) {
	tx := NewTransaction()
	defer tx.Close()

	configurationVisibilityDataTB := &CONFIGURATION_DATA_VISIBILITY_TB{}
	if err := tx.db.Where("CHANNEL_PK = ?", pk).First(&configurationVisibilityDataTB).Error; err != nil {
		Logger().Error("UpdateConfigurationVisibilityInfo, CONFIGURATION_DATA_VISIBILITY_TB Where failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
	}

	if err := tx.db.Model(&configurationVisibilityDataTB).Updates(map[string]interface{}{
		"CHECK_HOST_NAME": checkHost, "CHECK_BLOCK_HEIGHT": checkBH, "CHECK_RESPONSE_TIME": checkRT,
		"CHECK_NODE_IP": checkIP, "CHECK_TRANSACTION": checkTX, "CHECK_LEADER": checkLeader}).Error; err != nil {
		Logger().Error("UpdateConfigurationVisibilityInfo, CONFIGURATION_DATA_VISIBILITY_TB Update failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
	}
}

// UpdateConfigurationMasterInfo update master configuration and update it other related tables.
func UpdateConfigurationMasterInfo(sessionTimeout int, language string) {
	tx := NewTransaction()
	defer tx.Close()

	configurationMasterTB := &CONFIGURATION_MASTER_TB{}

	if err := tx.db.Model(&configurationMasterTB).Updates(map[string]interface{}{
		"SESSION_TIMEOUT": sessionTimeout, "SET_LANGUAGE": language}).Error; err != nil {
		Logger().Error("UpdateConfigurationMasterInfo, CONFIGURATION_MASTER_TB Update failed!")
		Logger().Errorf("%v+", err)
		tx.Fail()
	}
}
