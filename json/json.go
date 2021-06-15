package json

import (
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

type Object struct {
	maps map[string]interface{}
}

type Array struct {
	arrays []interface{}
}

func NewJsonObject(str string) (*Object, error) {
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &m)
	return &Object{maps: m}, err
}

func NewJsonObjectWithMap(m map[string]interface{}) *Object {
	return &Object{maps: m}
}

func newJsonArray(m []interface{}) *Array {
	return &Array{arrays: m}
}

func (obj *Object) GetObject(key string) *Object {
	if o, ok := obj.maps[key]; ok {
		return NewJsonObjectWithMap(o.(map[string]interface{}))
	}
	return &Object{}
}

func (obj *Object) Get(key string) interface{} {
	if o, ok := obj.maps[key]; ok {
		return o
	}
	return ""
}

func (obj *Object) GetArray(key string) *Array {
	if o, ok := obj.maps[key]; ok {
		return newJsonArray(o.([]interface{}))
	}
	return &Array{arrays: make([]interface{}, 0)}
}

func (obj *Object) ForEach(fn func(string, *Object)) {
	for key, obj := range obj.maps {
		fn(key, NewJsonObjectWithMap(obj.(map[string]interface{})))
	}
}

func (obj *Object) ForEachArray(fn func(string, *Array)) {
	for key, obj := range obj.maps {
		fn(key, newJsonArray(obj.([]interface{})))
	}
}

func (obj *Object) GetByKeys(keyStr string) interface{} {
	jsonMap := obj.maps
	keys := strings.Split(keyStr, ".")
	for i, k := range keys {

		if i == len(keys)-1 {
			if value, ok := jsonMap[k]; ok {
				return value
			}
		}

		kk := k
		if strings.Contains(k, "[") {
			kk = k[:strings.Index(k, "[")]
		}

		value, ok := jsonMap[kk]
		if !ok || value == nil {
			return ""
		}

		mt := reflect.TypeOf(value)
		if strings.HasPrefix(mt.String(), "map[string]interface ") {
			jsonMap = jsonMap[kk].(map[string]interface{})
		} else if strings.HasPrefix(mt.String(), "map[string]string") {
			_jsonMap := jsonMap[kk].(map[string]string)
			jsonMap = map[string]interface{}{}
			for k, v := range _jsonMap {
				jsonMap[k] = v
			}
		} else if strings.HasPrefix(mt.String(), "[]interface ") {
			array := jsonMap[kk].([]interface{})
			index := 0
			if strings.Contains(k, "[") && strings.Contains(k, "]") {
				indexStr := k[strings.Index(k, "[")+1 : strings.Index(k, "]")]
				index, _ = strconv.Atoi(indexStr)
			}
			if len(array) > index {
				jsonMap = array[index].(map[string]interface{})
			} else {
				return ""
			}
		} else {
			return ""
		}
	}
	return ""
}

func (obj *Array) ForEach(fn func(int, *Object)) {
	for i, obj := range obj.arrays {
		fn(i, NewJsonObjectWithMap(obj.(map[string]interface{})))
	}
}

func (obj *Array) ForEachArray(fn func(int, *Array)) {
	for i, obj := range obj.arrays {
		fn(i, newJsonArray(obj.([]interface{})))
	}
}

func (obj *Array) GetObject(i int) *Object {
	return NewJsonObjectWithMap(obj.arrays[i].(map[string]interface{}))
}

func (obj *Array) GetArray(i int) *Array {
	return newJsonArray(obj.arrays[i].([]interface{}))
}

func (obj *Array) Length() int {
	if obj.arrays == nil {
		return 0
	}
	return len(obj.arrays)
}

func (obj *Array) ToString() string {
	str, _ := json.Marshal(obj.arrays)
	return string(str)
}

func (obj *Object) ToString() string {
	str, _ := json.Marshal(obj.maps)
	return string(str)
}
