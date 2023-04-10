package common

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func MapIsNulOrEmpty(target map[string]string) bool {
	return target == nil || len(target) == 0
}

func IsStringInSlice(target string, list []string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}

func StructToMap(input interface{}, toUpper bool) map[string]interface{} {
	result := make(map[string]interface{})
	inputValue := reflect.ValueOf(input)
	inputType := reflect.TypeOf(input)

	for i := 0; i < inputValue.NumField(); i++ {
		fieldValue := inputValue.Field(i)
		fieldType := inputType.Field(i)
		key := fieldType.Name

		if toUpper {
			key = strings.ToUpper(key)
		} else {
			key = strings.ToLower(key)
		}

		result[key] = fieldValue.Interface()
	}

	return result
}

func ConvertToString(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case int:
		return strconv.Itoa(v), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("unsupported value type: %T", value)
	}
}

func MapInterfaceToString(input map[string]interface{}) (map[string]string, error) {
	result := make(map[string]string)

	for key, value := range input {
		strValue, err := ConvertToString(value)
		if err != nil {
			return nil, err
		}
		result[key] = strValue
	}

	return result, nil
}

func IsImageURLIncludesTag(imageURL string) bool {
	return strings.Contains(imageURL, ":")
}
