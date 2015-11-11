package cmd

import (
	/*
		"bufio"
		"bytes"
		"crypto/md5"
		"encoding/base64"
		"errors"
		"fmt"
		"github.com/asiainfoLDP/datahub/utils"
		"io/ioutil"
		"os"
	*/
	"fmt"
)

func Help(login bool, args []string) (err error) {

	for _, v := range Cmd {
		fmt.Printf("%-16s  %s\n", v.Name, v.Desc)
	}
	return nil
}
