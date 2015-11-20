package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var (
	loginLogged   = false
	loginAuthStr  string
	gstrUsername  string
	DefaultServer = "http://54.223.58.0:8888"
)

type UserForJson struct {
	Username string `json:"username", omitempty`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := DefaultServer + "/" //r.URL.Path
	r.ParseForm()

	if _, ok := r.Header["Authorization"]; !ok {

		if !loginLogged {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	userjsonbody, _ := ioutil.ReadAll(r.Body)
	userforjson := UserForJson{}
	if err := json.Unmarshal(userjsonbody, &userforjson); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	gstrUsername = userforjson.Username

	log.Println("login to", url, "Authorization:", r.Header.Get("Authorization"))
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	log.Println("login return", resp.StatusCode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		type tk struct {
			Token string `json:"token"`
		}
		token := &tk{}
		if err = json.Unmarshal(body, token); err != nil {
			panic(err)
			w.WriteHeader(resp.StatusCode)
			w.Write(body)
			fmt.Println(resp.StatusCode, string(body))
			return
		} else {
			loginAuthStr = "Token " + token.Token
			loginLogged = true
			log.Println(loginAuthStr)
		}
	}
	w.WriteHeader(resp.StatusCode)
}

func commToServer(method, path string, buffer []byte, w http.ResponseWriter) (resp *http.Response, err error) {

	log.Println("daemon: connecting to", DefaultServer+path)
	req, err := http.NewRequest(strings.ToUpper(method), DefaultServer+path, bytes.NewBuffer(buffer))
	if len(loginAuthStr) > 0 {
		req.Header.Set("Authorization", loginAuthStr)
	}

	//req.Header.Set("User", "admin")
	if resp, err = http.DefaultClient.Do(req); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write(body)
	fmt.Println(resp.StatusCode, string(body))
	return
}
