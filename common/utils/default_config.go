package utils

import (
	"log"
	"reflect"
	"strconv"
)

//这里第一个参数不要传递的时候不要取地址(类似&val)，而第二个参数必须要取地址(&val)
func SetConfigDefaultValue(ct reflect.Type, cv reflect.Value) {
	if ct.Kind() != reflect.Struct {
		log.Printf("not support not-struct")
		return
	}
	elements := cv.Elem()
	for i := 0; i < ct.NumField(); i++ {
		filed := ct.Field(i)
		element := elements.Field(i)
		defaultValue := filed.Tag.Get("default")
		if defaultValue == "" || !element.CanSet() {
			if filed.Type.Kind() != reflect.Struct {
				continue
			}
		}
		switch filed.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value, _ := strconv.ParseInt(defaultValue, 10, 64)
			element.SetInt(value)
		case reflect.Float32, reflect.Float64:
			value, _ := strconv.ParseFloat(defaultValue, 64)
			element.SetFloat(value)
		case reflect.Bool:
			value, _ := strconv.ParseBool(defaultValue)
			element.SetBool(value)
		case reflect.String:
			element.SetString(defaultValue)
		case reflect.Struct:
			SetConfigDefaultValue(ct.Field(i).Type, element.Addr())
		default:
			log.Printf("not support this type: %s", filed.Type.Kind().String())
		}
	}
}
