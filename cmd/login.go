package cmd

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/utils"
	"os"
)

type UserForJson struct {
	Username string `json:"username", omitempty`
}

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
	userJson := UserForJson{Username: User.userName}
	jsondata, _ := json.Marshal(userJson)

	resp, err := commToDaemon("get", "/", jsondata) //users/auth
	if err != nil {

		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		Logged = true
		return
	} else {
		return fmt.Errorf("ERROR %d: login failed.", resp.StatusCode)
	}
	/*
		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			//fmt.Println(resp.StatusCode, ShowMsgResp(body, false))
			fmt.Println(resp.StatusCode)
		}

		if resp.StatusCode == 401 {
			return fmt.Errorf(string(body))
		}
		return fmt.Errorf("ERROR %d: login failed.", resp.StatusCode)
	*/
}
