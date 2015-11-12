package cmd

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/asiainfoLDP/datahub/utils"
	"io/ioutil"
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

	User.b64 = base64.StdEncoding.EncodeToString([]byte(User.userName + ":" + User.password))
	//fmt.Printf("%s\n%s:%s\n", User.b64, User.userName, User.password)

	//req.Header.Set("Authorization", "Basic "+os.Getenv("DAEMON_USER_AUTH_INFO"))
	resp, err := commToDaemon("get", "/users/auth", nil)
	if err == nil && resp.StatusCode == 200 {
		if err = os.Setenv("DAEMON_USER_AUTH_TOKEN", resp.Header.Get("Authorization")); err != nil {
			panic(err)
		} else {
			Logged = true
		}
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		//fmt.Println(resp.StatusCode, ShowMsgResp(body, false))
		fmt.Println(resp.StatusCode)
	}

	if resp.StatusCode == 401 {
		return fmt.Errorf(string(body))
	}
	return fmt.Errorf("ERROR %d: login failed.", resp.StatusCode)
}
