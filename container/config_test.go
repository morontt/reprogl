package container

import (
	"reflect"
	"testing"

	"github.com/gookit/ini/v2"
)

type testStruct struct {
	A string `app_config:"PARAM_A"`
	B int    `app_config:"PARAM_B"`
	C string
}

type testInvalidStruct struct {
	D string  `app_config:"PARAM_D"`
	E int     `app_config:"PARAM_E"`
	F float32 `app_config:"PARAM_F"`
}

type testInvalidParam struct {
	D string `app_config:"PARAM_D"`
	Z int    `app_config:"PARAM_Z"`
}

var iniStr = `PARAM_A = test_a
PARAM_B=42
PARAM_D=test_d
PARAM_E=568
PARAM_F=0.301
`

func TestParseField(t *testing.T) {
	ta := testStruct{}
	fields := parseFields(reflect.TypeOf(ta))

	if !testContainField("A", fields) {
		t.Error("testStruct does not contain the \"A\" field")
	}

	if testContainField("C", fields) {
		t.Error("testStruct contain the \"C\" field")
	}
}

func TestSetupField(t *testing.T) {
	ini.LoadStrings(iniStr)

	ta := testStruct{}
	fields := parseFields(reflect.TypeOf(ta))

	rv := reflect.ValueOf(&ta)
	err := setupFields(rv, fields)
	if err != nil {
		t.Error(err.Error())
	}

	if ta.A != "test_a" {
		t.Error("incorrect value for PARAM_A")
	}
	if ta.B != 42 {
		t.Error("incorrect value for PARAM_B")
	}
	if ta.C != "" {
		t.Error("incorrect value for \"C\" field")
	}
}

func TestSetupFieldWithInvalidStruct(t *testing.T) {
	ini.LoadStrings(iniStr)

	ta := testInvalidStruct{}
	fields := parseFields(reflect.TypeOf(ta))

	rv := reflect.ValueOf(&ta)
	err := setupFields(rv, fields)
	if err != nil {
		if err.Error() != "config: undefined type float32 of field F" {
			t.Error("unknown error message for undefined type")
		}
	} else {
		t.Error("ignoring unknown type")
	}
}

func TestInvalidParameter(t *testing.T) {
	ini.LoadStrings(iniStr)

	ta := testInvalidParam{}
	fields := parseFields(reflect.TypeOf(ta))

	rv := reflect.ValueOf(&ta)
	err := setupFields(rv, fields)
	if err != nil {
		if err.Error() != "app.ini: Undefined parameter \"PARAM_Z\"" {
			t.Error("unknown error message for undefined parameter")
		}
	} else {
		t.Error("ignoring unknown parameter")
	}
}

func testContainField(fieldName string, fields []structField) bool {
	for _, el := range fields {
		if el.name == fieldName {
			return true
		}
	}

	return false
}
