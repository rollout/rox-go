package properties

type CustomPropertyType struct {
	Type         string
	ExternalType string
}

var (
	CustomPropertyTypeString = &CustomPropertyType{"string", "String"}
	CustomPropertyTypeBool   = &CustomPropertyType{"bool", "Boolean"}
	CustomPropertyTypeInt    = &CustomPropertyType{"int", "Number"}
	CustomPropertyTypeFloat  = &CustomPropertyType{"double", "Number"}
	CustomPropertyTypeSemver = &CustomPropertyType{"semver", "Semver"}
	CustomPropertyTypeTime   = &CustomPropertyType{"time", "DateTime"}
)
