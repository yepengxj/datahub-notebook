package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
	"io/ioutil"
	"net/http"
)

var (
	loginLogged  = false
	loginAuthStr string
)

func subsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Path)

}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := "http://10.1.51.32:8080/subscriptions/login"
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

func repoHandler(rw http.ResponseWriter, req *http.Request) {
	url := "http://10.1.235.96:8080/DataItem/"
	req.ParseForm()

	d := ds.FormatRepoList{}
	reqData, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(reqData, &d)

	if len(d.ItemID) > 0 {
		url = "http://10.1.235.96:8080/Subscriptions/" + d.ItemID + "?show_tag=1"
	}

	req, err := http.NewRequest(req.Method, url, req.Body)

	resp, err := http.DefaultClient.Do(req)

	if err != nil || (resp != nil && resp.StatusCode != 200) {
		if resp != nil {
			fmt.Println("http status code:", resp.StatusCode, err)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("response Body:", string(body))
		}

		fmt.Fprintln(rw, resp)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	fmt.Fprintln(rw, string(body))

}
