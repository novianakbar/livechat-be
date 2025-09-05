package utils

import (
	"database/sql"
	"reflect"
	"strings"
	"time"
)

// StructToMap converts a struct to map[string]interface{} with GORM column mapping
// Skips zero values and handles special types like sql.NullString
func StructToMap(s interface{}) map[string]interface{} {
	val := reflect.ValueOf(s)

	// Dereference if the passed value is a pointer
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Ensure we are now working with a struct
	if val.Kind() != reflect.Struct {
		panic("StructToMap only accepts struct or pointer to struct")
	}

	typ := val.Type()
	result := make(map[string]interface{})

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get column name from GORM tag or use field name
		columnName := getColumnName(fieldType)
		if columnName == "" || columnName == "-" {
			continue
		}

		// Handle different field types
		value := handleFieldValue(field)
		if value != nil {
			result[columnName] = value
		}
	}

	return result
}

// getColumnName extracts column name from GORM tag
func getColumnName(field reflect.StructField) string {
	gormTag := field.Tag.Get("gorm")
	if gormTag == "" {
		return strings.ToLower(field.Name)
	}

	// Parse GORM tag to find column name
	parts := strings.Split(gormTag, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "column:") {
			return strings.TrimPrefix(part, "column:")
		}
	}

	// If no column tag found, use snake_case of field name
	return toSnakeCase(field.Name)
}

// handleFieldValue processes different field types
func handleFieldValue(field reflect.Value) interface{} {
	if !field.IsValid() {
		return nil
	}

	// Handle pointers to basic types (like *string, *int, etc.)
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return nil
		}
		// Dereference the pointer and return the value
		return field.Elem().Interface()
	}

	// Handle sql.NullString
	if field.Type() == reflect.TypeOf(sql.NullString{}) {
		nullStr := field.Interface().(sql.NullString)
		if nullStr.Valid && nullStr.String != "" {
			return nullStr.String
		}
		return nil
	}

	// Handle sql.NullTime
	if field.Type() == reflect.TypeOf(sql.NullTime{}) {
		nullTime := field.Interface().(sql.NullTime)
		if nullTime.Valid {
			return nullTime.Time
		}
		return nil
	}

	// Skip zero values for basic types
	if field.IsZero() {
		return nil
	}

	// Handle time.Time
	if field.Type() == reflect.TypeOf(time.Time{}) {
		t := field.Interface().(time.Time)
		if t.IsZero() {
			return nil
		}
		return t
	}

	// Handle string
	if field.Kind() == reflect.String {
		str := field.String()
		if str == "" {
			return nil
		}
		return str
	}

	// Return the value for other types
	return field.Interface()
}

// toSnakeCase converts PascalCase to snake_case
func toSnakeCase(str string) string {
	// Handle special cases
	switch str {
	case "ID":
		return "id"
	case "UUID":
		return "uuid"
	}

	var result []rune
	for i, r := range str {
		// Add underscore before uppercase letters (except the first character)
		// But not if previous character was also uppercase (to handle abbreviations)
		if i > 0 && r >= 'A' && r <= 'Z' {
			// Check if previous character was also uppercase
			prevRune := []rune(str)[i-1]
			if !(prevRune >= 'A' && prevRune <= 'Z') {
				result = append(result, '_')
			}
		}
		// Convert to lowercase
		if r >= 'A' && r <= 'Z' {
			result = append(result, r-'A'+'a')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// UpdateRequest represents a generic update request
type UpdateRequest struct {
	ID string `json:"id" gorm:"column:id"`
}

// ToUpdateMap converts struct to map excluding ID and other non-updatable fields
func ToUpdateMap(s interface{}, excludeFields ...string) map[string]interface{} {
	updates := StructToMap(s)

	// Default excluded fields
	defaultExcludes := []string{
		"id", "created_at", "deleted_at", "ticket_code", "access_token",
		"created_via", "created_by_id", "updated_by", // Exclude audit fields
	}

	// Combine with user-provided excludes
	allExcludes := append(defaultExcludes, excludeFields...)

	// Remove excluded fields
	for _, field := range allExcludes {
		delete(updates, field)
	}

	// Handle special field mappings for UpdateTicketRequest
	if assignedTo, exists := updates["assignedto"]; exists {
		updates["assigned_to_id"] = assignedTo
		delete(updates, "assignedto")
	}

	// Always add updated_at
	updates["updated_at"] = time.Now()

	return updates
}
