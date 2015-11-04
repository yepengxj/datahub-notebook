package cmd

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/asiainfoLDP/datahub-client/utils"
	"os"
)

func Login(login bool, args []string) (err error) {
	fmt.Printf("login: ")
	reader := bufio.NewReader(os.Stdin)
	//loginName, _ := reader.ReadString('\n')
	loginName, _ := reader.ReadBytes('\n')

	loginName = bytes.TrimRight(loginName, "\r\n")
	fmt.Printf("password: ")
	pass := utils.GetPasswd(true) // Silent, for *'s use gopass.GetPasswdMasked()
	//fmt.Printf("[%s]:[%s]\n", string(loginName), string(pass))

	User.userName = string(loginName)
	//User.password = string(pass)
	User.password = fmt.Sprintf("%x", md5.Sum(pass))
	/*
		User.b64 = base64.StdEncoding.EncodeToString([]byte(User.password + ":" + User.password))
		fmt.Printf("%s\n%s:%s\n", User.b64, User.userName, User.password)
	*/
	Logged = true
	return nil
}
