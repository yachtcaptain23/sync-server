package utils_test

import (
	"database/sql"
	"testing"

	"github.com/brave-experiments/sync-server/utils"
	"github.com/stretchr/testify/assert"
)

func TestStringOrNull(t *testing.T) {
	// Passing nil should returns invalid sql.NullString.
	assert.False(t, utils.StringOrNull(nil).Valid)

	// Passing a normal string pointer should returns a valid sql.NullString.
	v := "test"
	assert.Equal(t, utils.StringOrNull(&v), sql.NullString{String: v, Valid: true})
}

func TestStringPtr(t *testing.T) {
	// Invalid sql.NullString should returns nil.
	nullString := &sql.NullString{String: "", Valid: false}
	assert.Nil(t, utils.StringPtr(nullString))

	// Valid sql.NullString should returns a pointer to the string value.
	v := "test"
	nullString = &sql.NullString{String: v, Valid: true}
	out := utils.StringPtr(nullString)
	assert.NotNil(t, out)
	assert.Equal(t, *out, v)
}
