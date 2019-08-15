package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/rollout/rox-go/core/consts"
)

func GenerateMD5(properties map[string]string, generatorList []consts.PropertyType, extraValues []string) string {
	var values []string

	for _, pt := range generatorList {
		value := properties[pt.Name]
		if value != "" {
			values = append(values, value)
		}
	}

	if len(extraValues) > 0 {
		// TODO
	}

	valueBytes := []byte(strings.Join(values, "|"))
	hashBytes := md5.Sum(valueBytes)

	return strings.ToUpper(hex.EncodeToString(hashBytes[:]))
}
