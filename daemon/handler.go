package daemon

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	loginLogged   = false
	loginAuthStr  string
	DefaultServer = "http://10.1.235.98:8080"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := DefaultServer + "/subscriptions/login"
	r.ParseForm()

	if _, ok := r.Header["Authorization"]; !ok {

		if !loginLogged {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if resp.StatusCode == 200 {
		loginAuthStr = r.Header.Get("Authorization")
		loginLogged = true
	}

	w.WriteHeader(resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write(body)
	return
	/*
		if err != nil || (resp != nil && resp.StatusCode != 200) {
			if resp != nil {
				w.WriteHeader(resp.StatusCode)
				fmt.Println("http status code:", resp.StatusCode, err)
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("response Body:", string(body))
			}

			//fmt.Fprintln(w, resp)
			return
		} else {
			loginAuthStr = r.Header.Get("Authorization")
			loginLogged = true
		}
		w.WriteHeader(resp.StatusCode)
	*/
}

func commToServer(method, path string, buffer []byte, w http.ResponseWriter) (resp *http.Response, err error) {

	fmt.Println("daemon: connecting to", DefaultServer+path)
	req, err := http.NewRequest(method, DefaultServer+path, bytes.NewBuffer(buffer))
	if len(loginAuthStr) > 0 {
		req.Header.Set("Authorization", loginAuthStr)
	}

	if resp, err = http.DefaultClient.Do(req); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	w.Write(body)
	return
}
