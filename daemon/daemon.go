package daemon

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asiainfoLDP/datahub/cmd"
	"github.com/asiainfoLDP/datahub/daemon/daemonigo"
	"github.com/asiainfoLDP/datahub/ds"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var g_ds = new(ds.Ds)

const g_dbfile string = "/var/run/datahub.db"

const g_strDpPath string = "/var/lib/datahub/"

type StoppableListener struct {
	*net.UnixListener          //Wrapped listener
	stop              chan int //Channel used only to indicate listener should shutdown
}

type strc_dp struct {
	Dpid   int
	Dptype string
}

func dbinit() {
	//log.Println("connect to db sqlite3")
	db, err := sql.Open("sqlite3", g_dbfile)
	//defer db.Close()
	chk(err)
	g_ds.Db = db

	g_ds.Create(ds.Create_dh_dp)
	g_ds.Create(ds.Create_dh_dp_repo_ditem_map)
	g_ds.Create(ds.Create_dh_repo_ditem_tag_map)

}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}
func get(err error) {
	if err != nil {
		log.Println(err)
	}
}

func New(l net.Listener) (*StoppableListener, error) {
	tcpL, ok := l.(*net.UnixListener)

	if !ok {
		return nil, errors.New("Cannot wrap listener")
	}

	retval := &StoppableListener{}
	retval.UnixListener = tcpL
	retval.stop = make(chan int)

	return retval, nil
}

var StoppedError = errors.New("Listener stopped")
var sl = new(StoppableListener)

func (sl *StoppableListener) Accept() (net.Conn, error) {

	for {
		//Wait up to one second for a new connection
		sl.SetDeadline(time.Now().Add(time.Second))

		newConn, err := sl.UnixListener.Accept()

		//Check for the channel being closed
		select {
		case <-sl.stop:
			return nil, StoppedError
		default:
			//If the channel is still open, continue as normal
		}

		if err != nil {
			netErr, ok := err.(net.Error)

			//If this is a timeout, then continue to wait for
			//new connections
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
		}

		return newConn, err
	}
}

func (sl *StoppableListener) Stop() {
	close(sl.stop)
}

func helloHttp(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	rw.WriteHeader(http.StatusOK)
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Fprintf(rw, "%s Hello HTTP!\n", req.URL.Path)
	fmt.Fprintf(rw, "%s \n", string(body))
}

func stopHttp(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Hello HTTP!\n")
	sl.Close()
	fmt.Println("connect close")
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}
	panic("not reached")
}
func isFileExists(file string) bool {
	fi, err := os.Stat(file)
	if err == nil {
		fmt.Println("exist", file)
		return !fi.IsDir()
	}
	return os.IsExist(err)
}

func RunDaemon() {
	fmt.Println("run daemon..")
	// Daemonizing echo server application.
	switch isDaemon, err := daemonigo.Daemonize(); {
	case !isDaemon:
		return
	case err != nil:
		log.Fatalf("main(): could not start daemon, reason -> %s", err.Error())
	}
	//fmt.Printf("server := http.Server{}\n")

	dbinit()

	if false == isDirExists(g_strDpPath) {
		err := os.MkdirAll(g_strDpPath, 0755)
		if err != nil {
			fmt.Printf("mkdir %s error! %v ", g_strDpPath, err)
		}

	}
	os.Chdir(g_strDpPath)
	originalListener, err := net.Listen("unix", cmd.UnixSock)
	if err != nil {
		panic(err)
	}

	sl, err = New(originalListener)
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	router.GET("/", helloHttp)
	router.POST("/datapools", dpPostOneHandler)
	router.GET("/datapools", dpGetAllHandler)
	router.GET("/datapools/:dpname", dpGetOneHandler)
	router.DELETE("/datapools/:dpname", dpDeleteOneHandler)

	router.GET("/subscriptions/:repo/:item", subsDetailHandler)
	router.GET("/subscriptions", subsHandler)
	router.POST("/subscriptions/:repo/:item/pull", pullHandler)

	http.Handle("/", router)
	http.HandleFunc("/stop", stopHttp)
	http.HandleFunc("/Repository", repoHandler)
	http.HandleFunc("/login", loginHandler)

	server := http.Server{}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		server.Serve(sl)
	}()

	//p2p server
	router_p2p := httprouter.New()
	router_p2p.GET("/", sayhello)
	router_p2p.GET("/pull/:repo/:dataitem/:tag", p2p_pull)
	go func() {
		wg.Add(1)
		defer wg.Done()
		http.ListenAndServe(":35800", router_p2p)
	}()

	fmt.Printf("Serving HTTP\n")
	select {
	case signal := <-stop:
		fmt.Printf("Got signal:%v\n", signal)
	}
	fmt.Printf("Stopping listener\n")
	sl.Stop()
	fmt.Printf("Waiting on server\n")
	wg.Wait()
	daemonigo.UnlockPidFile()
	g_ds.Db.Close()

}

/*pull parses filename and target IP from HTTP GET method, and start downloading routine. */
func p2p_pull(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("p2p pull...")
	r.ParseForm()
	//file := r.Form.Get("file")
	sRepoName := ps.ByName("repo")
	sDataItem := ps.ByName("dataitem")
	sTag := ps.ByName("tag")
	fmt.Println(sRepoName, sDataItem, sTag)
	var irpdmid, idpid int
	var stagdetail, sdpname, sdpconn string
	msg := &ds.MsgResp{}
	msg.Msg = "OK."

	sSqlGetRpdmidDpid := fmt.Sprintf(`SELECT DPID, RPDMID FROM DH_DP_RPDM_MAP 
    	WHERE REPOSITORY = '%s' AND DATAITEM = '%s'`, sRepoName, sDataItem)
	row, err := g_ds.QueryRow(sSqlGetRpdmidDpid)
	if err != nil {
		msg.Msg = err.Error()
	}
	row.Scan(&idpid, &irpdmid)
	fmt.Println("dpid", idpid, "rpdmid", irpdmid)

	sSqlGetTagDetail := fmt.Sprintf(`SELECT DETAIL FROM DH_RPDM_TAG_MAP 
        WHERE RPDMID = '%d' AND TAGNAME = '%s'`, irpdmid, sTag)
	tagrow, err := g_ds.QueryRow(sSqlGetTagDetail)
	if err != nil {
		msg.Msg = err.Error()
	}
	tagrow.Scan(&stagdetail)
	fmt.Println("tagdetail", stagdetail)

	sSqlGetDpconn := fmt.Sprintf(`SELECT DPNAME, DPCONN FROM DH_DP WHERE DPID='%d'`, idpid)
	dprow, err := g_ds.QueryRow(sSqlGetDpconn)
	if err != nil {
		msg.Msg = err.Error()
	}
	dprow.Scan(&sdpname, &sdpconn)
	fmt.Println("dpname", sdpname, "dpconn", sdpconn)

	filepathname := "/" + sdpconn + "/" + sdpname + "/" + sRepoName + "/" + sDataItem + "/" + stagdetail
	fmt.Println(" filename:", filepathname)
	if exists := isFileExists(filepathname); !exists {
		filepathname = "/" + sdpconn + "/" + stagdetail
		if exists := isFileExists(filepathname); !exists {
			filepathname = "/var/lib/datahub/" + sTag
			if exists := isFileExists(filepathname); !exists {
				fmt.Println(" filename:", filepathname)
				//http.NotFound(rw, r)
				msg.Msg = "tag not found"
				resp, _ := json.Marshal(msg)
				respStr := string(resp)
				fmt.Fprintln(rw, respStr)
				return
			}
		}
	}
	//rw.Header().Set("Content-Type", "file")
	http.ServeFile(rw, r, filepathname)

	resp, _ := json.Marshal(msg)
	respStr := string(resp)
	fmt.Fprintln(rw, respStr)
	return
}

func sayhello(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	rw.WriteHeader(http.StatusOK)
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Fprintf(rw, "%s Hello p2p HTTP !\n", req.URL.Path)
	fmt.Fprintf(rw, "%s \n", string(body))
}
