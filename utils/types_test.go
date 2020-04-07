package utils

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringOrNull(t *testing.T) {
	// Passing nil should returns invalid sql.NullString.
	assert.False(t, StringOrNull(nil).Valid)

	// Passing a normal string pointer should returns a valid sql.NullString.
	v := "test"
	assert.Equal(t, StringOrNull(&v), sql.NullString{String: v, Valid: true})
}

func TestStringPtr(t *testing.T) {
	// Invalid sql.NullString should returns nil.
	nullString := &sql.NullString{String: "", Valid: false}
	assert.Nil(t, StringPtr(nullString))

	// Valid sql.NullString should returns a pointer to the string value.
	v := "test"
	nullString = &sql.NullString{String: v, Valid: true}
	out := StringPtr(nullString)
	assert.NotNil(t, out)
	assert.Equal(t, *out, v)
}
