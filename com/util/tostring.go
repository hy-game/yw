package util

import (
	"reflect"
	"strconv"
)

func Show(s interface{}) string {
	return s.(string)
}

func ToString(val interface{}) (ret string) {
	switch val.(type) {
	case int:
		ret = strconv.FormatInt(int64(val.(int)), 10)
	case int8:
		ret = strconv.FormatInt(int64(val.(int8)), 10)
	case int16:
		ret = strconv.FormatInt(int64(val.(int16)), 10)
	case int32:
		ret = strconv.FormatInt(int64(val.(int32)), 10)
	case int64:
		ret = strconv.FormatInt(val.(int64), 10)
	case float32:
		ret = strconv.FormatFloat(float64(val.(float32)), 'f', -1, 64)
	case float64:
		ret = strconv.FormatFloat(val.(float64), 'f', -1, 64)
	case uint8:
		ret = strconv.FormatUint(uint64(val.(uint8)), 10)
	case uint16:
		ret = strconv.FormatUint(uint64(val.(uint16)), 10)
	case uint32:
		ret = strconv.FormatUint(uint64(val.(uint32)), 10)
	case uint64:
		ret = strconv.FormatUint(val.(uint64), 10)
	case bool:
		if val.(bool) {
			ret = "true"
		} else {
			ret = "false"
		}
	case string:
		ret = val.(string)
	default:
		if reflect.TypeOf(val).Kind() == reflect.Int32 {
			ret = strconv.FormatInt(reflect.ValueOf(val).Int(), 10)
		}
	}
	return
}
