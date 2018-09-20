package consts

type PropertyType struct {
	Value int
	Name  string
}

var (
	PropertyTypeCacheMissURL       = &PropertyType{1, "cache_miss_url"}
	PropertyTypePackageName        = &PropertyType{2, "package_name"}
	PropertyTypeVersionName        = &PropertyType{3, "version_name"}
	PropertyTypeLibVersion         = &PropertyType{4, "lib_version"}
	PropertyTypeRolloutBuild       = &PropertyType{5, "rollout_build"}
	PropertyTypeAPIVersion         = &PropertyType{6, "api_version"}
	PropertyTypeBuid               = &PropertyType{7, "buid"}
	PropertyTypeBuidGeneratorsList = &PropertyType{8, "buid_generators_list"}
	PropertyTypeAppVersion         = &PropertyType{9, "app_version"}
	PropertyTypeAppRelease         = &PropertyType{10, "app_release"}
	PropertyTypeDistinctID         = &PropertyType{11, "distinct_id"}
	PropertyTypeAppKey             = &PropertyType{12, "app_key"}
	PropertyTypeFeatureFlags       = &PropertyType{13, "feature_flags"}
	PropertyTypeRemoteVariables    = &PropertyType{14, "remote_variables"}
	PropertyTypeCustomProperties   = &PropertyType{15, "custom_properties"}
	PropertyTypePlatform           = &PropertyType{16, "platform"}
)
