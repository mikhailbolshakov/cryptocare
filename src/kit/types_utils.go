package kit

import (
	"encoding/json"
	"github.com/iancoleman/strcase"
	"reflect"
)

func MapsEqual(m1, m2 map[string]interface{}) bool {
	return reflect.DeepEqual(m1, m2)
}

func MapToLowerCamelKeys(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return nil
	}
	r := make(map[string]interface{}, len(m))
	for k, v := range m {
		if vMap, ok := v.(map[string]interface{}); ok && len(vMap) > 0 {
			r[strcase.ToLowerCamel(k)] = MapToLowerCamelKeys(vMap)
		} else {
			r[strcase.ToLowerCamel(k)] = v
		}
	}
	return r
}

func MapInterfacesToBytes(m map[string]interface{}) []byte {
	bytes, _ := json.Marshal(m)
	return bytes
}

func BytesToMapInterfaces(bytes []byte) map[string]interface{} {
	mp := make(map[string]interface{})
	_ = json.Unmarshal(bytes, &mp)
	return mp
}

func StringsToInterfaces(sl []string) []interface{} {
	if sl == nil {
		return nil
	}
	res := make([]interface{}, len(sl))
	for index, value := range sl {
		res[index] = value
	}

	return res
}
