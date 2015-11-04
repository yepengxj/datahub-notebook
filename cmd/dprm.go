package cmd

import (
	"fmt"
)

func DpRm(login bool, args []string) (err error) {
	fmt.Println("DpRm called", args)
	return nil
}
