package daemon

import (
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
)

func repoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println(r.URL.Path, "(repo)")
	reqBody, _ := ioutil.ReadAll(r.Body)
	commToServer("get", r.URL.Path, reqBody, w)

	return

}
func repoDetailHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println(r.URL.Path, "(subsdetail)")
	reqBody, _ := ioutil.ReadAll(r.Body)
	commToServer("get", r.URL.Path, reqBody, w)

	return

}
