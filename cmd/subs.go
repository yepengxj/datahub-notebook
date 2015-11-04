package cmd

import (
	"fmt"
)

func Subs(login bool, args []string) (err error) {
	fmt.Println("Subs called", args)
	return nil
}
