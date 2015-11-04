package cmd

import (
	"fmt"
	"os"
)

func Pull(login bool, args []string) (err error) {
	fmt.Println(args)

	if len(args) != 2 {
		fmt.Println("invalid argument..")
		pullUsage()
	}

	return nil
}

func pullUsage() {
	fmt.Printf("usage: %s pull [[REPO]/[ITEM][:TAG]] [DATAPOOL]\n", os.Args[0])
}
