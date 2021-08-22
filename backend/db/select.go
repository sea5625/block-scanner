package db

// GetUserPK return user PK.
func GetUserPK(userid string) string {
	userInfoTB := &USER_INFO_TB{}
	DBgorm().Where("USER_ID = ?", userid).First(&userInfoTB)

	return userInfoTB.PK
}

// GetUserInfo return user info struct.
func GetUserInfo() []USER_INFO_TB {
	var userInfoTB []USER_INFO_TB
	DBgorm().Find(&userInfoTB)

	return userInfoTB
}

// GetUserInfoByUserID return user info struct.
func GetUserInfoByUserID(userid string) *USER_INFO_TB {
	userInfoTB := &USER_INFO_TB{}
	DBgorm().Where("USER_ID = ?", userid).First(&userInfoTB)

	return userInfoTB
}

// GetUserInfoByPK return user info struct.
func GetUserInfoByPK(pk string) *USER_INFO_TB {
	userInfoTB := &USER_INFO_TB{}
	DBgorm().Where("PK = ?", pk).First(&userInfoTB)

	return userInfoTB
}

// GetUserPermissionChannels return user channel permission struct.
func GetUserPermissionChannels(pk string) []CHANNEL_USER_MAPPING_TB {
	var channelUserMappingTB []CHANNEL_USER_MAPPING_TB
	DBgorm().Where("PK = ?", pk).Find(&channelUserMappingTB)

	return channelUserMappingTB
}

// GetUserPermissionChannels return user channel permission struct.
func GetUserPermissionChannelsTable() []CHANNEL_USER_MAPPING_TB {
	var channelUserMappingTB []CHANNEL_USER_MAPPING_TB
	DBgorm().Find(&channelUserMappingTB)

	return channelUserMappingTB
}

// GetUserPermissions return user permission struct.
func GetUserPermissions(pk string) []PERMISSION_USER_MAPPING_TB {
	var permissionUserMappingTB []PERMISSION_USER_MAPPING_TB
	DBgorm().Where("PK = ?", pk).Find(&permissionUserMappingTB)

	return permissionUserMappingTB
}

// GetUserPermissionTable return user permission struct.
func GetUserPermissionTable() []PERMISSION_USER_MAPPING_TB {
	var permissionUserMappingTB []PERMISSION_USER_MAPPING_TB
	DBgorm().Find(&permissionUserMappingTB)

	return permissionUserMappingTB
}

// GetPermissionInfo return permission table info.
func GetPermissionInfo(index int) *PERMISSION_INFO_TB {
	permissionInfoTB := &PERMISSION_INFO_TB{}
	DBgorm().Where("PERMISSION_INDEX = ?", index).First(&permissionInfoTB)

	return permissionInfoTB
}

// GetPermissionTable return permission table info.
func GetPermissionTable() []PERMISSION_INFO_TB {
	var permissionInfoTB []PERMISSION_INFO_TB
	DBgorm().Find(&permissionInfoTB)

	return permissionInfoTB
}

// GetRoleGroupInfo return group info.
func GetRoleGroupInfo(pk string) *ROLE_GROUP_TB {
	roleGroupTB := &ROLE_GROUP_TB{}
	DBgorm().Where("GROUP_PK = ?", pk).First(&roleGroupTB)

	return roleGroupTB
}

// GetRoleGroupTable return group table struct.
func GetRoleGroupTable() []ROLE_GROUP_TB {
	var roleGroupTB []ROLE_GROUP_TB
	DBgorm().Find(&roleGroupTB)

	return roleGroupTB
}

// GetGroupPermissionChannels return group permission struct.
func GetGroupPermissionChannels(pk string) []CHANNEL_GROUP_MAPPING_TB {
	var channelGroupMappingTB []CHANNEL_GROUP_MAPPING_TB
	DBgorm().Where("GROUP_PK = ?", pk).Find(&channelGroupMappingTB)

	return channelGroupMappingTB
}

// GetGroupPermissionChannels return group permission struct.
func GetGroupPermissionChannelsTable() []CHANNEL_GROUP_MAPPING_TB {
	var channelGroupMappingTB []CHANNEL_GROUP_MAPPING_TB
	DBgorm().Find(&channelGroupMappingTB)

	return channelGroupMappingTB
}

// GetGroupPermissions return group permission struct.
func GetGroupPermissions(pk string) []PERMISSION_GROUP_MAPPING_TB {
	var permissionGroupMappingTB []PERMISSION_GROUP_MAPPING_TB
	DBgorm().Where("GROUP_PK = ?", pk).Find(&permissionGroupMappingTB)

	return permissionGroupMappingTB
}

// GetAlertConfigInfoByPK return Configuration of alert data.
func GetAlertConfigInfoByPK(pk string) *CONFIGURATION_DATA_ALERT_TB {
	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{}
	DBgorm().Where("CHANNEL_PK = ?", pk).First(&configurationAlertDataTB)

	return configurationAlertDataTB
}

// GetAlertConfigInfoByName return Configuration of alert data.
func GetAlertConfigInfoByName(name string) *CONFIGURATION_DATA_ALERT_TB {
	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{}
	DBgorm().Where("CHANNEL_NAME = ?", name).First(&configurationChannelTB)

	configurationAlertDataTB := &CONFIGURATION_DATA_ALERT_TB{}
	DBgorm().Where("CHANNEL_PK = ?", configurationChannelTB.CHANNEL_PK).First(&configurationAlertDataTB)

	return configurationAlertDataTB
}

// GetVisibilityConfigInfo return Configuration of visibility data.
func GetVisibilityConfigInfo(pk string) *CONFIGURATION_DATA_VISIBILITY_TB {
	configurationVisibilityDataTB := &CONFIGURATION_DATA_VISIBILITY_TB{}
	DBgorm().Where("CHANNEL_PK = ?", pk).First(&configurationVisibilityDataTB)

	return configurationVisibilityDataTB
}

// GetAlertConfigTable returns table data of alert configuration.
func GetAlertConfigTable() []CONFIGURATION_DATA_ALERT_TB {
	var configurationAlertDataTB []CONFIGURATION_DATA_ALERT_TB
	DBgorm().Find(&configurationAlertDataTB)

	return configurationAlertDataTB
}

// GetVisibilityConfigTable returns table data of visibility configuration.
func GetVisibilityConfigTable() []CONFIGURATION_DATA_VISIBILITY_TB {
	var configurationVisibilityDataTB []CONFIGURATION_DATA_VISIBILITY_TB
	DBgorm().Find(&configurationVisibilityDataTB)

	return configurationVisibilityDataTB
}

// GetConfigurationChannelTable return Configuration channel table.
func GetConfigurationChannelTable() []CONFIGURATION_CHANNEL_TB {
	var configurationChannelTB []CONFIGURATION_CHANNEL_TB
	DBgorm().Find(&configurationChannelTB)

	return configurationChannelTB
}

// GetConfigurationChannelInfo return Configuration channel Info.
func GetConfigurationChannelInfo(pk string) *CONFIGURATION_CHANNEL_TB {
	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{}
	DBgorm().Where("CHANNEL_PK = ?", pk).First(&configurationChannelTB)

	return configurationChannelTB
}

// GetConfigurationChannelInfo return Configuration channel Info by channel list.
func GetConfigurationChannelInfoByList(pk []string) []CONFIGURATION_CHANNEL_TB {
	var configurationChannelTB []CONFIGURATION_CHANNEL_TB
	DBgorm().Where(pk).Find(&configurationChannelTB)

	return configurationChannelTB
}

// GetConfigurationChannelInfo return Configuration channel Info by name.
func GetConfigurationChannelInfoByName(name string) *CONFIGURATION_CHANNEL_TB {
	configurationChannelTB := &CONFIGURATION_CHANNEL_TB{}
	DBgorm().Where("CHANNEL_NAME = ?", name).First(&configurationChannelTB)

	return configurationChannelTB
}

// GetChannelPK return channel PK.
func GetChannelPK(channel string) string {
	var channelTB CONFIGURATION_CHANNEL_TB
	DBgorm().Where("CHANNEL_NAME = ?", channel).First(&channelTB)

	return channelTB.CHANNEL_PK
}

// GetChannelPermissionNodesTable return node and channel mapping table.
func GetChannelPermissionNodesTable() []NODE_CHANNEL_MAPPING_TB {
	var nodeChannelMappingTB []NODE_CHANNEL_MAPPING_TB
	DBgorm().Find(&nodeChannelMappingTB)

	return nodeChannelMappingTB
}

// GetChannelPermissionNodes return Configuration channel Info.
func GetChannelPermissionNodes(pk string) []NODE_CHANNEL_MAPPING_TB {
	var nodeChannelMappingTB []NODE_CHANNEL_MAPPING_TB
	DBgorm().Where("CHANNEL_PK = ?", pk).Find(&nodeChannelMappingTB)

	return nodeChannelMappingTB
}

// GetConfigurationMasterTable return Configuration master table.
func GetConfigurationMasterTable() *CONFIGURATION_MASTER_TB {
	configurationMasterTB := &CONFIGURATION_MASTER_TB{}
	DBgorm().First(&configurationMasterTB)

	return configurationMasterTB
}

// GetConfigurationNodeInfoByNodePK return Configuration node Info.
func GetConfigurationNodeInfoByNodePK(pk string) *CONFIGURATION_NODE_TB {
	configurationNodeTB := &CONFIGURATION_NODE_TB{}
	DBgorm().Where("NODE_PK = ?", pk).First(&configurationNodeTB)

	return configurationNodeTB
}

// GetConfigurationNodeInfoByNodeName return Configuration node Info.
func GetConfigurationNodeInfoByNodeName(nodeName string) *CONFIGURATION_NODE_TB {
	configurationNodeTB := &CONFIGURATION_NODE_TB{}
	DBgorm().Where("NODE_NAME = ?", nodeName).First(&configurationNodeTB)

	return configurationNodeTB
}

// GetConfigurationNodeTable return Configuration node table Info.
func GetConfigurationNodeTable() []CONFIGURATION_NODE_TB {
	var configurationNodeTB []CONFIGURATION_NODE_TB
	DBgorm().Find(&configurationNodeTB)

	return configurationNodeTB
}