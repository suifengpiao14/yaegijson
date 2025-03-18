package yaegijson

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// GetValuesFromJson 从json字符串中获取指定路径的值
func GetValuesFromJson(jsonStr string, paths ...string) (values []string) {
	values = make([]string, len(paths))
	if jsonStr == "" {
		return values
	}
	results := gjson.GetMany(jsonStr, paths...)
	values = make([]string, len(paths))
	for i := 0; i < len(paths); i++ {
		values[i] = results[i].String()
	}
	return values
}
func SetValueToJson(jsonStr string, path string, value any) (newJsonStr string, err error) {
	jsonStr, err = sjson.Set(jsonStr, path, value)
	if err != nil {
		return "", err
	}
	return jsonStr, nil
}
