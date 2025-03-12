package yaegijson

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

//go:generate go install github.com/traefik/yaegi/cmd/yaegi
//go:generate yaegi extract github.com/tidwall/gjson
//go:generate yaegi extract github.com/tidwall/sjson
//go:generate yaegi extract github.com/spf13/cast
//github.com/suifengpiao14/httpraw
//github.com/suifengpiao14/gjsonmodifier

var Symbols = stdlib.Symbols

type DynamicExtension struct {
	SourceCode   string `json:"sourceCode"`
	SourcePath   string `json:"sourcePath"`
	_interpreter *interp.Interpreter
}

func NewDynamicExtension(sourceCode string, sourcePath string) *DynamicExtension {
	return &DynamicExtension{
		SourceCode: sourceCode,
		SourcePath: sourcePath,
	}
}

func (extension *DynamicExtension) _Eval() (err error) {
	if extension._interpreter != nil {
		return nil
	}

	if extension.SourceCode == "" && extension.SourcePath == "" {
		return errors.New("sourceCode or sourcePath required")
	}

	// 解析动态脚本
	interpreter := interp.New(interp.Options{})
	interpreter.Use(stdlib.Symbols)
	interpreter.Use(Symbols) //注册当前包结构体

	if extension.SourceCode != "" {
		_, err = interpreter.Eval(extension.SourceCode)
		if err != nil {
			err = errors.WithMessagef(err, "compile dynamic go sourceCode: %s", extension.SourcePath)
			return err
		}
	}

	if extension.SourcePath != "" {
		_, err = interpreter.EvalPath(extension.SourcePath)
		if err != nil {
			err = errors.WithMessagef(err, "compile dynamic go sourcePath: %s", extension.SourcePath)
			return err
		}
	}
	extension._interpreter = interpreter

	return nil
}

// GetDestFuncImpl 获取动态脚本中定义的函数实现，并赋值给dstFn
func (extension DynamicExtension) GetDestFuncImpl(funcName string, dstFn any) (err error) {
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
