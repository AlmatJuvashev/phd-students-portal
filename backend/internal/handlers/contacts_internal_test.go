package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactsHelpers(t *testing.T) {
	t.Run("contactNullableString", func(t *testing.T) {
		assert.Nil(t, contactNullableString(""))
		assert.Nil(t, contactNullableString("   "))
		assert.Equal(t, "val", contactNullableString("val"))
	})

	t.Run("contactNullablePtr", func(t *testing.T) {
		assert.Nil(t, contactNullablePtr(nil))
		s := ""
		assert.Nil(t, contactNullablePtr(&s))
		s = "  "
		assert.Nil(t, contactNullablePtr(&s))
		s = "val"
		assert.Equal(t, "val", contactNullablePtr(&s))
	})

	t.Run("toJSON", func(t *testing.T) {
		assert.Nil(t, toJSON(nil))
		assert.Nil(t, toJSON(map[string]string{}))
		
		m := map[string]string{"k": "v"}
		j := toJSON(m)
		assert.JSONEq(t, `{"k":"v"}`, string(j.([]byte)))
	})
}
