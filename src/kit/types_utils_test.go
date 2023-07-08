package kit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MapsEqual(t *testing.T) {
	tests := []struct {
		name  string
		m1    map[string]interface{}
		m2    map[string]interface{}
		equal bool
	}{
		{
			name:  "Both nils",
			m1:    nil,
			m2:    nil,
			equal: true,
		},
		{
			name:  "Nil vs empty",
			m1:    make(map[string]interface{}),
			m2:    nil,
			equal: false,
		},
		{
			name: "contains second",
			m1: map[string]interface{}{
				"k": "v",
			},
			m2: map[string]interface{}{
				"k": "v",
				"v": "k",
			},
			equal: false,
		},
		{
			name: "contains first",
			m1: map[string]interface{}{
				"k": "v",
				"v": "k",
			},
			m2: map[string]interface{}{
				"k": "v",
			},
		},
		{
			name:  "Both nils",
			m1:    nil,
			m2:    nil,
			equal: true,
		},
		{
			name:  "Nil vs empty",
			m1:    make(map[string]interface{}),
			m2:    nil,
			equal: false,
		},
		{
			name: "Single values",
			m1: map[string]interface{}{
				"k": "v",
			},
			m2: map[string]interface{}{
				"k": "v",
			},
			equal: true,
		},
		{
			name: "Complex values",
			m1: map[string]interface{}{
				"k": struct {
					Value string
				}{
					Value: "value",
				},
			},
			m2: map[string]interface{}{
				"k": struct {
					Value string
				}{
					Value: "value",
				},
			},
			equal: true,
		},
		{
			name: "Multiple values",
			m1: map[string]interface{}{
				"k1": "v1",
				"k2": 100,
				"k3": true,
			},
			m2: map[string]interface{}{
				"k1": "v1",
				"k2": 100,
				"k3": true,
			},
			equal: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.equal, MapsEqual(tt.m1, tt.m2))
		})
	}
}

func Test_MapToLowerCamelKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "nil maps",
			input:    nil,
			expected: nil,
		},
		{
			name:     "empty maps",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name:     "one level map",
			input:    map[string]interface{}{"Key": "value"},
			expected: map[string]interface{}{"key": "value"},
		},
		{
			name:     "multi words key",
			input:    map[string]interface{}{"VeryComplexKey": "value"},
			expected: map[string]interface{}{"veryComplexKey": "value"},
		},
		{
			name:     "multi level map",
			input:    map[string]interface{}{"Key": map[string]interface{}{"AnotherKey": "value", "SecondKey": map[string]interface{}{"Key": "value"}}},
			expected: map[string]interface{}{"key": map[string]interface{}{"anotherKey": "value", "secondKey": map[string]interface{}{"key": "value"}}},
		},
		{
			name:     "only keys",
			input:    map[string]interface{}{"Key": "VALUE"},
			expected: map[string]interface{}{"key": "VALUE"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, MapToLowerCamelKeys(tt.input))
		})
	}
}

func Test_MapInterfacesToBytesAndBack(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]interface{}
	}{
		{
			name: "nil maps",
			m:    nil,
		}, {
			name: "empty maps",
			m:    map[string]interface{}{},
		}, {
			name: "map one value",
			m:    map[string]interface{}{"key1": "value1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := MapInterfacesToBytes(tt.m)
			m := BytesToMapInterfaces(bytes)
			assert.Equal(t, tt.m, m)
		})
	}
}

func Test_MapInterfacesToBytesNestedTypesAndBack(t *testing.T) {
	tests := []struct {
		name string
		m    map[string]interface{}
	}{
		{
			name: "diff type values",
			m: map[string]interface{}{
				"key1": "value1",
				"key2": float64(10),
				"key3": 98.2,
			},
		}, {
			name: "diff type values with map value",
			m: map[string]interface{}{
				"key1": "value1",
				"key2": float64(10),
				"key3": 98.2,
				"key4": map[string]interface{}{"key4internal1": float64(10), "key4internal2": "value2"}},
		}, {
			name: "diff type values with map value",
			m: map[string]interface{}{
				"key1": "value1",
				"key2": float64(10),
				"key3": 98.2,
				"key4": map[string]interface{}{
					"key4internal1": float64(10),
					"key4internal2": map[string]interface{}{
						"key4internal1": float64(10),
						"key4internal2": "value2"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := MapInterfacesToBytes(tt.m)
			m := BytesToMapInterfaces(bytes)
			assertMap(t, tt.m, m)
		})
	}
}

func Test_StringsToInterfaces(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		expected []interface{}
	}{
		{
			name:     "nil slice",
			slice:    nil,
			expected: nil,
		}, {
			name:     "empty slice",
			slice:    []string{},
			expected: []interface{}{},
		}, {
			name:     "two value",
			slice:    []string{"value1", "value2"},
			expected: []interface{}{"value1", "value2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sl := StringsToInterfaces(tt.slice)
			assert.Equal(t, tt.expected, sl)
		})
	}

}

func assertMap(t *testing.T, expectedM, actualM map[string]interface{}) {
	assert.Equal(t, len(expectedM), len(actualM))
	for k, v := range expectedM {
		if internalV, ok := v.(map[string]interface{}); ok {
			assertMap(t, internalV, actualM[k].(map[string]interface{}))
		} else {
			assert.Equal(t, v, actualM[k])
		}
	}
}
