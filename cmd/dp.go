package cmd

import (
	"fmt"
)

func Dp(login bool, args []string) (err error) {
	fmt.Println("Dp called", args)
	return nil
}
