package config

import (
	"fmt"
	"github.com/Excoriate/stiletto/internal/common"
	"github.com/Excoriate/stiletto/internal/errors"
	"github.com/spf13/viper"
	"os"
	"reflect"
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
	GetFromViperOrDefault(key string, defaultValue interface{}) (CfgValue, error)
	GetFromEnvVars(key string) (CfgValue, error)
	GetFromAny(key string) (CfgValue, error)
	IsRunningInVendorAutomation() bool
	ValidateCfgKey(key string) (string, error)
	GetStringSliceFromViper(key string) (CfgValue, error)
	GetStringInterfaceMapFromViper(key string) (CfgValue, error)
}

func (c *Cfg) ValidateCfgKey(key string) (string, error) {
	var keyToSeek string

	if key == "" {
		keyToSeek = c.key
	} else {
		keyToSeek = key
	}

	if keyToSeek == "" {
		return "", errors.NewInternalPipelineError(fmt.Sprintf(
			"Failed to get config value (from viper) for key: %s. It's passed empty", keyToSeek))
	}

	keyNormalised := common.NormaliseNoSpaces(keyToSeek)
	return keyNormalised, nil
}

func (c *Cfg) GetStringSliceFromViper(key string) (CfgValue, error) {
	keyNormalised, err := c.ValidateCfgKey(key)
	if err != nil {
		return CfgValue{}, err
	}

	value := viper.GetStringSlice(keyNormalised)
	if len(value) == 0 {
		return CfgValue{Key: keyNormalised, Value: []string{}}, nil
	}

	return CfgValue{Key: keyNormalised, Value: value}, nil
}

func (c *Cfg) GetFromViperOrDefault(key string, defaultValue interface{}) (CfgValue, error) {
	keyNormalised, err := c.ValidateCfgKey(key)
	if err != nil {
		return CfgValue{}, err
	}

	value := viper.Get(keyNormalised)
	if value == nil {
		return CfgValue{Key: keyNormalised, Value: defaultValue}, nil
	}

	switch v := value.(type) {
	case string:
		if v == "" {
			return CfgValue{Key: keyNormalised, Value: defaultValue}, nil
		}
	case map[string]interface{}:
		if len(v) == 0 {
			return CfgValue{Key: keyNormalised, Value: defaultValue}, nil
		}
	case []string:
		if len(v) == 0 {
			return CfgValue{Key: keyNormalised, Value: defaultValue}, nil
		}
	default:
		// Check if value and defaultValue are of the same type
		if reflect.TypeOf(value) != reflect.TypeOf(defaultValue) {
			return CfgValue{}, errors.NewPipelineConfigurationError(
				"type mismatch between value and defaultValue", nil)
		}
	}

	return CfgValue{Key: keyNormalised, Value: value}, nil
}
func (c *Cfg) GetFromViper(key string) (CfgValue, error) {
	keyNormalised, err := c.ValidateCfgKey(key)
	if err != nil {
		return CfgValue{}, err
	}

	value := viper.Get(keyNormalised)

	if value == nil {
		return CfgValue{}, errors.NewInternalPipelineError(fmt.Sprintf(
			"Failed to get config value (from viper) for key: %s. It is not found.", keyNormalised))
	}

	if common.IsNotNilAndNotEmpty(value) {
		return CfgValue{Key: keyNormalised, Value: value}, nil
	}

	return CfgValue{
		Key:   keyNormalised,
		Value: value,
	}, nil
}

func (c *Cfg) GetStringInterfaceMapFromViper(key string) (CfgValue, error) {
	keyNormalised, err := c.ValidateCfgKey(key)
	if err != nil {
		return CfgValue{}, err
	}

	value := viper.GetStringMap(keyNormalised)
	if len(value) == 0 {
		return CfgValue{Key: keyNormalised, Value: map[string]interface{}{}}, nil
	}

	return CfgValue{Key: keyNormalised, Value: value}, nil
}

func (c *Cfg) GetFromEnvVars(key string) (CfgValue, error) {
	keyNormalised, err := c.ValidateCfgKey(key)
	if err != nil {
		return CfgValue{}, err
	}

	value := os.Getenv(keyNormalised)
	if common.IsNotNilAndNotEmpty(value) {
		return CfgValue{Key: keyNormalised, Value: value}, nil
	}

	return CfgValue{}, errors.NewInternalPipelineError(fmt.Sprintf("Failed to get config ("+
		"from env vars) value for key: %s. It is not found.", keyNormalised))
}

func (c *Cfg) GetFromAny(key string) (CfgValue, error) {
	keyNormalised, err := c.ValidateCfgKey(key)
	if err != nil {
		return CfgValue{}, err
	}

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
