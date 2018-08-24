package client

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/rollout/rox-go/core/consts"
	"github.com/rollout/rox-go/core/model"
	"sort"
	"strings"
)

var buidGenerators = []*consts.PropertyType{
	consts.PropertyTypePlatform,
	consts.PropertyTypeAppKey,
	consts.PropertyTypeLibVersion,
	consts.PropertyTypeAPIVersion,
	consts.PropertyTypeCustomProperties,
	consts.PropertyTypeFeatureFlags,
	consts.PropertyTypeRemoteVariables,
}

type buid struct {
	sdkSettings              model.SdkSettings
	deviceProperties         model.DeviceProperties
	flagRepository           model.FlagRepository
	customPropertyRepository model.CustomPropertyRepository
	buid                     string
}

func NewBUID(sdkSettings model.SdkSettings, deviceProperties model.DeviceProperties, flagRepository model.FlagRepository, customPropertyRepository model.CustomPropertyRepository) model.BUID {
	return &buid{
		sdkSettings:              sdkSettings,
		deviceProperties:         deviceProperties,
		flagRepository:           flagRepository,
		customPropertyRepository: customPropertyRepository,
	}
}

func (b *buid) GetValue() string {
	properties := b.deviceProperties.GetAllProperties()
	var values []string
	for _, generator := range buidGenerators {
		if value, ok := properties[generator.Name]; ok {
			values = append(values, value)
		}
	}

	values = append(values, b.serializeFeatureFlags())
	values = append(values, b.serializeCustomProperties())

	valueBytes := []byte(strings.Join(values, "|"))
	hasher := md5.New()
	hasher.Write(valueBytes)
	hashBytes := hex.EncodeToString(hasher.Sum(nil))

	b.buid = strings.ToUpper(hashBytes)
	return b.buid
}

func (b *buid) GetQueryStringParts() map[string]string {
	generators := make([]string, 0, len(buidGenerators))
	for _, generator := range buidGenerators {
		generators = append(generators, generator.Name)
	}

	return map[string]string{
		consts.PropertyTypeBuid.Name:               b.GetValue(),
		consts.PropertyTypeBuidGeneratorsList.Name: strings.Join(generators, ","),
		consts.PropertyTypeFeatureFlags.Name:       b.serializeFeatureFlags(),
		consts.PropertyTypeRemoteVariables.Name:    "[]",
		consts.PropertyTypeCustomProperties.Name:   b.serializeCustomProperties(),
	}
}
func (b *buid) serializeFeatureFlags() string {
	var flags []jsonFlag
	allFlags := b.flagRepository.GetAllFlags()
	sort.Slice(allFlags, func(i, j int) bool {
		return allFlags[i].Name() < allFlags[j].Name()
	})
	for _, f := range allFlags {
		flags = append(flags, jsonFlag{f.Name(), f.DefaultValue(), f.Options()})
	}
	result, _ := json.Marshal(flags)
	return string(result)
}

func (b *buid) serializeCustomProperties() string {
	var properties []jsonProperty
	customProperties := b.customPropertyRepository.GetAllCustomProperties()
	sort.Slice(customProperties, func(i, j int) bool {
		return customProperties[i].Name < customProperties[j].Name
	})
	for _, p := range customProperties {
		properties = append(properties, jsonProperty{p.Name, p.Type.Type, p.Type.ExternalType})
	}
	result, _ := json.Marshal(properties)
	return string(result)
}

type jsonProperty struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	ExternalType string `json:"externalType"`
}

type jsonFlag struct {
	Name         string   `json:"name"`
	DefaultValue string   `json:"defaultValue"`
	Options      []string `json:"options"`
}
