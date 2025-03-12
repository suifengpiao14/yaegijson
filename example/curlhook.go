package curlhook

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func BeforeFn(input string) (output string, err error) {
	timestamps := gjson.Get(input, "body._head._timestamps").String()
	_ = timestamps
	input, err = sjson.Set(input, "body._head._timestamps", "1111111111111111")
	if err != nil {
		return "", err
	}
	return input, nil
}
func AfterFn(input string) (output string, err error) {
	return input, nil
}
