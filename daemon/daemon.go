package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"database/sql"
	"github.com/asiainfoLDP/datahub-client/cmd"
	"github.com/asiainfoLDP/datahub-client/daemon/daemonigo"
	"github.com/asiainfoLDP/datahub-client/utils/httprouter"
	"github.com/asiainfoLDP/datahub-client/ds"
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

type strc_dp struct{
    	Dpid int
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
	fmt.Fprintf(rw, "Hello HTTP!\n")
}

func stopHttp(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
	fmt.Fprintf(rw, "Hello HTTP!\n")
	sl.Close()
	fmt.Println("connect close")
}

func dpPostOneHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	rw.WriteHeader(http.StatusOK)

	if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		reqJson := cmd.FormatDpCreate{}
		err := json.Unmarshal(result, &reqJson)
		if err != nil {
			fmt.Printf("%T\n%s\n%#v\n", err, err, err)
			fmt.Println(rw, "invalid argument.")
		}
		if len(reqJson.Name) == 0 {
			fmt.Fprintln(rw, "invalid argument.")
		} else {
			msg := &cmd.MsgResp{}
			if err := os.Mkdir(reqJson.Name, 0755); err != nil {
				msg.Msg = err.Error()
			} else {
				msg.Msg = "OK."
				sql_dp_insert := fmt.Sprintf(`insert into DH_DP (DPID, DPNAME, DPTYPE, DPCONN, STATUS)
					values (null, '%s', '%s', '%s', 'A')`, reqJson.Name, reqJson.Type, reqJson.Conn)
			    if _,err := g_ds.Insert(sql_dp_insert); err != nil {
			    	os.Remove(reqJson.Name)
			    	msg.Msg = err.Error()
			    }
			}
			resp, _ := json.Marshal(msg)
			respStr := string(resp)
			fmt.Fprintln(rw, respStr)
		}

	}

}

func dpGetAllHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	rw.WriteHeader(http.StatusOK)

	result, _ := ioutil.ReadAll(r.Body)
	reqJson := cmd.FormatDp{}
    err := json.Unmarshal(result, &reqJson)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		fmt.Println(rw, "invalid argument.")
	}

	msg := &cmd.MsgResp{}

	msg.Msg = "OK."
	wrtocli := cmd.FormatDp{}
	sql_dp := fmt.Sprintf(`SELECT DPNAME, DPTYPE FROM DH_DP WHERE STATUS = 'A'`)
	rows,err := g_ds.QueryRows(sql_dp)
	if err != nil {
		 msg.Msg = err.Error()
	}
	defer rows.Close()
	bresultflag := false
	for rows.Next() {
		bresultflag = true
		rows.Scan(&wrtocli.Name, &wrtocli.Type)
		resp, _ := json.Marshal(wrtocli)
	    respStr := string(resp)
	    fmt.Fprintln(rw, respStr)
	}
	if bresultflag == false{
        msg.Msg = "There isn't any datapool."
		resp, _ := json.Marshal(msg)
	    respStr := string(resp)
	    fmt.Fprintln(rw, respStr)
	}

}

func dpGetOneHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	rw.WriteHeader(http.StatusOK)
	dpname := ps.ByName("dpname")

    //In future, we need to get dptype in Json to surpport FILE\ DB\ SDK\ API datapool
	result, _ := ioutil.ReadAll(r.Body)
	reqJson := cmd.FormatDp{}
    err := json.Unmarshal(result, &reqJson)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
		fmt.Println(rw, "invalid argument.")
	}

	msg := &cmd.MsgResp{}
	msg.Msg = "OK."

    sql_total := fmt.Sprintf(`SELECT COUNT(*) FROM DH_DP 
		WHERE STATUS = 'A' AND DPNAME = '%s'`, string(dpname))
    row,err := g_ds.QueryRow(sql_total)
    if err != nil{
    	msg.Msg = err.Error()
    	resp, _ := json.Marshal(msg)
	    respStr := string(resp)
	    fmt.Fprintln(rw, respStr)
	    return
    }
    var total int
    row.Scan(&total)
    if total == 0 {
    	msg.Msg = fmt.Sprintf("Datapool %v not found.", dpname)
    	resp, _ := json.Marshal(msg)
	    respStr := string(resp)
	    fmt.Fprintln(rw, respStr)
	    return
    }

	sql_dp := fmt.Sprintf(`SELECT DPID, DPNAME, DPTYPE, DPCONN FROM DH_DP 
		WHERE STATUS = 'A' AND DPNAME = '%s'`, string(dpname))
	rows,err := g_ds.QueryRows(sql_dp)
	if err != nil {
		 msg.Msg = err.Error()
	}
	defer rows.Close()
	var dpid int
	//I use queryrows because  there will be distinct dataitems in the same datapool
	for rows.Next() {
		onedp := cmd.FormatDp_dpname{}
	    onedp.Items = make([]cmd.Item, 0,16)
		rows.Scan(&dpid, &onedp.Name, &onedp.Type, &onedp.Conn)
		if dpid > 0 {
			//Use "left out join" to get repository/dataitem records, whether it has tags or not.
			sql_tag := fmt.Sprintf(`SELECT A.REPOSITORY, A.DATAITEM, B.TAGNAME, strftime(B.CREATE_TIME), A.PUBLISH 
				FROM DH_DP_RPDM_MAP A LEFT JOIN DH_RPDM_TAG_MAP B
				ON (A.RPDMID = B.RPDMID)
				WHERE A.DPID = %v`, dpid)
			tagrows,err := g_ds.QueryRows(sql_tag)
			if err != nil {
		        msg.Msg = err.Error()
		    }
		    
	        fmt.Println(msg)
		    defer tagrows.Close()
		    for tagrows.Next(){
		    	fmt.Println(tagrows)
		    	item := cmd.Item{}
		    	tagrows.Scan(&item.Repository, &item.DataItem, &item.Tag, &item.Time, &item.Publish)
                onedp.Items = append(onedp.Items, item)
		    }
	    }
		
		resp, _ := json.Marshal(onedp)
	    respStr := string(resp)
	    fmt.Fprintln(rw, respStr)
	}

}

func dpDeleteOneHandler(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	rw.WriteHeader(http.StatusOK)

	dpname := ps.ByName("dpname")
	msg := &cmd.MsgResp{}

	sql_dp_rm := fmt.Sprintf(`SELECT DPID, DPTYPE FROM DH_DP WHERE DPNAME ='%s'`, dpname)
	dprows, err := g_ds.QueryRows(sql_dp_rm)
	if err != nil {
		 msg.Msg = err.Error()
	}
    
    bresultflag := false

    dpid_type := make([]strc_dp, 0, 8)
    strcone := strc_dp{}
    for dprows.Next(){
        dprows.Scan(&strcone.Dpid, &strcone.Dptype)
        dpid_type = append(dpid_type, strcone)
    }
    dprows.Close()

    for _, v := range dpid_type {
    	var dpid = v.Dpid
    	var dptype = v.Dptype
        bresultflag = true
    	//dprow.Scan(&dpid, &dptype)
    	sql_dp_item := fmt.Sprintf("SELECT PUBLISH FROM DH_DP_RPDM_MAP WHERE DPID = %v ", dpid)
    	row, err := g_ds.QueryRow(sql_dp_item)
    	if err != nil {
		     msg.Msg = err.Error()
	    }
	    //time.Sleep(60*time.Second)
	    var sPublish string
	    row.Scan(&sPublish)
	    if sPublish == "Y" {
	    	msg.Msg = fmt.Sprintf(`Datapool %s with type:%s can't be removed , it contains published DataItem !`, dpname, dptype)
	    }else{
	    	sql_update := fmt.Sprintf("UPDATE DH_DP SET STATUS = 'N' WHERE DPID = %v", dpid)
	    	_, err := g_ds.Update(sql_update)
	    	if err != nil {
		        msg.Msg = err.Error()
	        }else {
	            msg.Msg = fmt.Sprintf("Datapool %s with type:%s removed successfully!", dpname, dptype)
	        }
	    }
	    resp, _ := json.Marshal(msg)
	    respStr := string(resp)
	    fmt.Fprintln(rw, respStr)
    }
    if bresultflag == false {
    	msg.Msg = fmt.Sprintf("Datapool %s not found.\n", dpname)
    	resp, _ := json.Marshal(msg)
		respStr := string(resp)
		fmt.Fprintln(rw, respStr)
    }
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
	router.POST("/datapool", dpPostOneHandler)
	router.GET("/datapool", dpGetAllHandler)
	router.GET("/datapool/:dpname", dpGetOneHandler)
	router.DELETE("/datapool/:dpname", dpDeleteOneHandler)

	http.Handle("/", router)

	//http.HandleFunc("/", helloHttp)
	http.HandleFunc("/stop", stopHttp)
	//http.HandleFunc("/datapool", dpHttp)
	http.HandleFunc("/Repository", repoHandler)
	//http.HandleFunc("/subscriptions", subHttp)

	server := http.Server{}

	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		server.Serve(sl)
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
