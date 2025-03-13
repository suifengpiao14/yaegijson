package yaegijson

import "github.com/tidwall/gjson"

// GetValueFromJson 从json字符串中获取指定路径的值
func GetValueFromJson(jsonStr string, paths ...string) (values []string) {
	results := gjson.GetMany(jsonStr, paths...)
	values = make([]string, len(paths))
	for i := 0; i < len(paths); i++ {
		values[i] = results[i].String()
	}
	return values
}
