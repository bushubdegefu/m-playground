package utils

import (
	"fmt"
	"reflect"
)

// Return Unique values in list
func UniqueSlice(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Return Unique values in list
func CheckValueExistsInSlice(slice []string, role_test string) bool {
	for _, role := range slice {
		if role == role_test || role == "superuser" {
			return true
		}
	}
	return false
}

// Struct to Map conversion function
func StructToMap(input any) (map[string]any, error) {
	// Create an empty map
	result := make(map[string]any)

	// Get the reflect value of the struct
	val := reflect.ValueOf(input)

	// Ensure that the input is a pointer to a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check if the input is a struct
	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("input is not a struct")
	}

	// Loop through the struct fields
	for i := 0; i < val.NumField(); i++ {
		// Get the field and its name
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name

		// Insert the field name and value into the map
		result[fieldName] = field.Interface()
	}

	return result, nil
}
