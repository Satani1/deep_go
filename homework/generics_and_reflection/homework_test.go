package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	val := reflect.ValueOf(person)
	typ := val.Type()
	lines := make([]string, 0, typ.NumField())

	for idx := 0; idx < typ.NumField(); idx++ {
		field := typ.Field(idx)
		key, opts := parseFieldTag(field)
		if key == "" {
			continue
		}

		fieldValue := val.Field(idx)
		if shouldSkipField(opts, fieldValue) {
			continue
		}

		valueStr := formatValue(fieldValue)
		lines = append(lines, key+"="+valueStr)
	}

	return strings.Join(lines, "\n")
}

func parseFieldTag(field reflect.StructField) (key string, opts []string) {
	tag := field.Tag.Get("properties")
	if tag == "" {
		return "", nil
	}
	parts := strings.Split(tag, ",")
	return parts[0], parts[1:]
}

func shouldSkipField(opts []string, val reflect.Value) bool {
	for _, opt := range opts {
		if opt == "omitempty" && val.IsZero() {
			return true
		}
	}
	return false
}

func formatValue(val reflect.Value) string {
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", val.Int())
	case reflect.Bool:
		return fmt.Sprintf("%t", val.Bool())
	default:
		return fmt.Sprintf("%v", val.Interface())
	}
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
