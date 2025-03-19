package yaegijson_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suifengpiao14/yaegijson"
)

var dynamic = `
package curlhook_source

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

`
var path = "example/curlhook.go"
var input = `{"body":{"_head":{"_timestamps":"1234567890"}}}`
var dstFn = func(input string) (output string, err error) { return "to implement", nil }

func TestDynamicScript(t *testing.T) {
	pathFuncName := "curlhook.BeforeFn"
	sourceFuncName := "curlhook_source.BeforeFn"
	t.Run("from string", func(t *testing.T) {
		dynamicExtension := yaegijson.NewExtension().WithSouceCode(dynamic)
		err := dynamicExtension.GetDestFuncImpl(sourceFuncName, &dstFn)
		require.NoError(t, err)
		output, err := dstFn(input)
		require.NoError(t, err)
		require.NotEmpty(t, output)
		require.JSONEq(t, `{"body":{"_head":{"_timestamps":"1111111111111111"}}}`, output)
	})
	t.Run("from path", func(t *testing.T) {
		dynamicExtension := yaegijson.NewExtension().WithSourcePath(path)
		err := dynamicExtension.GetDestFuncImpl(pathFuncName, &dstFn)
		require.NoError(t, err)
		output, err := dstFn(input)
		require.NoError(t, err)
		require.NotEmpty(t, output)
		require.JSONEq(t, `{"body":{"_head":{"_timestamps":"1111111111111111"}}}`, output)
	})
	t.Run("from string and path ,use path", func(t *testing.T) {
		dynamicExtension := yaegijson.NewExtension().WithSouceCode(dynamic).WithSourcePath(path)
		err := dynamicExtension.GetDestFuncImpl(pathFuncName, &dstFn)
		require.NoError(t, err)
		output, err := dstFn(input)
		require.NoError(t, err)
		require.NotEmpty(t, output)
		require.JSONEq(t, `{"body":{"_head":{"_timestamps":"1111111111111111"}}}`, output)
	})
	t.Run("from string and path,use source code", func(t *testing.T) {
		dynamicExtension := yaegijson.NewExtension().WithSouceCode(dynamic).WithSourcePath(path)
		err := dynamicExtension.GetDestFuncImpl(sourceFuncName, &dstFn)
		require.NoError(t, err)
		output, err := dstFn(input)
		require.NoError(t, err)
		require.NotEmpty(t, output)
		require.JSONEq(t, `{"body":{"_head":{"_timestamps":"1111111111111111"}}}`, output)
	})
}
