package utils

import (
	"database/sql"
)

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// Int32 returns a pointer to the int32 value passed in.
func Int32(v int32) *int32 {
	return &v
}

// StringOrNull returns a sql.NullString from a string pointer.
func StringOrNull(v *string) sql.NullString {
	if v != nil {
		return sql.NullString{String: *v, Valid: true}
	}
	return sql.NullString{String: "", Valid: false}
}

// StringPtr returns a pointer to string from a sql.NullString.
func StringPtr(v *sql.NullString) *string {
	if v.Valid {
		return &v.String
	}
	return nil
}
