package db

import "motherbear/backend/logger"

// DeleteUserInfo delete user info and delete it other related tables.
func DeleteUserInfo(pk string) error {
	tx := NewTransaction()
	defer tx.Close()

	if err := tx.db.Delete(USER_INFO_TB{}, "PK = ?", pk).Error; err != nil {
		logger.Error("DeleteUserInfo, USER_INFO_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
		return err
	}
	if err := tx.db.Delete(CHANNEL_USER_MAPPING_TB{}, "PK = ?", pk).Error; err != nil {
		logger.Error("DeleteUserInfo, CHANNEL_USER_MAPPING_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
		return err
	}
	if err := tx.db.Delete(PERMISSION_USER_MAPPING_TB{}, "PK = ?", pk).Error; err != nil {
		logger.Error("DeleteUserInfo, PERMISSION_USER_MAPPING_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
		return err
	}

	return nil
}

// DeleteConfigurationChannel delete channel info and delete it other related tables.
func DeleteConfigurationChannel(pk string) error {
	tx := NewTransaction()
	defer tx.Close()

	if err := tx.db.Delete(CONFIGURATION_CHANNEL_TB{}, "CHANNEL_PK = ?", pk).Error; err != nil {
		logger.Error("DeleteConfigurationChannel, CONFIGURATION_CHANNEL_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
		return err
	}
	if err := tx.db.Delete(CONFIGURATION_DATA_ALERT_TB{}, "CHANNEL_PK = ?", pk).Error; err != nil {
		logger.Error("DeleteConfigurationChannel, CONFIGURATION_DATA_ALERT_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
		return err
	}
	if err := tx.db.Delete(CONFIGURATION_DATA_VISIBILITY_TB{}, "CHANNEL_PK = ?", pk).Error; err != nil {
		logger.Error("DeleteConfigurationChannel, CONFIGURATION_DATA_VISIBILITY_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
		return err
	}
	if err := tx.db.Delete(NODE_CHANNEL_MAPPING_TB{}, "CHANNEL_PK = ?", pk).Error; err != nil {
		logger.Error("DeleteConfigurationChannel, NODE_CHANNEL_MAPPING_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
		return err
	}
	if err := tx.db.Delete(CHANNEL_USER_MAPPING_TB{}, "CHANNEL_PK = ?", pk).Error; err != nil {
		logger.Error("DeleteConfigurationChannel, CHANNEL_USER_MAPPING_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
		return err
	}
	if err := tx.db.Delete(CHANNEL_GROUP_MAPPING_TB{}, "CHANNEL_PK = ?", pk).Error; err != nil {
		logger.Error("DeleteConfigurationChannel, CHANNEL_GROUP_MAPPING_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()
		return err
	}

	return nil
}

// DeleteConfigurationNodeInfo delete node info and delete it other related tables.
func DeleteConfigurationNodeInfo(pk string) error {
	tx := NewTransaction()
	defer tx.Close()

	if err := tx.db.Delete(CONFIGURATION_NODE_TB{}, "NODE_PK = ?", pk).Error; err != nil {
		logger.Error("DeleteConfigurationNodeInfo, CONFIGURATION_NODE_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()

		return err
	}
	if err := tx.db.Delete(NODE_CHANNEL_MAPPING_TB{}, "NODE_PK = ?", pk).Error; err != nil {
		logger.Error("DeleteConfigurationNodeInfo, NODE_CHANNEL_MAPPING_TB delete failed!")
		logger.Errorf("%v+", err)
		tx.Fail()

		return err
	}

	return nil
}
