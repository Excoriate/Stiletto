package config

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/spf13/viper"
	"os"
)

type CfgValue struct {
	Key   string
	Value interface{}
}

type Cfg struct {
	key string
}

type CfgRetriever interface {
	GetFromViper(key string) (CfgValue, error)
	GetFromEnvVars(key string) (CfgValue, error)
	GetFromAny(key string) (CfgValue, error)
	IsRunningInVendorAutomation() bool
}

func (c *Cfg) GetFromViper(key string) (CfgValue, error) {
	var keyToSeek string

	if key == "" {
		keyToSeek = c.key
	} else {
		keyToSeek = key
	}

	if keyToSeek == "" {
		return CfgValue{}, errors.NewInternalPipelineError(fmt.Sprintf(
			"Failed to get config value (from viper) for key: %s. It's passed empty", keyToSeek))
	}

	keyNormalised := common.NormaliseNoSpaces(keyToSeek)

	value := viper.Get(keyNormalised)

	if value == nil {
		return CfgValue{}, errors.NewInternalPipelineError(fmt.Sprintf(
			"Failed to get config value (from viper) for key: %s. It is not found.", keyNormalised))
	}

	if common.IsNotNilAndNotEmpty(value.(string)) {
		return CfgValue{Key: keyNormalised, Value: value}, nil
	}

	return CfgValue{}, errors.NewInternalPipelineError(fmt.Sprintf("Failed to get config ("+
		"from viper) value for key: %s. It is not found.", keyNormalised))
}

func (c *Cfg) GetFromEnvVars(key string) (CfgValue, error) {
	var keyToSeek string

	if key == "" {
		keyToSeek = c.key
	} else {
		keyToSeek = key
	}

	if keyToSeek == "" {
		return CfgValue{}, errors.NewInternalPipelineError(fmt.Sprintf(
			"Failed to get config (from env vars) value for key: %s. It's passed empty", keyToSeek))
	}

	keyNormalised := common.NormaliseNoSpaces(keyToSeek)

	value := os.Getenv(keyNormalised)
	if common.IsNotNilAndNotEmpty(value) {
		return CfgValue{Key: keyNormalised, Value: value}, nil
	}

	return CfgValue{}, errors.NewInternalPipelineError(fmt.Sprintf("Failed to get config ("+
		"from env vars) value for key: %s. It is not found.", keyNormalised))
}

func (c *Cfg) GetFromAny(key string) (CfgValue, error) {
	var keyToSeek string

	if key == "" {
		keyToSeek = c.key
	} else {
		keyToSeek = key
	}

	if keyToSeek == "" {
		return CfgValue{}, errors.NewInternalPipelineError(fmt.Sprintf(
			"Failed to get config (from any) value for key: %s. It's passed empty", keyToSeek))
	}

	keyNormalised := common.NormaliseNoSpaces(keyToSeek)

	value := viper.Get(keyNormalised)

	if common.IsNotNilAndNotEmpty(value) {
		return CfgValue{Key: keyNormalised, Value: value}, nil
	}

	value = os.Getenv(keyNormalised)
	if common.IsNotNilAndNotEmpty(value) {
		return CfgValue{Key: keyNormalised, Value: value}, nil
	}

	return CfgValue{}, errors.NewInternalPipelineError(fmt.Sprintf("Failed to get config ("+
		"from any) value for key: %s. It is not found.", keyNormalised))
}

func (c *Cfg) IsRunningInVendorAutomation() bool {
	runInVendor := viper.Get("run-in-vendor")
	if runInVendor == nil {
		return false
	}

	return runInVendor.(bool)
}
