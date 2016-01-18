package validate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidate(t *testing.T) {

	test := Validate{
		Input: "valid",
		Max:   8,
		Min:   3,
	}

	assert.False(t, test.IsEmpty(), "Should not be empty")

	assert.False(t, test.MaxLength(), "Should not be greater than max")

	assert.False(t, test.MinLength(), "Should be greater than min")

	bad := Validate{
		Input: "",
		Max:   8,
		Min:   3,
	}

	assert.True(t, bad.IsEmpty(), "Should  be empty")

	assert.False(t, bad.MaxLength(), "Should not be greater than max")

	assert.False(t, bad.MinLength(), "Should be less than min")

	long := Validate{
		Input: "waytoolongstring",
		Max:   8,
		Min:   3,
	}

	assert.True(t, long.MaxLength(), "Should not be greater than max")

	assert.False(t, long.MinLength(), "Should be more than min")

	short := Validate{
		Input: "hi",
		Max:   8,
		Min:   3,
	}

	assert.False(t, short.MaxLength(), "Should not be greater than max")

	assert.True(t, short.MinLength(), "Should be less than min")

}

func TestClamp(t *testing.T) {

	max := Clamp(10, 8, 3)

	assert.Equal(t, uint(8), max, "Should be clamped to max value")

	min := Clamp(2, 8, 3)

	assert.Equal(t, uint(3), min, "Should be clamped to min value")

	value := Clamp(6, 8, 3)

	assert.Equal(t, uint(6), value, "Should be actual value")

}
