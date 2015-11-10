package cmd

import (
	"testing"
)

func Testdp(t *testing.T) {
	Dp(false, nil)
	args := []string{"dp1"}
	Dp(false, args)
}
