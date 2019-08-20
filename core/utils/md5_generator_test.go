package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/rollout/rox-go/core/consts"
	"github.com/stretchr/testify/assert"
)

func HashToString(string) string {
	value := md5.Sum([]byte("value"))
	return strings.ToUpper(hex.EncodeToString(value[:]))
}

func TestWillCheckMD5UsesRightProps(t *testing.T) {
	props := make(map[string]string)
	props[consts.PropertyTypePlatform.Name] = "value"
	md5Computed := GenerateMD5(props, []consts.PropertyType{*consts.PropertyTypePlatform})
	md5Manual := HashToString("value")

	assert.Equal(t, md5Computed, md5Manual)
}

func TestWillCheckMD5NotUsingAllProps(t *testing.T) {
	props := make(map[string]string)
	props[consts.PropertyTypeDevModeSecret.Name] = "dev"
	props[consts.PropertyTypePlatform.Name] = "value"
	md5Computed := GenerateMD5(props, []consts.PropertyType{*consts.PropertyTypePlatform})
	md5Manual := HashToString("value")

	assert.Equal(t, md5Computed, md5Manual)
}

func TestWillCheckMD5UsingAllProps(t *testing.T) {
	props := make(map[string]string)
	props[consts.PropertyTypeDevModeSecret.Name] = "dev"
	props[consts.PropertyTypePlatform.Name] = "value"
	md5Computed := GenerateMD5(props, []consts.PropertyType{*consts.PropertyTypePlatform})
	md5Manual := HashToString("dev|value")

	assert.Equal(t, md5Computed, md5Manual)
}
