package utils

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// Int32 returns a pointer to the int32 value passed in.
func Int32(v int32) *int32 {
	return &v
}
