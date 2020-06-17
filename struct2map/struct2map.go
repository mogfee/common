package struct2map

import (
	"reflect"
	"strings"
)

func GetMapFromStruct(v interface{}, columns []string) map[string]interface{} {
	dataValue := reflect.ValueOf(v).Elem()
	dataType := reflect.TypeOf(v).Elem()
	allow := make(map[string]bool)
	for _, v := range columns {
		allow[v] = true
	}
	//fmt.Println(dataValue.FieldByName("Id").Interface())
	mp := make(map[string]interface{})
	for i := 0; i < dataType.NumField(); i++ {
		v := dataType.Field(i).Tag.Get("gorm")
		v1 := strings.Split(v, "column:")
		if len(v1) < 1 {
			continue
		}
		if len(v1) > 1 {
			v = strings.Split(v1[1], ";")[0]
		}
		if allow[v] {
			mp[v] = dataValue.Field(i).Interface()
		}
	}
	return mp
}
