package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/asiainfoLDP/datahub-client/ds"
	"io/ioutil"
	"net/http"
)

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
