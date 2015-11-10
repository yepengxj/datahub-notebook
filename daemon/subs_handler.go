package daemon

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

func subsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println(r.URL.Path, "(subs)")
	reqBody, _ := ioutil.ReadAll(r.Body)
	commToServer("get", r.URL.Path, reqBody, w)

	return

}

func subsDetailHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println(r.URL.Path, "(subsdetail)")
	reqBody, _ := ioutil.ReadAll(r.Body)
	commToServer("get", r.URL.Path, reqBody, w)

	return

}
