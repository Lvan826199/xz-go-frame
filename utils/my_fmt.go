/*
@Author: 梦无矶小仔
@Date:   2024/1/25 16:50
*/
package utils

import (
	"encoding/json"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"reflect"
)

// StructToMapJson 使用Json把struct转为map
func StructToMapJson(obj interface{}) map[string]any {
	var (
		data = make(map[string]interface{})
		buf  []byte
		err  error
	)
	buf, err = json.Marshal(obj)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(buf, &data)
	if err != nil {
		return nil
	}
	return data
}

// 使用反射把struct转化为map
func StructToMapReflect(obj interface{}) map[string]any {
	// 获取变量的类型
	obj1 := reflect.TypeOf(obj)
	// 获取变量的值
	obj2 := reflect.ValueOf(obj)

	data := make(map[string]any)
	// NumField()找到里面字段的数量
	for i := 0; i < obj1.NumField(); i++ {
		if obj1.Field(i).Tag.Get("mapstructure") != "" {
			data[obj1.Field(i).Tag.Get("mapstructure")] = obj2.Field(i).Interface()
		} else {
			data[obj1.Field(i).Name] = obj2.Field(i).Interface()
		}
	}
	return data
}

// MaheHump 将字符串转换为驼峰命名
func MaheHump(s string) string {
	humpName := cases.Title(language.English)
	s2 := humpName.String(s)
	return s2
}
