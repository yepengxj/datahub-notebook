package cmd

import (
	"fmt"
<<<<<<< HEAD
=======
	"os"
>>>>>>> master
)

func Pull(login bool, args []string) (err error) {
	fmt.Println(args)
<<<<<<< HEAD
	return nil
}
=======

	if len(args) != 2 {
		fmt.Println("invalid argument..")
		pullUsage()
	}

	return nil
}

func pullUsage() {
	fmt.Printf("usage: %s pull [[REPO]/[ITEM][:TAG]] [DATAPOOL]\n", os.Args[0])
}
>>>>>>> master
