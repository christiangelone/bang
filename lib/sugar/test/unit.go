package sugar_test

import (
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/christiangelone/bang/lib/meta"
	"github.com/christiangelone/bang/lib/ux/print"
	"github.com/smartystreets/goconvey/convey"
)

func DescribeFunc(fn interface{}, closure func(), t *testing.T) {
	funcStr := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	funcNameParts := strings.Split(funcStr, "/")
	describe := print.Sprint("Describe function:", funcNameParts[len(funcNameParts)-1])
	convey.Convey(describe, t, closure)
}

func DescribeStruct(aStruct meta.MetaStruct, closure func(), t *testing.T) {
	convey.Convey("Describe struct: "+aStruct.GetStructName(), t, closure)
}

func DescribeMethod(fn interface{}, closure func()) {
	methodStr := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	methodParts := strings.Split(methodStr, ".")
	methodName := methodParts[len(methodParts)-1]
	describe := print.Sprint("Describe method:", methodName[:len(methodName)-3])
	convey.Convey(describe, closure)
}

func Context(text string, closure func()) {
	convey.Convey("Context: "+text, closure)
}

func It(text string, closure func()) {
	convey.Convey("It "+text, closure)
}

func AssertThat(actual interface{}, assert func(actual interface{}, expected ...interface{}) string, expected ...interface{}) {
	convey.So(actual, assert, expected...)
}
