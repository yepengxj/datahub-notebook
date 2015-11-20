package daemon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Beatbody struct {
	Daemonid   string   `json:"daemonid ,omitempty"`
	EntryPoint []string `json:"entrypoint, omitempty"`
}

func HeartBeat() {
	for {
		heartbeatbody := Beatbody{}
		heartbeatbody.EntryPoint = append(heartbeatbody.EntryPoint, DefaultServer)
		jsondata, err := json.Marshal(heartbeatbody)
		url := DefaultServer + "/heartbeats/"
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsondata))
		if len(loginAuthStr) > 0 {
			req.Header.Set("Authorization", loginAuthStr)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("HeartBeat http statuscode:%v,  http body:%s\n", resp.StatusCode, body)

		time.Sleep(5 * time.Second)
	}
}
