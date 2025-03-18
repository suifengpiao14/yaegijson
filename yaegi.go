package yaegijson

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	_ "github.com/suifengpiao14/gjsonmodifier" // 引入自定义的gjson扩展包后会执行init()方法，注册自定义的gjson扩展函数

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

//go:generate go install github.com/traefik/yaegi/cmd/yaegi
//go:generate yaegi extract github.com/tidwall/gjson
//go:generate yaegi extract github.com/tidwall/sjson
//go:generate yaegi extract github.com/spf13/cast
//go:generate yaegi extract github.com/pkg/errors
//github.com/suifengpiao14/gjsonmodifier

var Symbols = stdlib.Symbols

type DynamicExtension struct {
	ExtensionCode string `json:"extensionCode"`
	ExtensionPath string `json:"extensionPath"`
	_interpreter  *interp.Interpreter
	symbols       map[string]map[string]reflect.Value `json:"-"`
}

func (extension *DynamicExtension) Withsymbols(symbols map[string]map[string]reflect.Value) *DynamicExtension {
	extension.symbols = symbols
	return extension
}

func NewDynamicExtension(sourceCode string, sourcePath string) *DynamicExtension {
	return &DynamicExtension{
		ExtensionCode: sourceCode,
		ExtensionPath: sourcePath,
	}
}

func (extension *DynamicExtension) _Eval() (err error) {
	if extension._interpreter != nil {
		return nil
	}

	if extension.ExtensionCode == "" && extension.ExtensionPath == "" {
		return errors.New("sourceCode or sourcePath required")
	}

	// 解析动态脚本
	interpreter := interp.New(interp.Options{})
	interpreter.Use(Symbols)           //注册当前包结构体
	interpreter.Use(extension.symbols) //注册外部包结构体

	if extension.ExtensionCode != "" {
		_, err = interpreter.Eval(extension.ExtensionCode)
		if err != nil {
			err = errors.WithMessagef(err, "compile dynamic go sourceCode: %s", extension.ExtensionPath)
			return err
		}
	}

	if extension.ExtensionPath != "" {
		_, err = interpreter.EvalPath(extension.ExtensionPath)
		if err != nil {
			err = errors.WithMessagef(err, "compile dynamic go sourcePath: %s", extension.ExtensionPath)
			return err
		}
	}
	extension._interpreter = interpreter

	return nil
}

// GetDestFuncImpl 获取动态脚本中定义的函数实现，并赋值给dstFn
func (extension DynamicExtension) GetDestFuncImpl(funcName string, dstFn any) (err error) {
	if funcName == "" {
		return nil
	}
	if dstFn == nil {
		err = errors.New("dstFn is nil")
		return err
	}
	if reflect.TypeOf(dstFn).Kind() != reflect.Pointer {
		err = errors.New("dstFn must be pointer")
		return err
	}
	err = extension._Eval()
	if err != nil {
		return err
	}
	err = getFn(extension._interpreter, funcName, dstFn)
	if err != nil {
		return err
	}
	return nil
}

var Error_not_found_func = errors.New("dynamic fun not found")

// getFn 从动态脚本中获取特定函数
func getFn(interpreter *interp.Interpreter, funcName string, dstFn any) (err error) {
	fnV, err := interpreter.Eval(funcName)

	if err != nil {
		err = errors.WithMessage(err, funcName)
		return err
	}
	rv := reflect.Indirect(reflect.ValueOf(dstFn))
	rt := rv.Type()
	if !rv.CanSet() {
		err = errors.Errorf("dstFn must can set, but %s", rt.String())
		return err
	}
	if fnV.IsNil() {
		err = errors.WithMessagef(Error_not_found_func, "func name: %s", funcName)
		return err
	}
	if !fnV.CanConvert(rt) {
		err = errors.Errorf("dynamic func %s ,must can convert to %s", funcName, fmt.Sprintf("%s.%s", rt.PkgPath(), rt.Name()))
		return err
	}
	rv.Set(fnV.Convert(rt))
	return nil
}
