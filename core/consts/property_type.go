package consts

type PropertyType struct {
	Value int
	Name  string
}

var (
	PropertyTypeCacheMissRelativeURL   = &PropertyType{1, "cache_miss_relative_url"}
	PropertyTypeLibVersion             = &PropertyType{4, "lib_version"}
	PropertyTypeRolloutBuild           = &PropertyType{5, "rollout_build"}
	PropertyTypeAPIVersion             = &PropertyType{6, "api_version"}
	PropertyTypeBuid                   = &PropertyType{7, "buid"}
	PropertyTypeBuidGeneratorsList     = &PropertyType{8, "buid_generators_list"}
	PropertyTypeAppRelease             = &PropertyType{10, "app_release"}
	PropertyTypeDistinctID             = &PropertyType{11, "distinct_id"}
	PropertyTypeAppKey                 = &PropertyType{12, "app_key"}
	PropertyTypeFeatureFlags           = &PropertyType{13, "feature_flags"}
	PropertyTypeRemoteVariables        = &PropertyType{14, "remote_variables"}
	PropertyTypeCustomProperties       = &PropertyType{15, "custom_properties"}
	PropertyTypePlatform               = &PropertyType{16, "platform"}
	PropertyTypeDevModeSecret          = &PropertyType{17, "devModeSecret"}
	PropertyTypeStateMD5               = &PropertyType{18, "state_md5"}
	PropertyTypeFeatureFlagsString     = &PropertyType{19, "feature_flags_string"}
	PropertyTypeRemoteVariablesString  = &PropertyType{20, "remote_variables_string"}
	PropertyTypeCustomPropertiesString = &PropertyType{21, "custom_properties_string"}
)
