package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Cast2Float64(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	t.Run("int64 with factor", func(t *testing.T) {
		expectedResult := 1e-6
		var given int64 = 1
		actualResult := Cast2Float64(given, 1e-6)
		a.Equal(expectedResult, actualResult)
	})
	t.Run("float64 with factor", func(t *testing.T) {
		expectedResult := 1e-6
		var given float64 = 1
		actualResult := Cast2Float64(given, 1e-6)
		a.Equal(expectedResult, actualResult)
	})
	t.Run("int64 without factor", func(t *testing.T) {
		var expectedResult float64 = 1
		var given float64 = 1
		actualResult := Cast2Float64(given, 1)
		a.Equal(expectedResult, actualResult)
	})
}
